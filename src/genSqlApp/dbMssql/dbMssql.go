// See License.txt in main repository directory

// dbMssql is the plugin for Microsoft SQL and  contains the
// data and functions specific for Microsoft SQL to generate
// table and field data for html forms, handlers and
// table sql i/o for a specific database.

// https://github.com/denisenkom/go-mssqldb  <== Driver being used
// https://godoc.org/github.com/denisenkom/go-mssqldb

// The sqlserver driver uses normal MS SQL Server syntax and expects
// parameters in the sql query to be in the form of either:
// 			@Name (using sql.Named)
//					or
// 			@p1 to @pN (ordinal position)


/***				T-SQL Examples

	// Create a database.
	USE master;
	GO
	IF DB_ID (N'mytest') IS NOT NULL
	DROP DATABASE mytest;
	GO
	CREATE DATABASE mytest;
	GO
	USE mytest;
	GO

 ***/

package dbMssql

import (
	"../../shared"
	"../../util"
	"../dbJson"
	"../dbPlugin"
	"../dbType"
	"fmt"
	"log"
	"strings"
)

// Notes:
//	* We are now using a Decimal Package for support of decimal operations including
//		monetary calculations via https://github.com/ericlagergren/decimal
var tds	= dbType.TypeDefns {
	{Name:"date", 		Html:"date", 		Sql:"DATE", 		Go:"time.Time",	DftLen:0,},
	{Name:"datetime",	Html:"datetime",	Sql:"DATETIME",		Go:"time.Time",	DftLen:0,},
	{Name:"email", 		Html:"email", 		Sql:"VARCHAR", 		Go:"string",	DftLen:50,},
	{Name:"dec", 		Html:"number",		Sql:"DEC",			Go:"float64",	DftLen:0,},
	{Name:"decimal", 	Html:"number",		Sql:"DEC",			Go:"float64",	DftLen:0,},
	{Name:"int", 		Html:"number",		Sql:"INT",			Go:"int64",		DftLen:0,},
	{Name:"integer", 	Html:"number",		Sql:"INT",			Go:"int64",		DftLen:0,},
	{Name:"money", 		Html:"number",		Sql:"DEC",			Go:"float64",	DftLen:0,},
	{Name:"number", 	Html:"number",		Sql:"INT",			Go:"int64",		DftLen:0,},
	{Name:"tel", 		Html:"tel",			Sql:"VARCHAR",		Go:"string",	DftLen:19,},	//+nnn (nnn) nnn-nnnn
	{Name:"text", 		Html:"text",		Sql:"NVARCHAR",		Go:"string",	DftLen:0,},
	{Name:"time", 		Html:"time",		Sql:"TIME",			Go:"time.Time",	DftLen:0,},
	{Name:"url", 		Html:"url",			Sql:"TEXT",			Go:"string",	DftLen:50,},
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
func (pd *Plugin) CreateDatabase() bool {
	return true
}

// DefaultDatabase returns default database name.
func (pd *Plugin) DefaultDatabase(db *dbJson.Database) string {
	return db.TitledName()
}

// DefaultPort returns default docker port.
func (pd *Plugin) DefaultPort() string {
	return "1401"
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
	return "sa"
}

// DockerName returns docker name used to pull the image.
func (pd *Plugin) DockerName() string {
	return "mcr.microsoft.com/mssql/server"
}

// DockerTag returns docker tag used to pull the image.
func (pd *Plugin) DockerTag() string {
	return "2017-latest"
}

// DriverName returns the name to be used on pkg database sql.Open calls
func (pd *Plugin) DriverName() string {
	return "mssql"
}

func (pd *Plugin) GenDatabaseCreateStmt(db *dbJson.Database) string {
	var str			util.StringBuilder

	str.WriteStringf("create database %s;\\n", db.TitledName())
	//str.WriteString( "go")

	return str.String()
}

func (pd *Plugin) GenDatabaseDeleteStmt(db *dbJson.Database) string {
	var str			util.StringBuilder

	str.WriteStringf("IF DB_ID (N'%s') IS NOT NULL\\n", db.TitledName())
	str.WriteStringf("DROP DATABASE %s ;\\n", db.TitledName())
	//str.WriteString("GO")

	return str.String()
}

func (pd *Plugin) GenDatabaseUseStmt(db *dbJson.Database) string {
	var str			util.StringBuilder

	str.WriteStringf("USE %s;\\n", db.TitledName())

	return str.String()
}

func (pd *Plugin) GenExecErrorCheck(db *dbJson.Database) string {
	var str			util.StringBuilder

	str.WriteString("if err != nil {\n")
	str.WriteString("\t\textra, ok := err.(ErrorWithExtraInfo)\n")
	str.WriteString("\t\tif ok {\n")
	str.WriteString("\t\t\tlineNo = int(extra.SQLErrorLineNo())\n")
	str.WriteString("\t\t}\n")
	str.WriteString("\t}\n")

	return str.String()
}

// GenFlagArgDefns generates a string that defines the various CLI options to allow the
// user to modify the connection string parameters for the Database connection.
func (pd *Plugin) GenFlagArgDefns(name string) string {
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

	str.WriteString("type ErrorWithExtraInfo interface {\n")
	str.WriteString("\tSQLErrorLineNo() int32\n")
	str.WriteString("\tSQLErrorNumber() int32\n")
	str.WriteString("}\n")

	return str.String()
}

// GenImportString returns the Database driver import string for this
// plugin.
func (pd *Plugin) GenImportString() string {
	return "\"github.com/denisenkom/go-mssqldb\""
}

// GenKeySearchPlaceHolder generates the string for multiple keys when an expression
// is involved such as used in RowFind(). The expression will always be '=' and will
// apply to all keys in the table. Example: "key1 = $1 AND key2 = $2"
func (pd *Plugin) GenKeySearchPlaceHolder(tb *dbJson.DbTable, rel string) string {

	insertStr := ""
	keys, _ := tb.Keys()
	for i, _ := range keys {
		cm := " AND "
		if i == len(keys) - 1 {
			cm = ""
		}
		insertStr += fmt.Sprintf("%s %s ?%s", keys[i], rel, cm)
	}

	return insertStr
}

func (pd *Plugin) GenRowFirstStmt(t *dbJson.DbTable) string {
	var str			util.StringBuilder

	db := t.DB

	// ORDER BY xx [OFFSET n ROWS [FETCH NEXT n ROWS ONLY]]
	str.WriteStringf("SELECT * FROM %s%s ORDER BY %s %s %s;\\n",
		db.Schema, t.TitledName(), t.KeysList("", " ASC"),
		pd.GenRowOffset(t, "0"), pd.GenRowLimit(t, "1"))

	return str.String()
}

func (pd *Plugin) GenRowLastStmt(t *dbJson.DbTable) string {
	var str			util.StringBuilder

	db := t.DB

	// ORDER BY xx [OFFSET n ROWS [FETCH NEXT n ROWS ONLY]]
	str.WriteStringf("SELECT * FROM %s%s ORDER BY %s %s;\\n",
		db.Schema, t.TitledName(), t.KeysList("", " DESC"),
		pd.GenRowOffset(t, "0"), pd.GenRowLimit(t, "1"))

	return str.String()
}

// GenRowLimit defines the interface for generating the LIMIT n option on
// SELECT.  LIMIT is used in general SQL, but not supported by T-SQL (Microsoft).
func (pd *Plugin) GenRowLimit(t *dbJson.DbTable, n string) string {
	var str			util.StringBuilder

	str.WriteStringf("FETCH NEXT %s ROWS ONLY", n)

	return str.String()
}

func (pd *Plugin) GenRowNextStmt(t *dbJson.DbTable) string {
	var str			util.StringBuilder

	db := t.DB

	// ORDER BY xx [OFFSET n ROWS [FETCH NEXT n ROWS ONLY]]
	str.WriteStringf("SELECT * FROM %s%s WHERE %s ORDER BY %s %s;\\n",
		db.Schema, t.TitledName(), pd.GenKeySearchPlaceHolder(t, ">"), t.KeysList("", " ASC"),
		pd.GenRowOffset(t, "0"), pd.GenRowLimit(t, "1"))

	return str.String()
}

// GenRowOffset defines the interface for generating the OFFSET n option on
// SELECT.  OFFSET has a slightly different grammar on T-SQL (Microsoft).
func (pd *Plugin) GenRowOffset(t *dbJson.DbTable, n string) string {
	var str			util.StringBuilder

	str.WriteStringf("OFFSET %s ROWS", n)

	return str.String()
}

func (pd *Plugin) GenRowPageStmt(t *dbJson.DbTable) string {
	var str			util.StringBuilder

	db := t.DB

	// ORDER BY xx [OFFSET n ROWS [FETCH NEXT n ROWS ONLY]]
	str.WriteStringf("SELECT * FROM %s%s ORDER BY %s %s %s;\\n",
		db.Schema, t.TitledName(), t.KeysList("", " ASC"),
		pd.GenRowOffset(t, "?"), pd.GenRowLimit(t, "?"))

	return str.String()
}

func (pd *Plugin) GenRowPrevStmt(t *dbJson.DbTable) string {
	var str			util.StringBuilder

	db := t.DB

	// ORDER BY xx [OFFSET n ROWS [FETCH NEXT n ROWS ONLY]]
	str.WriteStringf("SELECT * FROM %s%s WHERE %s ORDER BY %s %s;\\n",
		db.Schema, t.TitledName(), pd.GenKeySearchPlaceHolder(t, "<"), t.KeysList("", " DESC"),
		pd.GenRowOffset(t, "0"), pd.GenRowLimit(t, "1"))

	return str.String()
}

// GenSqlOpen generates the code to issue sql.Open() which is unique
// for each database server.
func (pd *Plugin) GenSqlOpen(dbSql,dbServer,dbPort,dbUser,dbPW,dbName string) string {
	var strs		util.StringBuilder

	/**************
	//strs.WriteString("connStr := fmt.Sprintf(\"sqlserver://%s:%s@%s:%s?database=master&connection+timeout=30\",")
	strs.WriteString("connStr := fmt.Sprintf(\"sqlserver://%s:%s@%s:%s?connection+timeout=30\",")
	strs.WriteString(dbUser)
	strs.WriteString(",")
	strs.WriteString(dbPW)
	strs.WriteString(",")
	strs.WriteString(dbServer)
	strs.WriteString(",")
	strs.WriteString(dbPort)
	strs.WriteString(")\n")
	 *************/

	strs.WriteString("\tquery := url.Values{}\n")
	strs.WriteString("\tquery.Add(\"connection+timeout\", \"30\")\n")
	strs.WriteString("\tu := &url.URL{\n")
	strs.WriteString("\t\tScheme:\t\t\"sqlserver\",\n")
	strs.WriteString("\t\tUser:\t\turl.UserPassword(")
	strs.WriteString(dbUser)
	strs.WriteString(", ")
	strs.WriteString(dbPW)
	strs.WriteString("),\n")
	strs.WriteString("\t\tHost:\t\tfmt.Sprintf(\"%s:%s\", ")
	strs.WriteString(dbServer)
	strs.WriteString(", ")
	strs.WriteString(dbPort)
	strs.WriteString("),\n")
	strs.WriteString("\t\tRawQuery:\tquery.Encode(),\n")
	strs.WriteString("\t}\n")
	strs.WriteString("\tconnStr := u.String()\n")

	if sharedData.GenDebugging() {
		strs.WriteStringf("\tlog.Printf(\"\\tConnecting to %s using %%s\\n\", connStr)\n", pd.DriverName())
	}
	strs.WriteStringf("\t%s, err = sql.Open(\"%s\", connStr)\n", dbSql, pd.DriverName())

	return strs.String()
}

func (pd *Plugin) GenTableCreateStmt(t *dbJson.DbTable) string {
	var str			strings.Builder
	var hasKeys		bool

	db := t.DB

	str.WriteString(fmt.Sprintf("CREATE TABLE %s%s (\\n", db.Schema, t.TitledName()))
	for i, _ := range t.Fields {
		var cm  		string
		var f			*dbJson.DbField
		var ft			string
		var nl			string
		var pk			string
		var sp			string

		f = &t.Fields[i]
		cm = ""
		if i != (len(t.Fields) - 1) {
			cm = ","
		} else {
			if hasKeys {
				cm = ","
			}
		}

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
		if f.KeyNum > 0 {
			hasKeys = true
		}
		sp = ""
		if len(f.SQLParms) > 0 {
			sp = " " + f.SQLParms
		}

		str.WriteString(fmt.Sprintf("\\t%s\\t%s%s%s%s%s\\n", f.Name, ft, nl, pk, cm, sp))
	}
	if hasKeys {
		wrk := fmt.Sprintf("\\tCONSTRAINT PK_%s PRIMARY KEY(%s)\\n", t.TitledName(), t.KeysList("", ""))
		str.WriteString(wrk)
	}
	str.WriteString(")")
	if len(t.SQLParms) > 0 {
		str.WriteString(",\\n")
		for _, l := range t.SQLParms {
			str.WriteString(fmt.Sprintf("%s\\n", l))
		}
	}
	str.WriteString(";\\n")

	return str.String()
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
func (pd *Plugin) Name() string {
	return "mssql"
}

// NeedUse indicates if the Database needs a USE
// SQL Statement before it can be used.
func (pd *Plugin) NeedUse() bool {
	return true
}

// SchemaName simply returns the external name that this plugin is known by
// or supports.
// Required method
func (pd *Plugin) SchemaName() string {
	return "dbo."
}

// Types returns the TypeDefn table for this plugin to the caller as defined in dbPlugin.
// Required method
func (pd *Plugin) Types() *dbType.TypeDefns {
	return &tds
}

//----------------------------------------------------------------------------
//							Global Support Functions
//----------------------------------------------------------------------------

var plug		*Plugin
var pluginData	*dbPlugin.PluginData

func init() {
	log.Printf("\tRegistering MS SQL\n")
	plug = &Plugin{}
	pluginData = &dbPlugin.PluginData{Name:"mssql", Types:&tds, Plugin:plug}
	dbPlugin.Register("mssql", *pluginData)
}

