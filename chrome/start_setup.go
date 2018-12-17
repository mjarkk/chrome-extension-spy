package chrome

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/mjarkk/chrome-extension-spy/fs"
	"github.com/mjarkk/chrome-extension-spy/types"
)

// Setup sets up the chrome part
func Setup(extTmpDir chan string, testMode bool) (map[string]*types.FullAndSmallExt, string, error) {
	returnExts := make(map[string]*types.FullAndSmallExt)
	chromeCommand, err := GetLocation()
	if err != nil {
		return returnExts, chromeCommand, err
	}
	chromeLocation := Location(chromeCommand)

	tempDir, err := ioutil.TempDir("", "chrome-extension-spy")
	if err != nil {
		return returnExts, chromeCommand, err
	}

	extTmpDir <- tempDir

	err = os.Chmod(tempDir, 0777)
	if err != nil {
		return returnExts, "", err
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
			return returnExts, chromeCommand, err
		}
		if len(versions) < 1 {
			return returnExts, chromeCommand, errors.New("Extension " + extension.Pkg + " has no version folder")
		}
		from = path.Join(from, versions[0].Name())

		to := path.Join(tempDir, extension.Pkg)

		err = os.Mkdir(to, 0777)
		if err != nil {
			return returnExts, chromeCommand, err
		}

		fs.CopyDir(from, to, []string{})

		err = EditExtension(to, extension, fullExtension)
		if err != nil {
			return returnExts, chromeCommand, err
		}
	}

	return returnExts, chromeCommand, nil
}
