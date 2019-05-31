// See License.txt in main repository directory

// dbSqlite is the plugin for SQLite and  contains the
// data and functions specific for SQLite to generate
// table and field data for html forms, handlers and
// table sql i/o for a specific database.

package dbSqlite

import (
	"../dbPlugin"
	"fmt"
	"log"
	"strings"
)

// Notes:
//	* We are now using a Decimal Package for support of decimal operations including
//		monetary calculations via https://github.com/ericlagergren/decimal
var tds	= dbPlugin.TypeDefns {
	{Name:"date", 		Html:"date", 		Sql:"DATE", 		Go:"string",	DftLen:0,},
	{Name:"datetime",	Html:"datetime",	Sql:"DATETIME",		Go:"string",	DftLen:0,},
	{Name:"email", 		Html:"email", 		Sql:"NVARCHAR", 	Go:"string",	DftLen:50,},
	{Name:"dec", 		Html:"number",		Sql:"DEC",			Go:"string",	DftLen:0,},
	{Name:"decimal", 	Html:"number",		Sql:"DEC",			Go:"string",	DftLen:0,},
	{Name:"int", 		Html:"number",		Sql:"INT",			Go:"int64",		DftLen:0,},
	{Name:"integer", 	Html:"number",		Sql:"INT",			Go:"int64",		DftLen:0,},
	{Name:"money", 		Html:"number",		Sql:"DEC",			Go:"string",	DftLen:0,},
	{Name:"number", 	Html:"number",		Sql:"INT",			Go:"int64",		DftLen:0,},
	{Name:"tel", 		Html:"tel",			Sql:"NVARCHAR",		Go:"string",	DftLen:19,},	//+nnn (nnn) nnn-nnnn
	{Name:"text", 		Html:"text",		Sql:"NVARCHAR",		Go:"string",	DftLen:0,},
	{Name:"time", 		Html:"time",		Sql:"TIME",			Go:"string",	DftLen:0,},
	{Name:"url", 		Html:"url",			Sql:"NVARCHAR",		Go:"string",	DftLen:50,},
}

func FlagsString(name string) string {
	var str			strings.Builder
	var wk			string

	wk = fmt.Sprintf("\tflag.StringVar(&db_name,\"dbName\",\"%s.db\",\"the database path\")\n", name)
	str.WriteString(wk)

	return str.String()
}

func ImportString() string {
	return "\"github.com/mattn/go-sqlite3\""
}

func init() {
	log.Printf("\tRegistering SQLite\n")
	pd :=  dbPlugin.PluginData{
		Name:"sqlite", Types:&tds, FlagsString:FlagsString, ImportString:ImportString,
		AddGo:false, CreateDB:false,
	}
	dbPlugin.Register(&pd)
}


