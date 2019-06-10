// See License.txt in main repository directory

// dbSql provides the functions to generate the go statements
// necessary to access and manipulate the SQL databases defined
// by the user. The problem that it tries to solve is that while
// SQL is supposed to be a universal language. It unfortunately
// is not and each type of database manager must be handled slightly
// differently.

// We give this package access to user defined JSON and the ap-
// propriate plugin for the data being processed. Between those
// two resources, it must generate the go code.

package dbSql

import (
	"../../shared"
	"../dbJson"
	"fmt"
	"log"
	"strconv"
	"strings"
	"unsafe"
)

//============================================================================
//                        	Interface Support
//============================================================================

// dbSql uses interfaces to determine what a plugin can do or not do and when it
// should be called.  If the plugin does not support a particular interface, then
// dbSql will perform default logic to handle the situation.
//
// The reason for all this is that even though Go uses a "common" interface for
// accessing SQL Servers. The SQL, itself, can vary.  Although SQL is supposed to
// to be a standard, it is not consistently implemented unforturnately.
//
// Functions that return a full SQL statement must return a slice of strings even
// if there is only one statement ever generated.  That is because some servers
// such as Microsoft's SQL Server may not do anything until an additional statement
// is issued such as "go".

//----------------------------------------------------------------------------
//                        	Database Interface Support
//----------------------------------------------------------------------------

type GenDatabaseCreateStmter interface {
	GenDatabaseCreateStmts(db *dbJson.Database) []string
}

type GenDatabaseDeleteStmter interface {
	GenDatabaseDeleteStmts(db *dbJson.Database) []string
}

//----------------------------------------------------------------------------
//                        	Table Interface Support
//----------------------------------------------------------------------------

type GenTableCreateStmter interface {
	GenTableCreateStmts(tb dbJson.DbTable) []string
}

type GenTableDeleteStmter interface {
	GenTableDeleteStmts(tb dbJson.DbTable) []string
}

//----------------------------------------------------------------------------
//                        	Row Interface Support
//----------------------------------------------------------------------------

type GenRowDeleteStmter interface {
	GenRowfDeleteStmts(tb dbJson.DbTable) []string
}

type GenRowFindStmter interface {
	GenRowFindStmts(tb dbJson.DbTable) []string
}

type GenRowFirstStmter interface {
	GenRowFirstStmts(tb dbJson.DbTable) []string
}

type GenRowInsertStmter interface {
	GenRowInsertStmts(tb dbJson.DbTable) []string
}

type GenRowLastStmter interface {
	GenRowLastStmts(tb dbJson.DbTable) []string
}

type GenRowNextStmter interface {
	GenRowNextStmts(tb dbJson.DbTable) []string
}

type GenRowPageStmter interface {
	GenRowPageStmts(tb dbJson.DbTable) []string
}

type GenRowPrevStmter interface {
	GenRowPrevStmts(tb dbJson.DbTable) []string
}

type GenRowUpdateStmter interface {
	GenRowUpdateStmts(tb dbJson.DbTable) []string
}


//============================================================================
//                        Type Definition Support
//============================================================================

type Field		struct {
	F			*dbJson.DbField
	Plg			interface{}
}

func (f *Field) CreateSql(cm string) string {
	var str			strings.Builder
	var ft			string
	var nl			string
	var pk			string
	var sp			string

	td := f.F.Typ
	if td == nil {
		log.Fatalln("Error - Could not find Type definition for field,",
			f.F.Name,"type:",f.F.TypeDefn)
	}
	tdd := td.Sql

	if f.F.Len > 0 {
		if f.F.Dec > 0 {
			ft = fmt.Sprintf("%s(%d,%d)", tdd, f.F.Len, f.F.Dec)
		} else {
			ft = fmt.Sprintf("%s(%d)", tdd, f.F.Len)
		}
	} else {
		ft = tdd
	}
	nl = " NOT NULL"
	if f.F.Nullable {
		nl = ""
	}
	pk = ""
	if f.F.PrimaryKey {
		pk = " PRIMARY KEY"
	}
	sp = ""
	if len(f.F.SQLParms) > 0 {
		sp = fmt.Sprintf(" %s", f.F.SQLParms)
	}

	str.WriteString(fmt.Sprintf("\\t%s\\t%s%s%s%s%s\\n", f.F.Name, ft, nl, pk, cm, sp))

	return str.String()
}

func (f *Field) CreateStruct() string {
	var str			strings.Builder

	tdd := f.GoType()
	str.WriteString(fmt.Sprintf("\t%s\t%s\n", strings.Title(f.F.Name),tdd))

	return str.String()
}

func (f *Field) FormInput() string {
	var str			strings.Builder
	var lbl			string
	var m			string

	td := f.F.Typ
	if td == nil {
		log.Fatalln("Error - Could not find Type definition for field,",
			f.F.Name,"type:",f.F.TypeDefn)
	}

	tdd := td.Html
	if len(f.F.Label) > 0 {
		lbl = strings.Title(f.F.Label)
	} else {
		lbl = strings.Title(f.F.Name)
	}
	switch td.Go {
	case "float64":
		m = "m=\"0\" step=\"0.01\" "
	default:
		m = ""
	}

	if f.F.Hidden {
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
func (f *Field) GenFromStringArray(dn,sn string, n int) string {
	var str			string
	var src			string

	src = sn + "[" + strconv.Itoa(n) + "]"
	str = f.GenFromString(dn, src)

	return str
}

// GenFromString generates the code to go from a string (sn) to
// a field (dn).  sn and dn are variable names.
func (f *Field) GenFromString(dn,sn string) string {
	var str			string

	td := f.F.Typ
	if td == nil {
		log.Fatalln("Error - Could not find Type definition for field,",
			f.F.Name,"type:",f.F.TypeDefn)
	}

	switch td.Go {
	case "int":
		fallthrough
	case "int32":
		fallthrough
	case "int64":
		{
			wrk := "\t%s.%s, err = strconv.ParseInt(%s,0,64)\n"
			str = fmt.Sprintf(wrk, dn, f.TitledName(), sn )
		}
	case "float64":
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
func (f *Field) GenToString(v string, st string) string {
	var str			string

	td := f.F.Typ
	if td == nil {
		log.Fatalln("Error - Could not find Type definition for field,",
			f.F.Name,"type:",f.F.TypeDefn)
	}

	tdd := td.Go
	switch tdd {
	case "int":
		fallthrough
	case "int32":
		fallthrough
	case "int64":
		str = fmt.Sprintf("\t%s = fmt.Sprintf(\"%%d\", %s.%s)\n", v, st, f.TitledName())
	case "float32":
		fallthrough
	case "float64":
		str = fmt.Sprintf("\t{\n")
		str += fmt.Sprintf("\t\ts := fmt.Sprintf(\"%s.4f\", %s.%s)\n", "%", st, f.TitledName())
		str += fmt.Sprintf("\t\t%s = strings.TrimRight(strings.TrimRight(s, \"0\"), \".\")\n", v)
		str += fmt.Sprintf("\t}\n")
	default:
		str = fmt.Sprintf("\t%s = %s.%s\n", v, st, f.TitledName())
	}

	return str
}

func (f *Field) GoType() string {

	td := f.F.Typ
	if td == nil {
		log.Fatalln("Error - Could not find Type definition for field,",
			f.F.Name,"type:",f.F.TypeDefn)
	}

	tdd := td.Go

	return tdd
}

func (f *Field) IsDec() bool {

	if f.F.TypeDefn == "dec" {
		return true
	}
	if f.F.TypeDefn == "decimal" {
		return true
	}
	if f.F.TypeDefn == "money" {
		return true
	}

	return false
}

func (f *Field) IsFloat() bool {

	tdd := f.GoType()
	if tdd == "float64" {
		return true
	}

	return false
}

func (f *Field) IsInteger() bool {

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

func (f *Field) IsText() bool {

	if f.TypeDefn == "text" {
		return true
	}

	return false
}

func (f *Field) RValueToStruct(dn string) string {
	var str			string

	td := f.Typ
	if td == nil {
		log.Fatalln("Error - Could not find Type definition for field,",
			f.Name,"type:",f.TypeDefn)
	}

	tdd := td.Go
	switch tdd {
	case "int":
		fallthrough
	case "int32":
		fallthrough
	case "int64":
		{
			wrk := "\twrk = r.FormValue(\"%s\")\n" +
				"\t%s.%s, err = strconv.ParseInt(wrk,0,64)\n"
			str = fmt.Sprintf(wrk, f.TitledName(), dn, f.TitledName())
		}
	case "float":
		fallthrough
	case "float32":
		fallthrough
	case "float64":
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

func (f *Field) TitledName( ) string {
	return strings.Title(f.Name)
}


type Table		dbJson.DbTable

// CreateInsertStr() creates a string of all the field
// names which can be used in SQL INSERT statements.
func (t *Table) CreateInsertStr() string {
	return t.ScanFields("")
}

func (t *Table) CreateSql() string {
	var str			strings.Builder
	var ff			*Field

	str.WriteString(fmt.Sprintf("CREATE TABLE %s (\\n", t.Name))
	for i, f := range t.Fields {
		var cm  		string

		cm = ""
		if i != (len(t.Fields) - 1) {
			cm = ","
		}
		ff = (*Field)(unsafe.Pointer(&f))
		str.WriteString(fmt.Sprintf("%s\\n", ff.CreateSql(cm)))
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

func (t *Table) CreateStruct( ) string {
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
func (t *Table) CreateValueStr() string {

	insertStr := ""
	for i, _ := range t.Fields {
		cm := ", "
		if i == len(t.Fields) - 1 {
			cm = ""
		}
		insertStr += fmt.Sprintf("?%s", cm)
		//insertStr += fmt.Sprintf("$%d%s", i+1, cm)
	}
	return insertStr
}

func (t *Table) DeleteSql() string {
	var str			strings.Builder

	str.WriteString(fmt.Sprintf("DROP TABLE IF EXISTS %s;\\n", t.Name))
	if dbStruct.SqlType == "mssql" {
		str.WriteString("GO\\n")
	}
	return str.String()
}

func (t *Table) ForFields(f func(f *Field) ) {
	for i,_ := range t.Fields {
		f(&t.Fields[i])
	}
}

func (t *Table) FieldIndex(n string) int {
	for i, f := range t.Fields {
		if f.Name == n {
			return i
		}
	}
	return -1
}

// HasDec returns true if any of the fields are a
// decimal type which will need string conversion
func (t *Table) HasDec() bool {

	for i,_ := range t.Fields {
		if t.Fields[i].IsDec() {
			return true
		}
	}
	return false
}

// HasFloat returns true if any of the fields are a
// float which will need float to string conversion
func (t *Table) HasFloat() bool {

	for i,_ := range t.Fields {
		f := (*Field)(unsafe.Pointer(&t.Fields[i]))
		if f.IsFloat() {
			return true
		}
	}
	return false
}

// HasInteger returns true if any of the fields are a
// integers which will need float to string conversion
func (t *Table) HasInteger() bool {

	for i,_ := range t.Fields {
		f := (*Field)(unsafe.Pointer(&t.Fields[i]))
		if f.IsInteger() {
			return true
		}
	}
	return false
}

// PrimaryKey returns the first field that it finds
// that is marked as a primary key.
func (t *Table) PrimaryKey() *Field {

	for i, f := range t.Fields {
		ff := (*Field)(unsafe.Pointer(&t.Fields[i]))
		if f.PrimaryKey {
			return ff
		}
	}
	return nil
}

// ScanFields returns struct fields to be used in
// a row.Scan.  It assumes that the struct's name
// is "data"
func (t *Table) ScanFields(prefix string) string {
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

func (t *Table) TitledName( ) string {
	return strings.Title(t.Name)
}


type Database	dbJson.Database

func (d *Database) HasFloat( ) bool {
	for i,_ := range d.Tables {
		t := (*Table)(unsafe.Pointer(&d.Tables[i]))
		if t.HasFloat() {
			return true
		}
	}
	return false
}

func (d *Database) HasMoney( ) bool {
	for i,_ := range d.Tables {
		t := (*Table)(unsafe.Pointer(&d.Tables[i]))
		if t.HasFloat() {
			return true
		}
	}
	return false
}


// SqlWork, The type definition struct, defines one acceptable type accepted in
// the JSON defining the Database Structure.  There must be a TypeDefn for each
// type accepted in each plugin.

type SqlWork struct {
	name		string				`json:"Name,omitempty"`		// Type Name
	db			dbJson.Database
}

func (s SqlWork) DB() dbJson.Database {
	return s.db
}

func (s SqlWork) SetDB(db dbJson.Database) {
	s.db = db
}

func (s SqlWork) Name() string {
	return s.name
}

func (s SqlWork) SetName(n string) {
	s.name = n
}

// SqlWorks provides a convenient way of defining a SQL Work Table.
type SqlWorks	[]SqlWork

//----------------------------------------------------------------------------
//						Global/Internal Object Functions
//----------------------------------------------------------------------------

func (t SqlWorks) FindDefn(name string) *SqlWork {
	for i, v := range t {
		if name == v.Name() {
			return &t[i]
		}
	}
	return nil
}

func (s *SqlWork) GenDatabaseCreateStmts() []string {
	var str			strings.Builder
	var strs  		[]string

	str.WriteString(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;\\n", s.db.TitledName))
	if s.db.SqlType == "mssql" {
		str.WriteString("GO\\n")
	}
	str.WriteString(fmt.Sprintf("USE %s;\\n", s.db.TitledName))
	if s.db.SqlType == "mssql" {
		str.WriteString("GO\\n")
	}

	strs = append(strs, str.String())
	return strs
}

func (s *SqlWork) GenTableDeleteStmts(tb *dbJson.DbTable) []string {
	var str			strings.Builder
	var strs  		[]string
	var intr		GenTableDeleteStmter
	var ok			bool

	intr, ok = tb.DB.Plugin.(GenTableDeleteStmter)
	if ok {
		return intr.GenTableDeleteStmts(tb)
	}

	str.WriteString(fmt.Sprintf("DROP DATABASE IF EXISTS %s;\\n", s.db.TitledName))
	if s.db.SqlType == "mssql" {
		str.WriteString("GO\\n")
	}

	strs = append(strs, str.String())
	return strs
}

//----------------------------------------------------------------------------
//						Global Support Functions
//----------------------------------------------------------------------------

func GenDatabaseCreateStmts(db *dbJson.Database) []string {
	var str			strings.Builder
	var strs  		[]string
	var intr		GenDatabaseCreateStmter
	var ok			bool

	intr, ok = db.Plugin.(GenDatabaseCreateStmter)
	if ok {
		return intr.GenDatabaseCreateStmts(db)
	}

	strs = append(strs, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;\\n", db.TitledName))
	if db.SqlType == "mssql" {
		strs = append(strs, "GO\\n")
	}
	strs = append(strs, fmt.Sprintf("USE %s;\\n", db.TitledName))
	if db.SqlType == "mssql" {
		strs = append(strs, "GO\\n")
	}

	strs = append(strs, str.String())
	return strs
}

func GenTableDeleteStmts(db *dbJson.Database) []string {
	var str			strings.Builder
	var strs  		[]string
	var intr		GenDatabaseDeleteStmter
	var ok			bool

	intr, ok = db.Plugin.(GenDatabaseDeleteStmter)
	if ok {
		return intr.GenDatabaseDeleteStmts(db)
	}

	str.WriteString(fmt.Sprintf("DROP DATABASE IF EXISTS %s;\\n", s.db.TitledName))
	if db.SqlType == "mssql" {
		str.WriteString("GO\\n")
	}

	strs = append(strs, str.String())
	return strs
}

// init() is called before main(). Here we define the functions that will be
// used in the templates.
func init() {
	sharedData.SetFunc("GenDatabaseCreateStmts", GenDatabaseCreateStmts)
}

// New provides a factory method to create an Sql Object.
func New(db dbJson.Database) (*SqlWork) {
	sw := &SqlWork{}
	sw.db = db
	return sw
}

