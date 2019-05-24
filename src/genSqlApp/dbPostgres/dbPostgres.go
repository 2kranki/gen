// See License.txt in main repository directory

// dbPostgres is the plugin for PostgreSQL and  contains the
// data and functions specific for PostgreSQL to generate
// table and field data for html forms, handlers and
// table sql i/o for a specific database.

package dbPostgres

import (
	"../dbData"
)

var tds	= dbData.TypeDefns {
	{Name:"date", 		Html:"date", 		Sql:"DATE", 		Go:"string",	DftLen:0,},
	{Name:"datetime",	Html:"datetime",	Sql:"DATETIME",		Go:"string",	DftLen:0,},
	{Name:"email", 		Html:"email", 		Sql:"VARCHAR", 		Go:"string",	DftLen:50,},
	{Name:"dec", 		Html:"number",		Sql:"DEC",			Go:"float64",	DftLen:0,},
	{Name:"decimal", 	Html:"number",		Sql:"DEC",			Go:"float64",	DftLen:0,},
	{Name:"int", 		Html:"number",		Sql:"INT",			Go:"int",		DftLen:0,},
	{Name:"integer", 	Html:"number",		Sql:"INT",			Go:"int",		DftLen:0,},
	{Name:"money", 		Html:"number",		Sql:"DEC",			Go:"float64",	DftLen:0,},
	{Name:"number", 	Html:"number",		Sql:"INT",			Go:"int",		DftLen:0,},
	{Name:"tel", 		Html:"tel",			Sql:"VARCHAR",		Go:"string",	DftLen:19,},	//+nnn (nnn) nnn-nnnn
	{Name:"text", 		Html:"text",		Sql:"VARCHAR",		Go:"string",	DftLen:0,},
	{Name:"time", 		Html:"time",		Sql:"TIME",			Go:"string",	DftLen:0,},
	{Name:"url", 		Html:"url",			Sql:"VARCHAR",		Go:"string",	DftLen:50,},
}

func ImportString() string {
	return "\"github.com/lib/pq\""
}

func init() {
	pd :=  dbData.Plugin_Data{"postgres", &tds, ImportString, false, false}
	dbData.Register(&pd)
}

