// See License.txt in main repository directory

// https://github.com/lib/pq  			<== Driver being used
// https://godoc.org/github.com/lib/pq

// dbPostgres is the plugin for PostgreSQL and  contains the
// data and functions specific for PostgreSQL to generate
// table and field data for html forms, handlers and
// table sql i/o for a specific database.

package dbPostgres

import (
	"../../shared"
	"../dbJson"
	"../dbPlugin"
	"../dbType"
	"fmt"
	"log"
	"strings"
)

const(
	extName="postgres"
)

// Notes:
//	* We are now using a Decimal Package for support of decimal operations including
//		monetary calculations via https://github.com/ericlagergren/decimal
var tds	= dbType.TypeDefns {
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

//----------------------------------------------------------------------------
//								Plugin Data and Methods
//----------------------------------------------------------------------------

// PluginData defines some of the data for the plugin.  Data within this package may also be
// used.  However, we use methods based off the PluginData to supply the data or other
// functionality.
type	Plugin struct {}

// CreateDatabase indicatess if the Database needs to be
// created before it can be used.
func (pd Plugin) CreateDatabase() bool {
	return false
}

// DefaultDatabase returns default database name.
func (pd *Plugin) DefaultDatabase(db *dbJson.Database) string {
	return db.TitledName()
}

// DefaultPort returns default docker port.
func (pd *Plugin) DefaultPort() string {
	return "5430"
}

// DefaultPW returns default docker password.
func (pd *Plugin) DefaultPW() string {
	return "Passw0rd!"
}

// DefaultServer returns default docker server name.
func (pd *Plugin) DefaultServer() string {
	return "localhost"
}

// DefaultUser returns default docker user.
func (pd *Plugin) DefaultUser() string {
	return "postgres"
}

// DockerName returns docker name used to pull the image.
func (pd Plugin) DockerName() string {
	return "postgres"
}

// DockerTag returns docker tag used to pull the image.
func (pd Plugin) DockerTag() string {
	return "11.3"
}

// DriverName returns the name to be used on pkg database sql.Open calls
func (pd *Plugin) DriverName() string {
	return "postgres"
}

// GenFlagArgDefns generates a string that defines the various CLI options to allow the
// user to modify the connection string parameters for the Database connection.
func (pd Plugin) GenFlagArgDefns(name string) string {
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

// GenImportString returns the Database driver import string for this
// plugin.
func (pd Plugin) GenImportString() string {
	return "\"github.com/lib/pq\""
}

// GenSqlOpen generates the code to issue sql.Open() which is unique
// for each database server.
func (pd Plugin) GenSqlOpen(dbSql,dbServer,dbPort,dbUser,dbPW,dbName string) string {
	var str			strings.Builder

	str.WriteString("\tconnStr := fmt.Sprintf(\"user=%s password='%s' host=%s port=%s \", ")
	str.WriteString(dbUser)
	str.WriteString(", ")
	str.WriteString(dbPW)
	str.WriteString(", ")
	str.WriteString(dbServer)
	str.WriteString(", ")
	str.WriteString(dbPort)
	str.WriteString(")\n")
	str.WriteString("\tif len(")
	str.WriteString(dbName)
	str.WriteString(") > 0 {\n")
	str.WriteString("\t\tconnStr += fmt.Sprintf(\"dbname='%s' \", ")
	str.WriteString(dbName)
	str.WriteString(")\n")
	str.WriteString("\t}\n")
	str.WriteString("\tconnStr += \"sslmode=disable\"\n")
	if sharedData.GenDebugging() {
		str.WriteString("\tlog.Printf(\"\\tConnecting to postgres using %s\\n\", connStr)\n")
	}
	str.WriteString("\t")
	str.WriteString(dbSql)
	str.WriteString(", err = sql.Open(\"postgres\", connStr)\n")

	return str.String()
}

// Name simply returns the external name that this plugin is known by
// or supports.
// Required method
func (pd Plugin) Name() string {
	return extName
}

// NeedUse indicates if the Database needs a USE
// SQL Statement before it can be used.
func (pd Plugin) NeedUse() bool {
	return false
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
	log.Printf("\tRegistering Postgres\n")
	plug = &Plugin{}
	pluginData = &dbPlugin.PluginData{Name:extName, Types:&tds, Plugin:plug}
	dbPlugin.Register(extName, *pluginData)
}

