// See License.txt in main repository directory

// Template Functions used in generation
// See the specific template files for how the functions
// and data are used.

package mainData

import (
	"genapp/pkg/sharedData"
	"testing"
)

func TestReadJsonFileMain(t *testing.T) {
	var err error

	sharedData.SetDebug(true)
	sharedData.SetMainPath("./test/main.json.txt")
	if err = ReadJsonFileMain(sharedData.MainPath()); err != nil {
		t.Errorf("ReadJsonFile() Reading Main JSON failed: %s'\n", sharedData.MainPath())
	}

	if len(mainStruct.Flags) != 1 {
		t.Errorf("ReadJsonFile() failed: should be 1 flags but is %d\n", len(mainStruct.Flags))
	}
	if mainStruct.Flags[0].Name != "exec" {
		t.Errorf("ReadJsonFile() failed: should be 5 flags but is %s\n", mainStruct.Flags[0].Name)
	}

	if len(mainStruct.Usage.Notes) != 3 {
		t.Errorf("ReadJsonFile() failed: len(Notes) should be 3 but is '%d'\n", len(mainStruct.Usage.Notes))
	}

	//t.Log(logData.String())
	t.Log("Successfully, completed: TestReadAppJson")

}
