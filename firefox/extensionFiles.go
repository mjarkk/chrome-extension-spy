package firefox

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/mjarkk/chrome-extension-spy/funs"
)

// GetRawExts exstracts all extensions and puts it inside folders
func (f *FF) GetRawExts() {
	if f.Err() {
		return
	}
	extensionZips, err := ioutil.ReadDir(path.Join(f.UserProfileLocation, "extensions"))
	if err != nil {
		f.HasErr = err
		return
	}

	unpackExts, err := ioutil.TempDir("", "FFExtensions")
	if err != nil {
		f.HasErr = err
		return
	}
	f.TmpDirs.UnpackExts = unpackExts

	for _, zip := range extensionZips {
		name := zip.Name()
		l := len(name)
		if string(name[l-4:l]) != ".xpi" {
			continue
		}
		extDir := path.Join(unpackExts, name)
		err := os.MkdirAll(extDir, os.ModePerm)
		if err != nil {
			fmt.Println("Skipping", name, "can't create folder, ERROR:", err.Error())
			continue
		}
		funs.Unzip(path.Join(f.UserProfileLocation, "extensions", name), extDir)
	}
}

// PackExtensions packs the extensions
func (f *FF) PackExtensions() {
	if f.Err() {
		return
	}

	dirs, err := ioutil.ReadDir(f.TmpDirs.UnpackExts)
	if err != nil {
		f.HasErr = err
		return
	}

	ffExtDir := path.Join(f.TmpDirs.Profile, "extensions")
	os.MkdirAll(ffExtDir, os.FileMode(int(0700)))
	for _, dir := range dirs {
		name := dir.Name()
		currentDir := path.Join(f.TmpDirs.UnpackExts, name)
		err := funs.ZipFiles(path.Join(ffExtDir, name), ListFilesInDir(currentDir), currentDir)
		if err != nil {
			fmt.Println("NOTE: can't create zip from:", name)
		}
	}
	fmt.Println()
}

// ListFilesInDir returns a list of all files in a dir this is also the fils in sub folders
func ListFilesInDir(startDir string) []string {
	toReturn := []string{}
	fileObjects, err := ioutil.ReadDir(startDir)
	if err != nil {
		return toReturn
	}
	for _, obj := range fileObjects {
		name := obj.Name()
		pathTo := path.Join(startDir, name)
		if obj.IsDir() {
			toReturn = append(toReturn, ListFilesInDir(pathTo)...)
			continue
		}
		toReturn = append(toReturn, pathTo)
	}
	return toReturn
}
