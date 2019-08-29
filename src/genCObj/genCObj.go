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
	"flag"
	"fmt"
	"log"

	"github.com/2kranki/go_util"
)

// FileDefns controls what files are generated.
var FileDefs1 	[]genCmn.FileDefn = []genCmn.FileDefn{
	{"obj_int_h.txt",
		[]string{"src"},
		"${Name}_internal.h",
		"text",
		0644,
		"",
		0,
	},
	{"obj_obj_c.txt",
		[]string{"src"},
		"${Name}_object.c",
		"text",
		0644,
		"",
		0,
	},
	{"obj_c.txt",
		[]string{"src"},
		"${Name}.c",
		"text",
		0644,
		"",
		0,
	},
	{"obj_h.txt",
		[]string{"src"},
		"${Name}.h",
		"text",
		0644,
		"",
		0,
	},
	{"obj_test_c.txt",
		[]string{"tests"},
		"test_${Name}.c",
		"text",
		0644,
		"",
		0,
	},
}

func init() {

}

//----------------------------------------------------------------------------
//								createOutputDir
//----------------------------------------------------------------------------

// CreateOutputDir creates the output directory on disk given a
// subdirectory (dir).
func CreateOutputDir(g *genCmn.GenData, dir []string) error {
	var err 	error
	var outPath *util.Path

	mapper := func(placeholderName string) string {
		switch placeholderName {
		case "Name":
			if len(dbStruct.Name) > 0 {
				return dbStruct.Name
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

// createOutputDir creates the output directory on disk given a
// subdirectory (dir).
func CreateOutputDirs(g *genCmn.GenData) error {
	var err 	error
	var outDir	*util.Path

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
	err = CreateOutputDir(g, []string{"src"})
	if err != nil {
		return err
	}
	err = CreateOutputDir(g, []string{"tests"})
	if err != nil {
		return err
	}

	return err
}

//----------------------------------------------------------------------------
//								CreateOutputFilePath
//----------------------------------------------------------------------------

func CreateOutputFilePath(name string, dir []string, fn string) (*util.Path, error) {
	var outPath 	*util.Path

	mapper := func(varSub string) string {
		switch varSub {
		case "Name":
			return name
		}
		return ""
	}

	outPath = util.NewPath(sharedData.OutDir())
	for _, d := range dir {
		outPath = outPath.Append(d)
	}
	outPath = outPath.Append(fn)
	outPath = outPath.Expand(mapper)

	if outPath.IsPathRegularFile() {
		if !sharedData.Force() {
			return outPath, fmt.Errorf("Over-write error of: %s\n", outPath)
		}
	}

	return outPath, nil
}

//----------------------------------------------------------------------------
//							readJsonFileData
//----------------------------------------------------------------------------

// ReadJsonFileData reads in the Data JSON file(s) that define the
// application to be generated.
func ReadJsonFileData(g *genCmn.GenData) error {
	var err error

	if err = ReadJsonFile(sharedData.DataPath()); err != nil {
		return fmt.Errorf("Error: Reading Data Json Input:%s %s\n",
			sharedData.DataPath(), err.Error())
	}
	g.TmplData.Data = DbStruct()


	return nil
}

//----------------------------------------------------------------------------
//							SetupFile
//----------------------------------------------------------------------------

// SetupFile sets up the task data defining what is to be done and
// pushes it on the work queue.
func SetupFile(g *genCmn.GenData, fd genCmn.FileDefn, wrk *util.WorkQueue) error {
	var err		error

	data := &genCmn.TaskData{}
	data.FD = &fd
	data.TD = DbStruct()
	data.Data = DbStruct()

	// Create the input model file path.
	data.PathIn, err = g.CreateModelPath(fd.ModelName)
	if err != nil {
		return fmt.Errorf("Error: %s: %s\n", data.PathIn.String(), err.Error())
	}
	if sharedData.Debug() {
		log.Println("\t\tmodelPath=", data.PathIn.String())
	}

	// Create the output path
	data.PathOut, err = CreateOutputFilePath(dbStruct.Name, fd.FileDir, fd.FileName)
	if err != nil {
		log.Fatalln(err)
	}
	if sharedData.Debug() {
		log.Println("\t\t outPath=", data.PathOut)
	}

	// Generate the file.
	wrk.PushWork(data)

	return nil
}

//============================================================================
//								GenCObj
//============================================================================

func GenCObj(inDefns map[string]interface{}) error {
	var genData		genCmn.GenData

	genData.Name = "cobj"
	genData.Mapper = func(varSub string) string {
		switch varSub {
		case "Name":
			return dbStruct.Name
		}
		return ""
	}
	genData.FileDefs1 = &FileDefs1
	genData.CreateOutputDirs = CreateOutputDirs
	genData.ReadJsonData = ReadJsonFileData
	genData.SetupFile = SetupFile
	genData.TmplData.Data = DbStruct()

	if sharedData.Debug() {
		log.Println("GenCObj: In Debug Mode...")
		log.Printf("\t  args: %q\n", flag.Args())
	}

	genData.GenOutput()

	return nil
}
