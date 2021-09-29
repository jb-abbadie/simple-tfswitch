package pkg

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// RenameFile : rename file name
func RenameFile(src string, dest string) {
	if err := os.Rename(src, dest); err != nil {
		log.Errorln("Error while renaming", err)
	}
}

// RemoveFiles : remove file
func RemoveFiles(src string) {
	files, err := filepath.Glob(src)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
}

// CheckFileExist : check if file exist in directory
func CheckFileExist(file string) bool {
	_, err := os.Stat(file)

	return err == nil
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		err := handleZipFile(f, dest)
		if err != nil {
			return err
		}
	}

	return nil
}

// handle 1 zip file
func handleZipFile(f *zip.File, dest string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// Store filename/path for returning and using later on
	fpath := filepath.Join(dest, f.Name)

	if f.FileInfo().IsDir() {
		// Make Folder
		err := os.MkdirAll(fpath, os.ModePerm)
		if err != nil {
			return err
		}

		return nil
	}

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

	return err
}

// CreateDirIfNotExist : create directory if directory does not exist
func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Debugf("Creating directory for terraform binary at: %v", dir)
		err = os.MkdirAll(dir, 0o755)
		if err != nil {
			log.Errorf("Unable to create directory for terraform binary at: %v", dir)
			panic(err)
		}
	}
}

// CheckDirExist : check if directory exist
// dir=path to file
// return bool
func CheckDirExist(dir string) bool {
	_, err := os.Stat(dir)

	return os.IsNotExist(err)
}

// Path : returns path of directory
// value=path to file
func Path(value string) string {
	return filepath.Dir(value)
}
