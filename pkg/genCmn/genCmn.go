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

package genCmn

import (
	"errors"
	"flag"
	"fmt"
	"genapp/pkg/mainData"
	"genapp/pkg/sharedData"
	"log"
	"os"

	"github.com/2kranki/go_util"
)

//============================================================================
//							File Definition
//============================================================================

// FileDefn gives the parameters needed to generate a file.  The fields of
// the struct have been simplified to allow for easy json encoding/decoding.
type FileDefn struct {
	ModelName string      `json:"ModelName,omitempty"`
	FileDir   []string    `json:"FileDir,omitempty"`   // Output File Directory
	FileName  string      `json:"FileName,omitempty"`  // Output File Name
	FileType  string      `json:"Type,omitempty"`      // copy, text, sql, html
	FilePerms os.FileMode `json:"FilePerms,omitempty"` // Output File Permissions
	Class     string      `json:"Class,omitempty"`     // single, table
	PerGrp    int         `json:"PerGrp,omitempty"`    // 0 == generate one file
	//															// 1 == generate one file for a database
	// 															// 2 == generate one file for a table
}

// See tasks.go for TaskData

//============================================================================
//							Template Data
//============================================================================

// TmplData is used to centralize all the inputs to the generators.
// We maintain generic JSON structures for the templating system
// which does not support structs.
type TmplData struct {
	Main   interface{}
	Data   interface{}
	Extra1 interface{}
	Extra2 interface{}
}

//============================================================================
//							Generate Data
//============================================================================

// GenData is used to centralize all the inputs to the generators.  We
// maintain generic JSON structures for the templating system which does
// not support structs.  (Not certain why yet.) We also maintain the data
// in structs for easier access by the generation functions.
type GenData struct {
	// Name indicates the type of generation being performed. Valid values
	// are:
	//				"cobj"
	//				"sqlappgo"
	// Note: This must be kept consistent with the directory structure
	// 			of "models".
	Name string
	// Mapper is an optional function which translates a string to another
	// string.  It can be used for translating external names to structs
	// or fields within this system.
	Mapper func(string) string
	// TmplData contains the primary data which is passed to the templating
	// system, to various plugins and generally used for generation.
	TmplData TmplData
	// fileDefs1 define the files that will make up the application
	// fill the various sub-directories of that application
	FileDefs1 *[]FileDefn
	// fileDefs2 define the sub-directories that will be copied into
	// the application's sub-directories.
	FileDefs2        *[]FileDefn
	CreateOutputDirs func(g *GenData) error
	// GenFile should actually generate the file generally using the
	// tasks portion of genCmn.
	GenFile func(g *GenData, d *TaskData)
	// ReadJsonFileData reads in the Data JSON file(s) that define the
	// application to be generated.
	ReadJsonData func(g *GenData) error
	// SetupFile sets up the task data (TaskData) defining what is to
	// be done and pushes it on the work queue.
	SetupFile func(g *GenData, fd FileDefn, work *util.WorkQueue) error
	// SetupTmplData
	SetupTmplData func(g *GenData) error
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
	path = path.Append(g.Name)
	path = path.Append(fn)
	if !path.IsPathRegularFile() && !path.IsPathDir() {
		return nil, fmt.Errorf("Error: %s is not a directory or a file!\n", path.String())
	}

	return path, nil
}

//----------------------------------------------------------------------------
//								CreateOutputPath
//----------------------------------------------------------------------------

// CreateOutputPath creates an output path from our models and verifies that it
// exists.  Regular files and directories that are found in the "models"
// directory are acceptable. genDir is the model subdirectory for this type
// of generation.
func (g *GenData) CreateOutputPath(mapper func(string) string, dir []string, fn string) (*util.Path, error) {
	var err error
	var outPath string
	var path *util.Path

	// Create output path with possible embedded substitutions.
	outPath = sharedData.OutDir()
	outPath += string(os.PathSeparator)
	for _, d := range dir {
		outPath += d
		outPath += string(os.PathSeparator)
	}
	outPath += fn

	// Perform substitutions in path.
	outPath = os.Expand(outPath, mapper)

	// Create a util.path for our finalized path.
	path = util.NewPath(outPath)
	if path == nil {
		return nil, fmt.Errorf("Error - out of memory\n")
	}

	// Create the directory portion of the path if it doesn't exist.
	oDirPath := util.NewPath(path.Dir())
	if !oDirPath.IsPathDir() {
		if sharedData.Noop() {
			log.Printf("genCmn::CreateOutputPath would have created %s\n", oDirPath.String())
		} else {
			oDirPath.CreateDir()
		}
	}

	// Check if we would over-write the file.
	if path.IsPathRegularFile() {
		if !sharedData.Replace() {
			err = fmt.Errorf("Over-write error of %s!\n", outPath)
		}
	}

	return path, err
}

//----------------------------------------------------------------------------
//							Read JSON Files
//----------------------------------------------------------------------------

// readJsonFileMain reads in the Main JSON file that define the
// application arguments.
func (g *GenData) ReadJsonFileMain() error {
	var err error

	if err = mainData.ReadJsonFileMain(sharedData.MainPath()); err != nil {
		return errors.New(fmt.Sprintln("Error: Reading Main Json Input:", sharedData.MainPath(), err))
	}
	g.TmplData.Main = mainData.MainJson()

	return nil
}

//----------------------------------------------------------------------------
//							Generate the Output
//----------------------------------------------------------------------------

func (g *GenData) GenFiles(fd *[]FileDefn) error {
	var err error
	var pathIn *util.Path
	var work *util.WorkQueue

	// Setup the worker queue.
	work = util.NewWorkQueue(
		func(a interface{}, cmn interface{}) {
			var data *TaskData
			var ok bool

			data, ok = a.(*TaskData)
			if ok {
				data.genFile()
			} else {
				panic(fmt.Sprintf("FATAL: Invalid TaskData Type of %T!", a))
			}
		},
		nil,
		0)

	// Run the first phase of file generation.
	for _, def := range *fd {

		if !sharedData.Quiet() {
			log.Println("Setting up file:", def.ModelName, "generating:", def.FileName, "...")
		}

		// Create the input model file path.
		if pathIn, err = g.CreateModelPath(def.ModelName); err != nil {
			return fmt.Errorf("Error: %s: %s\n", pathIn.String(), err.Error())
		}
		if sharedData.Debug() {
			log.Println("\t\tmodelPath=", pathIn)
		}

		// Now setup to generate the file pushing the setup info
		// onto the work queue.
		if g.SetupFile != nil {
			g.SetupFile(g, def, work)
		} else {
			panic("FATAL: Missing GenData::SetupFile!\n")
		}
	}
	work.CloseAndWaitForCompletion()

	return err
}

func (g *GenData) GenOutput() error {
	var err error

	if sharedData.Debug() {
		log.Println("\t genOutput: In Debug Mode")
		log.Printf("\t    args: %q\n", flag.Args())
		log.Printf("\t  mdldir: %s\n", sharedData.MdlDir())
	}

	// Read the JSON files.
	if err = g.ReadJsonFileMain(); err != nil {
		log.Fatalln(err)
	}
	if g.ReadJsonData != nil {
		if err = g.ReadJsonData(g); err != nil {
			log.Fatalln(err)
		}
	}

	// Set up template data
	if g.SetupTmplData != nil {
		err = g.SetupTmplData(g)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if sharedData.OutDir() == "" {
		return fmt.Errorf("Error - 'libPath' cli argument is required!\n\n\n")
	}

	// Set up the output directory structure
	if g.CreateOutputDirs != nil {
		if err = g.CreateOutputDirs(g); err != nil {
			return err
		}
	}

	g.GenFiles(g.FileDefs1)
	if g.FileDefs2 != nil {
		g.GenFiles(g.FileDefs2)
	}

	return nil
}

//----------------------------------------------------------------------------
//						Create a New GenData Structure
//----------------------------------------------------------------------------

func NewGenData() *GenData {
	gd := &GenData{}
	return gd
}
