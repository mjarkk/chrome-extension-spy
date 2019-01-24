package funs

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// RightPad make a string a fixed size
func RightPad(s string, padStr string, overallLen int) string {
	var padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}

// PrintErr check if there is a error and if so it returns that error and exits the program
func PrintErr(err error) {
	if err != nil {
		errMsg := color.New(color.FgRed, color.Bold)
		errMsg.Println("Program stopped because of an error:")
		fmt.Println(err)
		os.Exit(1)
	}
}

// RemoveTmpDirs removes a list of input dirs
func RemoveTmpDirs(dirs []string) {
	for _, dir := range dirs {
		if len(dir) > 0 {
			os.RemoveAll(dir)
		}
	}
}

// Flags are the flags the program can show
type Flags struct {
	ForceFF     bool
	ForceChrome bool
	IsInfo      bool
}

// GetFlags returns the program setted flags
func GetFlags() Flags {
	isInfo := flag.Bool("info", false, "Get info about this application")
	isFF := flag.Bool("isFF", false, "Force using firefox")
	isChrome := flag.Bool("isChrome", false, "Force using chrome")

	flag.Parse()
	return Flags{
		IsInfo:      *isInfo,
		ForceChrome: *isChrome,
		ForceFF:     *isFF,
	}
}
