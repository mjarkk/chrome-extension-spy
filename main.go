package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	output, err := getChromeLocation()
	printErr(err)
	fullpath := chromeLoc(output)
	extensions, fullExtension := getExtensions(fullpath)
	_, ext, fullExt := selectExtensionToUse(extensions, fullExtension)
	// create a temp dir to store the extension
	tempDir, err := ioutil.TempDir("", "chrome-extension-spy")
	defer os.RemoveAll(tempDir)
	printErr(err)
	err = copyFullExtension(ext.fullPkgURL, tempDir, []string{})
	printErr(err)
	err = editExtension(tempDir, ext, fullExt)
	printErr(err)

	var waitForGinAndChrome sync.WaitGroup
	waitForGinAndChrome.Add(2)

	var waitForExit = make(chan struct{})

	var ginErr error
	go func() {
		ginErr = startWebServer(waitForExit)
		waitForGinAndChrome.Done()
	}()
	var chromeErr error
	go func() {
		chromeErr = launchChrome(tempDir, waitForExit)
		waitForGinAndChrome.Done()
	}()
	go func() {
		waitForExitInput()
		waitForExit <- struct{}{}
	}()

	waitForGinAndChrome.Wait()
	printErr(chromeErr)
	printErr(ginErr)
}

func waitForExitInput() {
	var input string
	fmt.Print("Type exit to exit the program")
	fmt.Scanf("%s", &input)
	if input != "exit" {
		waitForExitInput()
	}
}

func launchChrome(extPath string, closeChrome chan struct{}) error {
	tempDir, err := ioutil.TempDir("", "chrome-data")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)
	var cmd *exec.Cmd
	go func() {
		<-closeChrome
		fmt.Println("killing chrome process")
		cmd.Process.Kill()
	}()
	cmd = exec.Command(
		"google-chrome-unstable",
		"--user-data-dir="+tempDir,
		"--disable-background-networking",
		"--disable-background-timer-throttling",
		"--disable-backgrounding-occluded-windows",
		"--disable-breakpad",
		"--disable-client-side-phishing-detection",
		"--disable-default-apps",
		"--disable-dev-shm-usage",
		"--disable-extensions",
		"--disable-features=site-per-process",
		"--disable-hang-monitor",
		"--disable-ipc-flooding-protection",
		"--disable-popup-blocking",
		"--disable-prompt-on-repost",
		"--disable-renderer-backgrounding",
		"--disable-sync",
		"--disable-translate",
		"--metrics-recording-only",
		"--no-first-run",
		"--safebrowsing-disable-auto-update",
		"--enable-automation",
		"--password-store=basic",
		"--use-mock-keychain",
	)
	_, err = cmd.Output()
	return err
}

func editExtension(tmpDir string, ext chromeExtension, fullExt extensionManifest) error {
	thisFileDir, err := os.Executable()
	if err != nil {
		return err
	}
	injectable, err := ioutil.ReadFile(path.Join(thisFileDir, "../web_static/extension_inject.js"))
	if err != nil {
		return err
	}
	for _, srcItem := range fullExt.Background.Scripts {
		fullFileDir := path.Join(tmpDir, srcItem)
		file, err := ioutil.ReadFile(fullFileDir)
		if err != nil {
			return err
		}
		toWrite := string(injectable) + string(file)
		ioutil.WriteFile(fullFileDir, []byte(toWrite), 0666)
	}
	return nil
}

func printErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func proxyHandeler(c *gin.Context, reqType string) {
	rawURL := c.Param("url")
	parsedURL, err := url.PathUnescape(rawURL)
	if err != nil {
		c.String(http.StatusConflict, "")
	}

	hc := http.Client{}
	req, err := http.NewRequest(reqType, parsedURL, nil)

	if reqType == "POST" {
		req.Body = c.Request.Body
	}

	for key, value := range c.Request.Header {
		req.Header.Add(key, value[0])
	}

	rs, err := hc.Do(req)
	if err != nil {
		c.String(400, "")
		return
	}

	for key, item := range rs.Header {
		c.Header(key, item[0])
	}

	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		body = []byte("")
	}

	c.Data(rs.StatusCode, rs.Header.Get("Content-Type"), body)
}

func proxyHandelerPost(c *gin.Context) {
	proxyHandeler(c, "POST")
}

func proxyHandelerGet(c *gin.Context) {
	proxyHandeler(c, "GET")
}

func startWebServer(closeServer chan struct{}) error {
	gin.SetMode("release")
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))
	r.GET("/proxy/:url", proxyHandelerGet)
	r.POST("/proxy/:url", proxyHandelerPost)
	r.Static("/web_static", "./web_static")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	<-closeServer
	fmt.Println("Stopping server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	return nil
}

func copyFullExtension(baseDir string, tempDir string, extensionDir []string) error {
	extensionDirPath := strings.Join(extensionDir, "/")
	fullExtensionDirPath := path.Join(baseDir, extensionDirPath)
	files, err := ioutil.ReadDir(fullExtensionDirPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		name := file.Name()

		if file.IsDir() {
			// create a dir and loop over that dir
			os.MkdirAll(path.Join(tempDir, extensionDirPath, name), 0777)
			copyFullExtension(baseDir, tempDir, append(extensionDir, name))
		} else {
			// copy a file over
			from, err := os.Open(path.Join(fullExtensionDirPath, name))
			if err != nil {
				return err
			}
			to, err := os.Create(path.Join(tempDir, extensionDirPath, name))
			if err != nil {
				return err
			}
			_, err = io.Copy(to, from)
			if err != nil {
				return err
			}
			from.Close()
			to.Close()
		}
	}
	return nil
}

func selectExtensionToUse(exts []chromeExtension, fullExts []extensionManifest) (int64, chromeExtension, extensionManifest) {
	printExtensions(exts)
	fmt.Println("------------------------------")
	fmt.Println("Type the id you want to spy on")
	i := askForNum(int64(len(exts) - 1))
	return i, exts[i], fullExts[i]
}

func askForNum(max int64) int64 {
	var input string
	fmt.Print("> ")
	fmt.Scanf("%s", &input)
	i, err := strconv.ParseInt(input, 10, 64)
	if err != nil || i > max {
		fmt.Println("Not a valid input")
		i = askForNum(max)
	}
	return i
}

func printExtensions(exts []chromeExtension) {
	maxNameLen := 0
	maxShortNameLen := 0
	maxPkgVersionLen := 7
	for _, ext := range exts {
		if len(ext.name) > maxNameLen {
			maxNameLen = len(ext.name)
		}
		if len(ext.shortName) > maxShortNameLen {
			maxShortNameLen = len(ext.name)
		}
		if len(ext.pkgVersion) > maxPkgVersionLen {
			maxPkgVersionLen = len(ext.pkgVersion)
		}
	}
	fmt.Printf("%s\t%s%s%s%s\n", "id", rightPad("name", " ", maxNameLen+1), rightPad("short name", " ", maxShortNameLen+1), rightPad("version", " ", maxPkgVersionLen+1), "homepage")
	for id, ext := range exts {
		name := ext.name
		if len(name) == 0 {
			name = "-"
		}
		shortName := ext.shortName
		if len(shortName) == 0 {
			shortName = "-"
		}
		homepageURL := ext.homepageURL
		if len(homepageURL) == 0 {
			homepageURL = "-"
		}
		pkgVersion := ext.pkgVersion
		if len(pkgVersion) == 0 {
			pkgVersion = "-"
		}
		fmt.Printf(
			"%v\t%s%s%s%s\n",
			id,
			rightPad(name, " ", maxNameLen+1),
			rightPad(shortName, " ", maxShortNameLen+1),
			rightPad(pkgVersion, " ", maxPkgVersionLen+1),
			homepageURL,
		)
	}
}

type chromeExtension struct {
	pkg         string // the chrome extension folder name
	pkgVersion  string // the package version
	fullPkgURL  string // full url to extension
	name        string // extension name
	shortName   string // extension longname
	homepageURL string // homepage url
}

func getExtensions(extensionsPath string) ([]chromeExtension, []extensionManifest) {
	toReturn := []chromeExtension{}
	toReturnFull := []extensionManifest{}
	files, err := ioutil.ReadDir(extensionsPath)
	if err != nil {
		return toReturn, toReturnFull
	}
	for _, f := range files {
		fName := f.Name()
		if len(fName) == 32 {
			extensionPath := path.Join(extensionsPath, fName)
			files, err := ioutil.ReadDir(extensionPath)
			if err != nil {
				return toReturn, toReturnFull
			}
			version := ""
			for _, versionDir := range files {
				version = versionDir.Name()
			}
			dat, err := ioutil.ReadFile(path.Join(extensionPath, version, "/manifest.json"))
			if err == nil {
				var manifest extensionManifest
				var addToReturnValue chromeExtension
				json.Unmarshal(dat, &manifest)
				addToReturnValue.name = manifest.Name
				addToReturnValue.homepageURL = manifest.HomepageURL
				addToReturnValue.pkg = fName
				addToReturnValue.pkgVersion = version
				addToReturnValue.shortName = manifest.ShortName
				addToReturnValue.fullPkgURL = path.Join(extensionPath, version, "/")
				toReturn = append(toReturn, addToReturnValue)
				toReturnFull = append(toReturnFull, manifest)
			}
		}
	}
	return toReturn, toReturnFull
}

func chromeLoc(version string) string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("can't get home directory")
		os.Exit(1)
	}
	return path.Join(usr.HomeDir, "/.config/", version, "/Default/Extensions")
}

func getChromeLocation() (string, error) {
	if _, err := os.Stat(chromeLoc("chromium")); !os.IsNotExist(err) {
		return "chromium", nil
	}
	if _, err := os.Stat(chromeLoc("google-chrome")); !os.IsNotExist(err) {
		return "google-chrome", nil
	}
	if _, err := os.Stat(chromeLoc("google-chrome-beta")); !os.IsNotExist(err) {
		return "google-chrome-beta", nil
	}
	if _, err := os.Stat(chromeLoc("google-chrome-dev")); !os.IsNotExist(err) {
		return "google-chrome-dev", nil
	}
	if _, err := os.Stat(chromeLoc("google-chrome-unstable")); !os.IsNotExist(err) {
		return "google-chrome-unstable", nil
	}
	if _, err := os.Stat(chromeLoc("google-chrome-canary")); !os.IsNotExist(err) {
		return "google-chrome-canary", nil
	}
	return "", errors.New("Chrome location not found")
}

type extensionManifest struct {
	Applications struct {
		Gecko struct {
			ID               string `json:"id"`
			StrictMinVersion string `json:"strict_min_version"`
		} `json:"gecko"`
	} `json:"applications"`
	App struct {
		Launch struct {
			Container string `json:"container"`
			WebURL    string `json:"web_url"`
		} `json:"launch"`
		WebContent struct {
			Enabled bool   `json:"enabled"`
			Origin  string `json:"origin"`
		} `json:"web_content"`
		Urls []string `json:"urls"`
	} `json:"app"`
	Commands struct {
		ExecutePageAction struct {
			SuggestedKey struct {
				Default string `json:"default"`
			} `json:"suggested_key"`
		} `json:"_execute_page_action"`
		DevtoolsBottom struct {
			Description string `json:"description"`
		} `json:"devtools-bottom"`
		DevtoolsLeft struct {
			Description string `json:"description"`
		} `json:"devtools-left"`
		DevtoolsRemote struct {
			Description string `json:"description"`
		} `json:"devtools-remote"`
		DevtoolsRight struct {
			Description string `json:"description"`
		} `json:"devtools-right"`
	} `json:"commands"`
	Background struct {
		Persistent    bool     `json:"persistent"`
		AllowJsAccess bool     `json:"allow_js_access"`
		Scripts       []string `json:"scripts"`
	} `json:"background"`
	ChromeURLOverrides struct {
		Newtab string `json:"newtab"`
	} `json:"chrome_url_overrides"`
	DisplayInLauncher   bool `json:"display_in_launcher"`
	DisplayInNewTabPage bool `json:"display_in_new_tab_page"`
	BrowserAction       struct {
		DefaultIcon  string `json:"default_icon"`
		DefaultTitle string `json:"default_title"`
		DefaultPopup string `json:"default_popup"`
	} `json:"browser_action"`
	ContentScripts []struct {
		AllFrames       bool     `json:"all_frames"`
		CSS             []string `json:"css"`
		Js              []string `json:"js"`
		MatchAboutBlank bool     `json:"match_about_blank"`
		ExcludeGlobs    []string `json:"exclude_globs"`
		Matches         []string `json:"matches"`
		RunAt           string   `json:"run_at"`
	} `json:"content_scripts"`
	ContentSecurityPolicy string `json:"content_security_policy"`
	Description           string `json:"description"`
	Icons                 struct {
		Num16  string `json:"16"`
		Num32  string `json:"32"`
		Num48  string `json:"48"`
		Num64  string `json:"64"`
		Num128 string `json:"128"`
	} `json:"icons"`
	Key                   string      `json:"key"`
	Author                interface{} `json:"author"`
	Incognito             string      `json:"incognito"`
	DevtoolsPage          string      `json:"devtools_page"`
	HomepageURL           string      `json:"homepage_url"`
	ManifestVersion       int         `json:"manifest_version"`
	DefaultLocale         string      `json:"default_locale"`
	OfflineEnabled        bool        `json:"offline_enabled"`
	MinimumChromeVersion  string      `json:"minimum_chrome_version"`
	MinimumOperaVersion   string      `json:"minimum_opera_version"`
	ShortName             string      `json:"short_name"`
	Name                  string      `json:"name"`
	ExternallyConnectable struct {
		Ids     []string `json:"ids"`
		Matches []string `json:"matches"`
	} `json:"externally_connectable"`
	Storage struct {
		ManagedSchema string `json:"managed_schema"`
	} `json:"storage"`
	OptionsPage string        `json:"options_page"`
	Permissions []interface{} `json:"permissions"`
	UpdateURL   string        `json:"update_url"`
	Version     string        `json:"version"`
	Oauth2      struct {
		AutoApprove bool     `json:"auto_approve"`
		ClientID    string   `json:"client_id"`
		Scopes      []string `json:"scopes"`
	} `json:"oauth2"`
	OptionsUI struct {
		ChromeStyle bool   `json:"chrome_style"`
		OpenInTab   bool   `json:"open_in_tab"`
		Page        string `json:"page"`
	} `json:"options_ui"`
	Sandbox struct {
		ContentSecurityPolicy string   `json:"content_security_policy"`
		Pages                 []string `json:"pages"`
	} `json:"sandbox"`
	URLHandlers struct {
		PostmanCollection struct {
			Matches []string `json:"matches"`
			Title   string   `json:"title"`
		} `json:"postman_collection"`
	} `json:"url_handlers"`
	PageAction struct {
		DefaultIcon  string `json:"default_icon"`
		DefaultPopup string `json:"default_popup"`
		DefaultTitle string `json:"default_title"`
	} `json:"page_action"`
	WebAccessibleResources []string `json:"web_accessible_resources"`
}

func rightPad(s string, padStr string, overallLen int) string {
	var padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}
