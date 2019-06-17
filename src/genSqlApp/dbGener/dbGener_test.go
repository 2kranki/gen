// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Test Generate dbGener Application Generator

// See dbJson_test.go for details on how to read the test file.

package dbGener

import (
	"../../shared"
	"../dbJson"
	"../dbPlugin"
	"fmt"
	"log"
	"testing"
	// Include the various Database Plugins so that they will register
	// with dbPlugin.
	_ "../dbMariadb"
	_ "../dbMssql"
	_ "../dbMysql"
	_ "../dbPostgres"
	_ "../dbSqlite"
)

const jsonTestPath = "../../../misc/test01/db.json.txt"

var jsonData	*dbJson.Database

//----------------------------------------------------------------------------
//								ReadJsonFile
//----------------------------------------------------------------------------

func ReadJsonFile(t *testing.T) {
	var err			error

	t.Logf("ReadJsonFile()..\n")
	sharedData.SetDebug(true)
	sharedData.SetMainPath(jsonTestPath)
	err = dbJson.ReadJsonFile(sharedData.MainPath())
	if err != nil {
		t.Fatalf("ReadJsonFile() Reading JSON tables failed: %s: %s\n", sharedData.MainPath(), err.Error())
	}
	jsonData = dbJson.DbStruct()
	if jsonData == nil {
		t.Fatalf("ReadJsonFile() Returned JSON tables are nil for: %s\n", sharedData.MainPath())
	}

	t.Logf("ReadJsonFile: end of test\n")

}


//----------------------------------------------------------------------------
//								GenFormDataDisplay
//----------------------------------------------------------------------------

func TestGenFormDataDisplay(t *testing.T) {
	var strs		[]string
	var dataTest	= []string{"<table>\n",
		"\t<tr><td><label>Num</label></td> <td><input type=\"number\" name=\"Num\" id=\"Num\" value=\"{{.Rcd.Num}}\"></td></tr>\n",
		"\t<tr><td><label>Name</label></td> <td><input type=\"text\" name=\"Name\" id=\"Name\" value=\"{{.Rcd.Name}}\"></td></tr>\n",
		"\t<tr><td><label>Addr1</label></td> <td><input type=\"text\" name=\"Addr1\" id=\"Addr1\" value=\"{{.Rcd.Addr1}}\"></td></tr>\n",
		"\t<tr><td><label>Addr2</label></td> <td><input type=\"text\" name=\"Addr2\" id=\"Addr2\" value=\"{{.Rcd.Addr2}}\"></td></tr>\n",
		"\t<tr><td><label>City</label></td> <td><input type=\"text\" name=\"City\" id=\"City\" value=\"{{.Rcd.City}}\"></td></tr>\n",
		"\t<tr><td><label>State</label></td> <td><input type=\"text\" name=\"State\" id=\"State\" value=\"{{.Rcd.State}}\"></td></tr>\n",
		"\t<tr><td><label>Zip</label></td> <td><input type=\"text\" name=\"Zip\" id=\"Zip\" value=\"{{.Rcd.Zip}}\"></td></tr>\n",
		"\t<tr><td><label>CurBal</label></td> <td><input type=\"number\" name=\"CurBal\" id=\"CurBal\" value=\"{{.Rcd.CurBal}}\"></td></tr>\n",
		"</table>\n",
		"<input type=\"hidden\" id=\"key0\" name=\"key0\"value=\"{{.Rcd.Num}}\">\n",
	}

	log.Printf("TestGenFormDataDisplay()..\n")
	sharedData.SetDebug(true)

	// Read the test JSON Tables
	ReadJsonFile(t)

	strs = GenFormDataDisplay(&jsonData.Tables[0])
	t.Log("===")
	t.Log(strs)
	t.Log("===")
	if len(strs) != len(dataTest) {
		t.Fatalf("TestGenFormDataDisplay() generated data did not match saved data - length\n")
	}
	for i, v := range strs {
		if v != dataTest[i] {
			t.Errorf(" gen: %s", v)
			t.Errorf("data: %s", dataTest[i])
			t.Fatalf("TestGenFormDataDisplay() generated data did not match saved data - line %d\n", i)
		}
	}

	t.Log("TestGenFormDataDisplay: end of test\n")

}

func TestGenFormDataKeyGet(t *testing.T) {
	var strs		[]string
	var dataTest	= []string{"\t\t\tkey0 = document.getElementById(\"key0\").value\n",}

	log.Printf("TestGenFormDataKeyGet()..\n")
	sharedData.SetDebug(true)

	// Read the test JSON Tables
	ReadJsonFile(t)

	strs = GenFormDataKeyGet(&jsonData.Tables[0])
	t.Log("===")
	t.Log(strs)
	t.Log("===")
	if len(strs) != len(dataTest) {
		t.Fatalf("TestGenFormDataKeyGet() generated data did not match saved data - length\n")
	}
	for i, v := range strs {
		if v != dataTest[i] {
			t.Errorf(" gen: %s", v)
			t.Errorf("data: %s", dataTest[i])
			t.Fatalf("TestGenFormDataKeyGet() generated data did not match saved data - line %d\n", i)
		}
	}

	t.Log("TestGenFormDataKeyGet: end of test\n")

}

func TestGenFormDataKeys(t *testing.T) {
	var str			string
	var dataTest	= "\"?\"+\"key0=\"+key0"

	log.Printf("TestGenFormDataKeys()..\n")
	sharedData.SetDebug(true)

	// Read the test JSON Tables
	ReadJsonFile(t)

	str = GenFormDataKeys(&jsonData.Tables[0])
	t.Log("===")
	t.Log(str)
	t.Log("===")
	if len(str) != len(dataTest) {
		t.Fatalf("TestGenFormDataKeys() generated data did not match saved data - length\n")
	}
	if str != dataTest {
		t.Errorf(" str: %s", str)
		t.Errorf("data: %s", dataTest)
		t.Fatalf("TestGenFormDataKeys() generated data did not match saved data\n")
	}

	t.Log("TestGenFormDataKeys: end of test\n")

}

func TestGenTableCreateStmt(t *testing.T) {
	var str			string
	var dataTest	= "\"CREATE TABLE IF NOT EXISTS Customer (\n\tNum\tINT NOT NULL,\n\tName\tNVARCHAR(30),\n\tAddr1\tNVARCHAR(30),\n\tAddr2\tNVARCHAR(30),\n\tCity\tNVARCHAR(20),\n\tState\tNVARCHAR(10),\n\tZip\tNVARCHAR(15),\n\tCurBal\tDEC(15,2)\n);\n\""
	var err		    error
	var plg			dbPlugin.PluginData

	log.Printf("TestGenTableCreateStmt()..\n")
	sharedData.SetDebug(true)

	// Read the test JSON Tables
	ReadJsonFile(t)
	jsonData.SqlType = "sqlite"
	if plg, err = dbPlugin.FindPlugin(jsonData.SqlType); err != nil {
		t.Fatalf(fmt.Sprintf("Error: Can't find plugin for %s!\n\n\n", jsonData.SqlType))
	}
	if err = jsonData.SetupPlugin(plg); err != nil {
		t.Fatalf(fmt.Sprintf("Error: Plugin setup failure for %s!\n\n\n", jsonData.SqlType))
	}

	str = GenTableCreateStmt(&jsonData.Tables[0])
	t.Log("===")
	t.Log(str)
	t.Log("===")
	if len(str) != len(dataTest) {
		t.Fatalf("TestGenTableCreateStmt() generated data did not match saved data - length\n")
	}
	if str != dataTest {
		t.Errorf(" str: %s", str)
		t.Errorf("data: %s", dataTest)
		t.Fatalf("TestGenTableCreateStmt() generated data did not match saved data\n")
	}

	t.Log("TestGenTableCreateStmt: end of test\n")

}

