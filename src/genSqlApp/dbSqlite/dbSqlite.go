// See License.txt in main repository directory

// dbSqlite is the plugin for SQLite and  contains the
// data and functions specific for SQLite to generate
// table and field data for html forms, handlers and
// table sql i/o for a specific database.

// SQLite is a fairly simple interface and easy to work
// with.  The database name becomes also the file path
// suffixed with ".db".

// Remarks:
//	*	The Create Database statement is irrelevant for
//		this server since the database name is used as
//		part of the SQLite file path.

package dbSqlite

import (
	"../dbJson"
	"../dbPlugin"
	"../dbType"
	"fmt"
	"log"
	"strings"

	"github.com/2kranki/go_util"
)

const(
	extName="sqlite"
)

// Notes:
//	* SQLite does We will be using a Decimal Package for support of decimal operations including
//		monetary calculations via https://github.com/ericlagergren/decimal
//    Currently, we just have decimal as text. SQLite does not handle decimal only
//		text, real64 and int64.
var tds	= dbType.TypeDefns {
	{Name:"date", 		Html:"date", 		Sql:"DATE", 		Go:"time.Time",	DftLen:0,},
	{Name:"datetime",	Html:"datetime",	Sql:"DATETIME",		Go:"time.Time",	DftLen:0,},
	{Name:"email", 		Html:"email", 		Sql:"VARCHAR", 		Go:"string",	DftLen:50,},
	{Name:"dec", 		Html:"number",		Sql:"TEXT",			Go:"string",	DftLen:0,},
	{Name:"decimal", 	Html:"number",		Sql:"TEXT",			Go:"string",	DftLen:0,},
	{Name:"int", 		Html:"number",		Sql:"INTEGER",		Go:"int64",		DftLen:0,},
	{Name:"integer", 	Html:"number",		Sql:"INTEGER",		Go:"int64",		DftLen:0,},
	{Name:"money", 		Html:"number",		Sql:"TEXT",			Go:"string",	DftLen:0,},
	{Name:"number", 	Html:"number",		Sql:"INT",			Go:"int64",		DftLen:0,},
	{Name:"tel", 		Html:"tel",			Sql:"VARCHAR",		Go:"string",	DftLen:19,},	//+nnn (nnn) nnn-nnnn
	{Name:"text", 		Html:"text",		Sql:"VARCHAR",		Go:"string",	DftLen:0,},
	{Name:"time", 		Html:"time",		Sql:"TIME",			Go:"time.Time",	DftLen:0,},
	{Name:"url", 		Html:"url",			Sql:"VARCHAR",		Go:"string",	DftLen:50,},
}

//----------------------------------------------------------------------------
//								Plugin Data and Methods
//----------------------------------------------------------------------------

// PluginData defines some of the data for the plugin.  Data within this package may also be
// used.  However, we use methods based off the PluginData to supply the data or other
// functionality.
type	Plugin struct {}

// CreateDatabase indicates if the Database needs to be
// created before it can be used.
func (pd Plugin) CreateDatabase() bool {
	return false
}

// DefaultDatabase returns default database name.
func (pd *Plugin) DefaultDatabase(db *dbJson.Database) string {
	return db.TitledName()+".db"
}

// DefaultPort returns default docker port.
func (pd *Plugin) DefaultPort() string {
	return ""
}

// DefaultPW returns default docker password.
func (pd *Plugin) DefaultPW() string {
	return ""
}

// DefaultServer returns default docker server name.
func (pd *Plugin) DefaultServer() string {
	return ""
}

// DefaultUser returns default docker user.
func (pd *Plugin) DefaultUser() string {
	return ""
}

// DriverName returns the name to be used on pkg database sql.Open calls
func (pd *Plugin) DriverName() string {
	return "sqlite3"
}

// GenFlagArgDefns generates a string that defines the various CLI options to allow the
// user to modify the connection string parameters for the Database connection.
func (pd Plugin) GenFlagArgDefns(name string) string {
	var str			strings.Builder
	var wk			string

	wk = fmt.Sprintf("\tflag.StringVar(&db_name,\"dbName\",\"%s.db\",\"the database path\")\n", name)
	str.WriteString(wk)
	return str.String()
}

// GenHeader returns any header information needed for I/O.
// This is included in both Database I/O and Table I/O.
func (pd *Plugin) GenHeader() string {
	var str			util.StringBuilder

	return str.String()
}

// GenImportString returns the Database driver import string for this
// plugin.
func (pd Plugin) GenImportString() string {
	return "\"github.com/mattn/go-sqlite3\""
}

// GenSqlBuildConn generates the code to build the connection string that would be
// issued to sql.Open() which is unique for each database server.
func (pd *Plugin) GenSqlBuildConn(dbServer,dbPort,dbUser,dbPW,dbName string) string {
	var strs		util.StringBuilder

	strs.WriteString("\tconnStr := fmt.Sprintf(\"%s\", ")
	strs.WriteString(dbName)
	strs.WriteString(")\n")

	return strs.String()
}

// GenTrailer returns any trailer information needed for I/O.
// This is included in both Database I/O and Table I/O.
func (pd *Plugin) GenTrailer() string {
	var str			util.StringBuilder

	return str.String()
}

// Name simply returns the external name that this plugin is known by
// or supports.
// Required method
func (pd Plugin) Name() string {
	return extName
}

// SchemaName simply returns the external name that this plugin is known by
// or supports.
// Required method
func (pd *Plugin) SchemaName() string {
	return ""
}

// Types returns the TypeDefn table for this plugin to the caller as defined in dbPlugin.
// Required method
func (pd Plugin) Types() *dbType.TypeDefns {
	return &tds
}

//----------------------------------------------------------------------------
//							Global Support Functions
//----------------------------------------------------------------------------

var plug		*Plugin
var pluginData	*dbPlugin.PluginData

func init() {
	log.Printf("\tRegistering SQLite\n")
	plug = &Plugin{}
	pluginData = &dbPlugin.PluginData{Name:extName, Types:&tds, Plugin:plug}
	dbPlugin.Register(extName, *pluginData)
}


