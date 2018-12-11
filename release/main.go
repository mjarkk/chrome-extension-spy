package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	defer func() {

	}()
	fmt.Println("\nChecking...")
	Check()

	fmt.Println("\nPrepeare to compile...")
	Prepare()

	fmt.Println("\nCompiling the go code and build web_static files...")
	Compile()

	fmt.Println("\nZip all files...")
	ZipIt()

	fmt.Println("\nCleanup...")
	Cleanup()

	fmt.Println("\nSucsessfull created release zips, They are located under: ./release/builds/")
}

var currPath = ""
var projectDir = ""

// Check checks if the needed programs are installed
func Check() {
	fmt.Println("- checking if programs are installed")
	checks := map[string]string{
		"node": "node --version",
		"npm":  "npm --version",
		"yarn": "yarn --version",
		// "gox":  "gox -h", // TODO: find alternative to check this on windows
	}
	for program, command := range checks {
		err := Run(command, "")
		if err != nil {
			ExitErr(errors.New(program + " not found"))
		}
	}
	fmt.Println("- creating global variables")
	ex, err := os.Executable()
	ExitErr(err)
	currPath = filepath.Dir(ex)
	projectDir = path.Join(currPath, "..")
}

// Prepare for compiling
func Prepare() {
	var waitJobs sync.WaitGroup
	jobs := []func(){
		func() {
			fmt.Println("- removing web_static/build")
			os.RemoveAll(path.Join(projectDir, "web_static", "build"))
		},
		func() {
			fmt.Println("- removing web_static/node_modules")
			os.RemoveAll(path.Join(projectDir, "web_static", "node_modules"))
		},
		func() {
			fmt.Println("- removing contents of release/builds")
			os.RemoveAll(path.Join(currPath, "builds"))
			os.Mkdir(path.Join(currPath, "builds"), 0777)
		},
		func() {
			fmt.Println("- removing contents of release/tempgobuild")
			os.RemoveAll(path.Join(currPath, "tempGoBuild"))
			os.Mkdir(path.Join(currPath, "tempGoBuild"), 0777)
		},
	}
	waitJobs.Add(len(jobs))
	for _, job := range jobs {
		go func(task func()) {
			task()
			waitJobs.Done()
		}(job)
	}
	waitJobs.Wait()
}

// Compile compiles the program and builds the web_static files
func Compile() {
	var waitJobs sync.WaitGroup
	jobs := []func(){
		func() {
			// Install deps
			yarn := "yarn"
			fmt.Println("- running `" + yarn + "`")
			Run(yarn, path.Join(projectDir, "web_static"))

			// Build the web static files
			yarnBuild := "yarn build"
			fmt.Println("- running `" + yarnBuild + "`")
			Run(yarnBuild, path.Join(projectDir, "web_static"))
			waitJobs.Done()
		},
		func() {
			// Get golang libs
			goGet := "go get"
			fmt.Println("- running `" + goGet + "`")
			Run(goGet, projectDir)

			// Build the go program
			gox := "gox -arch !s390x -arch !mips -arch !mips64 -arch !mipsle -arch !ppc64 -arch !ppc64le -os windows -os linux -os darwin -output ./release/tempGoBuild/{{.Dir}}_{{.OS}}_{{.Arch}}"
			fmt.Println("- running `" + gox + "` (this might take a while)")
			Run(gox, projectDir)
			waitJobs.Done()
		},
	}
	waitJobs.Add(2)
	go jobs[0]()
	go jobs[1]()
	waitJobs.Wait()
}

// ZipIt zips puts everything inside a zip
func ZipIt() {
	filesToCompile, err := ioutil.ReadDir(path.Join(currPath, "tempGoBuild"))
	ExitErr(err)
	var waitForZips sync.WaitGroup
	waitForZips.Add(len(filesToCompile))
	for _, toInclude := range filesToCompile {
		go func(fileObj os.FileInfo) {
			fname := fileObj.Name()
			fullFname := path.Join(currPath, "tempGoBuild", fname)
			fmt.Println("- building:", fname)
			cut := strings.Split(strings.Replace(fname, "chrome-extension-spy_", "", 1), "_")
			platform := cut[0]
			arch := strings.Replace(cut[1], ".exe", "", 1)
			zipFile, err := os.Create(path.Join(currPath, "builds", "chrome-extension-spy_"+platform+"_"+arch+".zip"))
			ExitErr(err)
			filenameEnd := ""
			if strings.Contains(cut[1], ".exe") {
				filenameEnd = ".exe"
			}
			archive := zip.NewWriter(zipFile)

			filesToAdd := map[string]string{
				fullFname: "chrome-extension-spy" + filenameEnd,
				path.Join(projectDir, "web_static", "extension_inject.js"):              "web_static/extension_inject.js",
				path.Join(projectDir, "web_static", "build", "index.html"):              "web_static/build/index.html",
				path.Join(projectDir, "web_static", "build", "js", "bundel.js"):         "web_static/build/js/bundel.js",
				path.Join(projectDir, "web_static", "build", "js", "bundel.js.LICENSE"): "web_static/build/js/bundel.js.LICENSE",
			}

			for from, to := range filesToAdd {
				file, err := archive.Create(to)
				ExitErr(err)
				data, err := ioutil.ReadFile(from)
				ExitErr(err)
				file.Write(data)
			}

			archive.Close()
			zipFile.Close()
			waitForZips.Done()
		}(toInclude)
	}
	waitForZips.Wait()
}

// Cleanup removes junk that issn't needed anymore
func Cleanup() {
	os.RemoveAll(path.Join(currPath, "tempGoBuild"))
}

// ExitErr exit the program with the error message when there is an error
func ExitErr(err error) {
	if err != nil {
		fmt.Println("\nCan't make a releases\n---------------------")
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// Run runs a command
func Run(input string, executeingDir string) error {
	command := strings.Split(input, " ")

	cmd := exec.Command(command[0], command[1:]...)
	if len(executeingDir) > 0 {
		cmd.Dir = executeingDir
	}
	err := cmd.Run()
	return err
}
