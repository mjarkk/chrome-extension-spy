package chrome

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"
	"runtime"
)

// Location returns the config location of the installed chrome
func Location(version string) string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("can't get home directory")
		os.Exit(1)
	}
	if runtime.GOOS == "windows" {
		return path.Join(usr.HomeDir, "AppData", "Local", version, "User Data", "Default", "Extensions")
	}
	if runtime.GOOS == "darwin" {
		return path.Join(usr.HomeDir, "Library", "Application Support", version, "Default", "Extensions")
	}
	// use linux dir as fallback
	return path.Join(usr.HomeDir, ".config", version, "Default", "Extensions")
}

// GetLocation returns the config folder location of the current installed chrome
func GetLocation() (string, error) {
	checkLocation := func(input string) bool {
		_, err := os.Stat(Location("chromium"))
		return !os.IsNotExist(err)
	}
	if checkLocation("chromium") {
		return "chromium", nil
	}
	if checkLocation("google-chrome") {
		return "google-chrome", nil
	}
	if checkLocation("google-chrome-beta") {
		return "google-chrome-beta", nil
	}
	if checkLocation("google-chrome-dev") {
		return "google-chrome-dev", nil
	}
	if checkLocation("google-chrome-unstable") {
		return "google-chrome-unstable", nil
	}
	if checkLocation("google-chrome-canary") {
		return "google-chrome-canary", nil
	}
	return "", errors.New("Chrome location not found")
}
