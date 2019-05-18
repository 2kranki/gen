// See License.txt in main repository directory

// dbPkg contains the data and functions to generate
// table and field data for html forms, handlers and
// table sql i/o for a specific database.  Multiple
// databases should be handled with multiple ??? of
// this package.

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

const (
	DBTYPE_MARIABDB	= 1 << iota
	DBTYPE_MSSQL
	DBTYPE_MYSQL
	DBTYPE_POSTGRES
	DBTYPE_SQLITE
)

// StringsBuilder is a composition of strings.Builder so that
// we can add supplemental functions.
type StringBuilder struct {
	str		strings.Builder
}

func NewStringBuilder() *StringBuilder {
	sb := StringBuilder{}
	return &sb
}

func (s StringBuilder) String( ) string {
	return s.str.String()
}

// WriteString allows us to write a string to the buffer.
func (s StringBuilder) WriteString(format string) error {
	_, err := s.str.WriteString(format)
	return err
}

// WriteStringf allows us to write a formatted string.
func (s StringBuilder) WriteStringf(format string, a ...interface{}) error {
	str := fmt.Sprintf(format, a...)
	err := s.WriteString(str)
	return err
}

// DbProp defines an Objects Property.
type DbProp struct {
	Name		string		`json:"Name"`					// External Name
	Internal	string		`json:"Internal,omitempty"`		// Optional Internal Name
	Desc		string		`json:"Desc,omitempty"`			// Optional Description
	TypeDefn	string		`json:"TypeDef,omitempty"`		// Type Definition
	Init		string		`json: Init,omitempty`			// Initialization
	Object  	bool	    `json:"Object,omitempty"`
	Vis			string		`json:"Vis,omitempty"`			// Visibility: public,private,read-only,none
	Base		string		`json:"Base,omitempty"`			// Base Struct or Pointer for Offset
	Offset		int	    	`json:"Offset,omitempty"`		// Offset into Base in bytes
	Shift	    int		    `json:"Len,omitempty"`			// Shift Amount to put field in lowest bit
	Size	    int		    `json:"Dec,omitempty"`			// Size in bits
	Nullable	bool		`json:"Null,omitempty"`			// Allow NULL for this field
}

func (f *DbProp) GenInit( ) string {
	var str			*StringBuilder

	str = NewStringBuilder()
	str.WriteStringf("\tthis->%s = %s;\n",f.Name, f.Init)

	return str.String()
}

func (f *DbProp) GenDefn() string {
	var str			strings.Builder
	var name		string

	if len(f.Internal) > 0 {
		name = f.Internal
	} else {
		name = f.Name
	}
	str.WriteString(fmt.Sprintf("\t%s\t\t%s;\n", name, f.TypeDefn))

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
	Name		string		`json:"Name,omitempty"`
	Props		[]DbProp	`json:"Properties,omitempty"`
}

func (t *DbObject) ForFields(f func(f *DbProp) ) {
	for i, _ := range t.Props {
		f(&t.Props[i])
	}
}

func (t *DbObject) TitledName( ) string {
	return strings.Title(t.Name)
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

