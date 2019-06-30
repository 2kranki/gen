// See License.txt in main repository directory

// dbJson contains the database definitions as defined
// by the user.

package dbJson

import (
	"../../shared"
	"../../util"
	"../dbPlugin"
	"../dbType"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strconv"
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
	KeyNum  	int	    			`json:"KeyNum,omitempty"`		// If not a key field, then 0. Otherwise, 1 for
	//																// highest level key, 2 for 2nd highest, ...
	Hidden		bool	    		`json:"Hidden,omitempty"`		// Do not display in the browser
	Nullable	bool				`json:"Null,omitempty"`			// Add NULL for this field
	Unique		bool				`json:"Unique,omitempty"`		// Add UNIQUE to this field
	Incr		bool				`json:"Incr,omitempty"`			// true == Auto Increment Field
	SQLParms	string				`json:"SQLParms,omitempty"`		// Extra SQL Parameters
	List		bool	    		`json:"List,omitempty"`			// Include in List Report
	Tbl			*DbTable			`json:"-"`						// (ignored)  Filled in after JSON is parsed
	Typ			*dbType.TypeDefn	`json:"-"`						// (ignored) Filled in after JSON is parsed
}

func (f *DbField) CreateSql(cm string) string {
	var str			strings.Builder
	var ft			string
	var nl			string
	var pk			string
	var sp			string

	td := f.Typ
	if td == nil {
		log.Fatalln("Error - Could not find Type definition for field,",
			f.Name,"type:",f.TypeDefn)
	}
	tdd := f.Typ.SqlType()

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
	//FIXME: if f.PrimaryKey {
		//pk = " PRIMARY KEY"
	//}
	sp = ""
	if len(f.SQLParms) > 0 {
		sp = fmt.Sprintf(" %s", f.SQLParms)
	}

	str.WriteString(fmt.Sprintf("\\t%s\\t%s%s%s%s%s\\n", f.Name, ft, nl, pk, cm, sp))

	return str.String()
}

func (f *DbField) CreateStruct() string {
	var str			strings.Builder

	str.WriteString(fmt.Sprintf("\t%s\t%s\n", strings.Title(f.Name),f.GoType()))

	return str.String()
}

func (f *DbField) FormInput() string {
	var str			strings.Builder
	var lbl			string
	var m			string

	tdd := f.Typ.Html
	if len(f.Label) > 0 {
		lbl = strings.Title(f.Label)
	} else {
		lbl = strings.Title(f.Name)
	}
	switch f.Typ.GoType() {
	case "float64":
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
// a field of (dn).  sn and dn are variable names.
func (f *DbField) GenFromString(dn,sn string) string {
	var str			string

	switch f.Typ.GoType() {
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
	case "time.Time":
		{
			wrk := "\t%s.%s, err = time.Parse(time.RFC3339, %s)\n"
			str = fmt.Sprintf(wrk, dn, f.TitledName(), sn )
		}
	default:
		str = fmt.Sprintf("\t%s.%s = %s\n", dn, f.TitledName(), sn)
	}

	return str
}

// GenToString generates code to convert the struct st.f field to string in variable, v.
func (f *DbField) GenToString(v string, st string) string {
	var str			string

	switch f.Typ.GoType() {
	case "int":
		fallthrough
	case "int32":
		fallthrough
	case "int64":
		str = fmt.Sprintf("\t%s = fmt.Sprintf(\"%%d\", %s.%s)\n", v, st, f.TitledName())
	case "float32":
		fallthrough
	case "float64":
		str = "\t{\n"
		str += fmt.Sprintf("\t\ts := fmt.Sprintf(\"%s.4f\", %s.%s)\n", "%", st, f.TitledName())
		str += fmt.Sprintf("\t\t%s = strings.TrimRight(strings.TrimRight(s, \"0\"), \".\")\n", v)
		str += "\t}\n"
	case "time.Time":
		{
			wrk := "\t{\n\t\twrk, _ := %s.%s.MarshalText()\n" +
				"\t\t%s = wrk\n\t}\n"
			str = fmt.Sprintf(wrk, st, f.TitledName(), v)
		}
	default:
		str = fmt.Sprintf("\t%s = %s.%s\n", v, st, f.TitledName())
	}

	return str
}

func (f *DbField) GoType() string {
	return f.Typ.GoType()
}

func (f *DbField) IsDate() bool {

	if f.TypeDefn == "date" {
		return true
	}
	if f.TypeDefn == "datetime" {
		return true
	}

	return false
}

func (f *DbField) IsDec() bool {

	if f.TypeDefn == "dec" {
		return true
	}
	if f.TypeDefn == "decimal" {
		return true
	}
	if f.TypeDefn == "money" {
		return true
	}

	return false
}

func (f *DbField) IsFloat() bool {

	tdd := f.Typ.GoType()
	if tdd == "float64" {
		return true
	}

	return false
}

func (f *DbField) IsInteger() bool {

	tdd := f.Typ.GoType()
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

	tdd := f.Typ.GoType()
	if tdd == "string" {
		return true
	}

	return false
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
	Name			string		`json:"Name,omitempty"`
	Fields			[]DbField	`json:"Fields,omitempty"`
	SQLParms		[]string	`json:"SQLParms,omitempty"`		// Extra SQL Parameters
	DB				*Database	`json:"-"`
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

func (t *DbTable) DeleteSql() string {
	var str			strings.Builder

	str.WriteString(fmt.Sprintf("DROP TABLE IF EXISTS %s;\\n", t.Name))
	if dbStruct.SqlType == "mssql" {
		str.WriteString("GO\\n")
	}
	return str.String()
}

func (t *DbTable) FieldCount() int {
	return len(t.Fields)
}

func (t *DbTable) FieldIndex(name string) int {
	for i, f := range t.Fields {
		if f.Name == name {
			return i
		}
	}
	return -1
}

// FieldNameList returns struct fields separated by
// commas with an optional per field prefix.
func (t *DbTable) FieldNameList(prefix string) string {
	var str			strings.Builder

	for i,f := range t.Fields {
		cm := ", "
		if i == len(t.Fields) - 1 {
			cm = ""
		}
		if len(prefix) > 0 {
			str.WriteString(fmt.Sprintf("%s%s%s", prefix, f.Name, cm))
		} else {
			str.WriteString(fmt.Sprintf("%s%s", f.Name, cm))
		}
	}
	return str.String()
}

func (t *DbTable) FindField(name string) *DbField {
	for i, f := range t.Fields {
		if f.Name == name {
			return &t.Fields[i]
		}
	}
	return nil
}

func (t *DbTable) FindIndex(idx int) *DbField {
	if idx < len(t.Fields) && idx >= 0 {
		return &t.Fields[idx]
	}
	return nil
}

func (t *DbTable) ForFields(f func(f *DbField) ) {
	for i, _ := range t.Fields {
		f(&t.Fields[i])
	}
}

// HasDate returns true if any of the fields are a
// date or datetime type.
func (t *DbTable) HasDate() bool {

	for i,_ := range t.Fields {
		if t.Fields[i].IsDate() {
			return true
		}
	}
	return false
}

// HasDec returns true if any of the fields are a
// decimal type which will need string conversion
func (t *DbTable) HasDec() bool {

	for i,_ := range t.Fields {
		if t.Fields[i].IsDec() {
			return true
		}
	}
	return false
}

// HasFloat returns true if any of the fields are a
// float which will need float to string conversion
func (t *DbTable) HasFloat() bool {

	for _, f := range t.Fields {
		if f.IsFloat() {
			return true
		}
	}
	return false
}

// HasInteger returns true if any of the fields are a
// integers which will need float to string conversion
func (t *DbTable) HasInteger() bool {

	for _, f := range t.Fields {
		if f.IsInteger() {
			return true
		}
	}
	return false
}

// HasText returns true if any of the fields are a
// string types
func (t *DbTable) HasText() bool {

	for _, f := range t.Fields {
		if f.IsText() {
			return true
		}
	}
	return false
}

// KeyCount returns the number of key fields in the table.
func (t *DbTable) KeyCount() int {
	var count		int

	// accumulate the number of key fields
	for _, v := range t.Fields {
		if v.KeyNum > 0 {
			count++
		}
	}

	return count
}

// Keys returns the field names marked as keys in ascending order
// by KeyNum which is descending order of importance.
func (t *DbTable) Keys() ([]string, error) {
	var strs  	[]string
	var mapKeys	[]int

	// accumulate the keys
	keys := map[int]string{}
	for _, v := range t.Fields {
		if v.KeyNum > 0 {
			if nm, ok := keys[v.KeyNum]; ok && nm != v.Name {
				return nil, fmt.Errorf("Error: Duplicate Keys - %s %s\n", nm, v.Name)
			}
			keys[v.KeyNum] = v.Name
			mapKeys = append(mapKeys, v.KeyNum)
		}
	}

	// generate the keys in ascending order.
	sort.Ints(mapKeys)
	for i := 0; i < len(mapKeys); i++ {
		key := mapKeys[i]
		strs = append(strs, keys[key])
	}

	return strs, nil
}

// KeysList returns the table's keys in number order as
// a comma separated list.
func (t *DbTable) KeysList(prefix, suffix string) string {
	var str			strings.Builder
	var strs		[]string

	strs, _ = t.Keys()
	for i, fn := range strs {
		cm := ", "
		if i == len(strs) - 1 {
			cm = ""
		}
		pref := ""
		if len(prefix) > 0 {
			pref = prefix
		}
		suf := ""
		if len(suffix) > 0 {
			suf = suffix
		}
		str.WriteString(fmt.Sprintf("%s%s%s%s", pref, fn, suf, cm))
	}
	return str.String()
}

// KeysListStr returns the table's keys in number order as
// a comma separated list of strings.
func (t *DbTable) KeysListStr() string {
	var str			strings.Builder
	var strs		[]string

	strs, _ = t.Keys()
	for i, fn := range strs {
		cm := ", "
		if i == len(strs) - 1 {
			cm = ""
		}
		str.WriteString(fmt.Sprintf("\"%s\"%s", fn, cm))
	}
	return str.String()
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

func (d *Database) FindTable(name string) *DbTable {
	for i, t := range d.Tables {
		if t.Name == name {
			return &d.Tables[i]
		}
	}
	return nil
}

func (d *Database) ForTables(f func(t *DbTable) ) {
	for i, _ := range d.Tables {
		f(&d.Tables[i])
	}
}

func (d *Database) HasDate() bool {
	for _, t := range d.Tables {
		if t.HasDate() {
			return true
		}
	}
	return false
}

func (d *Database) HasFloat() bool {
	for _, t := range d.Tables {
		if t.HasFloat() {
			return true
		}
	}
	return false
}

// SetupPlugin sets up the plugin within the database.
func (d *Database) SetupPlugin(plg dbPlugin.PluginData) error {

	// Validate the Plugin if possible.
	if plg.Types == nil {
		return fmt.Errorf("Error: Plugin missing types for %s!\n\n\n", d.SqlType)
	}
	if plg.Plugin == nil {
		return fmt.Errorf("Error: Plugin missing support for %s!\n\n\n", d.SqlType)
	}

	// Save the plugin.
	d.Plugin = plg

	// Fix up the tables with back pointers that we do not store externally.
	for i, t := range d.Tables {
		for ii, _ := range t.Fields {
			t.Fields[ii].Tbl = &t
			t.Fields[ii].Typ = plg.Types.FindDefn(t.Fields[ii].TypeDefn)
			if t.Fields[ii].Typ == nil {
				return fmt.Errorf("Error: Invalid Field Type for %s:%s!\n\n\n", t.Name, t.Fields[ii].Name)
			}
		}
		d.Tables[i].DB = d
	}

	return nil
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
	var plg			dbPlugin.PluginData

	jsonPath, _ = filepath.Abs(fn)
	if sharedData.Debug() {
		log.Println("json path:", jsonPath)
	}

	// Read in the json file structurally
	if err = util.ReadJsonFileToData(jsonPath, &dbStruct); err != nil {
		return fmt.Errorf("Error: unmarshalling: %s : %s", jsonPath, err)
	}

	// Set up Plugin Support for this database type.
	if plg, err = dbPlugin.FindPlugin(dbStruct.SqlType); err != nil {
		return fmt.Errorf("Error: Can't find plugin for %s!\n\n\n", dbStruct.SqlType)
	}
	err = SetupPlugin(plg)
	if err != nil {
		return fmt.Errorf("Error: Plugin Setup failed - %s\n", err)
	}

	if err = ValidateData(); err != nil {
		return err
	}

	if sharedData.Debug() {
		log.Printf("\tdbStruct: %+v\n", dbStruct)
	}

	return nil
}

// SetupPlugin sets up the plugin within the database.
func SetupPlugin(plg dbPlugin.PluginData) error {

	// Validate the Plugin if possible.
	if plg.Types == nil {
		return fmt.Errorf("Error: Plugin missing types for %s!\n\n\n", dbStruct.SqlType)
	}
	if plg.Plugin == nil {
		return fmt.Errorf("Error: Plugin missing support for %s!\n\n\n", dbStruct.SqlType)
	}

	// Save the plugin.
	dbStruct.Plugin = plg

	// Fix up the tables with back pointers that we do not store externally.
	for i, t := range dbStruct.Tables {
		for ii, _ := range t.Fields {
			t.Fields[ii].Tbl = &t
			t.Fields[ii].Typ = plg.Types.FindDefn(t.Fields[ii].TypeDefn)
			if t.Fields[ii].Typ == nil {
				return fmt.Errorf("Error: Invalid Field Type for %s:%s!\n\n\n", t.Name, t.Fields[ii].Name)
			}
		}
		dbStruct.Tables[i].DB = &dbStruct
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

// ValidateData checks the JSON built structures for errors. Some of
// errors may be duplicates of the JSON Unmarshalling process which
// is ok, because this function can be used if the data is from a
// different source.
func ValidateData() error {
	var plg		dbPlugin.PluginData
	var err		error

	// Set up Plugin Support for this database type.
	if plg, err = dbPlugin.FindPlugin(dbStruct.SqlType); err != nil {
		return err
	}

	if dbStruct.Name == "" {
		return fmt.Errorf("Error: Database Name is missing!")
	}
	if dbStruct.SqlType == "" {
		return fmt.Errorf("Error: SQL Type is missing!")
	}
	if len(dbStruct.Tables) == 0 {
		return fmt.Errorf("There are no tables defined for %s!", dbStruct.Name)
	}
	for i, t := range dbStruct.Tables {
		if t.Name == "" {
			return fmt.Errorf("%d Table Name is missing!", i)
		}
		if len(t.Fields) == 0 {
			return fmt.Errorf("There are no fields defined for %s!", t.Name)
		}
		if _, err = t.Keys(); err != nil {
			return err
		}
		for j,f := range t.Fields {
			if f.Name == "" {
				return fmt.Errorf("%d Field Name is missing from table %s!", j, t.Name)
			}
			td := plg.Types.FindDefn(f.TypeDefn)
			if td == nil {
				fmt.Errorf("Error - Could not find Type definition for field: %s  type: %s",
					f.Name, f.TypeDefn)
			}
		}
	}

	return nil
}

