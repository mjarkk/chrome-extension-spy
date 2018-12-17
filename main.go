package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/mjarkk/chrome-extension-spy/chrome"
	"github.com/mjarkk/chrome-extension-spy/funs"
	"github.com/mjarkk/chrome-extension-spy/types"
	"github.com/mjarkk/chrome-extension-spy/webserver"
)

func main() {
	err := run()
	funs.PrintErr(err)
}

func run() error {

	isInfoMode := getFlags()
	if isInfoMode {
		getAppInfo()
		os.Exit(0)
	}

	var extTmpDir = make(chan string)
	var startWebServer = make(chan struct{})
	var chromeCommand = ""
	var extensions map[string]*types.FullAndSmallExt
	go func() {
		exts, chromeLaunchCommand, err := chrome.Setup(extTmpDir, isInfoMode)
		extensions = exts
		chromeCommand = chromeLaunchCommand
		funs.PrintErr(err)
		startWebServer <- struct{}{}
	}()
	tempDir := <-extTmpDir

	defer os.RemoveAll(tempDir)

	// Wait for chrome to complete it's tasks
	<-startWebServer

	var tasks sync.WaitGroup
	tasks.Add(2)

	var webserverErr error
	var chromeErr error

	forceClose := make(chan struct{})
	go func() {
		chromeErr = chrome.Launch(tempDir, chromeCommand, forceClose)
		forceClose <- struct{}{}
		tasks.Done()
	}()
	go func() {
		webserverErr = webserver.StartWebServer(tempDir, forceClose, extensions)
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

func getFlags() (isInfoMode bool) {
	isInfoPointer := flag.Bool("info", false, "Get info about this application")
	flag.Parse()
	return *isInfoPointer
}

func getAppInfo() error {
	chromeCMD, err := chrome.GetLocation()
	if err != nil {
		fmt.Println("- ERROR: chrome not found")
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
