package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"sync"

	"github.com/mjarkk/chrome-extension-spy/chrome"
	"github.com/mjarkk/chrome-extension-spy/fs"
	"github.com/mjarkk/chrome-extension-spy/types"
)

func main() {
	output, err := chrome.GetLocation()
	printErr(err)
	fullpath := chrome.Location(output)
	extensions, fullExtension := chrome.GetExtensions(fullpath)
	_, ext, fullExt := selectExtensionToUse(extensions, fullExtension)
	// create a temp dir to store the extension
	tempDir, err := ioutil.TempDir("", "chrome-extension-spy")
	defer os.RemoveAll(tempDir)
	printErr(err)
	err = fs.CopyFullExtension(ext.FullPkgURL, tempDir, []string{})
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

func editExtension(tmpDir string, ext types.ChromeExtension, fullExt types.ExtensionManifest) error {
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

func selectExtensionToUse(exts []types.ChromeExtension, fullExts []types.ExtensionManifest) (int64, types.ChromeExtension, types.ExtensionManifest) {
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

func printExtensions(exts []types.ChromeExtension) {
	maxNameLen := 0
	maxShortNameLen := 0
	maxPkgVersionLen := 7
	for _, ext := range exts {
		if len(ext.Name) > maxNameLen {
			maxNameLen = len(ext.Name)
		}
		if len(ext.ShortName) > maxShortNameLen {
			maxShortNameLen = len(ext.ShortName)
		}
		if len(ext.PkgVersion) > maxPkgVersionLen {
			maxPkgVersionLen = len(ext.PkgVersion)
		}
	}
	fmt.Printf("%s\t%s%s%s%s\n", "id", funs.RightPad("name", " ", maxNameLen+1), funs.RightPad("short name", " ", maxShortNameLen+1), funs.RightPad("version", " ", maxPkgVersionLen+1), "homepage")
	for id, ext := range exts {
		name := ext.Name
		if len(name) == 0 {
			name = "-"
		}
		shortName := ext.ShortName
		if len(shortName) == 0 {
			shortName = "-"
		}
		homepageURL := ext.HomepageURL
		if len(homepageURL) == 0 {
			homepageURL = "-"
		}
		pkgVersion := ext.PkgVersion
		if len(pkgVersion) == 0 {
			pkgVersion = "-"
		}
		fmt.Printf(
			"%v\t%s%s%s%s\n",
			id,
			funs.RightPad(name, " ", maxNameLen+1),
			funs.RightPad(shortName, " ", maxShortNameLen+1),
			funs.RightPad(pkgVersion, " ", maxPkgVersionLen+1),
			homepageURL,
		)
	}
}
