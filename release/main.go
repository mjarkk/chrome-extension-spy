package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("\nChecking the program...")
	Check()

	fmt.Println("\nPrepearing the program...")
	Prepare()

	fmt.Println("\nCompiling the go code and build web_static files...")
	Compile()

	fmt.Println("\nZip all files...")
	ZipIt()
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
		"gox":  "gox -h",
	}
	for program, command := range checks {
		err := Run(command, "", []string{})
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
	fmt.Println("- removing web_static/build")
	os.RemoveAll(path.Join(projectDir, "web_static", "build"))
	fmt.Println("- removing web_static/node_modules")
	os.RemoveAll(path.Join(projectDir, "web_static", "node_modules"))
	fmt.Println("- removing contents of release/builds")
	os.RemoveAll(path.Join(currPath, "builds"))
	os.Mkdir(path.Join(currPath, "builds"), 0777)
	fmt.Println("- removing contents of release/tempgobuild")
	os.Mkdir(path.Join(currPath, "tempGoBuild"), 0777)
}

// Compile compiles the program and builds the web_static files
func Compile() {
	yarn := "yarn"
	fmt.Println("- running `" + yarn + "`")
	Run(yarn, path.Join(projectDir, "web_static"), []string{})

	yarnBuild := "yarn build"
	fmt.Println("- running `" + yarnBuild + "`")
	Run(yarnBuild, path.Join(projectDir, "web_static"), []string{})

	gox := "gox -output ./release/tempGoBuild/{{.Dir}}_{{.OS}}_{{.Arch}}"
	fmt.Println("- running `" + gox + "` (this might take a while)")
	Run(gox, projectDir, []string{})

	fmt.Println("")
}

// ZipIt zips puts everything inside a zip
func ZipIt() {

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
func Run(input string, executeingDir string, envs []string) error {
	command := strings.Split(input, " ")
	cmd := exec.Command(command[0], command[1:]...)
	if len(executeingDir) > 0 {
		cmd.Dir = executeingDir
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, envs...)
	err := cmd.Run()
	return err
}
