package chrome

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/mjarkk/chrome-extension-spy/types"
)

// EditExtension injects the spy code into a extension
func EditExtension(tmpDir string, ext types.ChromeExtension, fullExt types.ExtensionManifest) error {
	thisFileDir, err := os.Executable()
	if err != nil {
		return err
	}
	injectable, err := ioutil.ReadFile(path.Join(thisFileDir, "../web_static/extension_inject.js"))
	if err != nil {
		return err
	}
	for _, srcItem := range fullExt.Background.Scripts {
		fullFileDir := path.Join(tmpDir, srcItem)
		file, err := ioutil.ReadFile(fullFileDir)
		if err != nil {
			return err
		}
		toWrite := string(injectable) + string(file)
		ioutil.WriteFile(fullFileDir, []byte(toWrite), 0777)
	}
	return nil
}
