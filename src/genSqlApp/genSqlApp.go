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

/* 				*** Package Structure ***
genSqlApp
	- Responsible for generating the application
	- Is controlled by internal table of files to generate
		and support routines to generate each file given its
		definition
	- Sub-packages
		- dbPlugin
			- Provides Managerial functions for all plugins:
				- Plugin registration
				- Plugin unregistration
				- Finding Registered plugins
			- Defines the base struct that all plugins must support
				and inherit from
		- dbJson
			- Responsible for importing, set up and validation of the
				user defined JSON definition file.
			- Constains JSON tables that define the databases, tables
				and fields providing access to same.
		- dbSql
			- Responsible for generating the SQL statements that will
				work in the appropriate database plugin
		- dbForm
			- Responsible for generating the HTML form data needed to
				access/manage the database
 */


package genSqlApp

import (
	"../mainData"
	"../shared"
	"../util"
	_ "./dbForm"
	"./dbJson"
	_ "./dbGener"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	// Include the various Database Plugins so that they will register
	// with dbPlugin.
	_ "./dbMariadb"
	_ "./dbMssql"
	_ "./dbMysql"
	_ "./dbPostgres"
	_ "./dbSqlite"
)


// FileDefn gives the parameters needed to generate a file.  The fields of
// the struct have been simplified to allow for easy json encoding/decoding.
type FileDefn struct {
	ModelName		string 		`json:"ModelName,omitempty"`
	FileDir			string		`json:"FileDir,omitempty"`		// Output File Directory
	FileName  		string 		`json:"FileName,omitempty"`		// Output File Name
	FileType  		string 		`json:"Type,omitempty"`  		// copy, text, sql, html
	FilePerms		os.FileMode	`json:"FilePerms,omitempty"`	// Output File Permissions
	Class     		string 		`json:"Class,omitempty"` 		// single, table
	PerGrp  		int			`json:"PerGrp,omitempty"` 		// 0 == generate one file
	//															// 1 == generate one file for a database
	// 															// 2 == generate one file for a table
}

// FileDefns controls what files are generated.
var FileDefns []FileDefn = []FileDefn{
	{"bld.sh.txt",
		"",
		"b.sh",
		"text",
		0755,
		"one",
		0,
	},
	{"static.txt",
		"/static",
		"README.txt",
		"copy",
		0644,
		"one",
		0,
	},
	{"tst.sh.txt",
		"/src",
		"t.sh",
		"copy",
		0755,
		"one",
		0,
	},
	{"tst.sh.txt",
		"/src/hndlr${DbName}",
		"t.sh",
		"copy",
		0755,
		"one",
		0,
	},
	{"tst.sh.txt",
		"/src/io${DbName}",
		"t.sh",
		"copy",
		0755,
		"one",
		0,
	},
	{"form.html",
		"/html",
		"form.html",
		"copy",
		0644,
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
	{"list.html.tmpl.txt",
		"/tmpl",
		"${DbName}.${TblName}.list.gohtml",
		"text",
		0644,
		"one",
		2,
	},
	{"main.menu.html.tmpl.txt",
		"/html",
		"${DbName}.menu.html",
		"text",
		0644,
		"one",
		0,
	},
	{"main.go.tmpl.txt",
		"/src",
		"main.go",
		"text",
		0644,
		"one",
		0,
	},
	{"mainExec.go.tmpl.txt",
		"/src",
		"mainExec.go",
		"text",
		0644,
		"single",
		0,
	},
	{"handlers.go.tmpl.txt",
		"/src/hndlr${DbName}",
		"hndlr${DbName}.go",
		"text",
		0644,
		"single",
		0,
	},
	{"handlers_test.go.tmpl.txt",
		"/src/hndlr${DbName}",
		"hndlr${DbName}_test.go",
		"text",
		0644,
		"single",
		0,
	},
	{"handlers.table.go.tmpl.txt",
		"/src/hndlr${DbName}",
		"${TblName}.go",
		"text",
		0644,
		"single",
		2,
	},
	{"io.go.tmpl.txt",
		"/src/io${DbName}",
		"io${DbName}.go",
		"text",
		0644,
		"single",
		0,
	},
	{"io_test.go.tmpl.txt",
		"/src/io${DbName}",
		"io${DbName}_test.go",
		"text",
		0644,
		"single",
		0,
	},
	{"io.table.go.tmpl.txt",
		"/src/io${DbName}",
		"${TblName}.go",
		"text",
		0644,
		"single",
		2,
	},
	{"util.go.txt",
		"/src/util",
		"util.go",
		"text",
		0644,
		"single",
		0,
	},
}

// TmplData is used to centralize all the inputs
// to the generators.  We maintain generic JSON
// structures for the templating system which does
// not support structs.  (Not certain why yet.)
// We also maintain the data in structs for easier
// access by the generation functions.
type TmplData struct {
	Data     	*dbJson.Database
	Main     	*mainData.MainData
	Table		*dbJson.DbTable
}

var tmplData TmplData

type TaskData struct {
	FD			*FileDefn
	TD			*TmplData
	Table		*dbJson.DbTable
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

//----------------------------------------------------------------------------
//								copyFile
//----------------------------------------------------------------------------

func copyFile(modelPath, outPath string) (int64, error) {
	var dst *os.File
	var err error
	var src *os.File

	if _, err = util.IsPathRegularFile(modelPath); err != nil {
		return 0, errors.New(fmt.Sprint("Error - model file does not exist:", modelPath, err))
	}

	if outPath, err = util.IsPathRegularFile(outPath); err == nil {
		if sharedData.Replace() {
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

//----------------------------------------------------------------------------
//								createModelPath
//----------------------------------------------------------------------------

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
	modelPath, err = util.IsPathRegularFile(modelPath)

	return modelPath, err
}

//----------------------------------------------------------------------------
//								createOutputPath
//----------------------------------------------------------------------------

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
		outPath = strings.Replace(outPath, "${DbName}", strings.Title(dn), -1)
	}
	if len(tn) > 0 {
		outPath = strings.Replace(outPath, "${TblName}", strings.Title(tn), -1)
	}
	if sharedData.Debug() && strings.Contains(outPath, "${DbName}") {
		log.Fatalf("Error: output path, %s, contains $DbName request!.  args: %q\n", outPath)
	}
	if sharedData.Debug() && strings.Contains(outPath, "${TblName}") {
		log.Fatalf("Error: output path, %s, contains $TblName request!.  args: %q\n", outPath)
	}
	outPath, err = util.IsPathRegularFile(outPath)
	if err == nil {
		if !sharedData.Replace() {
			return outPath, errors.New(fmt.Sprint("Over-write error of:", outPath))
		}
	}

	return outPath, nil
}

//----------------------------------------------------------------------------
//								readJsonFiles
//----------------------------------------------------------------------------

// readJsonFiles reads in the two JSON files that define the
// application to be generated.
func readJsonFiles() error {
	var err error

	if err = mainData.ReadJsonFileMain(sharedData.MainPath()); err != nil {
		return errors.New(fmt.Sprintln("Error: Reading Main Json Input:", sharedData.MainPath(), err))
	}

	if err = dbJson.ReadJsonFile(sharedData.DataPath()); err != nil {
		return errors.New(fmt.Sprintln("Error: Reading Data Json Input:", sharedData.DataPath(), err))
	}

    return nil
}

func GenSqlApp(inDefns map[string]interface{}) error {
	var err 	error
	var pathIn	string
	//var ok 		bool

	if sharedData.Debug() {
		log.Println("\t sql_app: In Debug Mode")
		log.Printf("\t  args: %q\n", flag.Args())
		log.Printf("\tmdldir: %s\n", sharedData.MdlDir())
	}

    // Read the JSON files.
    if err = readJsonFiles(); err != nil {
		log.Fatalln(err)
    }

	// Set up template data
	tmplData.Main = mainData.MainStruct()
	tmplData.Data = dbJson.DbStruct()

	// Set up the output directory structure
    if !sharedData.Noop() {
        tmpName := path.Clean(sharedData.OutDir())
        // We only delete main directory if forced to. Otherwise, we
        // will simply replace our files within it.
        if sharedData.Force() {
			if err = os.RemoveAll(tmpName); err != nil {
				log.Fatalln("Error: Could not remove output directory:", tmpName, err)
			}
		}
		tmpName = path.Clean(sharedData.OutDir() + "/html")
		if err = os.MkdirAll(tmpName, os.ModeDir+0777); err != nil {
			log.Fatalln("Error: Could not create output directory:", tmpName, err)
		}
		tmpName = path.Clean(sharedData.OutDir() + "/static")
		if err = os.MkdirAll(tmpName, os.ModeDir+0777); err != nil {
			log.Fatalln("Error: Could not create output directory:", tmpName, err)
		}
		tmpName = path.Clean(sharedData.OutDir() + "/style")
		if err = os.MkdirAll(tmpName, os.ModeDir+0777); err != nil {
			log.Fatalln("Error: Could not create output directory:", tmpName, err)
		}
        tmpName = path.Clean(sharedData.OutDir() + "/tmpl")
        if err = os.MkdirAll(tmpName, os.ModeDir+0777); err != nil {
            log.Fatalln("Error: Could not create output directory:", tmpName, err)
        }
        tmpName = path.Clean(sharedData.OutDir() + "/src/hndlr" + tmplData.Data.TitledName())
        if err = os.MkdirAll(tmpName, os.ModeDir+0777); err != nil {
            log.Fatalln("Error: Could not create output directory:", tmpName, err)
        }
        tmpName = path.Clean(sharedData.OutDir() + "/src/io" + tmplData.Data.TitledName())
        if err = os.MkdirAll(tmpName, os.ModeDir+0777); err != nil {
            log.Fatalln("Error: Could not create output directory:", tmpName, err)
        }
		tmpName = path.Clean(sharedData.OutDir() + "/src/util")
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
			dbJson.DbStruct().ForTables(
				func(v *dbJson.DbTable) {
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
