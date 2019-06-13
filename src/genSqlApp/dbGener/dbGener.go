// See License.txt in main repository directory

// dbGener provides the functions to generate the go statements
// necessary to access and manipulate the SQL databases defined
// by the user. The problem that it tries to solve is that while
// SQL is supposed to be a universal language. It unfortunately
// is not and each type of database manager must be handled slightly
// differently.

// We give this package access to user defined JSON and the ap-
// propriate plugin for the data being processed. Between those
// two resources, it must generate the go code.

package dbGener

import (
	"../../shared"
	"../dbJson"
	"fmt"
	"strings"
)

//============================================================================
//                        	Interface Support
//============================================================================

// dbSql uses interfaces to determine what a plugin can do or not do and when it
// should be called.  If the plugin does not support a particular interface, then
// dbSql will perform default logic to handle the situation.
//
// The reason for all this is that even though Go uses a "common" interface for
// accessing SQL Servers. The SQL, itself, can vary.  Although SQL is supposed to
// to be a standard, it is not consistently implemented unforturnately.
//
// Functions that return a full SQL statement must return a slice of strings even
// if there is only one statement ever generated.  That is because some servers
// such as Microsoft's SQL Server may not do anything until an additional statement
// is issued such as "go".

//----------------------------------------------------------------------------
//                        	Database Interface Support
//----------------------------------------------------------------------------

type GenDatabaseCreateStmter interface {
	GenDatabaseCreateStmt(db *dbJson.Database) []string
}

type GenDatabaseDeleteStmter interface {
	GenDatabaseDeleteStmt(db *dbJson.Database) []string
}

//----------------------------------------------------------------------------
//                        	Table Interface Support
//----------------------------------------------------------------------------

type GenTableCreateStmter interface {
	GenTableCreateStmt(tb *dbJson.DbTable) []string
}

type GenTableDeleteStmter interface {
	GenTableDeleteStmt(tb *dbJson.DbTable) []string
}

//----------------------------------------------------------------------------
//                        	Row Interface Support
//----------------------------------------------------------------------------

type GenRowDeleteStmter interface {
	GenRowDeleteStmt(tb *dbJson.DbTable) []string
}

type GenRowFindStmter interface {
	GenRowFindStmt(tb *dbJson.DbTable) []string
}

type GenRowFirstStmter interface {
	GenRowFirstStmt(tb *dbJson.DbTable) []string
}

type GenRowInsertStmter interface {
	GenRowInsertStmt(tb *dbJson.DbTable) []string
}

type GenRowLastStmter interface {
	GenRowLastStmt(tb *dbJson.DbTable) []string
}

type GenRowNextStmter interface {
	GenRowNextStmt(tb *dbJson.DbTable) []string
}

type GenRowPageStmter interface {
	GenRowPageStmt(tb *dbJson.DbTable) []string
}

type GenRowPrevStmter interface {
	GenRowPrevStmt(tb *dbJson.DbTable) []string
}

type GenRowUpdateStmter interface {
	GenRowUpdateStmt(tb *dbJson.DbTable) []string
}


//----------------------------------------------------------------------------
//						Global Database Support Functions
//----------------------------------------------------------------------------

func GenDatabaseCreateStmt(db *dbJson.Database) []string {
	var str			strings.Builder
	var strs  		[]string
	var intr		GenDatabaseCreateStmter
	var ok			bool

	intr, ok = db.Plugin.(GenDatabaseCreateStmter)
	if ok {
		return intr.GenDatabaseCreateStmt(db)
	}

	strs = append(strs, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;\\n", db.TitledName()))
	if db.SqlType == "mssql" {
		strs = append(strs, "GO\\n")
	}
	strs = append(strs, fmt.Sprintf("USE %s;\\n", db.TitledName()))
	if db.SqlType == "mssql" {
		strs = append(strs, "GO\\n")
	}

	strs = append(strs, str.String())
	return strs
}

func GenDatabaseDeleteStmt(db *dbJson.Database) []string {
	var strs  		[]string
	var intr		GenDatabaseDeleteStmter
	var ok			bool

	intr, ok = db.Plugin.(GenDatabaseDeleteStmter)
	if ok {
		return intr.GenDatabaseDeleteStmt(db)
	}

	strs = append(strs, fmt.Sprintf("DELETE DATABASE IF EXISTS %s;\\n", db.TitledName))
	if db.SqlType == "mssql" {
		strs = append(strs, "GO\\n")
	}

	return strs
}

//----------------------------------------------------------------------------
//						Global Table Support Functions
//----------------------------------------------------------------------------

func GenTableCreateStmt(t *dbJson.DbTable) []string {
	var str			strings.Builder
	var strs  		[]string
	var intr		GenTableCreateStmter
	var ok			bool

	db := t.DB
	intr, ok = db.Plugin.(GenTableCreateStmter)
	if ok {
		return intr.GenTableCreateStmt(t)
	}

	strs = append(strs, fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s;\\n", t.TitledName()))
	if db.SqlType == "mssql" {
		strs = append(strs, "GO\\n")
	}
	strs = append(strs, fmt.Sprintf("USE %s;\\n", t.TitledName()))
	if db.SqlType == "mssql" {
		strs = append(strs, "GO\\n")
	}

	strs = append(strs, str.String())
	return strs
}

func GenTableDeleteStmt(t *dbJson.DbTable) []string {
	var str			strings.Builder
	var strs  		[]string
	var intr		GenTableDeleteStmter
	var ok			bool

	db := t.DB
	intr, ok = db.Plugin.(GenTableDeleteStmter)
	if ok {
		return intr.GenTableDeleteStmt(t)
	}

	str.WriteString(fmt.Sprintf("DROP TABLE IF EXISTS %s;\\n", t.TitledName()))
	if db.SqlType == "mssql" {
		str.WriteString("GO\\n")
	}

	strs = append(strs, str.String())
	return strs
}

//----------------------------------------------------------------------------
//						Global Row Support Functions
//----------------------------------------------------------------------------

func GenRowDeleteStmt(t *dbJson.DbTable) []string {
	var str			strings.Builder
	var strs  		[]string
	var intr		GenRowDeleteStmter
	var ok			bool

	db := t.DB
	intr, ok = db.Plugin.(GenRowDeleteStmter)
	if ok {
		return intr.GenRowDeleteStmt(t)
	}

	str.WriteString(fmt.Sprintf("DROP DATABASE IF EXISTS %s;\\n", db.TitledName()))
	if db.SqlType == "mssql" {
		str.WriteString("GO\\n")
	}

	strs = append(strs, str.String())
	return strs
}

func GenRowFindStmt(t *dbJson.DbTable) []string {
	var str			strings.Builder
	var strs  		[]string
	var intr		GenRowFindStmter
	var ok			bool

	db := t.DB
	intr, ok = db.Plugin.(GenRowFindStmter)
	if ok {
		return intr.GenRowFindStmt(t)
	}

	str.WriteString(fmt.Sprintf("DROP DATABASE IF EXISTS %s;\\n", db.TitledName()))
	if db.SqlType == "mssql" {
		str.WriteString("GO\\n")
	}

	strs = append(strs, str.String())
	return strs
}

func GenRowFirstStmt(t *dbJson.DbTable) []string {
	var str			strings.Builder
	var strs  		[]string
	var intr		GenRowFirstStmter
	var ok			bool

	db := t.DB
	intr, ok = db.Plugin.(GenRowFirstStmter)
	if ok {
		return intr.GenRowFirstStmt(t)
	}

	str.WriteString(fmt.Sprintf("DROP DATABASE IF EXISTS %s;\\n", db.TitledName()))
	if db.SqlType == "mssql" {
		str.WriteString("GO\\n")
	}

	strs = append(strs, str.String())
	return strs
}

func GenRowInsertStmt(t *dbJson.DbTable) []string {
	var str			strings.Builder
	var strs  		[]string
	var intr		GenRowInsertStmter
	var ok			bool

	db := t.DB
	intr, ok = db.Plugin.(GenRowInsertStmter)
	if ok {
		return intr.GenRowInsertStmt(t)
	}

	str.WriteString(fmt.Sprintf("DROP DATABASE IF EXISTS %s;\\n", db.TitledName()))
	if db.SqlType == "mssql" {
		str.WriteString("GO\\n")
	}

	strs = append(strs, str.String())
	return strs
}

func GenRowLastStmt(t *dbJson.DbTable) []string {
	var str			strings.Builder
	var strs  		[]string
	var intr		GenRowLastStmter
	var ok			bool

	db := t.DB
	intr, ok = db.Plugin.(GenRowLastStmter)
	if ok {
		return intr.GenRowLastStmt(t)
	}

	str.WriteString(fmt.Sprintf("DROP DATABASE IF EXISTS %s;\\n", db.TitledName()))
	if db.SqlType == "mssql" {
		str.WriteString("GO\\n")
	}

	strs = append(strs, str.String())
	return strs
}

func GenRowNextStmt(t *dbJson.DbTable) []string {
	var str			strings.Builder
	var strs  		[]string
	var intr		GenRowNextStmter
	var ok			bool

	db := t.DB
	intr, ok = db.Plugin.(GenRowNextStmter)
	if ok {
		return intr.GenRowNextStmt(t)
	}

	str.WriteString(fmt.Sprintf("DROP DATABASE IF EXISTS %s;\\n", db.TitledName()))
	if db.SqlType == "mssql" {
		str.WriteString("GO\\n")
	}

	strs = append(strs, str.String())
	return strs
}

func GenRowPageStmt(t *dbJson.DbTable) []string {
	var str			strings.Builder
	var strs  		[]string
	var intr		GenRowPageStmter
	var ok			bool

	db := t.DB
	intr, ok = db.Plugin.(GenRowPageStmter)
	if ok {
		return intr.GenRowPageStmt(t)
	}

	str.WriteString(fmt.Sprintf("DROP DATABASE IF EXISTS %s;\\n", db.TitledName()))
	if db.SqlType == "mssql" {
		str.WriteString("GO\\n")
	}

	strs = append(strs, str.String())
	return strs
}

func GenRowPrevStmt(t *dbJson.DbTable) []string {
	var str			strings.Builder
	var strs  		[]string
	var intr		GenRowPrevStmter
	var ok			bool

	db := t.DB
	intr, ok = db.Plugin.(GenRowPrevStmter)
	if ok {
		return intr.GenRowPrevStmt(t)
	}

	str.WriteString(fmt.Sprintf("DROP DATABASE IF EXISTS %s;\\n", db.TitledName()))
	if db.SqlType == "mssql" {
		str.WriteString("GO\\n")
	}

	strs = append(strs, str.String())
	return strs
}

func GenRowUpdateStmt(t *dbJson.DbTable) []string {
	var str			strings.Builder
	var strs  		[]string
	var intr		GenRowUpdateStmter
	var ok			bool

	db := t.DB
	intr, ok = db.Plugin.(GenRowUpdateStmter)
	if ok {
		return intr.GenRowUpdateStmt(t)
	}

	str.WriteString(fmt.Sprintf("DROP DATABASE IF EXISTS %s;\\n", db.TitledName()))
	if db.SqlType == "mssql" {
		str.WriteString("GO\\n")
	}

	strs = append(strs, str.String())
	return strs
}

//----------------------------------------------------------------------------
//							Global Support Functions
//----------------------------------------------------------------------------

// init() is called before main(). Here we define the functions that will be
// used in the templates.
func init() {
	sharedData.SetFunc("GenDatabaseCreateStmt", GenDatabaseCreateStmt)
	sharedData.SetFunc("GenDatabaseDeleteStmt", GenDatabaseDeleteStmt)
	sharedData.SetFunc("GenTableCreateStmt", GenTableCreateStmt)
	sharedData.SetFunc("GenTableDeleteStmt", GenTableDeleteStmt)
	sharedData.SetFunc("GenRowDeleteStmt", GenRowDeleteStmt)
	sharedData.SetFunc("GenRowFindStmt", GenRowFindStmt)
	sharedData.SetFunc("GenRowFirstStmt", GenRowFirstStmt)
	sharedData.SetFunc("GenRowInsertStmt", GenRowInsertStmt)
	sharedData.SetFunc("GenRowLastStmt", GenRowLastStmt)
	sharedData.SetFunc("GenRowNextStmt", GenRowNextStmt)
	sharedData.SetFunc("GenRowPageStmt", GenRowPageStmt)
	sharedData.SetFunc("GenRowPrevStmt", GenRowPrevStmt)
	sharedData.SetFunc("GenRowUpdateStmt", GenRowUpdateStmt)
}

