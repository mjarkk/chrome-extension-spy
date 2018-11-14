package chrome

import (
	"encoding/json"
	"io/ioutil"
	"path"

	"github.com/mjarkk/chrome-extension-spy/types"
)

// GetExtensions gets all chrome extensions from google chrome
func GetExtensions(extensionsPath string) ([]types.ChromeExtension, []types.ExtensionManifest) {
	toReturn := []types.ChromeExtension{}
	toReturnFull := []types.ExtensionManifest{}
	files, err := ioutil.ReadDir(extensionsPath)
	if err != nil {
		return toReturn, toReturnFull
	}
	for _, f := range files {
		fName := f.Name()
		if len(fName) == 32 {
			extensionPath := path.Join(extensionsPath, fName)
			files, err := ioutil.ReadDir(extensionPath)
			if err != nil {
				return toReturn, toReturnFull
			}
			version := ""
			for _, versionDir := range files {
				version = versionDir.Name()
			}
			dat, err := ioutil.ReadFile(path.Join(extensionPath, version, "/manifest.json"))
			if err == nil {
				var manifest types.ExtensionManifest
				var addToReturnValue types.ChromeExtension
				json.Unmarshal(dat, &manifest)
				addToReturnValue.Name = manifest.Name
				addToReturnValue.HomepageURL = manifest.HomepageURL
				addToReturnValue.Pkg = fName
				addToReturnValue.PkgVersion = version
				addToReturnValue.ShortName = manifest.ShortName
				addToReturnValue.FullPkgURL = path.Join(extensionPath, version, "/")
				toReturn = append(toReturn, addToReturnValue)
				toReturnFull = append(toReturnFull, manifest)
			}
		}
	}
	return toReturn, toReturnFull
}
