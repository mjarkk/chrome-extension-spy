package fs

import (
	"io"
	"io/ioutil"
	"os"
	"path"
)

// CopyDir copys a full dir over to another direcotry
func CopyDir(from string, to string, extensionDirs []string) error {
	extensionDirPath := path.Join(extensionDirs...)
	fullExtensionDirPath := path.Join(from, extensionDirPath)
	files, err := ioutil.ReadDir(fullExtensionDirPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		name := file.Name()

		if file.IsDir() {
			// create a dir and loop over that dir
			os.MkdirAll(path.Join(to, extensionDirPath, name), 0777)
			CopyDir(from, to, append(extensionDirs, name))
		} else {
			// copy a file over
			from, err := os.Open(path.Join(fullExtensionDirPath, name))
			if err != nil {
				return err
			}
			to, err := os.Create(path.Join(to, extensionDirPath, name))
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
