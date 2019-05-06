// See License.txt in main repository directory

// Shared contains the shared variables and data used by
// main and the other packages.  This was created to remove
// circular references from main's data and it's being used
// in sub-packages which have the need of reference and
// sometimes manipulating or adding to that data.

package sharedData

import (
	"fmt"
	"time"
)

var cmd			string
var dataPath	string
var defns		map[string]interface{}
var	debug		bool
var force		bool
var funcs		map[string]interface{}
var mainPath	string
var	mdlDir		string
var	noop		bool
var	outDir		string
var quiet		bool
var timeNow		string

func init() {
	defns  = map[string]interface{}{}
	funcs  = map[string]interface{}{}
	mdlDir = "./models"
	outDir = "./test"
	timeNow = time.Now().Format("Mon Jan _2, 2006 15:04")
	//sharedData.SetFunc("Time", Time)	<== Causes import cycle, added to main
}

func Cmd() string {
	return cmd
}

func SetCmd(f string) {
	cmd = f
}

// DataPath is the path to the app json file.
func DataPath() string {
	return dataPath
}

func SetDataPath(f string) {
	dataPath = f
}

func Debug() bool {
	return debug
}

func SetDebug(f bool) {
	debug = f
}

func Defn(nm string) interface{} {
	switch nm {
	case "cmd":
		return cmd
	case "dataPath":
		return dataPath
	case "debug":
		return debug
	case "force":
		return force
	case "mainPath":
		return mainPath
	case "mdlDir":
		return mdlDir
	case "noop":
		return noop
	case "outDir":
		return outDir
	case "quiet":
		return quiet
	case "timeNow":
		return timeNow
	}
	 d, _ := defns[nm]
	return d
}

func SetDefn(nm string, d interface{}) {
	var ok		bool
	var sw		bool
	var str		string

	switch nm {
	case "cmd":
		if str, ok = d.(string); ok {
			cmd = str
		}
	case "dataPath":
		if str, ok = d.(string); ok {
			dataPath = str
		}
	case "debug":
		if sw, ok = d.(bool); ok {
			debug = sw
		}
	case "force":
		if sw, ok = d.(bool); ok {
			force = sw
		}
	case "mainPath":
		if str, ok = d.(string); ok {
			mainPath = str
		}
	case "mdlDir":
		if str, ok = d.(string); ok {
			mdlDir = str
		}
	case "noop":
		if sw, ok = d.(bool); ok {
			noop = sw
		}
	case "outDir":
		if str, ok = d.(string); ok {
			outDir = str
		}
	case "quiet":
		if sw, ok = d.(bool); ok {
			quiet = sw
		}
	case "timeNow":
		if str, ok = d.(string); ok {
			timeNow = str
		}
	default:
		defns[nm] = d
	}
}

func Force() bool {
	return force
}

func SetForce(f bool) {
	force = f
}

func Funcs() map[string]interface{} {
	return funcs
}

func FuncsSlice() []interface{} {
	var f = []interface{}{}

	for _, v := range funcs {
		f = append(f, v)
	}

	return f
}

func SetFunc(nm string, d interface{}) {
	funcs[nm] = d
}

// MainPath is the path to the main json file.
func MainPath() string {
	return mainPath
}

func SetMainPath(f string) {
	mainPath = f
}

func MdlDir() string {
	return mdlDir
}

func SetMdlDir(f string) {
	mdlDir = f
}

func Noop() bool {
	return noop
}

func SetNoop(f bool) {
	noop = f
}

func OutDir() string {
	return outDir
}

func SetOutDir(f string) {
	outDir = f
}

func Quiet() bool {
	return quiet
}

func SetQuiet(f bool) {
	quiet = f
}

// String returns a stringified version of the shared data
func String() string {
	s := "{"
	s += fmt.Sprintf("cmd:%q,",cmd)
	s += fmt.Sprintf("dataPath:%q,",dataPath)
	s += fmt.Sprintf("debug:%v,",debug)
	s += fmt.Sprintf("force:%v,",force)
	s += fmt.Sprintf("mainPath:%q,",mainPath)
	s += fmt.Sprintf("mdlDir:%q,",mdlDir)
	s += fmt.Sprintf("noop:%v,",noop)
	s += fmt.Sprintf("outDir:%q,",outDir)
	s += fmt.Sprintf("quiet:%v,",quiet)
	s += fmt.Sprintf("time:%q,",timeNow)
	s += "}"
	return s
}

// MainPath is the path to the main json file.
func Time() string {
	return timeNow
}

func SetTime(f string) {
	timeNow = f
}

