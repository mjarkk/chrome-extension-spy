package chrome

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/mjarkk/chrome-extension-spy/funs"

	"github.com/mjarkk/chrome-extension-spy/fs"
	"github.com/mjarkk/chrome-extension-spy/types"
)

// ChromeExts is a list of all extensions on the system
var ChromeExts = make(map[string]*types.FullAndSmallExt)

// Setup sets up the chrome part
func Setup(extTmpDir chan string, flags funs.Flags, useFF chan bool) (string, error) {
	returnExts := make(map[string]*types.FullAndSmallExt)
	defer func() {
		ChromeExts = returnExts
	}()
	if flags.ForceFF {
		useFF <- true
		return "", nil
	}
	chromeCommand, err := GetLocation()
	if err != nil && !flags.ForceChrome {
		useFF <- true
		return chromeCommand, err
	}
	useFF <- false
	chromeLocation := Location(chromeCommand)

	tempDir, err := ioutil.TempDir("", "chrome-extension-spy")
	if err != nil {
		return "", err
	}

	extTmpDir <- tempDir

	err = os.Chmod(tempDir, 0777)
	if err != nil {
		return "", err
	}

	extensions, fullExtensions := GetExtensions(chromeLocation)
	for id, fullExtension := range fullExtensions {
		extension := extensions[id]

		returnExts[extension.Pkg] = &types.FullAndSmallExt{
			Small: extension,
			Full:  fullExtension,
		}

		from := path.Join(chromeLocation, extension.Pkg)
		versions, err := ioutil.ReadDir(from)
		if err != nil {
			return chromeCommand, err
		}
		if len(versions) < 1 {
			return chromeCommand, errors.New("Extension " + extension.Pkg + " has no version folder")
		}
		from = path.Join(from, versions[0].Name())

		to := path.Join(tempDir, extension.Pkg)

		err = os.Mkdir(to, 0777)
		if err != nil {
			return chromeCommand, err
		}

		fs.CopyDir(from, to, []string{})

		err = EditExtension(to, extension, fullExtension)
		if err != nil {
			return chromeCommand, err
		}
	}

	return chromeCommand, nil
}
