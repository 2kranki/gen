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
var funcs		map[string]interface{}
var mainPath	string
var	mdlDir		string
var	outDir		string

func init() {
	defns  = map[string]interface{}{}
	funcs  = map[string]interface{}{}
	mdlDir = "./src/models"
	outDir = "./src/test"
	defns["Debug"] = false
	defns["Force"] = false
	defns["GenDebugging"] = false
	funcs["GenDebugging"] = GenDebugging
	defns["GenLogging"] = false
	funcs["GenLogging"] = GenLogging
	defns["Noop"] = false
	defns["Quiet"] = false
	defns["Replace"] = true
	defns["Time"] = time.Now().Format("Mon Jan _2, 2006 15:04")
	funcs["Time"] = Time
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
	return defns["Debug"].(bool)
}

func SetDebug(f bool) {
	defns["Debug"] = f
}

func Defn(nm string) interface{} {
	switch nm {
	case "cmd":
		return cmd
	case "dataPath":
		return dataPath
	case "mainPath":
		return mainPath
	case "mdlDir":
		return mdlDir
	case "outDir":
		return outDir
	}
	 d, _ := defns[nm]
	return d
}

func IsDefined(nm string) bool {
	x := Defn(nm)
	if x != nil {
		return true
	}
	return false
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
	case "Debug":
		if sw, ok = d.(bool); ok {
			defns["Debug"] = sw
		}
	case "Force":
		if sw, ok = d.(bool); ok {
			defns["Force"] = sw
		}
	case "mainPath":
		if str, ok = d.(string); ok {
			mainPath = str
		}
	case "mdlDir":
		if str, ok = d.(string); ok {
			mdlDir = str
		}
	case "Noop":
		if sw, ok = d.(bool); ok {
			defns["Noop"] = sw
		}
	case "outDir":
		if str, ok = d.(string); ok {
			outDir = str
		}
	case "Quiet":
		if sw, ok = d.(bool); ok {
			defns["Quiet"] = sw
		}
	case "Replace":
		if sw, ok = d.(bool); ok {
			defns["Replace"] = sw
		}
	case "Time":
		if str, ok = d.(string); ok {
			defns["Time"] = str
		}
	default:
		defns[nm] = d
	}
}

func Force() bool {
	return defns["Force"].(bool)
}

func SetForce(f bool) {
	defns["Force"] = f
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

func GenDebugging() bool {
	return defns["GenDebugging"].(bool)
}

func GenLogging() bool {
	return defns["GenLogging"].(bool)
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

// MergeFrom merges the given map into the shared
// data definitions optionally replacing any that
// already exist.
func MergeFrom(m map[string]interface{}, rep bool) {
	var ok			bool

	for k, v := range m {
		if rep {
			defns[k] = v
		} else {
			if _, ok = defns[k]; !ok {
				defns[k] = v
			}
		}
	}
}

// MergeTo merges the shared into the given map
// optionally replacing any in the given map.\
func MergeTo(m map[string]interface{}, rep bool) {
	var ok			bool

	for k, v := range defns {
		if rep {
			m[k] = v
		} else {
			if _, ok = m[k]; !ok {
				m[k] = v
			}
		}
	}
}

func Noop() bool {
	return defns["Noop"].(bool)
}

func SetNoop(f bool) {
	defns["Noop"] = f
}

func OutDir() string {
	return outDir
}

func SetOutDir(f string) {
	outDir = f
}

func Quiet() bool {
	return defns["Quiet"].(bool)
}

func SetQuiet(f bool) {
	defns["Quiet"] = f
}

func Replace() bool {
	return defns["Replace"].(bool)
}

func SetReplace(f bool) {
	defns["Replace"] = f
}

// String returns a stringified version of the shared data
func String() string {
	s := "{"
	s += fmt.Sprintf("cmd:%q,",cmd)
	s += fmt.Sprintf("dataPath:%q,",dataPath)
	s += fmt.Sprintf("Debug:%v,",defns["Debug"])
	s += fmt.Sprintf("Force:%v,",defns["Force"])
	s += fmt.Sprintf("mainPath:%q,",mainPath)
	s += fmt.Sprintf("mdlDir:%q,",mdlDir)
	s += fmt.Sprintf("Noop:%v,",defns["Noop"])
	s += fmt.Sprintf("outDir:%q,",outDir)
	s += fmt.Sprintf("Quiet:%v,",defns["Quiet"])
	s += fmt.Sprintf("Time:%q,",defns["Time"])
	s += "}"
	return s
}

// MainPath is the path to the main json file.
func Time() string {
	return defns["Time"].(string)
}

func SetTime(f string) {
	defns["Time"] = f
}

