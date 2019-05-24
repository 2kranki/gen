// See License.txt in main repository directory

// dbPkg contains the data and functions to generate
// table and field data for html forms, handlers and
// table sql i/o for a specific database.  Multiple
// databases should be handled with multiple ??? of
// this package.

package dbData

import (
	"../../shared"
	"../../util"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	DBTYPE_MARIABDB	= 1 << iota
	DBTYPE_MSSQL
	DBTYPE_MYSQL
	DBTYPE_POSTGRES
	DBTYPE_SQLITE
)

type Plugin_Data	struct {
	Name		string
	T			*TypeDefns
	ImportString func() string
	AddGo		bool			// Add "GO" after major sql statements
	CreateDB	bool
}

/**
func (pd Plugin_Data) ImportString() string {
	return ""
}
**/

var plugins		map[string]*Plugin_Data
// The index into Plugins is the database name as used in
// the JSON input (ie mariadb, mssql, mysql, postgresql, sqlite)

// Register() registers the given plugin into the map.
func Register(pd *Plugin_Data) error {
	if plugins == nil {
		plugins = map[string]*Plugin_Data{}
	}
	plugins[pd.Name] = pd
	return nil
}

// Unregister() unregisters a given plugin in the map.
func Unregister(name string) {
	if plugins == nil {
		return
	}
	if _, ok := plugins[name]; ok {
		delete(plugins, name)
	}
}

// Plugin returns the Plugin interface for a name if possible.
func Plugin(name string) *Plugin_Data {
	if plugins == nil {
		return nil
	}
	if _, ok := plugins[name]; ok {
		return plugins[name]
	}
	return nil
}

type TypeDefn struct {
	Name		string		`json:"Name,omitempty"`		// Type Name
	Html		string		`json:"Html,omitempty"`		// HTML Type
	Sql			string		`json:"Sql,omitempty"`		// SQL Type
	Go			string		`json:"Go,omitempty"`		// GO Type
	DftLen		int			`json:"DftLen,omitempty"`	// Default Length (used if length is not
	//													//	given)(0 == Max Length)
}

type TypeDefns []TypeDefn

func (t TypeDefns) DftLen(name string) int {
	tdd := t.FindDefn(name)
	if tdd != nil {
		return tdd.DftLen
	}
	return -1
}

func (t TypeDefns) FindDefn(name string) *TypeDefn {
	for i,v := range t {
		if name == v.Name {
			return &t[i]
		}
	}
	return nil
}

func (t TypeDefns) GoType(name string) string {
	tdd := t.FindDefn(name)
	if tdd != nil {
		return tdd.Go
	}
	return ""
}

func (t TypeDefns) HtmlType(name string) string {
	tdd := t.FindDefn(name)
	if tdd != nil {
		return tdd.Html
	}
	return ""
}

func (t TypeDefns) SqlType(name string) string {
	tdd := t.FindDefn(name)
	if tdd != nil {
		return tdd.Sql
	}
	return ""
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
	Hidden		bool	    `json:"Hidden,omitempty"`		// Do not display in the browser
	Nullable	bool		`json:"Null,omitempty"`			// Allow NULL for this field
	SQLParms	string		`json:"SQLParms,omitempty"`		// Extra SQL Parameters
	List		bool	    `json:"List,omitempty"`			// Include in List Report
	Tbl			*DbTable									// Filled in after JSON is parsed
}

func (f *DbField) CreateSql(cm string) string {
	var str			strings.Builder
	var ft			string
	var nl			string
	var pk			string
	var sp			string

	td := Plugin(dbStruct.SqlType).T.FindDefn(f.TypeDefn)
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
	sp = ""
	if len(f.SQLParms) > 0 {
		sp = fmt.Sprintf(" %s", f.SQLParms)
	}

	str.WriteString(fmt.Sprintf("\\t%s\\t%s%s%s%s%s\\n", f.Name, ft, nl, pk, cm, sp))

	return str.String()
}

func (f *DbField) CreateStruct() string {
	var str			strings.Builder

	tdd := f.GoType()
	str.WriteString(fmt.Sprintf("\t%s\t%s\n", strings.Title(f.Name),tdd))

	return str.String()
}

func (f *DbField) FormInput() string {
	var str			strings.Builder
	var lbl			string
	var m			string

	td := Plugin(dbStruct.SqlType).T.FindDefn(f.TypeDefn)
	if td == nil {
		log.Fatalln("Error - Could not find Type definition for field,",
			f.Name,"type:",f.TypeDefn)
	}

	tdd := td.Html
	if len(f.Label) > 0 {
		lbl = strings.Title(f.Label)
	} else {
		lbl = strings.Title(f.Name)
	}
	switch td.Name {
	case "money":
		m = "m=\"0\" step=\"0.01\" "
	default:
		m = ""
	}

	if f.Hidden {
		str.WriteString(fmt.Sprintf("\t<input type=\"hidden\" name=\"%s\" id=\"%s\" %svalue=\"{{.Rcd.%s}}\">\n",
			f.TitledName(), f.TitledName(), m, f.TitledName()))
	} else {
		str.WriteString(fmt.Sprintf("\t<label>%s: <input type=\"%s\" name=\"%s\" id=\"%s\" %svalue=\"{{.Rcd.%s}}\"></label>\n",
			lbl, tdd, f.TitledName(), f.TitledName(), m, f.TitledName()))
	}

	return str.String()
}

// GenFromStringArray generates the code to go from a string array
// (sn) element (n) to a field (dn).  sn and dn are variable names.
func (f *DbField) GenFromStringArray(dn,sn string, n int) string {
	var str			string
	var src			string

	src = sn + "[" + strconv.Itoa(n) + "]"
	str = f.GenFromString(dn, src)

	return str
}

// GenFromString generates the code to go from a string (sn) to
// a field (dn).  sn and dn are variable names.
func (f *DbField) GenFromString(dn,sn string) string {
	var str			string

	td := Plugin(dbStruct.SqlType).T.FindDefn(f.TypeDefn)
	if td == nil {
		log.Fatalln("Error - Could not find Type definition for field,",
			f.Name,"type:",f.TypeDefn)
	}

	switch td.Name {
	case "int":
		fallthrough
	case "integer":
		{
			wrk := "\t%s.%s, err = strconv.Atoi(%s)\n"
			str = fmt.Sprintf(wrk, dn, f.TitledName(), sn )
		}
	case "dec":
		fallthrough
	case "decimal":
		fallthrough
	case "money":
		{
			wrk := 	"\t{\n\t\twrk := r.FormValue(\"%s\")\n" +
				"\t\t%s.%s, err = strconv.ParseFloat(wrk, 64)\n\t}\n"
			str = fmt.Sprintf(wrk, f.TitledName(), dn, f.TitledName())
		}
	default:
		str = fmt.Sprintf("\t%s.%s = %s\n", dn, f.TitledName(), sn)
	}

	return str
}

// GenToString generates code to convert the struct st.f field to string in variable, v.
func (f *DbField) GenToString(v string, st string) string {
	var str			string

	td := Plugin(dbStruct.SqlType).T.FindDefn(f.TypeDefn)
	if td == nil {
		log.Fatalln("Error - Could not find Type definition for field,",
			f.Name,"type:",f.TypeDefn)
	}

	tdd := td.Name
	switch tdd {
	case "int":
		fallthrough
	case "integer":
		str = fmt.Sprintf("\t%s = strconv.Itoa(%s.%s)\n", v, st, f.TitledName())
	case "dec":
		fallthrough
	case "decimal":
		fallthrough
	case "money":
		str = fmt.Sprintf("\t{\n")
		str += fmt.Sprintf("\t\ts := fmt.Sprintf(\"%s.4f\", %s.%s)\n", "%", st, f.TitledName())
		str += fmt.Sprintf("\t\t%s = strings.TrimRight(strings.TrimRight(s, \"0\"), \".\")\n", v)
		str += fmt.Sprintf("\t}\n")
	default:
		str = fmt.Sprintf("\t%s = %s.%s\n", v, st, f.TitledName())
	}

	return str
}

func (f *DbField) GoType() string {

	td := Plugin(dbStruct.SqlType).T.FindDefn(f.TypeDefn)
	if td == nil {
		log.Fatalln("Error - Could not find Type definition for field,",
			f.Name,"type:",f.TypeDefn)
	}

	tdd := td.Go

	return tdd
}

func (f *DbField) IsFloat() bool {

	tdd := f.GoType()
	if tdd == "float64" {
		return true
	}

	return false
}

func (f *DbField) IsInteger() bool {

	tdd := f.GoType()
	if tdd == "int32" {
		return true
	}
	if tdd == "int64" {
		return true
	}
	if tdd == "int" {
		return true
	}

	return false
}

func (f *DbField) IsText() bool {

	if f.TypeDefn == "text" {
		return true
	}

	return false
}

func (f *DbField) RValueToStruct(dn string) string {
	var str			string

	td := Plugin(dbStruct.SqlType).T.FindDefn(f.TypeDefn)
	if td == nil {
		log.Fatalln("Error - Could not find Type definition for field,",
			f.Name,"type:",f.TypeDefn)
	}

	tdd := td.Name
	switch tdd {
	case "int":
		fallthrough
	case "integer":
		{
			wrk := "\twrk = r.FormValue(\"%s\")\n" +
				"\t%s.%s, err = strconv.Atoi(wrk)\n"
			str = fmt.Sprintf(wrk, f.TitledName(), dn, f.TitledName())
		}
	case "dec":
		fallthrough
	case "decimal":
		fallthrough
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
	SQLParms	[]string	`json:"SQLParms,omitempty"`		// Extra SQL Parameters
	DB			*Database
}

// CreateInsertStr() creates a string of all the field
// names which can be used in SQL INSERT statements.
func (t *DbTable) CreateInsertStr() string {
	return t.ScanFields("")
}

func (t *DbTable) CreateSql() string {
	var str			strings.Builder

	str.WriteString(fmt.Sprintf("CREATE TABLE %s (\\n", t.Name))
	for i, f := range t.Fields {
		var cm  		string

		cm = ""
		if i != (len(t.Fields) - 1) {
			cm = ","
		}
		str.WriteString(fmt.Sprintf("%s\\n", f.CreateSql(cm)))
	}
	if len(t.SQLParms) > 0 {
		str.WriteString(",\\n")
		for _, l := range t.SQLParms {
			str.WriteString(fmt.Sprintf("%s\\n", l))
		}
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
	str.WriteString("}\n\n")

	// I was generating some of the struct functions here.  It turned out to be a
	// mistake.  Using the template system and supplement it with small functions
	// is far easier making it a much better strategy.

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

func (t *DbTable) FieldIndex(n string) int {
	for i, f := range t.Fields {
		if f.Name == n {
			return i
		}
	}
	return -1
}

// HasFloat returns true if any of the fields are a
// float which will need float to string conversion
func (t *DbTable) HasFloat() bool {

	for i,_ := range t.Fields {
		if t.Fields[i].IsFloat() {
			return true
		}
	}
	return false
}

// HasInteger returns true if any of the fields are a
// integers which will need float to string conversion
func (t *DbTable) HasInteger() bool {

	for i,_ := range t.Fields {
		if t.Fields[i].IsInteger() {
			return true
		}
	}
	return false
}

// PrimaryKey returns the first field that it finds
// that is marked as a primary key.
func (t *DbTable) PrimaryKey() *DbField {

	for i, f := range t.Fields {
		if f.PrimaryKey {
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
	Name		string			`json:"Name,omitempty"`
	SqlType		string			`json:"SqlType,omitempty"`
	SQLParms	string			`json:"SQLParms,omitempty"`		// Extra SQL Parameters
	Server		string			`json:"Server,omitempty"`
	Port		string			`json:"Port,omitempty"`
	PW			string			`json:"PW,omitempty"`
	Tables  	[]DbTable		`json:"Tables,omitempty"`
	ImportStr	string
	Plugin		*Plugin_Data
}

func (d *Database) CreateSql() string {
	var str			strings.Builder

	str.WriteString(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;\\n", d.TitledName))
	if dbStruct.SqlType == "mssql" {
		str.WriteString("GO\\n")
	}
	str.WriteString(fmt.Sprintf("USE %s;\\n", d.TitledName))
	if dbStruct.SqlType == "mssql" {
		str.WriteString("GO\\n")
	}

	return str.String()
}

func (d *Database) DeleteSql() string {
	var str			strings.Builder

	str.WriteString(fmt.Sprintf("DROP DATABASE IF EXISTS %s;\\n", d.TitledName))
	if dbStruct.SqlType == "mssql" {
		str.WriteString("GO\\n")
	}
	return str.String()
}

func (d *Database) ForTables(f func(t *DbTable) ) {
	for i,_ := range d.Tables {
		f(&d.Tables[i])
	}
}

func (d *Database) HasFloat( ) bool {
	for i,_ := range d.Tables {
		if d.Tables[i].HasFloat() {
			return true
		}
	}
	return false
}

func (d *Database) HasMoney( ) bool {
	for i,_ := range d.Tables {
		if d.Tables[i].HasFloat() {
			return true
		}
	}
	return false
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

// ReadJsonFile reads the input JSON file for app
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

	if err = ValidateData(); err != nil {
		return err
	}

	// Fix up the tables with back pointers that we do not store externally.
	for i, v := range dbStruct.Tables {
		for ii, _ := range v.Fields {
			v.Fields[ii].Tbl = &v
		}
		dbStruct.Tables[i].DB = &dbStruct
	}
	if sharedData.Debug() {
		log.Printf("\tplugins: %d\n", len(plugins))
		for n, v := range plugins {
			log.Printf("\t\tplugin: %s %q\n", n, v)
		}
	}
	if plg := Plugin(dbStruct.SqlType); plg != nil {
		dbStruct.ImportStr = plg.ImportString()
		dbStruct.Plugin = plg
	} else {
		return errors.New(fmt.Sprintf("Error: Can't find import string for %s!\n\n\n", dbStruct.SqlType))
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
		if t.PrimaryKey() == nil {
			return errors.New(fmt.Sprintf("There is no key defined for %s!", t.Name))
		}
		for j,f := range t.Fields {
			if f.Name == "" {
				return errors.New(fmt.Sprintf("%d Field Name is missing from table %s!", j, t.Name))
			}
			td := dbStruct.Plugin.T.FindDefn(f.TypeDefn)
			if td == nil {
				log.Fatalln("Error - Could not find Type definition for field,",
					f.Name,"type:",f.TypeDefn)
			}
		}
	}

	return nil
}
