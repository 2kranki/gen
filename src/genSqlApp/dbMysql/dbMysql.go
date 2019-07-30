// See License.txt in main repository directory

// dbMysql is the plugin for MySQL and  contains the
// data and functions specific for SQLite to generate
// table and field data for html forms, handlers and
// table sql i/o for a specific database.

package dbMysql

import (
	"../../shared"
	"../../util"
	"../dbJson"
	"../dbPlugin"
	"../dbType"
	"log"
)

const(
	extName="mysql"
)

// Notes:
//	* We are now using a Decimal Package for support of decimal operations including
//		monetary calculations via https://github.com/ericlagergren/decimal
var tds	= dbType.TypeDefns {
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
	return "3306"
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
	return "root"
}

// DockerName returns docker name used to pull the image.
func (pd Plugin) DockerName() string {
	return "mysql"
}

// DockerTag returns docker tag used to pull the image.
func (pd Plugin) DockerTag() string {
	return "5.7"
}

// DriverName returns the name to be used on pkg database sql.Open calls
func (pd *Plugin) DriverName() string {
	return "mysql"
}

func (pd *Plugin) GenDatabaseUseStmt(db *dbJson.Database) string {
	var str			util.StringBuilder

	str.WriteStringf("USE %s;", db.TitledName())

	return str.String()
}

// GenFlagArgDefns generates a string that defines the various CLI options to allow the
// user to modify the connection string parameters for the Database connection.
func (pd Plugin) GenFlagArgDefns(name string) string {
	var str			util.StringBuilder

	str.WriteStringf("\tflag.StringVar(&db_pw,\"dbPW\",\"%s\",\"the database password\")\n", pd.DefaultPW())
	str.WriteStringf("\tflag.StringVar(&db_port,\"dbPort\",\"%s\",\"the database port\")\n", pd.DefaultPort())
	str.WriteStringf("\tflag.StringVar(&db_srvr,\"dbServer\",\"%s\",\"the database server\")\n", pd.DefaultServer())
	str.WriteStringf("\tflag.StringVar(&db_user,\"dbUser\",\"%s\",\"the database user\")\n", pd.DefaultUser())
	str.WriteStringf("\tflag.StringVar(&db_name,\"dbName\",\"%s\",\"the database name\")\n", name)
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
	return "\"github.com/go-sql-driver/mysql\""
}

// GenSqlOpen generates the code to issue sql.Open() which is unique
// for each database server.
func (pd Plugin) GenSqlOpen(dbSql,dbServer,dbPort,dbUser,dbPW,dbName string) string {
	var strs		util.StringBuilder

	strs.WriteString("connStr := fmt.Sprintf(\"%s:%s@tcp(%s:%s)/\",")
	strs.WriteString(dbUser)
	strs.WriteString(",")
	strs.WriteString(dbPW)
	strs.WriteString(",")
	strs.WriteString(dbServer)
	strs.WriteString(",")
	strs.WriteString(dbPort)
	strs.WriteString(")\n")
	if sharedData.GenDebugging() {
		strs.WriteStringf("\tlog.Printf(\"\\tConnecting to %s using %%s\\n\", connStr)\n", pd.DriverName())
	}
	strs.WriteStringf("\t%s, err = sql.Open(\"%s\", connStr)\n", dbSql, pd.DriverName())

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

// NeedUse indicates if the Database needs a USE
// SQL Statement before it can be used.
func (pd Plugin) NeedUse() bool {
	return true
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
	log.Printf("\tRegistering MySQL\n")
	plug = &Plugin{}
	pluginData = &dbPlugin.PluginData{Name:extName, Types:&tds, Plugin:plug}
	dbPlugin.Register(extName, *pluginData)
}

