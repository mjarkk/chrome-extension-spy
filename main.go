package main

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/mjarkk/chrome-extension-spy/chrome"
	"github.com/mjarkk/chrome-extension-spy/firefox"
	"github.com/mjarkk/chrome-extension-spy/funs"
	"github.com/mjarkk/chrome-extension-spy/webserver"
)

func main() {
	err := run()
	funs.PrintErr(err)
}

func run() error {

	flags := funs.GetFlags()
	if flags.IsInfo {
		err := getAppInfo()
		if err != nil {
			fmt.Println("- ERROR:", err.Error())
		}
		os.Exit(0)
	}

	var extTmpDir = make(chan string)
	var startWebServer = make(chan struct{})
	var useFF = make(chan bool)
	var chromeCommand = ""

	go func() {
		cmd, err := chrome.Setup(extTmpDir, flags, useFF)
		chromeCommand = cmd
		funs.PrintErr(err)
		startWebServer <- struct{}{}
	}()

	tmpDir := ""

	if <-useFF {
		fmt.Println("using firefox!")
		f := firefox.Setup()
		fmt.Println(f)
		if len(f.TmpDirs.Profile) > 0 {
			os.RemoveAll(f.TmpDirs.Profile)
		}
		if len(f.TmpDirs.UnpackExts) > 0 {
			os.RemoveAll(f.TmpDirs.UnpackExts)
		}
		os.Exit(1)
	} else {
		tmpDir = <-extTmpDir

		// Wait for chrome to complete it's tasks
		<-startWebServer
	}

	var tasks sync.WaitGroup
	tasks.Add(2)

	var webserverErr error
	var chromeErr error

	forceClose := make(chan struct{})
	chromeDirTmpDir := make(chan string)
	go func() {
		chromeErr = chrome.Launch(tmpDir, chromeDirTmpDir, chromeCommand, forceClose)
		forceClose <- struct{}{}
		tasks.Done()
	}()
	browserTmpDir := <-chromeDirTmpDir
	go func() {
		webserverErr = webserver.StartWebServer(tmpDir, browserTmpDir, forceClose)
		forceClose <- struct{}{}
		tasks.Done()
	}()
	go func() {
		waitForExitInput()
		forceClose <- struct{}{}
	}()

	tasks.Wait()

	if webserverErr != nil {
		return webserverErr
	}
	if chromeErr != nil {
		return chromeErr
	}
	return nil
}

func waitForExitInput() {
	var input string
	fmt.Print("Type exit to exit the program\n")
	fmt.Scanf("%s", &input)
	if input != "exit" {
		waitForExitInput()
	}
}

func getAppInfo() error {
	chromeCMD, err := chrome.GetLocation()
	if err != nil {
		return errors.New("chrome not found")
	}
	fmt.Println("- OK: chrome version:", chromeCMD)

	chromeLoc := chrome.Location(chromeCMD)
	fmt.Println("- OK: chrome extension path:", chromeLoc)

	out, _ := chrome.GetExtensions(chromeLoc)
	if len(out) == 0 {
		fmt.Println("- ERROR: no extensions found")
	} else {
		fmt.Println("- OK:", len(out), "extensions found")
	}

	fmt.Println("- OK: command to launch chorme:", chrome.ChromeLocation(chromeCMD))
	return nil
}
