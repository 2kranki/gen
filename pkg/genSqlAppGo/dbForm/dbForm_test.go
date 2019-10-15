// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Test Generate SQL Application Generator

package dbForm

import (
	"genapp/pkg/sharedData"
	"log"
	"testing"
)

func TestNewFormWork(t *testing.T) {
	var fw *FormWork

	log.Printf("dbForm::TestNewFormWork()..\n")
	sharedData.SetDebug(true)

	if fw = NewFormWork(nil); fw == nil {
		t.Errorf("Error: Could not allocate a FormWork object!\n")
	}

	t.Log("...end of dbForm::TestNewFormWork\n")

}
