// See License.txt in main repository directory

// Shared contains the shared variables and data used by
// main and the other packages.

package sharedData

import	"testing"

func TestDataPath(t *testing.T) {
	SetDataPath("xyzzy")
	if DataPath() != "xyzzy" {
		t.Errorf("TestDataPath() failed: should be 'xyzzy' but is %s\n", DataPath())
	}
	SetDataPath("abc")
	if DataPath() != "abc"{
		t.Errorf("TestDataPath() failed: should be 'abc' but is %s\n", DataPath())
	}
}

func TestDebug(t *testing.T) {
	SetDebug(true)
	if !Debug() {
		t.Errorf("TestDebug() failed: should be true but is false\n")
	}
	SetDebug(false)
	if Debug() {
		t.Errorf("TestDebug() failed: should be false but is true\n")
	}
}

func TestDefn(t *testing.T) {
	if Defn("mdlDir") != "./models" {
		t.Errorf("TestMainPath() failed: should be './models' but is %s\n", Defn("mdlDir"))
	}
	SetDefn("mdlDir","xyzzy")
	if Defn("mdlDir") != "xyzzy" {
		t.Errorf("TestMainPath() failed: should be 'xyzzy' but is %s\n", Defn("mdlDir"))
	}
	SetDefn("mdlDir","abc")
	if Defn("mdlDir") != "abc"{
		t.Errorf("TestMainPath() failed: should be 'abc' but is %s\n", Defn("mdlDir"))
	}
}

func TestForce(t *testing.T) {
	SetForce(true)
	if !Force() {
		t.Errorf("TestForce() failed: should be true but is false\n")
	}
	SetForce(false)
	if Force() {
		t.Errorf("TestForce() failed: should be false but is true\n")
	}
}

func TestMainPath(t *testing.T) {
	SetMainPath("xyzzy")
	if MainPath() != "xyzzy" {
		t.Errorf("TestMainPath() failed: should be 'xyzzy' but is %s\n", MainPath())
	}
	SetMainPath("abc")
	if MainPath() != "abc"{
		t.Errorf("TestMainPath() failed: should be 'abc' but is %s\n", MainPath())
	}
}

func TestMdlDir(t *testing.T) {
	SetMdlDir("xyzzy")
	if MdlDir() != "xyzzy" {
		t.Errorf("TestMdlDir() failed: should be 'xyzzy' but is %s\n", MdlDir())
	}
	SetMdlDir("abc")
	if MdlDir() != "abc"{
		t.Errorf("TestMdlDir() failed: should be 'abc' but is %s\n", MdlDir())
	}
}

func TestNoop(t *testing.T) {
	SetNoop(true)
	if !Noop() {
		t.Errorf("TestNoop() failed: should be true but is false\n")
	}
	SetNoop(false)
	if Noop() {
		t.Errorf("TestNoop() failed: should be false but is true\n")
	}
}

func TestOutDir(t *testing.T) {
	SetOutDir("xyzzy")
	if OutDir() != "xyzzy" {
		t.Errorf("TestOutDir() failed: should be 'xyzzy' but is %s\n", OutDir())
	}
	SetOutDir("abc")
	if OutDir() != "abc"{
		t.Errorf("TestOutDir() failed: should be 'abc' but is %s\n", OutDir())
	}
}

func TestQuiet(t *testing.T) {
	SetQuiet(true)
	if !Quiet() {
		t.Errorf("TestQuiet() failed: should be true but is false\n")
	}
	SetQuiet(false)
	if Quiet() {
		t.Errorf("TestQuiet() failed: should be false but is true\n")
	}
}

