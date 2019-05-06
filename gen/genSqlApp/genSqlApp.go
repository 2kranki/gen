// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Generate SQL Application programs in go

// Notes:
//	1.	The html and text templating systems require that
//		their data be separated since it is not identical.
//		So, we put them in separate files.
//	2.	The html and text templating systems access generic
//		structures with range, with, if.  They do not handle
//		structures well especially arrays of structures within
//		structures.

package genSqlApp

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
	"strings"
)


// FileDefn gives the parameters needed to generate a file.  The fields of
// the struct have been simplified to allow for easy json encoding/decoding.
type FileDefn struct {
	ModelName		string 		`json:"ModelName,omitempty"`
	FileDir			string		`json:"FileDir,omitempty"`		// Output File Directory
	FileName  		string 		`json:"FileName,omitempty"`		// Output File Name
	FileType  		string 		`json:"Type,omitempty"`  		// text, sql, html
	FilePerms		os.FileMode	`json:"FilePerms,omitempty"`	// Output File Permissions
	Class     		string 		`json:"Class,omitempty"` 		// single, table
	PerGrp  		int			`json:"PerGrp,omitempty"` 		// 0 == generate one file
	//															// 1 == generate one file for a database
	// 															// 2 == generate one file for a table
}

// FileDefns controls what files are generated.
var FileDefns []FileDefn = []FileDefn{
	{"base.html.tmpl.txt",
		"/tmpl",
		"base.html.tmpl",
		"copy",
		0644,
		"one",
		0,
	},
	{"bld.sh.txt",
		"",
		"bld.sh",
		"copy",
		0755,
		"one",
		0,
	},
	{"form.html.tmpl.txt",
		"/tmpl",
		"${DbName}.${TblName}.form.gohtml",
		"text",
		0644,
		"one",
		2,
	},
	{"main.go.tmpl.txt",
		"",
		"main.go",
		"text",
		0644,
		"one",
		0,
	},
	{"mainExec.go.tmpl.txt",
		"",
		"mainExec.go",
		"text",
		0644,
		"single",
		0,
	},
	{"handlers.go.tmpl.txt",
		"/hndlr${DbName}",
		"hndlr${DbName}.go",
		"text",
		0644,
		"single",
		0,
	},
	{"handlers_test.go.tmpl.txt",
		"/hndlr${DbName}",
		"hndlr${DbName}_test.go",
		"text",
		0644,
		"single",
		0,
	},
	{"handlers.table.go.tmpl.txt",
		"/hndlr${DbName}",
		"${TblName}.go",
		"text",
		0644,
		"single",
		2,
	},
	{"io.go.tmpl.txt",
		"/io${DbName}",
		"io${DbName}.go",
		"text",
		0644,
		"single",
		0,
	},
	{"io_test.go.tmpl.txt",
		"/io${DbName}",
		"io${DbName}_test.go",
		"text",
		0644,
		"single",
		0,
	},
	{"io.table.go.tmpl.txt",
		"/io${DbName}",
		"${TblName}.go",
		"text",
		0644,
		"single",
		2,
	},
}

// TmplData is used to centralize all the inputs
// to the generators.  We maintain generic JSON
// structures for the templating system which does
// not support structs.  (Not certain why yet.)
// We also maintain the data in structs for easier
// access by the generation functions.
type TmplData struct {
	Data     	*dbPkg.Database
	Main     	*mainData.MainData
	Table		*dbPkg.DbTable
}

var tmplData TmplData

type TaskData struct {
	FD			*FileDefn
	TD			*TmplData
	Table		*dbPkg.DbTable
	PathIn	  	string						// Input File Path
	PathOut	  	string						// Output File Path

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
				os.Chmod(t.PathOut, t.FD.FilePerms)
				if !sharedData.Quiet() {
					log.Printf("\tCopied %d bytes from %s to %s\n", amt, t.PathIn, t.PathOut)
				}
			} else {
				log.Fatalf("Error - Copied %d bytes from %s to %s with error %s\n",
					amt, t.PathIn, t.PathOut, err)
			}
		}
	case "html":
		if err = GenHtmlFile(t.PathIn, t.PathOut, t); err == nil {
			os.Chmod(t.PathOut, t.FD.FilePerms)
			if !sharedData.Quiet() {
				log.Printf("\tGenerated HTML from %s to %s\n", t.PathIn, t.PathOut)
			}
		} else {
			log.Fatalf("Error - Generated HTML from %s to %s with error %s\n",
				t.PathIn, t.PathOut, err)
		}
	case "text":
		if err = GenTextFile(t.PathIn, t.PathOut, t); err == nil {
			os.Chmod(t.PathOut, t.FD.FilePerms)
			if !sharedData.Quiet() {
				log.Printf("\tGenerated HTML from %s to %s\n", t.PathIn, t.PathOut)
			}
		} else {
			log.Fatalf("Error - Generated HTML from %s to %s with error %s\n",
				t.PathIn, t.PathOut, err)
		}
	default:
		log.Fatalln("Error: Invalid file type:", t.FD.FileType, "for", t.FD.ModelName, err)
	}


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

func createOutputPath(dir string, dn string, tn string, fn string) (string, error) {
	var outPath string
	var err error

	outPath = sharedData.OutDir()
	outPath += "/"
	if len(dir) > 0 {
		outPath += dir
		outPath += "/"
	}
	outPath += fn
	if len(dn) > 0 {
		outPath = strings.ReplaceAll(outPath, "${DbName}", strings.Title(dn))
	}
	if len(tn) > 0 {
		outPath = strings.ReplaceAll(outPath, "${TblName}", strings.Title(tn))
	}
	if sharedData.Debug() && strings.Contains(outPath, "${DbName}") {
		log.Fatalf("Error: output path, %s, contains $DbName request!.  args: %q\n", outPath)
	}
	if sharedData.Debug() && strings.Contains(outPath, "${TblName}") {
		log.Fatalf("Error: output path, %s, contains $TblName request!.  args: %q\n", outPath)
	}
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

func GenSqlApp(inDefns map[string]interface{}) error {
	var err 	error
	var pathIn	string
	//var ok 		bool

	if sharedData.Debug() {
		log.Println("\tsql_app: In Debug Mode")
		log.Printf("\t  args: %q\n", flag.Args())
	}

    // Read the JSON files.
    if err = readJsonFiles(); err != nil {
		log.Fatalln(err)
    }

	// Set up template data
	//if tmplData.MainJson, ok = mainData.MainJson().(map[string]interface{}); !ok {
		//log.Fatalln("Error - Could not type assert mainData.MainJson()")
	//}
	tmplData.Main = mainData.MainStruct()
	//if tmplData.DataJson, ok = appData.AppJson().(map[string]interface{}); !ok {
	//	log.Fatalln("Error - Could not type assert appData.AppJson()")
	//}
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
        tmpName = path.Clean(sharedData.OutDir() + "/hndlr" + tmplData.Data.TitledName())
        if err = os.MkdirAll(tmpName, os.ModeDir+0777); err != nil {
            log.Fatalln("Error: Could not create output directory:", tmpName, err)
        }
        tmpName = path.Clean(sharedData.OutDir() + "/io" + tmplData.Data.TitledName())
        if err = os.MkdirAll(tmpName, os.ModeDir+0777); err != nil {
            log.Fatalln("Error: Could not create output directory:", tmpName, err)
        }
    }

	// Setup the worker queue.
	done := make(chan bool)
	inputQueue := util.Workers(
					func(a interface{}) {
						var t		TaskData
						t = a.(TaskData)
						t.genFile()
					},
					func() {
						done <- true
					},
					5)

	for i, def := range FileDefns {

		if !sharedData.Quiet() {
			log.Println("Process file:", def.ModelName, "generating:", def.FileName, "...")
		}

		// Create the input model file path.
		if pathIn, err = createModelPath(def.ModelName); err != nil {
			return errors.New(fmt.Sprintln("Error:", pathIn, err))
		}
		if sharedData.Debug() {
			log.Println("\t\tmodelPath=", pathIn)
		}

		// Now generate the file.
		switch def.PerGrp {
		case 0:
			// Standard File
			data := TaskData{FD:&FileDefns[i], TD:&tmplData, PathIn:pathIn}
			// Create the output path
			if data.PathOut, err = createOutputPath(def.FileDir, tmplData.Data.Name, "", def.FileName); err != nil {
				log.Fatalln(err)
			}
			if sharedData.Debug() {
				log.Println("\t\t outPath=", data.PathOut)
			}
			// Generate the file.
			inputQueue <- data
		case 2:
			// Output File is Titled Table Name in Titled Database Name directory
			dbPkg.ForTables(
				func(v *dbPkg.DbTable) {
					data := TaskData{FD:&FileDefns[i], TD:&tmplData, Table:v, PathIn:pathIn}
					if data.PathOut, err = createOutputPath(def.FileDir, tmplData.Data.Name, v.Name, def.FileName); err != nil {
						log.Fatalln(err)
					}
					if sharedData.Debug() {
						log.Println("\t\t outPath=", data.PathOut)
					}
					// Generate the file.
					inputQueue <- data
				})
		default:
			log.Printf("Skipped %s because of type!\n", def.FileName)
		}
	}
	close(inputQueue)

	<-done

	return nil
}
