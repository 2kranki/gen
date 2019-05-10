// See License.txt in main repository directory

// Template Functions used in generation
// See the specific template files for how the functions
// and data are used.

package mainData

import (
	"../shared"
	"fmt"
	"testing"
)

func TestReadJsonFileMain(t *testing.T) {
	var err			error
	var json		map[string]interface{}
	var ok			bool
	var wa		    []interface{}
	var wi		    interface{}
	var map1		map[string]interface{}

	sharedData.SetMainPath("./test/main.json.txt")
	if err = ReadJsonFileMain(sharedData.MainPath()); err != nil {
		t.Errorf("ReadJsonFile() Reading Main JSON failed: %s'\n", sharedData.MainPath())
	}

	wi = MainJson()
	t.Log(fmt.Sprintf("MainJson() Type: %T",MainJson()))
	t.Log(fmt.Sprintf("*MainJson() Type: %T", wi))
	if json, ok = wi.(map[string]interface{}); !ok {
		t.Errorf("ReadJsonFile() Main JSON type assertion failed: %s'\n", sharedData.MainPath())
	}

	if wi, ok = json["Flags"]; !ok {
		t.Errorf("ReadJsonFile() failed: Could not find Flags\n")
	}
	t.Log(fmt.Sprintf("Flags Type: %T", wi))
	if wa, ok = wi.([]interface{}); !ok {
		t.Errorf("ReadJsonFile() failed: Flags type assertion failed\n")
	}
	wi = wa[0]
	if map1, ok = wi.(map[string]interface{}); !ok {
		t.Errorf("ReadJsonFile() Flags[0] type assertion failed: %q\n", wi)
	}
	if _, ok = map1["Name"]; !ok {
		t.Errorf("ReadJsonFile() Flags[0].Name not found: %q\n", map1)
	}

	if len(mainStruct.Usage.Notes) != 3 {
		t.Errorf("ReadJsonFile() failed: len(Notes) should be 3 but is '%d'\n", len(mainStruct.Usage.Notes))
	}
	if len(mainStruct.Flags) != 7 {
		t.Errorf("ReadJsonFile() failed: should be 7 flags but is %d\n", len(mainStruct.Flags))
	}

	//t.Log(logData.String())
	t.Log("Successfully, completed: TestReadAppJson")

}
