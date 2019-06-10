// See License.txt in main repository directory

// dbPlugin contains the data and functions that form
// the basis for a database plugin.  The various plugins
// inherit from this to provide spedific support for a
// Database System such MySQL, Microsoft SQL or SQLite.

package dbForm

import (
	"../dbJson"
)

//============================================================================
//                        Type Definition Support
//============================================================================

// SqlWork, The type definition struct, defines one acceptable type accepted in the JSON
// defining the Database Structure.  There must be a TypeDefn for each type
// accepted in each plugin.

type FormWork struct {
	Name		string		`json:"Name,omitempty"`		// Type Name
	j			dbJson.Database
}

//----------------------------------------------------------------------------
//						Global/Internal Object Functions
//----------------------------------------------------------------------------

/***
func (t FormWork) FindDefn(name string) *FormWork {
	for i, v := range t {
		if name == v.Name {
			return &t[i]
		}
	}
	return nil
}
 ***/

func (t FormWork) Write(p []byte) (int, error) {
	var n		int
	var err		error

	// Write p to somewhere.

	return n, err
}

func (t FormWork) WriteString(s string) (int, error) {
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
func New(j *dbJson.Database) (*FormWork) {
	fw := &FormWork{}
	fw.j = *j
	return fw
}

