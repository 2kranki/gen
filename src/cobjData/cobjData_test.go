// See License.txt in main repository directory

// Template Functions used in generation
// See the specific template files for how the functions
// and data are used.

package cobjData

import (
	"../shared"
	"fmt"
	"testing"
)

func TestReadJsonFileApp(t *testing.T) {
	var err			error
	var json		map[string]interface{}
	var ok			bool
	var wa		    []interface{}
	var wi		    interface{}
	var map1		map[string]interface{}

	sharedData.SetMainPath("./test/app.json.txt")
	if err = ReadJsonFileApp(sharedData.MainPath()); err != nil {
		t.Errorf("ReadJsonFileApp() Reading Main JSON failed: %s'\n", sharedData.MainPath())
	}

	wi = AppJson()
	t.Log(fmt.Sprintf("MainJson() Type: %T",AppJson()))
	t.Log(fmt.Sprintf("*MainJson() Type: %T", wi))
	if json, ok = wi.(map[string]interface{}); !ok {
		t.Errorf("ReadJsonFileApp() Main JSON type assertion failed: %s'\n", sharedData.MainPath())
	}

	if wi, ok = json["Tables"]; !ok {
		t.Errorf("ReadJsonFileApp() failed: Could not find Tables\n")
	}
	t.Log(fmt.Sprintf("Tables Type: %T", wi))
	if wa, ok = wi.([]interface{}); !ok {
		t.Errorf("ReadJsonFileApp() failed: Tables type assertion failed\n")
	}
	wi = wa[0]
	if map1, ok = wi.(map[string]interface{}); !ok {
		t.Errorf("ReadJsonFileApp() Tables[0] type assertion failed: %q\n", wi)
	}
	if _, ok = map1["Name"]; !ok {
		t.Errorf("ReadJsonFileApp() Tables[0].Name not found: %q\n", map1)
	}

	if len(appStruct.Tables) != 2 {
		t.Errorf("ReadJsonFileApp() failed: len(Tables) should be 2 but is '%d'\n", len(appStruct.Name))
	}
	if len(appStruct.Tables[0].Fields) != 8 {
		t.Errorf("ReadJsonFileApp() failed: should be 8 Tables but is %d\n", len(appStruct.Tables[0].Fields))
	}

	//t.Log(logData.String())
	t.Log("Successfully, completed: TestReadAppJson")

}
