package funs

import (
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
