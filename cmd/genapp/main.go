// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// The purpose of this program is to generate other programs and code
// using Golang's templating system.

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"genapp/pkg/genCObj"
	"genapp/pkg/genSqlAppGo"
	"genapp/pkg/sharedData"

	"github.com/2kranki/go_util"
)

var (
	debug         	bool
	execPath      	string
	force         	bool
	genDebugging  	bool
	genHttps    	bool			// Generate HTTPS support.
	genLogging    	bool
	genMuxWrapper 	bool
	jsonPath      	string
	mainPath      	string
	mdldir        	string
	noop          	bool
	outdir        	string
	quiet         	bool
	replace       	bool
)

var defns map[string]interface{}

type defineFlags []string

func (t *defineFlags) String() string {
	return ""
}

func (t *defineFlags) Set(value string) error {
	*t = append(*t, value)
	return nil
}

var defnFlags defineFlags

func ChkSetMdlDir(s string) {
	var path *util.Path

	if len(s) == 0 {
		s = "./models"
	}
	path = util.NewPath(s)
	if !path.IsPathDir() {
		path = util.NewPath(os.ExpandEnv("${GENAPP_MODELS}"))
		if !path.IsPathDir() {
			panic("Error - Missing Model Directory! Set with Env(GENAPP_MODELS) or -mdldir flag.")
		}
	}
	mdldir = path.String()
	sharedData.SetMdlDir(path.String())
}

// SetupShared combines several sources of program options into
// one shared package used throughout the program.
func SetupShared(execPath string, cmd string) error {
	var jsonOut interface{}
	var wrk interface{}
	var err error
	//var flag		bool
	var ok bool

	// Copy the default CLI flags to the shared data.
	defns = map[string]interface{}{}
	sharedData.SetCmd(cmd)
	sharedData.SetDebug(debug)
	sharedData.SetForce(force)
	sharedData.SetDefn("GenDebugging", genDebugging)
	sharedData.SetDefn("GenLogging", genLogging)
	sharedData.SetDefn("GenHttps", genHttps)
	sharedData.SetDefn("GenMuxWrapper", genMuxWrapper)
	ChkSetMdlDir(mdldir)
	sharedData.SetNoop(noop)
	if len(outdir) > 0 {
		sharedData.SetOutDir(outdir)
	}
	sharedData.SetQuiet(quiet)
	sharedData.SetReplace(replace)

	// Now merge in CLI defines. Defines are the <name>=<value>
	// pairs found on the command line.  Note it is fine to
	// specify CLI flags in the defines as long as the type
	// of the value is correct.
	for _, v := range defnFlags {
		s := strings.Split(v, "=")
		if len(s) > 1 {
			sharedData.SetDefn(s[0], s[1])
		}
	}

	// Now merge in the Exec JSON File if there is one.  It
	// uses the same names as the CLI flags.
	if len(execPath) > 0 {
		jsonOut, err = util.ReadJsonFile(execPath)
		if err != nil {
			return errors.New(fmt.Sprintln("Error: Exec JSON,", execPath, ", "+
				"file did not unmarshal properly:", err))
		}
		if debug {
			fmt.Println("\tData:", jsonOut)
		}
		m := jsonOut.(map[string]interface{})
		if m == nil {
			return errors.New("Error: Exec JSON file did not unmarshal properly!")
		}
		if wrk, ok = m["data"]; ok {
			sharedData.SetDataPath(wrk.(string))
		}
		if wrk, ok = m["debug"]; ok {
			sharedData.SetDebug(wrk.(bool))
		}
		if wrk, ok = m["force"]; ok {
			sharedData.SetForce(wrk.(bool))
		}
		if wrk, ok = m["main"]; ok {
			sharedData.SetMainPath(wrk.(string))
		}
		if wrk, ok = m["mdldir"]; ok {
			wrkDir := wrk.(string)
			ChkSetMdlDir(wrkDir)
		}
		if wrk, ok = m["noop"]; ok {
			sharedData.SetNoop(wrk.(bool))
		}
		if wrk, ok = m["outdir"]; ok {
			sharedData.SetOutDir(wrk.(string))
		}
		if wrk, ok = m["quiet"]; ok {
			sharedData.SetQuiet(wrk.(bool))
		}
		if wrk, ok = m["replace"]; ok {
			sharedData.SetReplace(wrk.(bool))
		}
		if wrk, ok = m["define"]; ok {
			s := strings.Split(wrk.(string), ",")
			for _, v := range s {
				ss := strings.Split(v, "=")
				if len(ss) > 1 {
					sharedData.SetDefn(ss[0], ss[1])
				}
			}
		}
		if wrk, ok = m["cmd"]; ok {
			sharedData.SetCmd(wrk.(string))
		}
	}

	return nil
}

func main() {
	var err error

	flag.Usage = usage
	flag.BoolVar(&debug, "debug", true, "enable debugging")
	flag.StringVar(&execPath, "exec", "", "exec json path (optional)")
	flag.StringVar(&execPath, "x", "", "exec json path (optional)")
	flag.BoolVar(&force, "force", false, "enable over-writes and deletions")
	flag.BoolVar(&force, "f", false, "enable over-writes and deletions")
	flag.BoolVar(&genDebugging, "genDebugging", true, "generate debugging output")
	flag.BoolVar(&genHttps, "genHttps", true, "generate HTTPS support")
	flag.BoolVar(&genLogging, "genLogging", true, "generate logging")
	flag.BoolVar(&genMuxWrapper, "genMuxWrapper", true, "generate a wrapper around the mux")
	flag.StringVar(&mainPath, "main", "", "set json main input path")
	flag.StringVar(&jsonPath, "json", "", "set json main input path")
	flag.StringVar(&mdldir, "mdldir", "./models", "set model input directory")
	flag.BoolVar(&noop, "noop", false, "execute program, but do not make real changes")
	flag.StringVar(&outdir, "outdir", "/tmp", "set output directory")
	flag.BoolVar(&quiet, "quiet", false, "enable quiet mode")
	flag.BoolVar(&replace, "replace", true, "overwrite existing files")
	flag.BoolVar(&quiet, "q", false, "enable quiet mode")
	flag.Var(&defnFlags, "define", "enter definitions (<name>=<string>)")
	flag.Var(&defnFlags, "d", "enter definitions (<name>=<string>)")
	flag.Parse()
	if debug {
		log.Println("\tIn Debug Mode...")
	}

	err = SetupShared(execPath, flag.Arg(0))
	if err != nil {
		log.Fatalln("Error: failed to set up main definitions:", err)
	}

	// Execute the command
	if debug {
		log.Println("\tcmd: '", sharedData.Cmd(), "'")
	}
	switch sharedData.Cmd() {
	case "cobj":
		err = genCObj.Generate(defns)
	case "sqlappgo":
		err = genSqlAppGo.Generate(defns)
	default:
		fmt.Println("\nError: command must be 'cobj' or 'sqlapp'")
	}
	if err != nil {
		log.Println("Error: generation failed:", err)
	}

	if !quiet {
		log.Println("\tEnd of Program")
	}
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "Usage:\n\tgen [options] (cobj | sqlapp)\n")
	fmt.Fprintf(flag.CommandLine.Output(), "\nOptions:\n")
	flag.PrintDefaults()
	fmt.Fprintf(flag.CommandLine.Output(), "\nNotes:\n")
	fmt.Fprintf(flag.CommandLine.Output(), "'exec json' is a file that defines the command line parameters \n")
	fmt.Fprintf(flag.CommandLine.Output(), "so that you can set them and then execute gen with -x or -exec\n")
	fmt.Fprintf(flag.CommandLine.Output(), "option.\n\n")
	fmt.Fprintf(flag.CommandLine.Output(), "'json path' is the json file that defines the data passed to the\n")
	fmt.Fprintf(flag.CommandLine.Output(), "template engine which controls data within the generated files.\n")
	fmt.Fprintf(flag.CommandLine.Output(), "'{{' and '}}' are not used in the basic templates.  Instead, '[['\n")
	fmt.Fprintf(flag.CommandLine.Output(), "and ']]' are used. This way, we can pass the generated text back\n")
	fmt.Fprintf(flag.CommandLine.Output(), "through the templating system at execution time.\n")
	fmt.Fprintf(flag.CommandLine.Output(), "Version: 0.1\n\n")
}
