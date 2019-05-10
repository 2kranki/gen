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
	"../dbPkg"
	"../mainData"
	"../shared"
	"../util"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
)

const (
	jsonDirCon = "./"
	// Merged from main.go
	cmdId     = "cmd"
	jsonDirId = "jsondir"
	nameId    = "name"
	timeId    = "time"
)

// FileDefn gives the parameters needed to generate a file.  The fields of
// the struct have been simplified to allow for easy json encoding/decoding.
type FileDefn struct {
	ModelName string `json:"ModelName,omitempty"`
	FileName  string `json:"FileName,omitempty"`
	FileType  string `json:"Type,omitempty"`  // text, sql, html
	Class     string `json:"Class,omitempty"` // single, table
}

// FileDefns controls what files are generated.
var FileDefns []FileDefn = []FileDefn{
	{"base.html.tmpl.txt",
		"/tmpl/base.html.tmpl",
		"copy",
		"one",
	},
	{"main.go.tmpl.txt",
		"main.go",
		"text",
		"one",
	},
	{"mainExec.go.tmpl.txt",
		"mainExec.go",
		"text",
		"single",
	},
	{"handlers.go.tmpl.txt",
		"/handlers/handlers.go",
		"text",
		"single",
	},
	{"tableio.go.tmpl.txt",
		"/tableio/tableio.go",
		"text",
		"single",
	},
}

// TmplData is used to centralize all the inputs
// to the generators.  We maintain generic JSON
// structures for the templating system which does
// not support structs.  (Not certain why yet.)
// We also maintain the data in structs for easier
// access by the generation functions.
type TmplData struct {
	Data     *dbPkg.Database
	Main     *mainData.MainData
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
	modelPath += "/sqlapp"
	modelPath += "/"
	modelPath += fn
	//modelPath = filepath.Clean(modelPath)				// Clean() is part of Abs()
	modelPath, err = filepath.Abs(modelPath)

	return modelPath, err
}

func createOutputPath(fn string) (string, error) {
	var outPath string
	var err error

	outPath = sharedData.OutDir()
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

// readJsonFiles reads in the two JSON files that define the
// application to be generated.
func readJsonFiles() error {
	var err error

	if err = mainData.ReadJsonFileMain(sharedData.MainPath()); err != nil {
		return errors.New(fmt.Sprintln("Error: Reading Main Json Input:", sharedData.MainPath(), err))
	}

	if err = dbPkg.ReadJsonFile(sharedData.DataPath()); err != nil {
		return errors.New(fmt.Sprintln("Error: Reading Main Json Input:", sharedData.DataPath(), err))
	}

    return nil
}

func GenCObj(inDefns map[string]interface{}) error {
	var err error

	if sharedData.Debug() {
		log.Println("\tsql_app: In Debug Mode")
		log.Printf("\t  args: %q\n", flag.Args())
	}

    // Read the JSON files.
    if err = readJsonFiles(); err != nil {
		log.Fatalln(err)
    }

	// Set up template data
	tmplData.Main = mainData.MainStruct()
	tmplData.Data = dbPkg.DbStruct()

	// Set up the output directory structure
    if !sharedData.Noop() {
        tmpName := path.Clean(sharedData.OutDir())
        if err = os.RemoveAll(tmpName); err != nil {
            log.Fatalln("Error: Could not remove output directory:", tmpName, err)
        }
        tmpName = path.Clean(sharedData.OutDir() + "/tmpl")
        if err = os.MkdirAll(tmpName, os.ModeDir+0777); err != nil {
            log.Fatalln("Error: Could not create output directory:", tmpName, err)
        }
        tmpName = path.Clean(sharedData.OutDir() + "/handlers")
        if err = os.MkdirAll(tmpName, os.ModeDir+0777); err != nil {
            log.Fatalln("Error: Could not create output directory:", tmpName, err)
        }
        tmpName = path.Clean(sharedData.OutDir() + "/tableio")
        if err = os.MkdirAll(tmpName, os.ModeDir+0777); err != nil {
            log.Fatalln("Error: Could not create output directory:", tmpName, err)
        }
    }

	// Now handle each FileDefn creating a file for it.
	for _, def := range (FileDefns) {
		var modelPath string
		var outPath string

		if !sharedData.Quiet() {
			log.Println("Process file:", def.ModelName, "generating:", def.FileName, "...")
		}

		//if def.PreprocSql {
		//PreprocessSql(json.(DatabaseData))
		//}

		// Create the input model file path.
		if modelPath, err = createModelPath(def.ModelName); err != nil {
			return errors.New(fmt.Sprintln("Error:", modelPath, err))
		}
		if sharedData.Debug() {
			log.Println("\t\tmodelPath=", modelPath)
		}

		// Create the output path
		if outPath, err = createOutputPath(def.FileName); err != nil {
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
