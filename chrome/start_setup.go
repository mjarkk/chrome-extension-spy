package chrome

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/mjarkk/chrome-extension-spy/fs"
)

// Setup sets up the chrome part
func Setup(extTmpDir chan string) (string, error) {
	chromeCommand, err := GetLocation()
	if err != nil {
		return chromeCommand, err
	}
	chromeLocation := Location(chromeCommand)

	tempDir, err := ioutil.TempDir("", "chrome-extension-spy")
	if err != nil {
		return chromeCommand, err
	}
	extTmpDir <- tempDir

	extensions, fullExtensions := GetExtensions(chromeLocation)
	for id, fullExtension := range fullExtensions {
		extension := extensions[id]

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
