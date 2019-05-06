// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Miscellaneous utility functions

package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/2kranki/jsonpreprocess"
	"os"
	"path/filepath"
	"strings"
)

// IsPathRegularFile cleans up the supplied file path
// and then checks the cleaned file path to see
// if it is an existing standard file. Return the
// cleaned up path and a potential error if it exists.
func IsPathRegularFile(fp string) (string, error) {
	var err error
	var path string

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

