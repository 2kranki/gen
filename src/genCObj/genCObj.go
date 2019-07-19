// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Generate C Object

// Notes:
//	1.	The html and text templating systems require that
//		their data be separated since it is not identical.
//		So, we put them in separate files.
//	2.	The html and text templating systems access generic
//		structures with range, with, if.  They do not handle
//		structures well especially arrays of structures within
//		structures.

package genCObj

import (
	"../genCmn"
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

const (
	jsonDirCon = "./"
	// Merged from main.go
	cmdId     = "cmd"
	jsonDirId = "jsondir"
	nameId    = "name"
	timeId    = "time"
)

// FileDefns controls what files are generated.
var FileDefns []genCmn.FileDefn = []genCmnFileDefn{
	{"obj_int_h.txt",
		"src",
		"${Name}_internal.h",
		"text",
		0644,
	},
	{"obj_obj_c.txt",
		"src",
		"${Name}_object.c",
		"text",
		0644,
	},
	{"obj_c.txt",
		"src",
		"${Name}.c",
		"text",
		0644,
	},
	{"obj_h.txt",
		"src",
		"${Name}.h",
		"text",
		0644,
	},
	{"obj_test_c.txt",
		"tests",
		"test_${Name}.c",
		"text",
		0644,
	},
}

// TmplData is used to centralize all the inputs
// to the generators.  We maintain generic JSON
// structures for the templating system which does
// not support structs.  (Not certain why yet.)
// We also maintain the data in structs for easier
// access by the generation functions.
type TmplData struct {
	Data     *DbObject
}

var tmplData TmplData

func init() {

}

func copyFile(modelPath, outPath string) (int64, error) {
	var dst *os.File
	var err error
	var src *os.File

	if _, err = util.IsPathRegularFile(modelPath); err != nil {
		return 0, errors.New(fmt.Sprint("Error - model file does not exist:", modelPath, err))
	}

	if outPath, err = util.IsPathRegularFile(outPath); err == nil {
		if sharedData.Force() {
			if err = os.Remove(outPath); err != nil {
				return 0, errors.New(fmt.Sprint("Error - could not delete:", outPath, err))
			}
		} else {
			return 0, errors.New(fmt.Sprint("Error - overwrite error of:", outPath))
		}
	}
	if dst, err = os.Create(outPath); err != nil {
		return 0, errors.New(fmt.Sprint("Error - could not create:", outPath, err))
	}
	defer dst.Close()

	if src, err = os.Open(modelPath); err != nil {
		return 0, errors.New(fmt.Sprint("Error - could not open model file:", modelPath, err))
	}
	defer src.Close()

	amt, err := io.Copy(dst, src)

	return amt, err
}

func createModelPath(fn string) (string, error) {
	var modelPath   string
	var err         error

	if modelPath, err = util.IsPathRegularFile(fn); err == nil {
		return modelPath, err
	}

	// Calculate the model path.
	modelPath = sharedData.MdlDir()
	modelPath += "/cobj/"
	modelPath += fn
	modelPath, err = util.IsPathRegularFile(modelPath)

	return modelPath, err
}

func createOutputPath(dn,fn string) (string, error) {
	var outPath 	string
	var err 		error

	if len(fn) > 0 {
		fn = strings.Replace(fn, "${Name}", tmplData.Data.Name, -1)
	}

	outPath = sharedData.OutDir()
	outPath += "/"
	outPath += dn
	outPath += "/"
	outPath += fn
	outPath, err = util.IsPathRegularFile(outPath)
	if err == nil {
		if !sharedData.Force() {
			return outPath, errors.New(fmt.Sprint("Over-write error of:", outPath))
		}
	}

	return outPath, nil
}

func GenCObj(inDefns map[string]interface{}) error {
	var err 	error

	if sharedData.Debug() {
		log.Println("GenCObj: In Debug Mode...")
		log.Printf("\t  args: %q\n", flag.Args())
	}

	// Set up template data
	tmplData.Data = DbStruct()

	// Read the JSON file.
	if sharedData.Debug() {
		log.Println("\tDataPath:", sharedData.DataPath())
	}
	if err = ReadJsonFile(sharedData.DataPath()); err != nil {
		log.Fatalln(errors.New(fmt.Sprintln("Error: Reading Main Json Input:", sharedData.DataPath(), err)))
	}

	if sharedData.OutDir() == "" {
		log.Fatalf("Error - 'libPath' cli argument is required!\n\n\n")
	}

	// Now handle each FileDefn creating a file for it.
	for _, def := range (FileDefns) {
		var modelPath string
		var outPath string

		if !sharedData.Quiet() {
			log.Println("Process file:", def.ModelName, "generating:", def.FileName, "...")
		}

		// Create the input model file path.
		if modelPath, err = createModelPath(def.ModelName); err != nil {
			return errors.New(fmt.Sprintln("Error:", modelPath, err))
		}
		if sharedData.Debug() {
			log.Println("\t\tmodelPath=", modelPath)
		}

		// Create the output path
		if outPath, err = createOutputPath(def.FileDir, def.FileName); err != nil {
			log.Fatalln(err)
		}
		if sharedData.Debug() {
			log.Println("\t\t outPath=", outPath)
		}

		// Now generate the file.
		switch def.FileType {
		case "copy":
			if sharedData.Noop() {
				if !sharedData.Quiet() {
					log.Printf("\tShould have copied from %s to %s\n", modelPath, outPath)
				}
			} else {
				if amt, err := copyFile(modelPath, outPath); err == nil {
					if !sharedData.Quiet() {
						log.Printf("\tCopied %d bytes from %s to %s\n", amt, modelPath, outPath)
					}
				} else {
					log.Fatalf("Error - Copied %d bytes from %s to %s with error %s\n",
						amt, modelPath, outPath, err)
				}
			}
		case "text":
			if err = GenTextFile(modelPath, outPath, tmplData); err == nil {
				if !sharedData.Quiet() {
					log.Printf("\tGenerated HTML from %s to %s\n", modelPath, outPath)
				}
			} else {
				log.Fatalf("Error - Generated HTML from %s to %s with error %s\n",
					modelPath, outPath, err)
			}
		default:
			return errors.New(fmt.Sprint("Error: Invalid file type:", def.FileType,
				"for", def.ModelName, err))
		}
	}

	return nil
}
