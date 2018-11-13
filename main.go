package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/mjarkk/chrome-extension-spy/chrome"
)

func main() {
	output, err := chrome.GetLocation()
	printErr(err)
	fullpath := chrome.Location(output)
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
		chromeErr = chrome.Launch(tempDir, waitForExit)
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

func rightPad(s string, padStr string, overallLen int) string {
	var padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}
