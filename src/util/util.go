// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Miscellaneous utility functions

// Some functions were taken from https://blog.kowalczyk.info/book/go-cookbook.html
// which was declared public domain at the time that the functions were taken.

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
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

//----------------------------------------------------------------------------
//                             Command Execution
//----------------------------------------------------------------------------

// os.Exec contains further details
type ExecCmd struct {
	cmd       	*exec.Cmd
}

func (c *ExecCmd) Cmd( ) *exec.Cmd {
	return c.cmd
}

func (c *ExecCmd) CommandString( ) string {
	n := len(c.cmd.Args)
	a := make([]string, n, n)
	for i := 0; i < n; i++ {
		a[i] = c.QuoteArgIfNeeded(i)
	}
	return strings.Join(a, " ")
}

func (c *ExecCmd) QuoteArgIfNeeded(n int) string {
	var s		string

	s = c.cmd.Args[n]
	if strings.Contains(s, " ") || strings.Contains(s, "\"") {
		s = strings.Replace(s, `"`, `\"`, -1)
		return `"` + s + `"`
	}
	return s
}

// Runt runs the previously set up command.
func (c *ExecCmd) Run( ) error {
	var err		error

	err = c.cmd.Run()

	return err
}

// RunWithOutput runs the previously set up command, gets the combined output
// of sysout and syserr, trims whitespace from it and returns if error free.
// If any error occurs, it is simply returned.
func (c *ExecCmd) RunWithOutput( ) (string, error) {
	var err		error

	outBytes, err := c.cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	s := string(outBytes)
	s = strings.TrimSpace(s)

	return s, nil
}

func NewExecCmd(name string, args ...string) *ExecCmd {
	ce := ExecCmd{}
	if len(name) > 0 {
		ce.cmd = exec.Command(name, args...)
	}
	return &ce
}


//----------------------------------------------------------------------------
//                             CopyDir
//----------------------------------------------------------------------------

// CopyDir copies from the given directory (src) and all of its files to the
// destination (dst).
func CopyDir(src, dst string) error {
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
		dir.Close()
		return err
	}
	dir.Close()

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

//----------------------------------------------------------------------------
//                             CopyFile
//----------------------------------------------------------------------------

// CopyFile copies a file given by its path (src) creating
// an output file given its path (dst)
func CopyFile(src, dst string) error {
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

//----------------------------------------------------------------------------
//                             ErrorString
//----------------------------------------------------------------------------

func ErrorString(err error) string {
	if err == nil {
		return "ok"
	} else {
		return err.Error()
	}
}

//----------------------------------------------------------------------------
//                             FileCompare
//----------------------------------------------------------------------------

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
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		log.Fatal(err)
	}
	defer f2.Close()

	b1 := make([]byte, 8192)
	b2 := make([]byte, 8192)
	for {
		c1, err1 := f1.Read(b1)
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

	return false
}

//----------------------------------------------------------------------------
//                             		FormatArgs
//----------------------------------------------------------------------------

func FormatArgs(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}
	format := args[0].(string)
	if len(args) == 1 {
		return format
	}
	return fmt.Sprintf(format, args[1:]...)
}

//----------------------------------------------------------------------------
//                             HomeDir
//----------------------------------------------------------------------------

// HomeDir returns $HOME diretory of the current user
func HomeDir() string {

	// user.Current() returns nil if cross-compiled e.g. on mac for linux
	if usr, _ := user.Current(); usr != nil {
		return usr.HomeDir
	}

	return os.Getenv("HOME")
}

//----------------------------------------------------------------------------
//                             IsPathDir
//----------------------------------------------------------------------------

// IsPathDir cleans up the supplied file path
// and then checks the cleaned file path to see
// if it is an existing standard directory. Return the
// cleaned up path and a potential error if it exists.
func IsPathDir(fp string) (string, error) {
	var err error
	var path string

	path = PathClean(fp)
	fi, err := os.Lstat(path)
	if err != nil {
		return path, errors.New("path not found")
	}
	if fi.Mode().IsDir() {
		return path, nil
	}
	return path, errors.New("path not regular file")
}

//----------------------------------------------------------------------------
//                            IsPathRegularFile
//----------------------------------------------------------------------------

// IsPathRegularFile cleans up the supplied file path
// and then checks the cleaned file path to see
// if it is an existing standard file. Return the
// cleaned up path and a potential error if it exists.
func IsPathRegularFile(fp string) (string, error) {
	var err error
	var path string

	path = PathClean(fp)
	fi, err := os.Lstat(path)
	if err != nil {
		return path, errors.New("path not found")
	}
	if fi.Mode().IsRegular() {
		return path, nil
	}
	return path, errors.New("path not regular file")
}

//----------------------------------------------------------------------------
//                             PanicIf
//----------------------------------------------------------------------------

func PanicIf(cond bool, args ...interface{}) {
	if !cond {
		return
	}
	s := FormatArgs(args...)
	if s == "" {
		s = "fatalIf: cond is false"
	}
	panic(s)
}

//----------------------------------------------------------------------------
//                             PanicIfErr
//----------------------------------------------------------------------------

func PanicIfErr(err error, args ...interface{}) {
	if err == nil {
		return
	}
	s := FormatArgs(args...)
	if s == "" {
		s = err.Error()
	}
	panic(s)
}

//----------------------------------------------------------------------------
//                             PathClean
//----------------------------------------------------------------------------

// PathClean cleans up the supplied file path.
func PathClean(fp string) string {
	var path string

	if strings.HasPrefix(fp, "~") {
		fp = HomeDir() + fp[1:]
	}
	fp = os.ExpandEnv(fp)
	fp = filepath.Clean(fp)
	path, _ = filepath.Abs(fp)

	return path
}

//----------------------------------------------------------------------------
//                            ReadJsonFile
//----------------------------------------------------------------------------

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

//----------------------------------------------------------------------------
//                            ReadJsonFileToData
//----------------------------------------------------------------------------

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

//----------------------------------------------------------------------------
//                            		Workers
//----------------------------------------------------------------------------

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

