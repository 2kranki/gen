// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Test Generate SQL Application Generator

package dbType

import (
	"testing"
	//"time"
	"genapp/pkg/sharedData"
)

var tds = TypeDefns{
	{Name: "date", Html: "date", Sql: "DATE", Go: "string", DftLen: 0},
	{Name: "datetime", Html: "datetime", Sql: "DATETIME", Go: "string", DftLen: 0},
	{Name: "email", Html: "email", Sql: "NVARCHAR", Go: "string", DftLen: 50},
	{Name: "dec", Html: "number", Sql: "DEC", Go: "string", DftLen: 0},
	{Name: "decimal", Html: "number", Sql: "DEC", Go: "string", DftLen: 0},
	{Name: "int", Html: "number", Sql: "INT", Go: "int64", DftLen: 0},
	{Name: "integer", Html: "number", Sql: "INT", Go: "int64", DftLen: 0},
	{Name: "money", Html: "number", Sql: "DEC", Go: "string", DftLen: 0},
	{Name: "number", Html: "number", Sql: "INT", Go: "int64", DftLen: 0},
	{Name: "tel", Html: "tel", Sql: "NVARCHAR", Go: "string", DftLen: 19}, //+nnn (nnn) nnn-nnnn
	{Name: "text", Html: "text", Sql: "NVARCHAR", Go: "string", DftLen: 0},
	{Name: "time", Html: "time", Sql: "TIME", Go: "string", DftLen: 0},
	{Name: "url", Html: "url", Sql: "NVARCHAR", Go: "string", DftLen: 50},
}

func TestFind(t *testing.T) {
	//var err			error
	//var str			string
	var typ *TypeDefn

	t.Logf("dbType::TestFind()..\n")
	sharedData.SetDebug(true)

	typ = tds.FindDefn("date")
	if typ == nil {
		t.Errorf("Error: Could not find 'date' entry\n")
	}

	typ = tds.FindDefn("url")
	if typ == nil {
		t.Errorf("Error: Could not find 'url' entry\n")
	}

	typ = tds.FindDefn("int")
	if typ == nil {
		t.Errorf("Error: Could not find 'int' entry\n")
	}

	t.Logf("dbType::TestFind: end of test\n")

}
