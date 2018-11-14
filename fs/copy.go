package fs

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// CopyFullExtension copys a exteion folder over to another direcotry
func CopyFullExtension(baseDir string, tempDir string, extensionDir []string) error {
	extensionDirPath := strings.Join(extensionDir, "/")
	fullExtensionDirPath := path.Join(baseDir, extensionDirPath)
	files, err := ioutil.ReadDir(fullExtensionDirPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		name := file.Name()

		if file.IsDir() {
			// create a dir and loop over that dir
			os.MkdirAll(path.Join(tempDir, extensionDirPath, name), 0777)
			CopyFullExtension(baseDir, tempDir, append(extensionDir, name))
		} else {
			// copy a file over
			from, err := os.Open(path.Join(fullExtensionDirPath, name))
			if err != nil {
				return err
			}
			to, err := os.Create(path.Join(tempDir, extensionDirPath, name))
			if err != nil {
				return err
			}
			_, err = io.Copy(to, from)
			if err != nil {
				return err
			}
			from.Close()
			to.Close()
		}
	}
	return nil
}
