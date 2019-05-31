// See License.txt in main repository directory

// dbPostgres is the plugin for PostgreSQL and  contains the
// data and functions specific for PostgreSQL to generate
// table and field data for html forms, handlers and
// table sql i/o for a specific database.

package dbPostgres

import (
	"fmt"
	"strings"
	"../dbPlugin"
)

// Notes:
//	* We are now using a Decimal Package for support of decimal operations including
//		monetary calculations via https://github.com/ericlagergren/decimal
var tds	= dbPlugin.TypeDefns {
	{Name:"date", 		Html:"date", 		Sql:"DATE", 		Go:"string",	DftLen:0,},
	{Name:"datetime",	Html:"datetime",	Sql:"DATETIME",		Go:"string",	DftLen:0,},
	{Name:"email", 		Html:"email", 		Sql:"VARCHAR", 		Go:"string",	DftLen:50,},
	{Name:"dec", 		Html:"number",		Sql:"DEC",			Go:"string",	DftLen:0,},
	{Name:"decimal", 	Html:"number",		Sql:"DEC",			Go:"string",	DftLen:0,},
	{Name:"int", 		Html:"number",		Sql:"INT",			Go:"int",		DftLen:0,},
	{Name:"integer", 	Html:"number",		Sql:"INT",			Go:"int",		DftLen:0,},
	{Name:"money", 		Html:"number",		Sql:"DEC",			Go:"string",	DftLen:0,},
	{Name:"number", 	Html:"number",		Sql:"INT",			Go:"int",		DftLen:0,},
	{Name:"tel", 		Html:"tel",			Sql:"VARCHAR",		Go:"string",	DftLen:19,},	//+nnn (nnn) nnn-nnnn
	{Name:"text", 		Html:"text",		Sql:"VARCHAR",		Go:"string",	DftLen:0,},
	{Name:"time", 		Html:"time",		Sql:"TIME",			Go:"string",	DftLen:0,},
	{Name:"url", 		Html:"url",			Sql:"VARCHAR",		Go:"string",	DftLen:50,},
}

func FlagsString(name string) string {
	var str			strings.Builder
	var wk			string

	str.WriteString("\tflag.StringVar(&db_pw,\"dbPW\",\"Passw0rd!\",\"the database password\")\n")
	str.WriteString("\tflag.StringVar(&db_port,\"dbPort\",\"5430\",\"the database port\")\n")
	str.WriteString("\tflag.StringVar(&db_srvr,\"dbServer\",\"localhost\",\"the database server\")\n")
	str.WriteString("\tflag.StringVar(&db_user,\"dbUser\",\"postgres\",\"the database user\")\n")
	wk = fmt.Sprintf("\tflag.StringVar(&db_name,\"dbName\",\"%s\",\"the database name\")\n", name)
	str.WriteString(wk)
	return str.String()
}

func ImportString() string {
	return "\"github.com/lib/pq\""
}

func init() {
	pd :=  dbPlugin.PluginData{
		Name:"postgres", Types:&tds, FlagsString:FlagsString, ImportString:ImportString,
		AddGo:false, CreateDB:false,
	}
	dbPlugin.Register(&pd)
}

