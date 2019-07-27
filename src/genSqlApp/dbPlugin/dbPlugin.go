// See License.txt in main repository directory

// dbPlugin contains the data and functions that form
// the basis for a database plugin as well as being the
// plugin manager.  The various plugins inherit from this
// to provide specific support for a Database System such
// MySQL, Microsoft SQL or SQLite. The plugins will register
// with this package to insure that the system is aware of
// them.  However, they have to be imported at a higher
// level package to get their init() function to be called
// by go. They need not be used within the highest level
// package just imported.

// So, every plugin must supply a PluginData Struct and
// build upon it. The struct defines the two responsibilities
// of a database plugin, provide relevant sql components
// as needed at the database, table, row and field levels.
// and to help control the generation of the sql.

package dbPlugin

import (
	"../dbType"
	"fmt"
	"sync"
)

//============================================================================
//								Interfaces
//============================================================================

// GenBaseStringer is the interface that must be defined minimally for each Database
// Server module supported.
type GenBaseStringer	interface {
	GenFlagArgDefns(name string) string
	GenImportString()	string
}

// GenTableCreateStringer is the interface that defines the Table Creation and Deletion
// methods.
type GenTableCreateStringer	interface {
	GenCreateTableSQL(table interface{}) string
	GenDeleteTableSQL(table interface{}) string
}

//============================================================================
//                        		Plugin Support
//============================================================================

// Plugin Support defines the basis for all the Database manager plugins. Each
// plugin must address all the definitions within this structure.  It also
// provides a consistent external interface used by the rest of the system.

type PluginData	struct {
	Name			string				// Name of database
	Types			*dbType.TypeDefns	// Type definitions useable for this database
	Plugin			interface{}			// Used to supply various interfaces declared above
}


//----------------------------------------------------------------------------
//					Global Plugin Support Functions
//----------------------------------------------------------------------------

// plugins provides the source of findong all registered functions. The index \
// into it is the database name as used in the JSON input (ie mariadb, mssql,
// mysql, postgresql, sqlite) Each plugin registers with this package at init()
// insuring that the package is available when needed.
var plugins		map[string]PluginData
var mtxPlugins	sync.Mutex

// FindPlugin returns the Plugin interface for a name if possible. NIL is
// returned if it is not found.
func FindPlugin(name string) (PluginData, error) {
	mtxPlugins.Lock()
	defer mtxPlugins.Unlock()
	if plugins == nil {
		return PluginData{}, fmt.Errorf("Error: Plugin, %s, not found!\n", name)
	}
	if _, ok := plugins[name]; ok {
		return plugins[name], nil
	}
	return PluginData{}, fmt.Errorf("Error: Plugin, %s, not found!\n", name)
}

// Register adds or replaces the given plugin in the plugins map.
func Register(name string, plg PluginData) {
	mtxPlugins.Lock()
	defer mtxPlugins.Unlock()
	if plugins == nil {
		plugins = map[string]PluginData{}
	}
	plugins[name] = plg
}

// Unregister() deletes a given plugin in the plugin map.
func Unregister(name string) {
	mtxPlugins.Lock()
	defer mtxPlugins.Unlock()
	if plugins == nil {
		return
	}
	if _, ok := plugins[name]; ok {
		delete(plugins, name)
	}
}

// UnregisterAll() removes all given plugins from the plugin map.
func UnregisterAll() {
	mtxPlugins.Lock()
	defer mtxPlugins.Unlock()
	if plugins == nil {
		return
	}
	for n, _ := range plugins {
		delete(plugins, n)
	}
}

