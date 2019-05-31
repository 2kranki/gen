// See License.txt in main repository directory

// dbSql provides the functions to generate the go statements
// necessary to access and manipulate the SQL databases defined
// by the user. The problem that it tries to solve is that while
// SQL is supposed to be a universal language. It unfortunately
// is not and each type of database manager must be handled slightly
// differently.

// We give this package access to user defined JSON and the ap-
// propriate plugin for the data being processed. Between those
// two resources, it must generate the go code.

package dbSql

import (
	"../dbJson"
	"fmt"
	"io"
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
	func GenDatabaseCreateStmts() []string
}

type GenDatabaseDeleteStmter interface {
	func GenDatabaseDeleteStmts() []string
}

//----------------------------------------------------------------------------
//                        	Table Interface Support
//----------------------------------------------------------------------------

type GenTableCreateStmter interface {
	func GenTableCreateStmts() []string
}

type GenTableDeleteStmter interface {
	func GenTableDeleteStmts() []string
}

//----------------------------------------------------------------------------
//                        	Row Interface Support
//----------------------------------------------------------------------------

type GenRowDeleteStmter interface {
	func GenRowfDeleteStmts() []string
}

type GenRowFindStmter interface {
	func GenRowFindStmts() []string
}

type GenRowInsertStmter interface {
	func GenRowInsertStmts() []string
}


//============================================================================
//                        Type Definition Support
//============================================================================

// SqlWork, The type definition struct, defines one acceptable type accepted in the JSON
// defining the Database Structure.  There must be a TypeDefn for each type
// accepted in each plugin.

type SqlWork struct {
	Name		string				`json:"Name,omitempty"`		// Type Name
	j			dbJson.Database
	w			io.Writer
}

// SqlWorks provides a convenient way of defining a SQL Work Table.
type SqlWorks	[]SqlWork

//----------------------------------------------------------------------------
//						Global/Internal Object Functions
//----------------------------------------------------------------------------

func (t SqlWorks) FindDefn(name string) *SqlWork {
	for i, v := range t {
		if name == v.Name {
			return &t[i]
		}
	}
	return nil
}

func (s *SqlWork) GenDatabaseCreateSqlStmts() []string {
	var str			strings.Builder

	str.WriteString(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;\\n", s.j.TitledName))
	if s.j.SqlType == "mssql" {
		str.WriteString("GO\\n")
	}
	str.WriteString(fmt.Sprintf("USE %s;\\n", s.j.TitledName))
	if s.j.SqlType == "mssql" {
		str.WriteString("GO\\n")
	}

	return str.String()
}

func (s *SqlWork) GenDeleteDatabasSql() string {
	var str			strings.Builder
	var intrfc		interface{}

	intffc = s.j.Plugin
	if obj, ok := intffc.(interface{GenDeleteDatabaseSql() string}); ok {
		str.WriteString(obj.GenDeleteDatabaseSql())
	}
	str.WriteString(fmt.Sprintf("DROP DATABASE IF EXISTS %s;\\n", s.j.TitledName))
	if s.j.SqlType == "mssql" {
		str.WriteString("GO\\n")
	}
	return str.String()
}

func (t SqlWork) Write(p []byte) (int, error) {
	var n		int
	var err		error

	// Write p to somewhere.

	return n, err
}

func (t SqlWork) WriteString(s string) (int, error) {
	var n		int
	var err		error
	var bs  	[]byte

	// Under the covers, a string is a read-only slice of bytes.
	// So, this conversion should not create much overhead, hopefully.
	bs = []byte(s)
	n, err = t.Write(bs)

	return n, err
}

//----------------------------------------------------------------------------
//						Global Support Functions
//----------------------------------------------------------------------------

// New provides a factory method to create an Sql Object.
func New(j *dbJson.Database, w io.Writer) (*SqlWork) {
	sw := &SqlWork{}
	sw.j = j
	sw.w = w
	return sw
}

