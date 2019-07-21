// See License.txt in main repository directory

// dbPkg contains the data and functions to generate
// table and field data for html forms, handlers and
// table sql i/o for a specific database.  Multiple
// databases should be handled with multiple ??? of
// this package.

// The template system does not allow nested if
// statements.  So, we are doing more generation
// here than we would like.

package genCObj

import (
	"../shared"
	"../util"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

// DbProp defines an Objects Property.
type DbProp struct {
	Name		string		`json:"name,omitempty"`			// External Name
	Internal	string		`json:"internal,omitempty"`		// Optional Internal Name
	External	string		`json:"external,omitempty"`		// Optional External Name
	Desc		string		`json:"desc,omitempty"`			// Optional Description
	TypeDefn	string		`json:"type,omitempty"`			// Type Definition
	Init		string		`json: init,omitempty`			// Initialization
	Object  	bool	    `json:"object,omitempty"`
	Vis			string		`json:"vis,omitempty"`			// Visibility: public,private,read-only,none
	Base		string		`json:"base,omitempty"`			// Base Struct or Pointer for Offset
	Offset		int	    	`json:"offset,omitempty"`		// Offset into Base in bytes
	Shift	    int		    `json:"len,omitempty"`			// Shift Amount to put field in lowest bit
	Size	    int		    `json:"Dec,omitempty"`			// Size in bits
}

func (f *DbProp) GenBody( ) string {
	var str			strings.Builder
	var name		string

	if f.Vis == "none" {
		return ""
	}
	if len(f.Internal) > 0 {
		name = f.Internal
	} else {
		name = f.Name
	}

	str.WriteString("\n\t//---------------------------------------------------------------\n")
	str.WriteString(fmt.Sprintf("\t//\t\t\t %s\n", f.Name))
	str.WriteString("\t//---------------------------------------------------------------\n\n")

	// Generate Get()
	if f.Object && f.TypeDefn != "OBJ_ID" {
		str.WriteString(fmt.Sprintf("\t%s *\t\t\t\t%s_get%s(\n", f.TypeDefn, dbStruct.Name, f.TitledName()))
	} else {
		str.WriteString(fmt.Sprintf("\t%s\t\t\t\t%s_get%s(\n", f.TypeDefn, dbStruct.Name, f.TitledName()))
	}
	str.WriteString(fmt.Sprintf("\t\t%s_DATA\t*this\n", dbStruct.UpperName()))
	str.WriteString("\t)\n")
	str.WriteString("\t{\n\n")
	str.WriteString("#ifdef NDEBUG\n")
	str.WriteString("#else\n")
	str.WriteString(fmt.Sprintf("\t\tif (%s_Validate(this)) {\n", dbStruct.Name))
	str.WriteString("\t\t\tDEBUG_BREAK();\n")
	if f.Object {
		str.WriteString("\t\t\treturn OBJ_NIL;\n")
	} else {
		str.WriteString("\t\t\treturn 0;\n")
	}
	str.WriteString("\t\t}\n")
	str.WriteString("#endif\n\n")
	str.WriteString(fmt.Sprintf("\t\treturn this->%s;\n", name))
	str.WriteString("\t}\n\n\n")

	// Generate Set()
	str.WriteString(fmt.Sprintf("\tbool\t\t\t\t%s_set%s(\n", dbStruct.Name, f.TitledName()))
	str.WriteString(fmt.Sprintf("\t\t%s_DATA\t*this,\n", dbStruct.UpperName()))
	if f.Object && f.TypeDefn != "OBJ_ID" {
		str.WriteString(fmt.Sprintf("\t\t%s\t\t*pValue,\n", f.TypeDefn))
	} else {
		value := "value"
		if f.TypeDefn == "OBJ_ID" {
			value = "pValue"
		}
		str.WriteString(fmt.Sprintf("\t\t%s\t\t%s,\n", f.TypeDefn, value))
	}
	str.WriteString("\t)\n")
	str.WriteString("\t{\n")
	str.WriteString("#ifdef NDEBUG\n")
	str.WriteString("#else\n")
	str.WriteString(fmt.Sprintf("\t\tif (%s_Validate(this)) {\n", dbStruct.Name))
	str.WriteString("\t\t\tDEBUG_BREAK();\n")
	str.WriteString("\t\t\treturn false;\n")
	str.WriteString("\t\t}\n")
	str.WriteString("#endif\n\n")
	if f.Object {
		str.WriteString("\t\tobj_Retain(pValue);\n")
		str.WriteString(fmt.Sprintf("\t\tif (this->%s) {\n", name))
		str.WriteString(fmt.Sprintf("\t\t\tobj_Release(this->%s)\n", name))
		str.WriteString("\t\t}\n")
		str.WriteString(fmt.Sprintf("\t\tthis->%s = pValue;\n", name))
	} else {
		str.WriteString(fmt.Sprintf("\t\tthis->%s = value;\n", name))
	}
	str.WriteString("\t\treturn true;\n")
	str.WriteString("\t}\n\n\n\n")

	return str.String()
}

func (f *DbProp) GenDefn() string {
	var str			strings.Builder

	if f.Vis != "none" && f.Vis != "private" {
		if len(f.Desc) > 0 {
			str.WriteString(fmt.Sprintf("//%s - %s\n", f.Name, f.Desc))
		}
		if f.Object && f.TypeDefn != "OBJ_ID" {
			str.WriteString(fmt.Sprintf("\t%s *\t\t\t\t%s_get%s(\n", f.TypeDefn, dbStruct.Name, f.TitledName()))
		} else {
			str.WriteString(fmt.Sprintf("\t%s\t\t\t\t%s_get%s(\n", f.TypeDefn, dbStruct.Name, f.TitledName()))
		}
		str.WriteString(fmt.Sprintf("\t\t%s_DATA\t*this\n", dbStruct.UpperName()))
		str.WriteString("\t);\n\n")
		if f.Vis == "public" {
			str.WriteString(fmt.Sprintf("\tbool\t\t\t\t%s_set%s(\n", dbStruct.Name, f.TitledName()))
			str.WriteString(fmt.Sprintf("\t\t%s_DATA\t*this,\n", dbStruct.UpperName()))
			if f.Object && f.TypeDefn != "OBJ_ID" {
				str.WriteString(fmt.Sprintf("\t\t%s\t\t*pValue\n", f.TypeDefn))
			} else {
				value := "value"
				if f.TypeDefn == "OBJ_ID" {
					value = "pValue"
				}
				str.WriteString(fmt.Sprintf("\t\t%s\t\t%s,\n", f.TypeDefn, value))
			}
			str.WriteString("\t);\n\n\n")
		} else {
			str.WriteString("\n")
		}
	}

	return str.String()
}

func (f *DbProp) GenDefnPrivate() string {
	var str			strings.Builder

	if f.Vis == "private" || f.Vis == "read-only" || f.Vis == "ro" {
		if len(f.Desc) > 0 {
			str.WriteString(fmt.Sprintf("\t//%s - %s\n", f.Name, f.Desc))
		}
	}
	if f.Vis == "private" {
		if f.Object && f.TypeDefn != "OBJ_ID" {
			str.WriteString(fmt.Sprintf("\t%s *\t\t\t\t%s_get%s(\n", f.TypeDefn, dbStruct.Name, f.TitledName()))
		} else {
			str.WriteString(fmt.Sprintf("\t%s\t\t\t\t%s_get%s(\n", f.TypeDefn, dbStruct.Name, f.TitledName()))
		}
		str.WriteString(fmt.Sprintf("\t\t%s_DATA\t*this\n", dbStruct.UpperName()))
		str.WriteString("\t);\n\n")
	}
	if f.Vis == "private" || f.Vis == "read-only" || f.Vis == "ro" {
		str.WriteString(fmt.Sprintf("\tbool\t\t\t\t%s_set%s(\n", dbStruct.Name, f.TitledName()))
		str.WriteString(fmt.Sprintf("\t\t%s_DATA\t*this,\n", dbStruct.UpperName()))
		if f.Object && f.TypeDefn == "OBJ_ID" {
			str.WriteString(fmt.Sprintf("\t\t%s\t\t*pValue\n", f.TypeDefn))
		} else {
			value := "value"
			if f.TypeDefn == "OBJ_ID" {
				value = "pValue"
			}
			str.WriteString(fmt.Sprintf("\t\t%s\t\t%s,\n", f.TypeDefn, value))
		}
		str.WriteString("\t);\n\n\n")
	}

	return str.String()
}

func (f *DbProp) GenInit( ) string {
	var str			strings.Builder
	var name		string

	if len(f.Internal) > 0 {
		name = f.Internal
	} else {
		name = f.Name
	}

	if len(f.Init) > 0 {
		str.WriteString(fmt.Sprintf("\t\tthis->%s = %s;\n", name, f.Init))
	}

	return str.String()
}

func (f *DbProp) GenProp( ) string {
	var str			strings.Builder
	var name		string

	if len(f.Internal) > 0 {
		name = f.Internal
	} else {
		name = f.Name
	}

	str.WriteString("\t{ ")
	str.WriteString(fmt.Sprintf("\"%s\",", f.Name))
	str.WriteString(fmt.Sprintf("\"%s\",", f.Internal))
	str.WriteString(fmt.Sprintf("\"%s\",", f.External))
	str.WriteString(fmt.Sprintf("\"%s\",", f.Desc))
	str.WriteString(fmt.Sprintf("\"%s\",", f.TypeDefn))
	str.WriteString(fmt.Sprintf("\"%s\",", f.Init))
	str.WriteString(fmt.Sprintf("\"%s\",", f.Vis))
	if len(f.Base) > 0 {
		str.WriteString(fmt.Sprintf("\"%s\",", f.Base))
	} else {
		str.WriteString(fmt.Sprintf("\"%s\",", dbStruct.DataName()))
	}
	if f.Offset < 0 {
		str.WriteString(fmt.Sprintf("offsetof(%s,%s),", dbStruct.DataName(), name))
	} else {
		str.WriteString(fmt.Sprintf("\"%s_DATA\",", dbStruct.UpperName()))
	}
	if f.Size == 0 {
		if f.Object && f.TypeDefn != "OBJ_ID" {
			str.WriteString(fmt.Sprintf("(sizeof(%s *) << 3),", f.TypeDefn))
		} else {
			str.WriteString(fmt.Sprintf("(sizeof(%s) << 3),", f.TypeDefn))
		}
	} else {
		str.WriteString(fmt.Sprintf("%d,", f.Size))
	}
	str.WriteString(fmt.Sprintf("%d", f.Shift))
	str.WriteString("\t},\n")

	return str.String()
}

func (f *DbProp) GenStruct() string {
	var str			strings.Builder
	var name		string

	if len(f.Internal) > 0 {
		name = f.Internal
	} else {
		name = f.Name
	}

	if f.Object && f.TypeDefn != "OBJ_ID" {
		str.WriteString(fmt.Sprintf("\t%s\t\t*%s;\n", f.TypeDefn, name))
	} else {
		str.WriteString(fmt.Sprintf("\t%s\t\t%s;\n", f.TypeDefn, name))
	}

	return str.String()
}

func (f *DbProp) TitledName( ) string {
	return strings.Title(f.Name)
}

// DbObject stands for Database Table and defines
// the make up of the SQL Table.
// Fields should be in the order in which they are to
// be displayed in the list form and the main form.
type DbObject struct {
	Name		string		`json:"name,omitempty"`
	Super		string		`json:"super,omitempty"`
	Props		[]DbProp	`json:"properties,omitempty"`
}

func (t *DbObject) DataName( ) string {
	name := strings.ToUpper(t.Name)
	return fmt.Sprintf("%s_DATA", name)
}

func (t *DbObject) DataSuper( ) string {
	var name	string
	if len(t.Super) > 0 {
		name = strings.ToUpper(t.Super)
	} else {
		name = "OBJ"
	}
	return fmt.Sprintf("%s_DATA", name)
}

func (t *DbObject) ForFields(f func(f *DbProp) ) {
	for i, _ := range t.Props {
		f(&t.Props[i])
	}
}

func (t *DbObject) PropCount( ) string {
	str := fmt.Sprintf("%d", len(t.Props))
	return str
}

func (t *DbObject) TitledName( ) string {
	return strings.Title(t.Name)
}

func (t *DbObject) TitledSuper( ) string {
	return strings.Title(t.Super)
}

func (t *DbObject) UpperName( ) string {
	return strings.ToUpper(t.Name)
}

var	dbStruct	DbObject

func DbStruct() *DbObject {
	return &dbStruct
}

// ReadJsonFileDb reads the input JSON file for app
// and stores the generic JSON Table as well as the
// decoded structs.
func ReadJsonFile(fn string) error {
	var err		    error
	var jsonPath	string

	jsonPath,_ = filepath.Abs(fn)
	if sharedData.Debug() {
		log.Println("json path:", jsonPath)
	}

	// Read in the json file structurally
	if err = util.ReadJsonFileToData(jsonPath, &dbStruct); err != nil {
		return errors.New(fmt.Sprintln("Error: unmarshalling", jsonPath, ", JSON input file:", err))
	}

	if sharedData.Debug() {
		log.Println("\tJson Struct:", dbStruct)
	}

	return nil
}

