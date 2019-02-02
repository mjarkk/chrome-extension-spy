package funs

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// RightPad make a string a fixed size
func RightPad(s string, padStr string, overallLen int) string {
	var padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}

// PrintErr check if there is a error and if so it returns that error and exits the program
func PrintErr(err error) {
	if err != nil {
		errMsg := color.New(color.FgRed, color.Bold)
		errMsg.Println("Program stopped because of an error:")
		fmt.Println(err)
		os.Exit(1)
	}
}

// RemoveTmpDirs removes a list of input dirs
func RemoveTmpDirs(dirs []string) {
	for _, dir := range dirs {
		if len(dir) > 0 {
			os.RemoveAll(dir)
		}
	}
}

// Flags are the flags the program can show
type Flags struct {
	ForceFF     bool
	ForceChrome bool
	IsInfo      bool
}

// GetFlags returns the program setted flags
func GetFlags() Flags {
	isInfo := flag.Bool("info", false, "Get info about this application")
	isFF := flag.Bool("isFF", false, "Force using firefox")
	isChrome := flag.Bool("isChrome", false, "Force using chrome")

	flag.Parse()
	return Flags{
		IsInfo:      *isInfo,
		ForceChrome: *isChrome,
		ForceFF:     *isFF,
	}
}

// CommandExsists returns true a command exsists
func CommandExsists(command string) bool {
	return exec.Command("which", command).Run() == nil
}

// Unzip unzips a zip file into a dir
func Unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			// Make File
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			_, err = io.Copy(outFile, rc)

			// Close the file without defer to close before next iteration of loop
			outFile.Close()

			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ZipFiles compresses one or many files into a single zip archive file.
// Param 1: filename is the output zip file's name.
// Param 2: files is a list of files to add to the zip.
func ZipFiles(filename string, files []string, toRemovePrefix string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file, toRemovePrefix); err != nil {
			return err
		}
	}
	return nil
}

// AddFileToZip adds a file to a zip object
func AddFileToZip(zipWriter *zip.Writer, filename, toRemovePrefix string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	shortName := strings.Replace(filename, toRemovePrefix, "", 1)
	if string(shortName[0]) == "/" {
		strings.Replace(shortName, "/", "", 1)
	}
	header.Name = shortName

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}
