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
genSqlAppGo
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

package genSqlAppGo

import (
	"flag"
	"fmt"
	"genapp/pkg/genCmn"
	_ "genapp/pkg/genSqlAppGo/dbForm"
	_ "genapp/pkg/genSqlAppGo/dbGener"
	"genapp/pkg/genSqlAppGo/dbJson"
	"genapp/pkg/genSqlAppGo/dbPlugin"
	sharedData "genapp/pkg/sharedData"
	"log"
	"os"
	"strings"
	// Include the various Database Plugins so that they will register
	// with dbPlugin.
	_ "genapp/pkg/genSqlAppGo/dbMariadb"
	_ "genapp/pkg/genSqlAppGo/dbMssql"
	_ "genapp/pkg/genSqlAppGo/dbMysql"
	_ "genapp/pkg/genSqlAppGo/dbPostgres"
	_ "genapp/pkg/genSqlAppGo/dbSqlite"

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
	{"jenkins_build.py.tmpl.txt",
		[]string{"jenkins", "build"},
		"build.py",
		"text",
		0755,
		"one",
		0,
	},
	{"jenkins_build_test.py.tmpl.txt",
		[]string{"jenkins", "build"},
		"build_test.py",
		"text",
		0755,
		"one",
		0,
	},
	{"jenkins_deploy.py.tmpl.txt",
		[]string{"jenkins", "deploy"},
		"deploy.py",
		"text",
		0755,
		"one",
		0,
	},
	{"jenkins_deploy_test.py.tmpl.txt",
		[]string{"jenkins", "deploy"},
		"deploy_test.py",
		"text",
		0755,
		"one",
		0,
	},
	{"jenkins_lint.py.tmpl.txt",
		[]string{"jenkins", "lint"},
		"lint.py",
		"text",
		0755,
		"one",
		0,
	},
	{"jenkins_lint_test.py.tmpl.txt",
		[]string{"jenkins", "lint"},
		"lint_test.py",
		"text",
		0755,
		"one",
		0,
	},
	{"jenkins_push.py.tmpl.txt",
		[]string{"jenkins", "push"},
		"push.py",
		"text",
		0755,
		"one",
		0,
	},
	{"jenkins_push_test.py.tmpl.txt",
		[]string{"jenkins", "push"},
		"push_test.py",
		"text",
		0755,
		"one",
		0,
	},
	{"jenkins_test.py.tmpl.txt",
		[]string{"jenkins", "test"},
		"test.py",
		"text",
		0755,
		"one",
		0,
	},
	{"jenkins_test_test.py.tmpl.txt",
		[]string{"jenkins", "test"},
		"test_test.py",
		"text",
		0755,
		"one",
		0,
	},
	{"util.py.tmpl.txt",
		[]string{},
		"util.py",
		"text",
		0755,
		"one",
		0,
	},
	{"util_test.py.tmpl.txt",
		[]string{},
		"util_test.py",
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
	{"tmpl.txt",
		[]string{"tmpl"},
		"README.txt",
		"copy",
		0644,
		"one",
		0,
	},
	{"tst.sh.txt",
		[]string{"cmd/${DbName}"},
		"t.sh",
		"text",
		0755,
		"one",
		0,
	},
	{"form.html",
		[]string{"static"},
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
		[]string{"tmpl"},
		"${DbName}.main.menu.gohtml",
		"text",
		0644,
		"one",
		0,
	},
	{"go.mod.tmpl.txt",
		[]string{""},
		"go.mod",
		"text",
		0644,
		"one",
		0,
	},
	{"dockerCli.py.tmpl.txt",
		[]string{""},
		"dockerCli.py",
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
	{"docker-compose.yaml.tmpl.txt",
		[]string{""},
		"docker-compose.yaml",
		"text",
		0644,
		"one",
		0,
	},
	{"main.go.tmpl.txt",
		[]string{"cmd/${DbName}"},
		"main.go",
		"text",
		0644,
		"one",
		0,
	},
	{"mainExec.go.tmpl.txt",
		[]string{"cmd/${DbName}"},
		"mainExec.go",
		"text",
		0644,
		"single",
		0,
	},
	{"handlers.go.tmpl.txt",
		[]string{"cmd/${DbName}"},
		"hndlr${DbName}.go",
		"text",
		0644,
		"single",
		0,
	},
	{"table.go.tmpl.txt",
		[]string{"cmd/${DbName}"},
		"${DbName}${TblName}.go",
		"text",
		0644,
		"single",
		2,
	},
	{"table.test.go.tmpl.txt",
		[]string{"cmd/${DbName}"},
		"${DbName}${TblName}_test.go",
		"text",
		0644,
		"single",
		2,
	},
	{"handlers.table.go.tmpl.txt",
		[]string{"cmd/${DbName}"},
		"hndlr${DbName}${TblName}.go",
		"text",
		0644,
		"single",
		2,
	},
	{"handlers.table.test.go.tmpl.txt",
		[]string{"cmd/${DbName}"},
		"hndlr${DbName}${TblName}_test.go",
		"text",
		0644,
		"single",
		2,
	},
	{"io.go.tmpl.txt",
		[]string{"cmd/${DbName}"},
		"io${DbName}.go",
		"text",
		0644,
		"single",
		0,
	},
	{"docker_test.go.tmpl.txt",
		[]string{"cmd/${DbName}"},
		"docker_test.go",
		"text",
		0644,
		"single",
		2,
	},
	{"io_test.go.tmpl.txt",
		[]string{"cmd/${DbName}"},
		"io${DbName}_test.go",
		"text",
		0644,
		"single",
		0,
	},
	{"io.table.go.tmpl.txt",
		[]string{"cmd/${DbName}"},
		"io${DbName}${TblName}.go",
		"text",
		0644,
		"single",
		2,
	},
	{"io.table.test.go.tmpl.txt",
		[]string{"cmd/${DbName}"},
		"io${DbName}${TblName}_test.go",
		"text",
		0644,
		"single",
		2,
	},
	{"vendor_notes.txt",
		[]string{"vendor"},
		"Notes.txt",
		"text",
		0644,
		"single",
		0,
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
	var err error
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
	var err error
	var outDir *util.Path

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
	// cmd is used for main.go
	err = CreateOutputDir([]string{"cmd/${DbName}"}, dn, "")
	if err != nil {
		return err
	}
	// pkg is used for application packages.
	err = CreateOutputDir([]string{"pkg"}, dn, "")
	if err != nil {
		return err
	}
	// Static is used for CSS, HTML, JPG and any other static data
	err = CreateOutputDir([]string{"static"}, dn, "")
	if err != nil {
		return err
	}
	// tmpl is used for web page templates normally named *.gohtml.
	err = CreateOutputDir([]string{"tmpl"}, dn, "")
	if err != nil {
		return err
	}
	// vendor is used for application dependencies which should
	// probably not be commited in git.
	err = CreateOutputDir([]string{"vendor"}, dn, "")
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

	oDirPath := util.NewPath(path.Dir())
	if !oDirPath.IsPathDir() {
		oDirPath.CreateDir()
	}

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
	var err error
	var pathIn *util.Path

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
	case 0: // Standard File
		data := &genCmn.TaskData{FD: &fd, TD: &g.TmplData, PathIn: pathIn, Data: dbJson.DbStruct()}
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
	case 2: // Output File is Titled Table Name in Titled Database
		// Name directory
		dbJson.DbStruct().ForTables(
			func(v *dbJson.DbTable) {
				data := &genCmn.TaskData{FD: &fd, TD: &g.TmplData, Table: v, PathIn: pathIn.Copy()}
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
// application to be generated and validates them.
func ReadJsonFileData(g *genCmn.GenData) error {
	var err error
	var db *dbJson.Database

	db = dbJson.DbStruct()
	if err = db.ReadJsonFile(sharedData.DataPath()); err != nil {
		return fmt.Errorf("Error: Reading Data Json Input: %s - %s\n",
			sharedData.DataPath(), err)
	}
	g.TmplData.Data = db

	if err = db.SetupPlugin(); err != nil {
		return err
	}
	if err = db.ValidatePlugin(); err != nil {
		return err
	}

	return nil
}

// SetupPlugin finds the plugin needed and sets it up within the database.
func SetupPlugin() error {
	var err error
	var intr dbPlugin.SchemaNamer
	var ok bool
	var plg dbPlugin.PluginData
	var db *dbJson.Database

	// Indicate the plugin needed.
	db = dbJson.DbStruct()
	if sharedData.Debug() {
		log.Printf("\t\tSqtype: %s\n", db.SqlType)
	}

	// Find the plugin for this database.
	if plg, err = dbPlugin.FindPlugin(db.SqlType); err != nil {
		return fmt.Errorf("Error: Can't find plugin for %s!\n\n\n", db.SqlType)
	}
	if sharedData.Debug() {
		log.Printf("\t\tPlugin Type: %T\n", plg)
		log.Printf("\t\tPlugin: %+v\n", plg)
		log.Printf("\t\tPlugin.Plugin: %+v\n", plg.Plugin)
	}

	// Validate the Plugin if possible.
	if plg.Types == nil {
		return fmt.Errorf("Error: Plugin missing types for %s!\n\n\n", db.SqlType)
	}

	// Save the plugin.
	db.Plugin = plg

	if len(db.Schema) == 0 {
		intr, ok = plg.Plugin.(dbPlugin.SchemaNamer)
		if ok {
			db.Schema = intr.SchemaName()
		}
	}

	// Set up the Table Fields so that point to the Plugin Field Type definition.
	for _, t := range db.Tables {
		for ii, _ := range t.Fields {
			t.Fields[ii].Typ = plg.Types.FindDefn(t.Fields[ii].TypeDefn)
			if t.Fields[ii].Typ == nil {
				return fmt.Errorf("Error: Invalid Field Type for %s:%s!\n\n\n", t.Name, t.Fields[ii].Name)
			}
		}
	}
	return nil
}

// ValidatePlugin checks the JSON built structures for errors with
// respect to the plugin. This assumes that the data was previously
// validated.
func validatePlugin() error {
	var err error
	var plg dbPlugin.PluginData

	// Set up Plugin Support for this database type.
	if plg, err = dbPlugin.FindPlugin(dbJson.DbStruct().SqlType); err != nil {
		return err
	}

	for _, t := range dbJson.DbStruct().Tables {
		for _, f := range t.Fields {
			td := plg.Types.FindDefn(f.TypeDefn)
			if td == nil {
				fmt.Errorf("Error - Could not find Type definition for field: %s  type: %s\n",
					f.Name, f.TypeDefn)
			}
		}
	}

	return nil
}

//============================================================================
//								GenSqlApp
//============================================================================

func Generate(inDefns map[string]interface{}) error {
	var genData genCmn.GenData

	// Set up genData.
	genData.Name = "sqlapp"
	genData.FileDefs1 = &FileDefs1
	genData.FileDefs2 = &FileDefs2
	genData.CreateOutputDirs = CreateOutputDirs
	genData.ReadJsonData = ReadJsonFileData
	genData.SetupFile = SetupFile
	genData.TmplData.Data = dbJson.DbStruct()

	if sharedData.Debug() {
		log.Println("\t sqlapp: In Debug Mode")
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
