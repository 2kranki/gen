// See License.txt in main repository directory

// dbJson contains the database definitions as defined
// by the user.

package dbJson

import (
	"../../shared"
	"../../util"
	"../dbPlugin"
	"../dbType"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

//============================================================================
//                        JSON Database Field Support
//============================================================================

// DbField defines a Table's field mostly in terms of
// SQL.
type DbField struct {
	Name		string				`json:"Name,omitempty"`			// Field Name
	Label		string				`json:"Label,omitempty"`		// Form Label
	TypeDefn	string				`json:"TypeDef,omitempty"`		// Type Definition
	Len		    int		    		`json:"Len,omitempty"`			// Data Maximum Length
	Dec		    int		    		`json:"Dec,omitempty"`			// Decimal Positions
	PrimaryKey  bool	    		`json:"PrimaryKey,omitempty"`
	Hidden		bool	    		`json:"Hidden,omitempty"`		// Do not display in the browser
	Nullable	bool				`json:"Null,omitempty"`			// Allow NULL for this field
	Incr		bool				`json:"Incr,omitempty"`			// true == Auto Increment Field
	SQLParms	string				`json:"SQLParms,omitempty"`		// Extra SQL Parameters
	List		bool	    		`json:"List,omitempty"`			// Include in List Report
	Tbl			*DbTable			`json:"-"`						// (ignored)  Filled in after JSON is parsed
	Typ			*dbType.TypeDefn	`json:"-"`						// (ignored) Filled in after JSON is parsed
}

func (f *DbField) TitledName( ) string {
	return strings.Title(f.Name)
}

//============================================================================
//                        JSON Database Table Support
//============================================================================

// DbTable stands for Database Table and defines
// the make up of the SQL Table.
// Fields should be in the order in which they are to
// be displayed in the list form and the main form.
type DbTable struct {
	Name		string		`json:"Name,omitempty"`
	Fields		[]DbField	`json:"Fields,omitempty"`
	SQLParms	[]string	`json:"SQLParms,omitempty"`		// Extra SQL Parameters
	DB			*Database	`json:"-"`
}

func (t *DbTable) ForFields(f func(f *DbField) ) {
	for i,_ := range t.Fields {
		f(&t.Fields[i])
	}
}

func (t *DbTable) TitledName( ) string {
	return strings.Title(t.Name)
}

//============================================================================
//                        	JSON Database Support
//============================================================================

type Database struct {
	Name		string			`json:"Name,omitempty"`
	SqlType		string			`json:"SqlType,omitempty"`
	SQLParms	string			`json:"SQLParms,omitempty"`		// Extra SQL Parameters
	Server		string			`json:"Server,omitempty"`
	Port		string			`json:"Port,omitempty"`
	PW			string			`json:"PW,omitempty"`
	Tables  	[]DbTable		`json:"Tables,omitempty"`
	// There can only be one Plugin per Database Definition.  Once we have decoded
	// the JSON, we will establish which plugin works with this JSON data if any.
	Plugin		interface{}		`json:"-"`
}

func (d *Database) ForTables(f func(t *DbTable) ) {
	for i,_ := range d.Tables {
		f(&d.Tables[i])
	}
}

func (d *Database) TitledName( ) string {
	return strings.Title(d.Name)
}

//----------------------------------------------------------------------------
//						Global Support Functions
//----------------------------------------------------------------------------

// New provides a factory method to create an Sql Object.
func New() (*Database) {
	db := &Database{}
	return db
}

var	dbStruct	Database

func DbStruct() *Database {
	return &dbStruct
}

func DefaultJsonFileName() string {
	return "db.json.txt"
}

// ReadJsonFile reads the input JSON file for app
// and stores the generic JSON Table as well as the
// decoded structs.
func ReadJsonFile(fn string) error {
	var err		    error
	var jsonPath	string
	var plg			interface{}

	jsonPath, _ = filepath.Abs(fn)
	if sharedData.Debug() {
		log.Println("json path:", jsonPath)
	}

	// Read in the json file structurally
	if err = util.ReadJsonFileToData(jsonPath, &dbStruct); err != nil {
		return errors.New(fmt.Sprintln("Error: unmarshalling:", jsonPath, ":", err))
	}

	// Set up Plugin Support for this database type.
	if plg = dbPlugin.FindPlugin(dbStruct.SqlType); plg == nil {
		return errors.New(fmt.Sprintf("Error: Can't find plugin for %s!\n\n\n", dbStruct.SqlType))
	}
	dbStruct.Plugin = plg

	// Fix up the tables with back pointers that we do not store externally.
	for i, v := range dbStruct.Tables {
		for ii, _ := range v.Fields {
			v.Fields[ii].Tbl = &v
		}
		dbStruct.Tables[i].DB = &dbStruct
	}

	if err = ValidateData(); err != nil {
		return err
	}

	if sharedData.Debug() {
		log.Printf("\tdbStruct: %+v\n", dbStruct)
	}

	return nil
}

func TableNames() []string {
	var list	[]string

	for _, v := range dbStruct.Tables {
		list = append(list, v.Name)
	}

	return list
}

func ValidateData() error {
	var plg		dbPlugin.Plugin

	// Set up Plugin Support for this database type.
	if plg = dbPlugin.FindPlugin(dbStruct.SqlType); plg == nil {
		return errors.New(fmt.Sprintf("Error: Can't find plugin for %s!\n\n\n", dbStruct.SqlType))
	}

	if dbStruct.Name == "" {
		return errors.New(fmt.Sprintf("Database Name is missing!"))
	}
	if dbStruct.SqlType == "" {
		return errors.New(fmt.Sprintf("SQL Type is missing!"))
	}
	if len(dbStruct.Tables) == 0 {
		return errors.New(fmt.Sprintf("There are no tables defined for %s!", dbStruct.Name))
	}
	for i, t := range dbStruct.Tables {
		if t.Name == "" {
			return errors.New(fmt.Sprintf("%d Table Name is missing!", i))
		}
		if len(t.Fields) == 0 {
			return errors.New(fmt.Sprintf("There are no fields defined for %s!", t.Name))
		}
		if t.PrimaryKey() == nil {
			return errors.New(fmt.Sprintf("There is no key defined for %s!", t.Name))
		}
		for j,f := range t.Fields {
			if f.Name == "" {
				return errors.New(fmt.Sprintf("%d Field Name is missing from table %s!", j, t.Name))
			}
			td := plg.Types().FindDefn(f.TypeDefn)
			if td == nil {
				log.Fatalln("Error - Could not find Type definition for field,",
					f.Name,"type:",f.TypeDefn)
			}
		}
	}

	return nil
}

