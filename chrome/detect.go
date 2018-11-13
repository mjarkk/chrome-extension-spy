package chrome

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"
)

// Location returns the config location of the installed chrome
func Location(version string) string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("can't get home directory")
		os.Exit(1)
	}
	return path.Join(usr.HomeDir, "/.config/", version, "/Default/Extensions")
}

// GetLocation returns the config folder location of the current installed chrome
func GetLocation() (string, error) {
	if _, err := os.Stat(Location("chromium")); !os.IsNotExist(err) {
		return "chromium", nil
	}
	if _, err := os.Stat(Location("google-chrome")); !os.IsNotExist(err) {
		return "google-chrome", nil
	}
	if _, err := os.Stat(Location("google-chrome-beta")); !os.IsNotExist(err) {
		return "google-chrome-beta", nil
	}
	if _, err := os.Stat(Location("google-chrome-dev")); !os.IsNotExist(err) {
		return "google-chrome-dev", nil
	}
	if _, err := os.Stat(Location("google-chrome-unstable")); !os.IsNotExist(err) {
		return "google-chrome-unstable", nil
	}
	if _, err := os.Stat(Location("google-chrome-canary")); !os.IsNotExist(err) {
		return "google-chrome-canary", nil
	}
	return "", errors.New("Chrome location not found")
}
