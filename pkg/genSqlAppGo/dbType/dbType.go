// See License.txt in main repository directory

// dbType contains the data and functions that form
// the basis for field definitions.

package dbType

import ()

//============================================================================
//								Interfaces
//============================================================================

//============================================================================
//                        Type Definition Support
//============================================================================

// TypeDefn, The type definition struct, defines one acceptable type accepted in the JSON
// defining a Database column (field) type.  There must be a TypeDefn for each type accepted
// by a plugin within it.

type TypeDefn struct {
	Name   string `json:"Name,omitempty"`   // Type Name
	Html   string `json:"Html,omitempty"`   // HTML Type
	Sql    string `json:"Sql,omitempty"`    // SQL Type
	Go     string `json:"Go,omitempty"`     // GO Type
	DftLen int    `json:"DftLen,omitempty"` // Default Length (used if length is not
	//													//	given)(0 == Max Length)
	GenLen bool `json:"GenLen,omitempty"` // If true, generate the length
	DftDec int  `json:"DftDec,omitempty"` // Default Decimal Positions (used if length is not
	//													//	given)(0 == Not Used)
	GenDec bool `json:"GenDec,omitempty"` // If true, generate the decimal positions
}

// CanGenDec indicates if Decimal Position should be generated in SQL
// and other
func (t TypeDefn) CanGenDec() bool {
	if !t.GenLen {
		return false
	}
	return t.GenDec
}

func (t TypeDefn) CanGenLen() bool {
	return t.GenLen
}

func (t TypeDefn) GoType() string {
	return t.Go
}

func (t TypeDefn) HtmlType() string {
	return t.Html
}

func (t TypeDefn) IsDec() bool {

	if t.Name == "dec" {
		return true
	}
	if t.Name == "decimal" {
		return true
	}
	if t.Name == "money" {
		return true
	}

	return false
}

func (t TypeDefn) IsFloat() bool {

	if t.Go == "float64" {
		return true
	}

	return false
}

func (t TypeDefn) IsInteger() bool {

	if t.Go == "int64" {
		return true
	}
	if t.Go == "int32" {
		return true
	}
	if t.Go == "int" {
		return true
	}

	return false
}

func (t TypeDefn) IsText() bool {

	if t.Name == "text" {
		return true
	}

	return false
}

func (t TypeDefn) SqlType() string {
	return t.Sql
}

//============================================================================
//                        Type Definition Table Support
//============================================================================

// TypeDefns provides a convenient way of defining a Type Definition Table which
// is used in the plugin environment.
type TypeDefns []TypeDefn

func (t TypeDefns) DftLen(name string) int {
	tdd := t.FindDefn(name)
	if tdd != nil {
		return tdd.DftLen
	}
	return -1
}

func (t TypeDefns) FindDefn(name string) *TypeDefn {
	for i, v := range t {
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

//============================================================================
//                        Default Type Definition Table
//============================================================================

// Notes:
//	* We are now using a Decimal Package for support of decimal operations including
//		monetary calculations via https://github.com/ericlagergren/decimal
//  * In this table, we pick the most common types which should be generic to any
//		SQL Server if possible
var DefaultTable = TypeDefns{
	{Name: "date", Html: "date", Sql: "DATE", Go: "string"},
	{Name: "datetime", Html: "datetime", Sql: "DATETIME", Go: "string"},
	{Name: "email", Html: "email", Sql: "VARCHAR", Go: "string", DftLen: 50, GenLen: true},
	{Name: "dec", Html: "number", Sql: "DEC", Go: "string"},
	{Name: "decimal", Html: "number", Sql: "DEC", Go: "string"},
	{Name: "int", Html: "number", Sql: "INT", Go: "int64"},
	{Name: "integer", Html: "number", Sql: "INT", Go: "int64"},
	{Name: "money", Html: "number", Sql: "DEC", Go: "string"},
	{Name: "number", Html: "number", Sql: "INT", Go: "int64"},
	{Name: "tel", Html: "tel", Sql: "VARCHAR", Go: "string", DftLen: 19, GenLen: true}, //+nnn (nnn) nnn-nnnn
	{Name: "text", Html: "text", Sql: "VARCHAR", Go: "string", DftLen: 0, GenLen: true},
	{Name: "time", Html: "time", Sql: "TIME", Go: "string", DftLen: 0, GenLen: false},
	{Name: "url", Html: "url", Sql: "VARCHAR", Go: "string", DftLen: 50, GenLen: true},
}
