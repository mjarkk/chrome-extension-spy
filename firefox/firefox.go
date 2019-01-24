package firefox

import (
	"fmt"
	"log"
	"os/user"
)

// FF is the main package
type FF struct {
	LaunchCMD           string // this is something like "firefox" on linux and "C:\Program Files\firefox\firefox.exe" on win
	UserProfileLocation string // the default user profile location
	TmpDirs             FfTmpDirs
}

// FfTmpDirs is the temp dirs struct for the FF type
type FfTmpDirs struct {
	UnpackExts string // The unpacked extensions will be here
	Profile    string // The created user profile
}

// Setup returns the default firefox struct
func Setup() FF {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(usr.HomeDir)

	return FF{
		LaunchCMD:           "firefox",
		UserProfileLocation: "",
		TmpDirs: FfTmpDirs{
			UnpackExts: "",
			Profile:    "",
		},
	}

}
