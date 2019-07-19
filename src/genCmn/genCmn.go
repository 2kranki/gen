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


//============================================================================
//							File Definition
//============================================================================

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

//============================================================================
//							Generate Data
//============================================================================

// TmplData is used to centralize all the inputs
// to the generators.  We maintain generic JSON
// structures for the templating system which does
// not support structs.  (Not certain why yet.)
// We also maintain the data in structs for easier
// access by the generation functions.
type GenData struct {
	genName		string
	TmplData   	interface{}
	// fileDefs1 define the files that will make up the application
	// fill the various sub-directories of that application
	fileDefs1	*[]FileDefn
	// fileDefs2 define the sub-directories that will be copied into
	// the application's sub-directories.
	fileDefs2	*[]FileDefn
}

func (g *GenData) GenName() string {
	return g.genName
}

func (g *GenData) SetGenName(s string) {
	g.genName = s
}

func (g *GenData) FileDefs1() *[]FileDefn {
	return g.fileDefs1
}

func (g *GenData) SetFileDefs1(f *[]FileDefn) {
	g.fileDefs1 = f
}

func (g *GenData) FileDefs2() *[]FileDefn {
	return g.fileDefs2
}

func (g *GenData) SetFileDefs2(f *[]FileDefn) {
	g.fileDefs2 = f
}

//----------------------------------------------------------------------------
//								CreateModelPath
//----------------------------------------------------------------------------

// CreateModelPath creates an input path from our models and verifies that it
// exists.  Regular files and directories that are found in the "models"
// directory are acceptable. genDir is the model subdirectory for this type
// of generation.
func (g *GenData) CreateModelPath(fn string) (*util.Path, error) {

	// Calculate the model path.
	path := util.NewPath(sharedData.MdlDir())
	path = path.Append(g.genName)
	path = path.Append(fn)
	if !path.IsPathRegularFile() && !path.IsPathDir() {
		return nil, fmt.Errorf("Error: %s is not a directory or a file!\n", path.String())
	}

	return path, nil
}

func (g *GenData) CreateOutputDirs() error {
	panic("Function, GenData.CreateOutputDirs, needs to be implemented!")
	return nil
}

//----------------------------------------------------------------------------
//								readJsonFiles
//----------------------------------------------------------------------------

// readJsonFileMain reads in the Main JSON file that define the
// application arguments.
func (g *GenData) readJsonFileMain( ) error {
	var err error

	if err = mainData.ReadJsonFileMain(sharedData.MainPath()); err != nil {
		return errors.New(fmt.Sprintln("Error: Reading Main Json Input:", sharedData.MainPath(), err))
	}

    return nil
}

// readJsonFileData reads in the Data JSON file(s) that define the
// application to be generated.
func (g *GenData) readJsonFileData( ) error {
	panic("Function, GenData.readJsonFileData, needs to be implemented!")
	return nil
}

//----------------------------------------------------------------------------
//						Set up various data elements
//----------------------------------------------------------------------------

func (g *GenData) SetupTmplData() (interface{}, error) {
	panic("Function, GenData.readJsonFileData, needs to be implemented!")
	return nil
}

// SetupGroup0File sets up the task data defining what is to be done and
// returning that so that it may be pushed on the work queue.
func (g *GenData) SetupGroup0File(fd FileDefn) {
	panic("Function, GenData.SetupGroup0File, needs to be implemented!")
}

// SetupGroup2File sets up the task data defining what is to be done and
// returning that so that it may be pushed on the work queue.
func (g *GenData) SetupGroup2File(fd FileDefn, wrk *chan interface{}) {
	panic("Function, GenData.SetupGroup2File, needs to be implemented!")
}

//----------------------------------------------------------------------------
//								Generate the Output
//----------------------------------------------------------------------------

// GenFile should actually generate the file generally using the
// tasks portion of genCmn.
func (g *GenData) GenFile(d interface{}) {
	panic("Function, GenData.GenFile, needs to be implemented!")
}


func (g *GenData) GenFiles(fd *[]FileDefn) error {
	var err error
	var pathIn *util.Path
	//var ok 		bool

	// Setup the worker queue.
	done := make(chan bool)
	inputQueue := util.Workers(
		func(a interface{}) {
			g.GenFile(a)
		},
		func() {
			done <- true
		},
		5)

	// Run the first phase of file generation.
	for i, def := range fd {

		if !sharedData.Quiet() {
			log.Println("Process file:", def.ModelName, "generating:", def.FileName, "...")
		}

		// Create the input model file path.
		if pathIn, err = g.CreateModelPath(def.ModelName); err != nil {
			return fmt.Errorf("Error: %s: %s\n", pathIn.String(), err.Error())
		}
		if sharedData.Debug() {
			log.Println("\t\tmodelPath=", pathIn)
		}

		// Now generate the file.
		switch def.PerGrp {
		case 0:
			// Standard File
			data := g.SetupGroup0File(def)
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

}

func (g *GenData) GenOutput(inDefns map[string]interface{}) error {
	var err 	error
	var pathIn	*util.Path
	//var ok 		bool

	if sharedData.Debug() {
		log.Println("\t genOutput: In Debug Mode")
		log.Printf("\t    args: %q\n", flag.Args())
		log.Printf("\t  mdldir: %s\n", sharedData.MdlDir())
	}

    // Read the JSON files.
    if err = g.ReadJsonFileMain(); err != nil {
		log.Fatalln(err)
    }
	if err = g.ReadJsonFileData(); err != nil {
		log.Fatalln(err)
	}

	// Set up template data
	g.TmplData, err = g.SetupTmplData()
	if err != nil {
		log.Fatalln(err)
	}

	// Set up the output directory structure
	if err = g.CreateOutputDirs( ); err != nil {
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
