// See License.txt in main repository directory

// Template Functions used in generation
// See the specific template files for how the functions
// and data are used.

package cobjData

import (
	"../shared"
	"../util"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

type DbField struct {
	Name		string		`json:"Name,omitempty"`
	Type		string		`json:"Type,omitempty"`
	Desc		string		`json:"Desc,omitempty"`
	Len		    int		    `json:"Len,omitempty"`
	Offset		int		    `json:"Offset,omitempty"`
	ShiftAmt	int		    `json:"ShiftAmt,omitempty"`
}

type Database struct {
	Name	    string			`json:"Name,omitempty"`
    Fields		[]DbField       `json:"Fields,omitempty"`
}

var	cobjStruct	Database

func CObjStruct() *Database {
	return &cobjStruct
}

// GenHeaderDefns
func GenHeaderDefns() string {
	var str			strings.Builder
	str.WriteString(fmt.Sprintf("\tfunc %sDeleteRow( ) {\n", t.Name))
	str.WriteString("\t}\n\n")
	str.WriteString(fmt.Sprintf("\tfunc %sInsertRow( ) {\n", t.Name))
	str.WriteString("\t}\n\n")
	str.WriteString(fmt.Sprintf("\tfunc %sSelect(sel string) ([]string, error) {\n", t.Name))
	/***
	  func {{title $t.Name }}Select(sel string) []string,error {
	      {{ if eq .Data.SqlType "mariadb" }}
	          ERROR - NOT IMPLEMENTED
	      {{ else if eq .Data.SqlType "mssql" }}
	      _ "github.com/2kranki/go-mssqldb"
	      {{ else if eq .Data.SqlType "mysql" }}
	          _ "github.com/go-sql-driver/mysql"
	      {{ else if eq .Data.SqlType "postgres" }}
	          _ "github.com/lib/pq"
	      {{ else if eq .Data.SqlType "sqlite" }}
	      _ "github.com/2kranki/go-sqlite3"
	      {{ end }}

	  }

	 */
	str.WriteString("\t}\n\n")
	str.WriteString(fmt.Sprintf("\tfunc %sSetupRow(r *http.Request) {\n", t.Name))
	/***
	    func {{ title $t.Name }}SetupRow(r *http.Request) {
	        data := interface{}
	        key := r.FormValue("{{$t.PrimaryKey}}")
		    if key == "" {
			    return data, errors.New("400. Bad Request.")
		    }
	        row := config.DB.QueryRow("SELECT * FROM {{$t.Name}} WHERE {{$t.PrimaryKey}} = $1", key)
	        err := row.Scan(
	                    &data.Isbn,
	                    &data.Title,
	                    &data.Author,
	                    &data.Price)
	        if err != nil {
	        	return data, err
	        }
	        	return data, nil
	    }
	*/
	str.WriteString("\t}\n\n")
	str.WriteString(fmt.Sprintf("\tfunc %sUpdateRow( ) {\n", t.Name))
	str.WriteString("\t}\n\n")
	return str.String()
}

// GenInternalDefn generates the internal field definitions.
func GenInternalDefn(f DbField) string {
	var str			strings.Builder

	str.WriteString("\t")
	str.WriteString(f.Type)
	str.WriteString("\t")
	str.WriteString(f.Name)
	str.WriteString(";\n")

	return str.String()
}

// GenInternalDefns generates the internal field definitions.
func GenInternalDefns() string {
	var str			strings.Builder

	for _, v := range cobjStruct.Fields {
		str.WriteString(GenInternalDefn(v))
	}
	str.WriteString("\n")

	return str.String()
}

// GenBody generates the getter/setter bodies.
func GenGSBody(f DbField) string {
	var str			strings.Builder

	str.WriteString(fmt.Sprintf("\tfunc %sDeleteRow( ) {\n", f.Name))
	str.WriteString("\t}\n\n")

	return str.String()
}

// GenGetSetBodies generates the getter/setter bodies.
func GenGetSetBodies() string {
	var str			strings.Builder

	for _, v := range cobjStruct.Fields {
		str.WriteString(GenInternalDefn(v))
	}

	return str.String()
}

// init() adds the functions needed for templating to
// shared data.
func init() {
	//sharedData.SetFunc("GenFlagSetup", GenFlagSetup)
	sharedData.SetFunc("GenGetSetBodies", GenGetSetBodies)
	sharedData.SetFunc("GenHeaderDefns", GenHeaderDefns)
	sharedData.SetFunc("GenInternalDefns", GenInternalDefns)
}

// ReadJsonFileApp reads the input JSON file for app
// and stores the generic JSON Table as well as the
// decoded structs.
func ReadJsonFileApp(fn string) error {
	var err		    error
	var jsonPath	string

	jsonPath,_ = filepath.Abs(fn)
	if sharedData.Debug() {
		log.Println("json path:", jsonPath)
	}

	// Read in the json file structurally
	if err = util.ReadJsonFileToData(jsonPath, &cobjStruct); err != nil {
		return errors.New(fmt.Sprintln("Error: unmarshalling", jsonPath, ", JSON input file:", err))
	}

	if sharedData.Debug() {
		log.Println("\tJson Struct:", cobjStruct)
	}

	return nil
}


