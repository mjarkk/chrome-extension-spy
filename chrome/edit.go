package chrome

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/mjarkk/chrome-extension-spy/types"
)

// EditExtension injects the spy code into a extension
func EditExtension(extDir string, ext types.ChromeExtension, fullExt types.ExtensionManifest) error {
	thisFileDir, err := os.Executable()
	if err != nil {
		return err
	}
	injectable, err := ioutil.ReadFile(path.Join(thisFileDir, "../web_static/extension_inject.js"))
	if err != nil {
		return err
	}
	for _, srcItem := range fullExt.Background.Scripts {
		fullFileDir := path.Join(extDir, srcItem)
		file, err := ioutil.ReadFile(fullFileDir)
		if err != nil {
			return err
		}
		toWrite := strings.Replace(string(injectable), "--EXT-APP-ID--", ext.Pkg, 1) + string(file)
		ioutil.WriteFile(fullFileDir, []byte(toWrite), 0777)
	}
	return nil
}
