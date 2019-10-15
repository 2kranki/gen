// See License.txt in main repository directory

// https://github.com/lib/pq  			<== Driver being used
// https://godoc.org/github.com/lib/pq

// dbPostgres is the plugin for PostgreSQL and  contains the
// data and functions specific for PostgreSQL to generate
// table and field data for html forms, handlers and
// table sql i/o for a specific database.

// Remarks:
//	*	Unquoted names are converted to lowercase. Quoted names are case-sensitive
//		(ie "MyTable" is not equal to "MYTABLE") and should be surrounded with
//		double-quotes.

package dbPostgres

import (
	"fmt"
	"genapp/pkg/genSqlAppGo/dbJson"
	"genapp/pkg/genSqlAppGo/dbPlugin"
	"genapp/pkg/genSqlAppGo/dbType"
	"github.com/2kranki/go_util"
	"log"
)

const (
	extName = "postgres"
)

// Notes:
//	* We are now using a Decimal Package for support of decimal operations including
//		monetary calculations via https://github.com/ericlagergren/decimal
var tds = dbType.TypeDefns{
	{Name: "date", Html: "date", Sql: "DATE", Go: "string", DftLen: 0},
	{Name: "datetime", Html: "datetime", Sql: "DATETIME", Go: "string", DftLen: 0},
	{Name: "email", Html: "email", Sql: "VARCHAR", Go: "string", DftLen: 50},
	{Name: "dec", Html: "number", Sql: "DEC", Go: "float64", DftLen: 0},
	{Name: "decimal", Html: "number", Sql: "DEC", Go: "float64", DftLen: 0},
	{Name: "int", Html: "number", Sql: "INT", Go: "int64", DftLen: 0},
	{Name: "integer", Html: "number", Sql: "INT", Go: "int64", DftLen: 0},
	{Name: "money", Html: "number", Sql: "DEC", Go: "float64", DftLen: 0},
	{Name: "number", Html: "number", Sql: "INT", Go: "int64", DftLen: 0},
	{Name: "tel", Html: "tel", Sql: "VARCHAR", Go: "string", DftLen: 19}, //+nnn (nnn) nnn-nnnn
	{Name: "text", Html: "text", Sql: "VARCHAR", Go: "string", DftLen: 0},
	{Name: "time", Html: "time", Sql: "TIME", Go: "string", DftLen: 0},
	{Name: "url", Html: "url", Sql: "VARCHAR", Go: "string", DftLen: 50},
}

//----------------------------------------------------------------------------
//								Plugin Data and Methods
//----------------------------------------------------------------------------

// PluginData defines some of the data for the plugin.  Data within this package may also be
// used.  However, we use methods based off the PluginData to supply the data or other
// functionality.
type Plugin struct{}

// CreateDatabase indicatess if the Database needs to be
// created before it can be used.
func (pd Plugin) CreateDatabase() bool {
	return true
}

// DefaultDatabase returns default database name.
func (pd *Plugin) DefaultDatabase(db *dbJson.Database) string {
	return db.TitledName()
}

// DefaultPort returns default docker port.
func (pd *Plugin) DefaultPort() string {
	return "5432"
}

// DefaultPW returns default docker password.
func (pd *Plugin) DefaultPW() string {
	return "Passw0rd"
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

// DockerPort returns docker port used by the image.
func (pd Plugin) DockerPort() string {
	return "5432"
}

// DockerTag returns docker tag used to pull the image.
func (pd Plugin) DockerTag() string {
	return "11.3"
}

// DriverName returns the name to be used on pkg database sql.Open calls
func (pd *Plugin) DriverName() string {
	return "postgres"
}

func (pd *Plugin) GenDatabaseCreateStmt(db *dbJson.Database) string {
	var str util.StringBuilder

	str.WriteString("\tstr.WriteStringf(\"CREATE DATABASE %s;\", dbName)\n")

	return str.String()
}

// GenEnvArgDefns generates a check for an environment variable over-ride and
// over-rides the parsed CLI option if the environment variable is present.
func (pd Plugin) GenEnvArgDefns(appName string) string {
	var str util.StringBuilder

	str.WriteStringf("\twrk = os.Getenv(\"%s_DB_PW\")\n", appName)
	str.WriteString("\tif len(wrk)>0 {\n")
	str.WriteString("\t\tdb_pw = wrk\n")
	str.WriteString("\t}\n")
	str.WriteStringf("\twrk = os.Getenv(\"%s_DB_PORT\")\n", appName)
	str.WriteString("\tif len(wrk)>0 {\n")
	str.WriteString("\t\tdb_port = wrk\n")
	str.WriteString("\t}\n")
	str.WriteStringf("\twrk = os.Getenv(\"%s_DB_SERVER\")\n", appName)
	str.WriteString("\tif len(wrk)>0 {\n")
	str.WriteString("\t\tdb_srvr = wrk\n")
	str.WriteString("\t}\n")
	str.WriteStringf("\twrk = os.Getenv(\"%s_DB_USER\")\n", appName)
	str.WriteString("\tif len(wrk)>0 {\n")
	str.WriteString("\t\tdb_user = wrk\n")
	str.WriteString("\t}\n")
	str.WriteStringf("\twrk = os.Getenv(\"%s_DB_NAME\")\n", appName)
	str.WriteString("\tif len(wrk)>0 {\n")
	str.WriteString("\t\tdb_name = wrk\n")
	str.WriteString("\t}\n")
	return str.String()
}

// GenFlagArgDefns generates a string that defines the various CLI options to allow the
// user to modify the connection string parameters for the Database connection.
func (pd Plugin) GenFlagArgDefns(name string) string {
	var str util.StringBuilder

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
	var str util.StringBuilder

	return str.String()
}

// GenImportString returns the Database driver import string for this
// plugin.
func (pd Plugin) GenImportString() string {
	return "\"github.com/lib/pq\""
}

func (pd Plugin) GenRowPageStmt(t *dbJson.DbTable) string {
	var str util.StringBuilder

	db := t.DB

	str.WriteStringf("SELECT * FROM %s%s ORDER BY %s LIMIT $1 OFFSET $2;\\n",
		db.Schema, t.Name, t.KeysList("", " ASC"))

	return str.String()
}

// GenSqlBuildConn generates the code to build the connection string that would be
// issued to sql.Open() which is unique for each database server.
func (pd *Plugin) GenSqlBuildConn(dbServer, dbPort, dbUser, dbPW, dbName string) string {
	var str util.StringBuilder

	str.WriteString("\tconnStr := fmt.Sprintf(\"user=%s password='%s' host=%s port=%s \", ")
	str.WriteString(dbUser)
	str.WriteString(", ")
	str.WriteString(dbPW)
	str.WriteString(", ")
	str.WriteString(dbServer)
	str.WriteString(", ")
	str.WriteString(dbPort)
	str.WriteString(")\n")
	//str.WriteStringf("\\tif len(%s) > 0 {\n", dbName)
	//str.WriteString("\\t\\tconnStr += fmt.Sprintf(\"dbname='%%s')\\n\", %s", dbName)
	//str.WriteString("\t}\n")
	str.WriteString("\tconnStr += \"sslmode=disable\"\n")

	return str.String()
}

// GenTrailer returns any trailer information needed for I/O.
// This is included in both Database I/O and Table I/O.
func (pd *Plugin) GenTrailer() string {
	var str util.StringBuilder

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
	return "public."
}

// Types returns the TypeDefn table for this plugin to the caller as defined in dbPlugin.
// Required method
func (pd Plugin) Types() *dbType.TypeDefns {
	return &tds
}

//----------------------------------------------------------------------------
//                        Miscellaneous Support
//----------------------------------------------------------------------------

// GenDataPlaceHolder generates the string for table columns when a list of them
// is involved such as used in RowInsert().  Example: "$1, $2, $3"
func (pd Plugin) GenDataPlaceHolder(tb *dbJson.DbTable) string {
	var str util.StringBuilder
	var cnt int

	// Accumulate field name count.
	for _, f := range tb.Fields {
		if !f.Incr {
			cnt++
		}
	}

	for i := 0; i < cnt; i++ {
		cm := ", "
		if i == cnt-1 {
			cm = ""
		}
		//str.WriteStringf("?%s", cm)
		str.WriteStringf("$%d%s", i+1, cm)
	}

	return str.String()
}

// GenKeySearchPlaceHolder generates the string for multiple keys when an expression
// is involved such as used in RowFind(). The expression will always be '=' and will
// apply to all keys in the table. Example: "key1 = $1 AND key2 = $2"
func (pd Plugin) GenKeySearchPlaceHolder(tb *dbJson.DbTable, rel string) string {

	insertStr := ""
	keys, _ := tb.Keys()
	for i, _ := range keys {
		cm := " AND "
		if i == len(keys)-1 {
			cm = ""
		}
		insertStr += fmt.Sprintf("%s %s $%d%s", keys[i], rel, i+1, cm)
	}

	return insertStr
}

// GenKeysPlaceHolder generates the string for multiple keys when a list of key
// is involved such as used in RowFind().  Example: "?, ?, ?"
func (pd Plugin) GenKeysPlaceHolder(tb *dbJson.DbTable) string {

	insertStr := ""
	keys, _ := tb.Keys()
	for i := 0; i < len(keys); i++ {
		cm := ", "
		if i == len(tb.Fields)-1 {
			cm = ""
		}
		//insertStr += fmt.Sprintf("?%s", cm)
		insertStr += fmt.Sprintf("$%d%s", i+1, cm)
	}
	return insertStr
}

//----------------------------------------------------------------------------
//							Global Support Functions
//----------------------------------------------------------------------------

var plug *Plugin
var pluginData *dbPlugin.PluginData

func init() {
	log.Printf("\tRegistering Postgres\n")
	plug = &Plugin{}
	pluginData = &dbPlugin.PluginData{Name: extName, Types: &tds, Plugin: plug}
	dbPlugin.Register(extName, *pluginData)
}
