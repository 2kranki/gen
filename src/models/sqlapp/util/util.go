// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Miscellaneous utility functions

// Some functions were taken from https://blog.kowalczyk.info/book/go-cookbook.html
// which was declared public domain at the time that the functions were taken.

package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/2kranki/jsonpreprocess"
	"golang.org/x/tools/go/gccgoexportdata"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"
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
func CopyDir(src, dst *Path) error {
	var err 	error

	//log.Printf("CopyDir: base: %s  last: %c\n", pathIn.Base(), dst[len(dst)-1])
	if dst.String()[len(dst.String())-1] == os.PathSeparator {
		dst = dst.Append(src.Base())
	}
	//log.Printf("CopyDir: %s -> %s\n", pathIn.String(), pathOut.String())

	if src.IsPathRegularFile() {
		return CopyFile(src, dst)
	}

	if !src.IsPathDir() {
		return fmt.Errorf("Error: CopyDir: %s is not a file or directory!\n", src.String())
	}

	si, err := os.Stat(src.Absolute())
	if err != nil {
		return err
	}
	mode := si.Mode() & 03777

	log.Printf("CopyDir: MkdirAll %s %o\n", dst.Absolute(), mode)
	err = os.MkdirAll(dst.Absolute(), mode)
	if err != nil {
		return err
	}
	if !dst.IsPathDir() {
		return fmt.Errorf("Error: %s could not be found!", dst.Absolute())
	}

	dir, err := os.Open(src.Absolute())
	if err != nil {
		return err
	}

	entries, err := dir.Readdir(-1)
	if err != nil {
		dir.Close()
		return err
	}
	dir.Close()

	for _, fi := range entries {
		srcNew := src.Append(fi.Name())
		dstNew := dst.Append(fi.Name())

		if fi.Mode().IsDir() {
			log.Printf("CopyDir: Dir: %s -> %s\n", srcNew.String(), dstNew.String())
			err = CopyDir(srcNew, dstNew)
			if err != nil {
				return err
			}
		} else if fi.Mode().IsRegular() {
			log.Printf("CopyDir: File: %s -> %s\n", srcNew.String(), dstNew.String())
			err = CopyFile(srcNew, dstNew)
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
func CopyFile(src, dst *Path) error {
	var err 	error

	log.Printf("CopyFile: %s -> %s\n", src.Absolute(), dst.Absolute())

	// Clean up the input file path and check for its existence.
	if !src.IsPathRegularFile() {
		return fmt.Errorf("Error: %s is not a file!\n", src.String())
	}

	// Open the input file.
	fileIn, err := os.Open(src.Absolute())
	if err != nil {
		return err
	}
	defer fileIn.Close()

	// Create the output file.
	fileOut, err := os.Create(dst.Absolute())
	if err != nil {
		return err
	}
	defer fileOut.Close()

	// Perform the copy and set the output file's privileges
	// to same as input file.
	_, err = io.Copy(fileOut, fileIn)
	if err == nil {
		si, err := os.Stat(src.Absolute())
		if err == nil {
			err = os.Chmod(dst.Absolute(), si.Mode())
		}
	}

	return err
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

// FileCompareEqual compares two files returning true
// if they are equal.
func FileCompareEqual(file1, file2 *Path) bool {
	var err 		error

	if !file1.IsPathRegularFile() {
		return false
	}

	if !file2.IsPathRegularFile() {
		return false
	}

	if file1.Size() != file2.Size() {
		return false
	}

	f1, err := os.Open(file1.Absolute())
	if err != nil {
		log.Fatal(err)
	}
	defer f1.Close()

	f2, err := os.Open(file2.Absolute())
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
//                             		Path
//----------------------------------------------------------------------------

// Path provides a centralized
type Path struct {
	str       	string
}

// Absolute returns the absolute file path for
// this path.
func (p *Path) Absolute( ) string {
	path := p.Clean()
	path, _ = filepath.Abs(path)
	return path
}

// Append a subdirectory or file name[.file extension] to the path.
// If the string is empty, then '/' will be appended.
func (p *Path) Append(s string) *Path {
	pth := Path{}
	pth.str = p.str + string(os.PathSeparator) + s
	pth.str = filepath.Clean(pth.str)
	return &pth
}

// Base returns the last component of the path. If the
// path is empty, "." is returned.
func (p *Path) Base( ) string {
	b := filepath.Base(p.str)
	return b
}

// Chmod changes the mode of the named file to the given mode.
func (p *Path) Chmod(mode os.FileMode) error {
	var err error

	b := p.Clean()
	if len(b) > 0 {
		err = os.Chmod(b, mode)
	}

	return err
}

// Clean cleans up the file path. It returns the absolute
// file path if needed.
func (p *Path) Clean( ) string {
	var path string

	if strings.HasPrefix(p.str, "~") {
		p.str = NewHomeDir().String() + string(os.PathSeparator) + p.str[1:]
	}
	p.str = os.ExpandEnv(p.str)
	p.str = filepath.Clean(p.str)
	path, _ = filepath.Abs(p.str)

	return path
}

// CreateDir assumes that this path represents a
// directory and creates it along with any parent
// directories needed as well.
func (p *Path) CreateDir( ) error {
	var err error

	b := p.Clean()
	if len(b) > 0 {
		err = os.MkdirAll(b, 0777)
	}

	return err
}

// DeleteFile assumes that this path represents a
// file and deletes it.
func (p *Path) DeleteFile( ) error {
	var err error

	pth := p.Clean()
	if len(pth) > 0 {
		fi, err := os.Lstat(pth)
		if err != nil {
			err = fmt.Errorf("Error: DeleteFile(): %s is not a file!\n", pth)
		} else {
			if fi.Mode().IsRegular() {
				err = os.Remove(pth)
			} else {
				err = fmt.Errorf("Error: DeleteFile(): %s is not a file!\n", pth)
			}
		}
	}

	return err
}

// IsPathDir cleans up the supplied file path
// and then checks the cleaned file path to see
// if it is an existing standard directory.
func (p *Path) IsPathDir( ) bool {
	var err error
	var pth string

	pth = p.Clean( )
	fi, err := os.Lstat(pth)
	if err != nil {
		return false
	}
	if fi.Mode().IsDir() {
		return true
	}
	return false
}

// IsPathRegularFile cleans up the supplied file path
// and then checks the cleaned file path to see
// if it is an existing standard file.
func (p *Path) IsPathRegularFile( ) bool {
	var err error
	var pth string

	pth = p.Clean()
	fi, err := os.Lstat(pth)
	if err != nil {
		return false
	}
	if fi.Mode().IsRegular() {
		return true
	}
	return false
}

func (p *Path) Mode( ) os.FileMode {
	var mode	os.FileMode

	si, err := os.Stat(p.Absolute())
	if err == nil {
		mode = si.Mode()
	}
	return mode
}

func (p *Path) ModTime( ) time.Time {
	var mod		time.Time

	si, err := os.Stat(p.Absolute())
	if err == nil {
		mod = si.ModTime()
	}
	return mod
}

// RemoveDir assumes that this path represents a
// directory and deletes it along with any parent
// directories that it can as well.
func (p *Path) RemoveDir( ) error {
	var err error

	b := p.Clean()
	if len(b) > 0 {
		err = os.RemoveAll(b)
	}

	return err
}

func (p *Path) SetStr(s string) {
	p.str = s
}

// Size returns length in bytes for regular files.
func (p *Path) Size( ) int64 {
	var size	int64

	si, err := os.Stat(p.Absolute())
	if err == nil {
		size = si.Size()
	}
	return size
}

func (p *Path) String( ) string {
	return p.str
}

// NewHomeDir returns the current working directory as a Path.
func NewHomeDir() *Path {
	p := Path{}

	// user.Current() returns nil if cross-compiled e.g. on mac for linux.
	// So, our backup is to get the home directory from the environment.
	if usr, _ := user.Current(); usr != nil {
		p.str = usr.HomeDir
	} else {
		p.str = os.Getenv("HOME")
	}

	return &p
}

func NewPath(s string) *Path {
	p := Path{}
	p.str = s
	return &p
}

// NewWorkDir returns the current working directory as a Path.
func NewWorkDir() *Path {
	p := Path{}
	p.str, _ = os.Getwd()
	return &p
}

// NewTempDir returns the temporary directory as a Path.
func NewTempDir() *Path {
	p := Path{}
	p.str = os.TempDir()
	return &p
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

//============================================================================
//                            		Workers
//============================================================================

type WorkQueue struct {
	queue		chan interface{}
	ack			chan bool
	done		chan bool
	doWork		func(interface{})
	complete	func()
	size		int
}

func (w *WorkQueue) CloseInput() {
	close(w.queue)
}

func (w *WorkQueue) complete() {
	w.done <- true
}

func (w *WorkQueue) PushWork(i interface{}) {
	w.queue <- i
}

func (w *WorkQueue) WaitForCompletion() {
	<-w.done
	w.Completed()
}

func NewWorkQueue(task func(interface{}), s int) *WorkQueue {
	// Set up the Work Queue.
	wq := &WorkQueue{}
	wq.doWork = task
	if s > 0 {
		wq.size = s
	} else {
		wq.size = runtime.NumCPU()
	}
	// Now set up the actual queues needed.
	wq.queue = make(chan interface{})
	wq.ack = make(chan bool)
	wq.done = make(chan bool)
	for i := 0; i < r; i++ {
		go func() {
			for {
				v, ok := <-input
				if ok {
					task(v)
				} else {
					wq.ack <- true
					return
				}
			}
		}()
	}
	go func() {
		for i := 0; i < r; i++ {
			<-wq.ack
		}
		wq.complete()
	}()

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

