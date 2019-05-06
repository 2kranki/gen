// See License.txt in main repository directory

// dbPkg contains the data and functions to generate
// table and field data for html forms, handlers and
// table sql i/o for a specific database.  Multiple
// databases should be handled with multiple ??? of
// this package.

package dbPkg

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

type TypeDefn struct {
	Name		string		`json:"Name,omitempty"`		// Type Name
	Html		string		`json:"Html,omitempty"`		// HTML Type
	Sql			string		`json:"Sql,omitempty"`		// SQL Type
	Go			string		`json:"Go,omitempty"`		// GO Type
	DftLen		int			`json:"DftLen,omitempty"`	// Default Length (used if length is not
	//													//	given)(0 == Max Length)
}

type TypeDefns []TypeDefn

func (t TypeDefns) FindDefn(Name string) *TypeDefn {
	for i,v := range t {
		if Name == v.Name {
			return &t[i]
		}
	}
	return nil
}

var tds	= TypeDefns {
	{Name:"date", 		Html:"date", 		Sql:"DATE", 		Go:"string",	DftLen:0,},
	{Name:"datetime",	Html:"datetime",	Sql:"DATETIME",		Go:"string",	DftLen:0,},
	{Name:"email", 		Html:"email", 		Sql:"NVARCHAR", 	Go:"string",	DftLen:50,},
	{Name:"dec", 		Html:"number",		Sql:"DEC",			Go:"float64",	DftLen:0,},
	{Name:"decimal", 	Html:"number",		Sql:"DEC",			Go:"float64",	DftLen:0,},
	{Name:"int", 		Html:"number",		Sql:"INT",			Go:"int",		DftLen:0,},
	{Name:"integer", 	Html:"number",		Sql:"INT",			Go:"int",		DftLen:0,},
	{Name:"money", 		Html:"number",		Sql:"DEC",			Go:"float64",	DftLen:0,},
	{Name:"number", 	Html:"number",		Sql:"INT",			Go:"int",		DftLen:0,},
	{Name:"tel", 		Html:"tel",			Sql:"NVARCHAR",		Go:"string",	DftLen:19,},	//+nnn (nnn) nnn-nnnn
	{Name:"text", 		Html:"text",		Sql:"NVARCHAR",		Go:"string",	DftLen:0,},
	{Name:"time", 		Html:"time",		Sql:"TIME",			Go:"string",	DftLen:0,},
	{Name:"url", 		Html:"url",			Sql:"NVARCHAR",		Go:"string",	DftLen:50,},
}

// DbField defines a Table's field mostly in terms of
// SQL.
type DbField struct {
	Name		string		`json:"Name,omitempty"`			// Field Name
	Label		string		`json:"Label,omitempty"`		// Form Label
	TypeDefn	string		`json:"TypeDef,omitempty"`		// Type Definition
	Len		    int		    `json:"Len,omitempty"`			// Data Maximum Length
	Dec		    int		    `json:"Dec,omitempty"`			// Decimal Positions
	PrimaryKey  bool	    `json:"PrimaryKey,omitempty"`
	Nullable	bool		`json:"Null,omitempty"`
	List		bool	    `json:"List,omitempty"`			// Include in List Report
}

func (f *DbField) CreateSql(cm string) string {
	var str			strings.Builder
	var ft			string
	var nl			string
	var pk			string

	td := tds.FindDefn(f.TypeDefn)
	if td == nil {
		log.Fatalln("Error - Could not find Type definition for field,",
			f.Name,"type:",f.TypeDefn)
	}
	tdd := td.Sql

	if f.Len > 0 {
		if f.Dec > 0 {
			ft = fmt.Sprintf("%s(%d,%d)", tdd, f.Len, f.Dec)
		} else {
			ft = fmt.Sprintf("%s(%d)", tdd, f.Len)
		}
	} else {
		ft = tdd
	}
	nl = " NOT NULL"
	if f.Nullable {
		nl = ""
	}
	pk = ""
	if f.PrimaryKey {
		pk = " PRIMARY KEY"
	}
	str.WriteString(fmt.Sprintf("\\t%s\\t%s%s%s%s\\n", f.Name, ft, nl, pk, cm))

	return str.String()
}

func (f *DbField) CreateStruct() string {
	var str			strings.Builder

	td := tds.FindDefn(f.TypeDefn)
	if td == nil {
		log.Fatalln("Error - Could not find Type definition for field,",
			f.Name,"type:",f.TypeDefn)
	}

	tdd := td.Go
	str.WriteString(fmt.Sprintf("\t%s\t%s\n", strings.Title(f.Name),tdd))

	return str.String()
}

func (f *DbField) FormInput() string {
	var str			strings.Builder

	td := tds.FindDefn(f.TypeDefn)
	if td == nil {
		log.Fatalln("Error - Could not find Type definition for field,",
			f.Name,"type:",f.TypeDefn)
	}

	tdd := td.Html
	str.WriteString(fmt.Sprintf("\t<label>%s: <input type=\"%s\" value={{%s}}></label>\n",
		strings.Title(f.Label), strings.Title(tdd), f.TitledName()))

	return str.String()
}

func (f *DbField) RValueToStruct(dn string) string {
	var str			string

	td := tds.FindDefn(f.TypeDefn)
	if td == nil {
		log.Fatalln("Error - Could not find Type definition for field,",
			f.Name,"type:",f.TypeDefn)
	}

	tdd := td.Name
	switch tdd {
	case "dec":
		fallthrough
	case "decimal":
		str = fmt.Sprintf("\t%s.%s = strconv.Atoi(r.FormValue(\"%s\"))\n", dn, f.TitledName(), f.TitledName())
	case "int":
		fallthrough
	case "integer":
		{
			wrk := "\twrk = r.FormValue(\"%s\")\n" +
				"\t%s.%s, err = strconv.Atoi(wrk)\n"
			str = fmt.Sprintf(wrk, f.TitledName(), dn, f.TitledName())
		}
	case "money":
		{
			wrk := 	"\twrk = r.FormValue(\"%s\")\n" +
					"\t%s.%s, err = strconv.ParseFloat(wrk, 64)\n"
			str = fmt.Sprintf(wrk, f.TitledName(), dn, f.TitledName())
		}
	default:
		str = fmt.Sprintf("\t%s.%s = r.FormValue(\"%s\")\n", dn, f.TitledName(), f.TitledName())
	}

	return str
}

func (f *DbField) TitledName( ) string {
	return strings.Title(f.Name)
}

// DbTable stands for Database Table and defines
// the make up of the SQL Table.
// Fields should be in the order in which they are to
// be displayed in the list form and the main form.
type DbTable struct {
	Name		string		`json:"Name,omitempty"`
	Fields		[]DbField	`json:"Fields,omitempty"`
}

// CreateInsertStr() creates a string of all the field
// names which can be used in SQL INSERT statements.
func (t *DbTable) CreateInsertStr() string {
	return t.ScanFields("")
}

func (t *DbTable) CreateSql() string {
	var str			strings.Builder

	str.WriteString(fmt.Sprintf("CREATE TABLE %s (\\n", t.Name))
	for i,f := range t.Fields {
		var cm  		string

		cm = ""
		if i != (len(t.Fields) - 1) {
			cm = ","
		}
		str.WriteString(f.CreateSql(cm))
	}
	str.WriteString(fmt.Sprintf(");\\n"))
	if dbStruct.SqlType == "mssql" {
		str.WriteString("GO\\n")
	}

	return str.String()
}

func (t *DbTable) CreateStruct( ) string {
	var str			strings.Builder

	str.WriteString(fmt.Sprintf("type %s struct {\n", t.TitledName()))
	for i,_ := range t.Fields {
		str.WriteString(t.Fields[i].CreateStruct())
	}
	str.WriteString("}\n")

	return str.String()
}

// CreateValueStr() creates a string of $nnn's
// which can be used in SQL INSERT VALUE statements.
func (t *DbTable) CreateValueStr() string {

	insertStr := ""
	for i, _ := range t.Fields {
		cm := ", "
		if i == len(t.Fields) - 1 {
			cm = ""
		}
		insertStr += fmt.Sprintf("$%d%s", i+1, cm)
	}
	return insertStr
}

func (t *DbTable) DeleteSql() string {
	var str			strings.Builder

	str.WriteString(fmt.Sprintf("DROP TABLE IF EXISTS %s;\\n", t.Name))
	if dbStruct.SqlType == "mssql" {
		str.WriteString("GO\\n")
	}
	return str.String()
}

func (t *DbTable) ForFields(f func(f *DbField) ) {
	for i,_ := range t.Fields {
		f(&t.Fields[i])
	}
}

func (t *DbTable) Form() string {
	var str			strings.Builder
	str.WriteString(fmt.Sprintf("<form id=\"%s\" method=\"post\">\n", t.Name))
	for _, v := range t.Fields {
		str.WriteString(fmt.Sprintf("<p>%s</p>\n",v.Name))
	}
	str.WriteString("<p/>\n<p/>\n<p/>\n")
	str.WriteString("\t<input type=submit onclick='onPrev()' value=\"Prev\">\n")
	str.WriteString("\t<input type=submit onclick='onAdd()' value=\"Add\">\n")
	str.WriteString("\t<input type=submit onclick='onDelete()' value=\"Delete\">\n")
	str.WriteString("\t<input type=submit onclick='onUpdate()' value=\"Update\">\n")
	str.WriteString("\t<input type=submit onclick='onNext()' value=\"Next\">\n")
	str.WriteString("\t<input type=reset onclick='onReset()' value=\"Reset\">\n")
	str.WriteString("</form>\n\n")
	str.WriteString("<script>\n")
	str.WriteString("\tfunction onAdd() {\n")
	str.WriteString(fmt.Sprintf("\t\tdocument.getElementById(\"%s\").action = \"/Create/Process\";\n", t.Name))
	str.WriteString("\t}\n")
	str.WriteString("\tfunction onDelete() {\n")
	str.WriteString(fmt.Sprintf("\t\tdocument.getElementById(\"%s\").action = \"/Delete/Process\";\n",t.Name))
	str.WriteString("\t}\n")
	str.WriteString("\tfunction onNext() {\n")
	str.WriteString(fmt.Sprintf("\t\tdocument.getElementById(\"%s\").action = \"/Next\";\n",t.Name))
	str.WriteString(fmt.Sprintf("\t\tdocument.getElementById(\"%s\").method = \"get\";\n",t.Name))
	str.WriteString("\t}\n")
	str.WriteString("\tfunction onPrev() {\n")
	str.WriteString(fmt.Sprintf("\t\tdocument.getElementById(\"%s\").action = \"/Prev\";\n",t.Name))
	str.WriteString(fmt.Sprintf("\t\tdocument.getElementById(\"%s\").method = \"get\";\n",t.Name))
	str.WriteString("\t}\n")
	str.WriteString("\tfunction onReset() {\n")
	str.WriteString("\t}\n")
	str.WriteString("\tfunction onUpdate() {\n")
	str.WriteString(fmt.Sprintf("\t\tdocument.getElementById(\"%s\").action = \"/Update/Process\";\n",t.Name))
	str.WriteString("\t}\n")
	str.WriteString("</script>\n")
	return str.String()
}

// PrimaryKey returns the first field that it finds
// that is marked as a primary key.
func (t *DbTable) PrimaryKey() *DbField {

	for i,_ := range t.Fields {
		if t.Fields[i].PrimaryKey {
			return &t.Fields[i]
		}
	}
	return nil
}

// ScanFields returns struct fields to be used in
// a row.Scan.  It assumes that the struct's name
// is "data"
func (t *DbTable) ScanFields(prefix string) string {
	var str			strings.Builder

	for i,f := range t.Fields {
		cm := ", "
		if i == len(t.Fields) - 1 {
			cm = ""
		}
		if len(prefix) > 0 {
			str.WriteString(fmt.Sprintf("%s.%s%s", prefix, f.Name, cm))
		} else {
			str.WriteString(fmt.Sprintf("%s%s", f.Name, cm))
		}
	}
	return str.String()
}

func (t *DbTable) TitledName( ) string {
	return strings.Title(t.Name)
}

type Database struct {
	Name	string			`json:"Name,omitempty"`
	SqlType	string			`json:"SqlType,omitempty"`
	Server	string			`json:"Server,omitempty"`
	Port	string			`json:"Port,omitempty"`
	PW		string			`json:"PW,omitempty"`
	Tables  []DbTable		`json:"Tables,omitempty"`
}

func (d *Database) ForTables(f func(t *DbTable) ) {
	for i,_ := range d.Tables {
		f(&d.Tables[i])
	}
}

func (d *Database) TitledName( ) string {
	return strings.Title(d.Name)
}

var	dbStruct	Database

func DbStruct() *Database {
	return &dbStruct
}

func DefaultJsonFileName() string {
	return "db.json.txt"
}

func InsertSql(t interface{}) string {
	//var Fields 	[]map[string] interface{}
	//var ok		bool
	var x		string

	//insertStr := ""
	return x
}

func ForTables(f func(*DbTable)) {
	for i,_ := range dbStruct.Tables {
		f(&dbStruct.Tables[i])
	}
}

func GenAccessFunc(t DbTable) string {
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

func GenAccessFuncs() string {
	var str			strings.Builder
	for _, v := range dbStruct.Tables {
		str.WriteString(GenAccessFunc(v))
	}
	return str.String()
}

func GenListField(f DbField) string {
	var str			strings.Builder

	if f.PrimaryKey {
		str.WriteString("<a href=\"\">")
	}

	return str.String()
}

func GenListBody(t *DbTable) string {
	var str			strings.Builder
	for _, v := range t.Fields {
		str.WriteString(GenListField(v))
	}
	return str.String()
}

// init() adds the functions needed for templating to
// shared data.
func init() {
	sharedData.SetFunc("GenAccessFuncs", GenAccessFuncs)
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

func TableNames() []string {
	var list	[]string

	for _, v := range dbStruct.Tables {
		list = append(list, v.Name)
	}

	return list
}

func ValidateData() error {

	switch dbStruct.SqlType {
	case "mariadb":
	case "mssql":
	case "mysql":
	case "postgres":
	case "sqlite":
	default:
		return errors.New(fmt.Sprintf("SqlType of %s is not supported!",dbStruct.SqlType))
	}
	if dbStruct.Name == "" {
		return errors.New(fmt.Sprintf("Database Name is missing!"))
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
		for j,f := range t.Fields {
			if f.Name == "" {
				return errors.New(fmt.Sprintf("%d Field Name is missing from table %s!", j, t.Name))
			}
			td := tds.FindDefn(f.TypeDefn)
			if td == nil {
				log.Fatalln("Error - Could not find Type definition for field,",
					f.Name,"type:",f.TypeDefn)
			}
		}
	}

	return nil
}
