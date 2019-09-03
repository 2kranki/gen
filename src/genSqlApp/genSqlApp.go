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
	"genapp/genCmn"
	"genapp/shared"
	_ "genapp/genSqlApp/dbForm"
	_ "genapp/genSqlApp/dbGener"
	"genapp/genSqlApp/dbJson"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	// Include the various Database Plugins so that they will register
	// with dbPlugin.
	_ "genapp/genSqlApp/dbMariadb"
	_ "genapp/genSqlApp/dbMssql"
	_ "genapp/genSqlApp/dbMysql"
	_ "genapp/genSqlApp/dbPostgres"
	_ "genapp/genSqlApp/dbSqlite"

	"github.com/2kranki/go_util"
)


// FileDefns controls what files are generated.
var FileDefs1 []genCmn.FileDefn = []genCmn.FileDefn{
	{"bld.sh.txt",
		[]string{},
		"b.sh",
		"text",
		0755,
		"one",
		0,
	},
	{"tst2.sh.txt",
		[]string{},
		"t.sh",
		"text",
		0755,
		"one",
		0,
	},
	{"static.txt",
		[]string{"static"},
		"README.txt",
		"copy",
		0644,
		"one",
		0,
	},
	{"tst.sh.txt",
		[]string{"src"},
		"t.sh",
		"text",
		0755,
		"one",
		0,
	},
	{"form.html",
		[]string{"html"},
		"form.html",
		"copy",
		0644,
		"one",
		0,
	},
	{"form.html.tmpl.txt",
		[]string{"tmpl"},
		"${DbName}.${TblName}.form.gohtml",
		"text",
		0644,
		"one",
		2,
	},
	{"list.html.tmpl.txt",
		[]string{"tmpl"},
		"${DbName}.${TblName}.list.gohtml",
		"text",
		0644,
		"one",
		2,
	},
	{"main.menu.html.tmpl.txt",
		[]string{"html"},
		"${DbName}.menu.html",
		"text",
		0644,
		"one",
		0,
	},
	{"go.mod.tmpl.txt",
		[]string{"src"},
		"go.mod",
		"text",
		0644,
		"one",
		0,
	},
	{"DockerBuild.sh.tmpl.txt",
		[]string{""},
		"DockerBuild.sh",
		"text",
		0755,
		"one",
		0,
	},
	{"Dockerfile.tmpl.txt",
		[]string{""},
		"Dockerfile",
		"text",
		0644,
		"one",
		0,
	},
	{"main.go.tmpl.txt",
		[]string{"src"},
		"main.go",
		"text",
		0644,
		"one",
		0,
	},
	{"mainExec.go.tmpl.txt",
		[]string{"src"},
		"mainExec.go",
		"text",
		0644,
		"single",
		0,
	},
	{"handlers.go.tmpl.txt",
		[]string{"src"},
		"hndlr${DbName}.go",
		"text",
		0644,
		"single",
		0,
	},
	{"table.go.tmpl.txt",
		[]string{"src"},
		"${DbName}${TblName}.go",
		"text",
		0644,
		"single",
		2,
	},
	{"table.test.go.tmpl.txt",
		[]string{"src"},
		"${DbName}${TblName}_test.go",
		"text",
		0644,
		"single",
		2,
	},
	{"handlers.table.go.tmpl.txt",
		[]string{"src"},
		"hndlr${DbName}${TblName}.go",
		"text",
		0644,
		"single",
		2,
	},
	{"handlers.table.test.go.tmpl.txt",
		[]string{"src"},
		"hndlr${DbName}${TblName}_test.go",
		"text",
		0644,
		"single",
		2,
	},
	{"io.go.tmpl.txt",
		[]string{"src"},
		"io${DbName}.go",
		"text",
		0644,
		"single",
		0,
	},
	{"docker_test.go.tmpl.txt",
		[]string{"src"},
		"docker_test.go",
		"text",
		0644,
		"single",
		2,
	},
	{"io_test.go.tmpl.txt",
		[]string{"src"},
		"io${DbName}_test.go",
		"text",
		0644,
		"single",
		0,
	},
	{"io.table.go.tmpl.txt",
		[]string{"src"},
		"io${DbName}${TblName}.go",
		"text",
		0644,
		"single",
		2,
	},
	{"io.table.test.go.tmpl.txt",
		[]string{"src"},
		"io${DbName}${TblName}_test.go",
		"text",
		0644,
		"single",
		2,
	},
}

var FileDefs2 []genCmn.FileDefn = []genCmn.FileDefn{
	{"dbs",
		[]string{""},
		"",
		"copyDir",
		0644,
		"single",
		0,
	},
}

//----------------------------------------------------------------------------
//								createOutputDir
//----------------------------------------------------------------------------

// CreateOutputDir creates the output directory on disk given a
// subdirectory (dir).
func CreateOutputDir(dir []string, dn string, tn string) error {
	var err 	error
	var outPath *util.Path

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

	outPath = util.NewPath(sharedData.OutDir())
	for _, d := range dir {
		if len(dir) > 0 {
			outPath = outPath.Append(d)
		}
	}
	outPath = outPath.Expand(mapper)

	if !outPath.IsPathDir() {
		log.Printf("\t\tCreating directory: %s...\n", outPath.String())
		err = outPath.CreateDir()
	}

	return err
}

//----------------------------------------------------------------------------
//								createOutputDirs
//----------------------------------------------------------------------------

// CreateOutputDir creates the output directory on disk given a
// subdirectory (dir).
func CreateOutputDirs(g *genCmn.GenData) error {
	var err 	error
	var outDir	*util.Path

	dn := dbJson.DbStruct().Name
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
	err = CreateOutputDir([]string{"html"}, dn, "")
	if err != nil {
		return err
	}
	err = CreateOutputDir([]string{"static"}, dn, "")
	if err != nil {
		return err
	}
	err = CreateOutputDir([]string{"style"}, dn, "")
	if err != nil {
		return err
	}
	err = CreateOutputDir([]string{"tmpl"}, dn, "")
	if err != nil {
		return err
	}
	err = CreateOutputDir([]string{"src"}, dn, "")
	if err != nil {
		return err
	}

	return err
}

//----------------------------------------------------------------------------
//								createOutputPath
//----------------------------------------------------------------------------

// CreateOutputPath creates an output path from a directory (dir),
// file name (fn), optional database name (dn) and optional table
// name (tn). The dn and tn are only used if "$(DbName}" or "${TblName}"
// are found in the file name.
func CreateOutputPath(dir []string, dn, tn, fn string) (*util.Path, error) {
	var err error
	var outPath string
	var path *util.Path

	outPath = sharedData.OutDir()
	outPath += string(os.PathSeparator)
	for _, d := range dir {
		outPath += d
		outPath += string(os.PathSeparator)
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
	path = util.NewPath(outPath)
	if path.IsPathRegularFile() {
		if !sharedData.Replace() {
			err = fmt.Errorf("Over-write error of %s!\n", outPath)
		}
	}

	return path, err
}

//----------------------------------------------------------------------------
//							SetupFile
//----------------------------------------------------------------------------

// SetupFile sets up the task data defining what is to be done and
// pushes it on the work queue.
func SetupFile(g *genCmn.GenData, fd genCmn.FileDefn, wrk *util.WorkQueue) error {
	var err		error
	var pathIn	*util.Path

	// Create the input model file path.
	pathIn, err = g.CreateModelPath(fd.ModelName)
	if err != nil {
		return fmt.Errorf("Error: %s: %s\n", pathIn.String(), err.Error())
	}
	if sharedData.Debug() {
		log.Println("\t\tmodelPath=", pathIn.String())
	}

	// Now set up to generate the file.
	switch fd.PerGrp {
	case 0:
		// Standard File
		data := &genCmn.TaskData{FD:&fd, TD:&g.TmplData, PathIn:pathIn, Data:dbJson.DbStruct()}
		// Create the output path
		data.PathOut, err = CreateOutputPath(fd.FileDir, dbJson.DbStruct().Name, "", fd.FileName)
		if err != nil {
			log.Fatalln(err)
		}
		if sharedData.Debug() {
			log.Println("\t\t outPath=", data.PathOut)
		}
		// Generate the file.
		wrk.PushWork(data)
	case 2:
		// Output File is Titled Table Name in Titled Database Name directory
		dbJson.DbStruct().ForTables(
			func(v *dbJson.DbTable) {
				data := &genCmn.TaskData{FD:&fd, TD:&g.TmplData, Table:v, PathIn:pathIn.Copy()}
				data.PathOut, err = CreateOutputPath(fd.FileDir, dbJson.DbStruct().Name, v.Name, fd.FileName)
				if err != nil {
					log.Fatalln(err)
				}
				if sharedData.Debug() {
					log.Println("\t\t outPath=", data.PathOut)
				}
				// Generate the file.
				wrk.PushWork(data)
			})
	default:
		log.Printf("Skipped %s because of type!\n", fd.FileName)
	}

	return nil
}

//----------------------------------------------------------------------------
//								readJsonFiles
//----------------------------------------------------------------------------

// ReadJsonFiles reads in the two JSON files that define the
// application to be generated.
func ReadJsonFileData(g *genCmn.GenData) error {
	var err error

	if err = dbJson.ReadJsonFile(sharedData.DataPath()); err != nil {
		return errors.New(fmt.Sprintln("Error: Reading Data Json Input:", sharedData.DataPath(), err))
	}
	g.TmplData.Data = dbJson.DbStruct()

    return nil
}

//============================================================================
//								GenSqlApp
//============================================================================

func GenSqlApp(inDefns map[string]interface{}) error {
	var genData		genCmn.GenData

	// Set up genData.
	genData.Name = "sqlapp"
	genData.FileDefs1 = &FileDefs1
	genData.FileDefs2 = &FileDefs2
	genData.CreateOutputDirs = CreateOutputDirs
	genData.ReadJsonData = ReadJsonFileData
	genData.SetupFile = SetupFile
	genData.TmplData.Data = dbJson.DbStruct()

	if sharedData.Debug() {
		log.Println("\t sql_app: In Debug Mode")
		log.Printf("\t  args: %q\n", flag.Args())
		log.Printf("\tmdldir: %s\n", sharedData.MdlDir())
	}

	genData.GenOutput()

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
