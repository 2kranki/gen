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
	_ "./dbGener"
	"./dbJson"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
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
	{"tst.sh.txt",
		"",
		"t.sh",
		"copy",
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
	{"handlers.test.go.tmpl.txt",
		"/src/hndlr${DbName}",
		"hndlr${DbName}_test.go",
		"text",
		0644,
		"single",
		0,
	},
	{"table.go.tmpl.txt",
		"/src/${DbName}${TblName}",
		"${DbName}${TblName}.go",
		"text",
		0644,
		"single",
		2,
	},
	{"table.test.go.tmpl.txt",
		"/src/${DbName}${TblName}",
		"${DbName}${TblName}_test.go",
		"text",
		0644,
		"single",
		2,
	},
	{"tst.sh.txt",
		"/src/${DbName}${TblName}",
		"t.sh",
		"copy",
		0755,
		"single",
		2,
	},
	{"handlers.table.go.tmpl.txt",
		"/src/hndlr${DbName}${TblName}",
		"${DbName}${TblName}.go",
		"text",
		0644,
		"single",
		2,
	},
	{"handlers.table.test.go.tmpl.txt",
		"/src/hndlr${DbName}${TblName}",
		"${DbName}${TblName}_test.go",
		"text",
		0644,
		"single",
		2,
	},
	{"handlers.table.test.fakedb.go.tmpl.txt",
		"/src/hndlr${DbName}${TblName}",
		"${DbName}${TblName}FakeDB.go",
		"text",
		0644,
		"single",
		2,
	},
	{"tst.sh.txt",
		"/src/hndlr${DbName}${TblName}",
		"t.sh",
		"copy",
		0755,
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
		"/src/io${DbName}${TblName}",
		"${TblName}.go",
		"text",
		0644,
		"single",
		2,
	},
	{"io.table.test.go.tmpl.txt",
		"/src/io${DbName}${TblName}",
		"${TblName}_test.go",
		"text",
		0644,
		"single",
		2,
	},
	{"tst.sh.txt",
		"/src/io${DbName}${TblName}",
		"t.sh",
		"copy",
		0755,
		"single",
		2,
	},
}

var FileDefns2 []FileDefn = []FileDefn{
	{"docker",
		"/src/",
		"",
		"copyDir",
		0644,
		"single",
		0,
	},
	{"util",
		"/src/",
		"",
		"copyDir",
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
	PathIn	  	util.Path					// Input File Path
	PathOut	  	util.Path					// Output File Path

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

func copyDir(modelPath, outPath util.Path) error {
	var err 	error
	var base	string
	var pathIn	util.Path
	var pathOut	util.Path

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

	err = util.CopyDir(pathIn.String(), pathOut.String())

	return err
}

//----------------------------------------------------------------------------
//								copyFile
//----------------------------------------------------------------------------

func copyFile(modelPath, outPath util.Path) (int64, error) {
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

//----------------------------------------------------------------------------
//								createModelPath
//----------------------------------------------------------------------------

// createModelPath creates an input path from our models and verifies that it
// exists.  Regular files and directories that are found in the "models"
// directory are acceptable.
func createModelPath(fn string) (util.Path, error) {
	var modelPath   string
	var err         error

	// Calculate the model path.
	modelPath = sharedData.MdlDir()
	modelPath += "/sqlapp"
	modelPath += "/"
	modelPath += fn
	modelPath, err = util.IsPathRegularFile(modelPath)

	if err != nil {
		modelPath2, err2 := util.IsPathDir(modelPath)
		if err2 == nil {
			modelPath = modelPath2
			err = nil
		}
	}

	return util.NewPath(modelPath), err
}

//----------------------------------------------------------------------------
//								createOutputDir
//----------------------------------------------------------------------------

// createOutputDir creates the output directory on disk given a
// subdirectory (dir).
func createOutputDir(dir, dn, tn string) error {
	var err error
	var outPath string
	var pathOut util.Path

	outPath = sharedData.OutDir()
	outPath += "/"
	if len(dir) > 0 {
		outPath += dir
	}
	mapper := func(placeholderName string) string {
		switch placeholderName {
		case "DbName":
			if len(dn) > 0 {
				return strings.Title(dn)
			}
		case "TblName":
			if len(tn) > 0 {
				return strings.Title(tn)
			}
		}
		return ""
	}
	outPath = os.Expand(outPath, mapper)

	pathOut = util.NewPath(outPath)
	if !pathOut.IsPathDir() {
		log.Printf("\t\tCreating directory: %s...\n", pathOut.String())
		err = pathOut.CreateDir()
	}

	return err
}

//----------------------------------------------------------------------------
//								createOutputDirs
//----------------------------------------------------------------------------

// createOutputDir creates the output directory on disk given a
// subdirectory (dir).
func createOutputDirs(dn string, Tables []dbJson.DbTable) error {
	var err 	error
	var outDir	util.Path

	if sharedData.Noop() {
		log.Printf("NOOP -- Skipping Creating directories\n")
		return nil
	}
	outDir = util.NewPath(sharedData.OutDir())

	// We only delete main directory if forced to. Otherwise, we
	// will simply replace our files within it.
	if sharedData.Force() {
		log.Printf("\tRemoving directory: %s...\n", outDir.String())
		if err = outDir.RemoveDir(); err != nil {
			return fmt.Errorf("Error: Could not remove output directory: %s: %s\n",
				outDir.String(), err.Error())
		}
	}

	// Create the main directory if needed.
	if !outDir.IsPathDir() {
		log.Printf("\tCreating directory: %s...\n", outDir.String())
		if err = outDir.CreateDir(); err != nil {
			return fmt.Errorf("Error: Could not crete output directory: %s: %s\n",
				outDir.String(), err.Error())
		}
	}

	log.Printf("\tCreating general directories...\n")
	err = createOutputDir("/html", dn, "")
	if err != nil {
		return err
	}
	err = createOutputDir("/static", dn, "")
	if err != nil {
		return err
	}
	err = createOutputDir("/style", dn, "")
	if err != nil {
		return err
	}
	err = createOutputDir("/tmpl", dn, "")
	if err != nil {
		return err
	}
	err = createOutputDir("/src", dn, "")
	if err != nil {
		return err
	}
	log.Printf("\tCreating src directories for %s...\n", dn)
	err = createOutputDir("/src/hndlr${DbName}", dn, "")
	if err != nil {
		return err
	}
	err = createOutputDir("/src/io${DbName}", dn, "")
	if err != nil {
		return err
	}
	for _, t := range Tables {
		tn := t.TitledName()
		log.Printf("\tCreating directories for Table: %s...\n", tn)
		err = createOutputDir("/src/${DbName}${TblName}", dn, tn)
		if err != nil {
			log.Printf("FAILED on creating /src/%s/%s!\n", dn, tn)
			return err
		}
		err = createOutputDir("/src/hndlr${DbName}${TblName}", dn, tn)
		if err != nil {
			log.Printf("FAILED on creating /src/hndlr%s%s!\n", dn, tn)
			return err
		}
		err = createOutputDir("/src/io${DbName}${TblName}", dn, tn)
		if err != nil {
			log.Printf("FAILED on creating /src/io%s%s!\n", dn, tn)
			return err
		}
	}

	return err
}

//----------------------------------------------------------------------------
//								createOutputPath
//----------------------------------------------------------------------------

// createOutputPath creates an output path from a directory (dir),
// file name (fn), optional database name (dn) and optional table
// name (tn). The dn and tn are only used if "$(DbName}" or "${TblName}"
// are found in the file name.
func createOutputPath(dir, dn, tn, fn string) (util.Path, error) {
	var outPath string
	var err error

	outPath = sharedData.OutDir()
	outPath += "/"
	if len(dir) > 0 {
		outPath += dir
		outPath += "/"
	}
	outPath += fn
	mapper := func(placeholderName string) string {
		switch placeholderName {
		case "DbName":
			if len(dn) > 0 {
				return strings.Title(dn)
			}
		case "TblName":
			if len(tn) > 0 {
				return strings.Title(tn)
			}
		}
		return ""
	}
	outPath = os.Expand(outPath, mapper)
	outPath, err = util.IsPathRegularFile(outPath)
	if err == nil {
		if !sharedData.Replace() {
			err = fmt.Errorf("Over-write error of %s!\n", outPath)
		}
	} else {
		err = nil
	}

	return util.NewPath(outPath), err
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
	var pathIn	util.Path
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
	dn := tmplData.Data.TitledName()

	// Set up the output directory structure
	if err = createOutputDirs(dn, tmplData.Data.Tables); err != nil {
		return err
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

	// Run the first phase of file generation.
	for i, def := range FileDefns {

		if !sharedData.Quiet() {
			log.Println("Process file:", def.ModelName, "generating:", def.FileName, "...")
		}

		// Create the input model file path.
		if pathIn, err = createModelPath(def.ModelName); err != nil {
			return fmt.Errorf("Error: %s: %s\n", pathIn.String(), err.Error())
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
			data.PathOut, err = createOutputPath(def.FileDir, tmplData.Data.Name, "", def.FileName)
			if err != nil {
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
					data.PathOut, err = createOutputPath(def.FileDir, tmplData.Data.Name, v.Name, def.FileName)
					if err != nil {
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

	// Run the second phase of file generation
	for i, def := range FileDefns2 {

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
			data := TaskData{FD:&FileDefns2[i], TD:&tmplData, PathIn:pathIn}
			// Create the output path
			data.PathOut, err = createOutputPath(def.FileDir, tmplData.Data.Name, "", def.FileName)
			if err != nil {
				log.Fatalln(err)
			}
			if sharedData.Debug() {
				log.Println("\t\t outPath=", data.PathOut)
			}
			// Generate the file.
			data.genFile()
		case 2:
			// Output File is Titled Table Name in Titled Database Name directory
			dbJson.DbStruct().ForTables(
				func(v *dbJson.DbTable) {
					data := TaskData{FD:&FileDefns2[i], TD:&tmplData, Table:v, PathIn:pathIn}
					data.PathOut, err = createOutputPath(def.FileDir, tmplData.Data.Name, v.Name, def.FileName)
					if err != nil {
						log.Fatalln(err)
					}
					if sharedData.Debug() {
						log.Println("\t\t outPath=", data.PathOut)
					}
					// Generate the file.
					data.genFile()
				})
		default:
			log.Printf("Skipped %s because of type!\n", def.FileName)
		}
	}

	if dbJson.DbStruct().SqlType == "sqlite" && dbJson.DbStruct().HasDec() {
		log.Printf("========================== WARNING ==========================\n")
		log.Printf("SQLite does not directly support DEC, DECIMAL or MONEY! It\n")
		log.Printf("internally converts decimal types to float. So, rounding\n")
		log.Printf("errors can creep in! So, we define decimal types as string\n")
		log.Printf("in Golang and recommend that you not do any calculations\n")
		log.Printf("in SQL, but rather pull the rows out, convert the data\n")
		log.Printf("to any of the IEEE 754R libraries and do your calculations\n")
		log.Printf("there.\n")
		log.Printf("=============================================================\n")
	}

	return nil
}
