// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Generate Tasks

// Notes:
//	1.	The html and text templating systems require that
//		their data be separated since it is not identical.
//		So, we put them in separate files.
//	2.	The html and text templating systems access generic
//		structures with range, with, if.  They do not handle
//		structures well especially arrays of structures within
//		structures.


package genCmn

import (
	"../mainData"
	"../shared"
	"../util"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)


//----------------------------------------------------------------------------
//								Task Data
//----------------------------------------------------------------------------

type TaskData struct {
	FD			*FileDefn
	TD			*TmplData
	PathIn	  	*util.Path					// Input File Path
	PathOut	  	*util.Path					// Output File Path
	Data		interface{}
	Table		interface{}

}

func (t *TaskData) genFile() {
	var err         error

	// Now generate the file.
	switch t.FD.FileType {
	case "copy":
		if sharedData.Noop() {
			if !sharedData.Quiet() {
				log.Printf("\tShould have copied from %s to %s\n", t.PathIn, t.PathOut)
			}
		} else {
			if amt, err := copyFile(t.PathIn, t.PathOut); err == nil {
				t.PathOut.Chmod(t.FD.FilePerms)
				if !sharedData.Quiet() {
					log.Printf("\tCopied %d bytes from %s to %s\n", amt, t.PathIn, t.PathOut)
				}
			} else {
				log.Fatalf("Error - Copied %d bytes from %s to %s with error %s\n",
					amt, t.PathIn, t.PathOut, err)
			}
		}
	case "copyDir":
		if sharedData.Noop() {
			if !sharedData.Quiet() {
				log.Printf("\tShould have copied directory from %s to %s\n", t.PathIn, t.PathOut)
			}
		} else {
			if err := copyDir(t.PathIn, t.PathOut); err == nil {
				if !sharedData.Quiet() {
					log.Printf("\tCopied from %s to %s\n", t.PathIn, t.PathOut)
				}
			} else {
				log.Fatalf("Error - Copied from %s to %s with error %s\n",
					t.PathIn, t.PathOut, err)
			}
		}
	case "html":
		if err = GenHtmlFile(t.PathIn, t.PathOut, t); err == nil {
			t.PathOut.Chmod(t.FD.FilePerms)
			if !sharedData.Quiet() {
				log.Printf("\tGenerated HTML from %s to %s\n", t.PathIn, t.PathOut)
			}
		} else {
			log.Fatalf("Error - Generated HTML from %s to %s with error %s\n",
				t.PathIn, t.PathOut, err)
		}
	case "text":
		if err = GenTextFile(t.PathIn, t.PathOut, t); err == nil {
			t.PathOut.Chmod(t.FD.FilePerms)
			if !sharedData.Quiet() {
				log.Printf("\tGenerated text from %s to %s\n", t.PathIn, t.PathOut)
			}
		} else {
			log.Fatalf("Error - Generated text from %s to %s with error %s\n",
				t.PathIn, t.PathOut, err)
		}
	default:
		log.Fatalln("Error: Invalid file type:", t.FD.FileType, "for", t.FD.ModelName, err)
	}


}

//----------------------------------------------------------------------------
//								copyDir
//----------------------------------------------------------------------------

func (t *TaskData) copyDir(modelPath, outPath *util.Path) error {
	var err 	error
	var base	string
	var pathIn	*util.Path
	var pathOut	*util.Path

	if !modelPath.IsPathDir( ) {
		return fmt.Errorf("Error - model directory, %s, does not exist!\n", pathIn.String())
	}
	base = modelPath.Base( )
	if len(base) == 0 {
		return fmt.Errorf("Error - model directory, %s, does not have base directory!\n", pathIn.String())
	}

	pathOut = outPath.Append(base)
	log.Printf("\tcopyDir:  inPath: %s\n", pathIn.String())
	log.Printf("\tcopyDir: outPath: %s base: %s\n", pathOut.String(), base)
	if outPath.IsPathDir( ) {
		if sharedData.Replace() {
			log.Printf("\tcopyDir: Removing %s\n", pathOut.String())
			if err = pathOut.RemoveDir( ); err != nil {
				return fmt.Errorf("Error - could not delete %s: %s\n", pathOut.String(), err.Error())
			}
		} else {
			return fmt.Errorf("Error - overwrite error of %s\n", pathOut.String())
		}
	}

	err = util.CopyDir(pathIn, pathOut)

	return err
}

//----------------------------------------------------------------------------
//								copyFile
//----------------------------------------------------------------------------

func (t *TaskData) copyFile(modelPath, outPath *util.Path) (int64, error) {
	var dst *os.File
	var err error
	var src *os.File

	if !modelPath.IsPathRegularFile( ) {
		return 0, fmt.Errorf("Error - model file does not exist for %s: %s\n", modelPath.String(), err.Error())
	}

	if outPath.IsPathRegularFile( ) {
		if !sharedData.Replace() {
			return 0, fmt.Errorf("Error - overwrite error of %s\n", outPath.String())
		}
	}
	if dst, err = os.Create(outPath.Absolute()); err != nil {
		return 0, fmt.Errorf("Error - could not create %s: %s\n", outPath.String(), err.Error())
	}
	defer dst.Close()

	if src, err = os.Open(modelPath.Absolute()); err != nil {
		return 0, fmt.Errorf("Error - could not open model file, %s: %s\n", modelPath.String(), err.Error())
	}
	defer src.Close()

	amt, err := io.Copy(dst, src)

	return amt, err
}


