package firefox

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/user"
	"path"
	"runtime"
	"strings"

	"github.com/mjarkk/chrome-extension-spy/funs"
)

// GetLaunchCMD detect the luanch command for to launch firefox and adds it to the input obj
func (f *FF) GetLaunchCMD() {
	if f.Err() {
		return
	}
	if runtime.GOOS == "windows" {
		f.LaunchCMD = "C:\\Program Files\\Mozilla Firefox\\"
	} else if runtime.GOOS == "linux" {
		options := []string{"firefox", "firefox-developer-edition"}
		for _, opt := range options {
			if funs.CommandExsists(opt) {
				f.LaunchCMD = opt
				break
			}
		}
	} else if runtime.GOOS == "darwin" {
		// TODO: detect the path on mac
	}

	if len(f.LaunchCMD) == 0 {
		f.HasErr = errors.New("Can't find firefox")
		return
	}

	fmt.Println("Found firefox:", f.LaunchCMD)
}

// GetUserLocation detect the user data location of firefox and adds it to the type
func (f *FF) GetUserLocation() {
	if f.Err() {
		return
	}

	if runtime.GOOS == "linux" {
		usr, err := user.Current()
		if err != nil {
			f.HasErr = err
			return
		}
		ffConfig := path.Join(usr.HomeDir, "/.mozilla/firefox/")
		contents, err := ioutil.ReadDir(ffConfig)
		if err != nil {
			f.HasErr = err
			return
		}

		possibleDirs := []string{}
		for _, item := range contents {
			name := item.Name()
			if item.IsDir() && !strings.Contains(name, " ") {
				if _, err := ioutil.ReadDir(path.Join(ffConfig, name, "extensions")); err != nil {
					continue
				}
				possibleDirs = append(possibleDirs, name)
			}
		}
		if len(possibleDirs) >= 1 {
			f.UserProfileLocation = path.Join(ffConfig, possibleDirs[0])
			return
		}
	} else if runtime.GOOS == "windows" {
		// TODO: find user data location in windows
	} else if runtime.GOOS == "darwin" {
		// TODO: find user data location for mac
	}

	err := errors.New("could not find user profile data")
	f.HasErr = err
	fmt.Println("ERROR:", err.Error())
}
