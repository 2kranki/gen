// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Miscellaneous utility functions

package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/2kranki/jsonpreprocess"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// CopyDir copies from the given directory (src)
// and all of its files to the destination (dst).
func CopyDir(src, dst string) (error) {
	var err error

	src, err = IsPathRegularFile(src)
	if err == nil {
		return CopyFile(src, dst)
	}

	src, err = IsPathDir(src)
	if err != nil {
		return err
	}

	si, err := os.Stat(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return err
	}

	dir, err := os.Open(src)
	if err != nil {
		return err
	}

	fis, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	for _, fi := range fis {

		srcpath := src + "/" + fi.Name()

		dstpath := dst + "/" + fi.Name()

		if fi.IsDir() {
			err = CopyDir(srcpath, dstpath)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(srcpath, dstpath)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

// CopyFile copies a file given by its path (src) creating
// an output file given its path (dst)
func CopyFile(src, dst string) (error) {
	var err error

	// Clean up the input file path and check for its existence.
	src, err = IsPathRegularFile(src)
	if err != nil {
		return err
	}

	// Open the input file.
	fileIn, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fileIn.Close()

	// Create the output file.
	dst, _ = IsPathRegularFile(dst)		// Clean up output file path
	fileOut, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fileOut.Close()

	// Perform the copy and set the output file's privileges
	// to same as input file.
	_, err = io.Copy(fileOut, fileIn)
	if err == nil {
		si, err := os.Stat(src)
		if err == nil {
			err = os.Chmod(dst, si.Mode())
		}
	}

	return nil
}

// FileCompare compares two files returning true
// if they are equal.
func FileCompare(file1, file2 string) bool {
	var err 		error
	var file1size	int64
	var file2size	int64


	file1, err = IsPathRegularFile(file1)
	if err != nil {
		return false
	}
	si, err := os.Stat(file1)
	if err == nil {
		file1size = si.Size()
	}

	file2, err = IsPathRegularFile(file2)
	if err != nil {
		return false
	}
	si, err = os.Stat(file2)
	if err == nil {
		file2size = si.Size()
	}

	if file1size != file2size {
		return false
	}

	f1, err := os.Open(file1)
	if err != nil {
		log.Fatal(err)
	}

	f2, err := os.Open(file2)
	if err != nil {
		log.Fatal(err)
	}

	for {
		b1 := make([]byte, 8192)
		c1, err1 := f1.Read(b1)

		b2 := make([]byte, 8192)
		c2, err2 := f2.Read(b2)

		if c1 != c2 {
			return false
		}

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			} else {
				log.Fatal(err1, err2)
			}
		}

		if !bytes.Equal(b1, b2) {
			return false
		}
	}

	return true
}

// IsPathDir cleans up the supplied file path
// and then checks the cleaned file path to see
// if it is an existing standard directory. Return the
// cleaned up path and a potential error if it exists.
func IsPathDir(fp string) (string, error) {
	var err error
	var path string

	fp = os.ExpandEnv(fp)
	fp = filepath.Clean(fp)
	path, err = filepath.Abs(fp)
	if err != nil {
		return path, errors.New(fmt.Sprint("Error getting absolute path for:", fp, err))
	}
	fi, err := os.Lstat(path)
	if err != nil {
		return path, errors.New("path not found")
	}
	if fi.Mode().IsDir() {
		return path, nil
	}
	return path, errors.New("path not regular file")
}

// IsPathRegularFile cleans up the supplied file path
// and then checks the cleaned file path to see
// if it is an existing standard file. Return the
// cleaned up path and a potential error if it exists.
func IsPathRegularFile(fp string) (string, error) {
	var err error
	var path string

	fp = os.ExpandEnv(fp)
	fp = filepath.Clean(fp)
	path, err = filepath.Abs(fp)
	if err != nil {
		return path, errors.New(fmt.Sprint("Error getting absolute path for:", fp, err))
	}
	fi, err := os.Lstat(path)
	if err != nil {
		return path, errors.New("path not found")
	}
	if fi.Mode().IsRegular() {
		return path, nil
	}
	return path, errors.New("path not regular file")
}

// ReadJsonFile preprocesses out comments and then unmarshals the data
// generically.
func ReadJsonFile(jsonPath string) (interface{}, error) {
	var err error
	var jsonOut interface{}

	// Open the input template file
	input, err := os.Open(jsonPath)
	if err != nil {
		return jsonOut, err
	}
	textBuf := strings.Builder{}
	err = jsonpreprocess.WriteMinifiedTo(&textBuf, input)
	if err != nil {
		return jsonOut, err
	}

	// Read and process the template file
	err = json.Unmarshal([]byte(textBuf.String()), &jsonOut)
	if err != nil {
		return jsonOut, err
	}

	return jsonOut, err
}

// ReadJsonFileToData preprocesses out comments and then unmarshals the data
// into a data structure previously defined.
func ReadJsonFileToData(jsonPath string, jsonOut interface{}) error {
	var err error

	// Open the input template file
	input, err := os.Open(jsonPath)
	if err != nil {
		return err
	}
	textBuf := strings.Builder{}
	err = jsonpreprocess.WriteMinifiedTo(&textBuf, input)
	if err != nil {
		return err
	}

	// Read and process the template file
	err = json.Unmarshal([]byte(textBuf.String()), jsonOut)
	if err != nil {
		return err
	}

	return err
}

// Workers allows us to perform r number of task(s) at a time until
// all tasks are completed.  The input channel to run the task is
// returned by this function. The caller of this function simply
// puts to the returned input channel until all tasks have been inputted.
// It then closes the channel indicating that there is no more data.
// The function, completed, will be called when all of the input
// has been processed.
// Thanks to Vignesh Sk for his blog post on this.
func Workers(task func(interface{}), completed func(), r int) chan interface{} {
	input := make(chan interface{})
	ack := make(chan bool)
	for i := 0; i < r; i++ {
		go func() {
			for {
				v, ok := <-input
				if ok {
					task(v)
				} else {
					ack <- true
					return
				}
			}
		}()
	}
	go func() {
		for i := 0; i < r; i++ {
			<-ack
		}
		completed()
	}()
	return input
}

