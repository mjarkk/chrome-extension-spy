package main

import (
	"fmt"

	"github.com/mjarkk/chrome-extension-spy/chrome"
	"github.com/mjarkk/chrome-extension-spy/funs"
)

func main() {
	var extTmpDir = make(chan string)
	var startWebServer = make(chan struct{})
	var chromeCommand = ""

	go func() {
		chromeLaunchCommand, err := chrome.Setup(extTmpDir)
		chromeCommand = chromeLaunchCommand
		funs.PrintErr(err)
		startWebServer <- struct{}{}
	}()
	tempDir := <-extTmpDir
	fmt.Println(tempDir)
	// defer os.RemoveAll(tempDir)

	// Wait for chrome to complete it's tasks
	<-startWebServer

	forceClose := make(chan struct{})
	go func() {
		chrome.Launch(tempDir, chromeCommand, forceClose)
		forceClose <- struct{}{}
	}()

	<-forceClose

	// var waitForGinAndChrome sync.WaitGroup
	// waitForGinAndChrome.Add(2)

	// var waitForExit = make(chan struct{})

	// var ginErr error
	// go func() {
	// 	ginErr = startWebServer(waitForExit)
	// 	waitForGinAndChrome.Done()
	// }()
	// var chromeErr error
	// go func() {
	// 	chromeErr = chrome.Launch(tempDir, waitForExit)
	// 	waitForGinAndChrome.Done()
	// }()
	// go func() {
	// 	waitForExitInput()
	// 	waitForExit <- struct{}{}
	// }()

	// waitForGinAndChrome.Wait()
	// funs.PrintErr(chromeErr)
	// funs.PrintErr(ginErr)
}

func waitForExitInput() {
	var input string
	fmt.Print("Type exit to exit the program")
	fmt.Scanf("%s", &input)
	if input != "exit" {
		waitForExitInput()
	}
}
