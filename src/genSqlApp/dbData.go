// See License.txt in main repository directory

// Template Functions used in generation
// See the specific template files for how the functions
// and data are used.

package genSqlApp

type DbField struct {
	Name		string		`json:"Name,omitempty"`
	Type		string		`json:"Type,omitempty"`
	Len		    int		    `json:"Len,omitempty"`
	PrimaryKey  bool	    `json:"PrimaryKey,omitempty"`
	List		bool	    `json:"List,omitempty"`			// Include in List Report
}

type DbTable struct {
	Name		string		`json:"Name,omitempty"`
	Fields		[]DbField
}

func (t *DbTable) CreateInsertStr() string {

	insertStr := ""
	for _, v := range t.Fields {
		insertStr += v.Name + ","
	}
	if len(insertStr) > 0 {
		insertStr = insertStr[0:len(insertStr)-1]
	}
	return insertStr
}

type Database struct {
	Name	string			`json:"Name,omitempty"`
	SqlType	string			`json:"SqlType,omitempty"`
	Tables  []DbTable
}

func CreateInsertSql(t interface{}) string {
	//var Fields 	[]map[string] interface{}
	//var ok		bool
	var x		string

	//insertStr := ""
	return x
}

func CreateTableStruct(t interface{}) string {
	//var Fields 	[]map[string] interface{}
	//var ok		bool
	var x		string

	//insertStr := ""
	return x
}

func dblClose() string {
	return "}}"
}

func dblOpen() string {
	return "{{"
}

// decodeDbField decodes a DbField from a generic
func decodeDbField(i interface{}) *DbField {
	var f		DbField
	var m		map[string]interface{}
	var ok		bool

	if m, ok = i.(map[string]interface{}); !ok {
		return nil
	}
	if f.Name, ok = m["Name"].(string); !ok {
		return nil
	}
	if f.Type, ok = m["Type"].(string); !ok {
		return nil
	}
	if f.Len, ok = m["Len"].(int); !ok {
		return nil
	}
	if f.PrimaryKey, ok = m["PrimaryKey"].(bool); !ok {
		return nil
	}
	if f.List, ok = m["List"].(bool); !ok {
		return nil
	}
	return &f
}

func decodeDbTable(i interface{}) *DbTable {
	var t		DbTable
	var m		map[string]interface{}
	//var fields	[]interface{}
	var fields	[]DbField
	var ok		bool

	if m, ok = i.(map[string]interface{}); !ok {
		return nil
	}
	if t.Name, ok = m["Name"].(string); !ok {
		return nil
	}
	if fields, ok = m["Fields"].([]DbField); !ok {
		return nil
	}
	for _, v := range fields {
		t.Fields = append(t.Fields, *decodeDbField(v))
	}
	return &t
}

func decodeDatabase(i interface{}) *Database {
	var t		Database
	var m		map[string]interface{}
	//var fields	[]interface{}
	var tables	[]DbTable
	var ok		bool

	if m, ok = i.(map[string]interface{}); !ok {
		return nil
	}
	if t.Name, ok = m["Name"].(string); !ok {
		return nil
	}
	if t.SqlType, ok = m["SqlType"].(string); !ok {
		return nil
	}
	if tables, ok = m["Tables"].([]DbTable); !ok {
		return nil
	}
	for _, v := range tables {
		t.Tables = append(t.Tables, *decodeDbTable(v))
	}
	return &t
}


