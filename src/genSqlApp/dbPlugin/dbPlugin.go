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

import "sync"

const (
	DBTYPE_MARIABDB	= 1 << iota
	DBTYPE_MSSQL
	DBTYPE_MYSQL
	DBTYPE_POSTGRES
	DBTYPE_SQLITE
)

//============================================================================
//                        Type Definition Support
//============================================================================

// TypeDefn, The type definition struct, defines one acceptable type accepted in the JSON
// defining the Database Structure.  There must be a TypeDefn for each type
// accepted in each plugin.

type TypeDefn struct {
	Name		string		`json:"Name,omitempty"`		// Type Name
	Html		string		`json:"Html,omitempty"`		// HTML Type
	Sql			string		`json:"Sql,omitempty"`		// SQL Type
	Go			string		`json:"Go,omitempty"`		// GO Type
	DftLen		int			`json:"DftLen,omitempty"`	// Default Length (used if length is not
	//													//	given)(0 == Max Length)
}

// TypeDefns provides a convenient way of defining a Type Definition Table which
// is used in the plugin environment.
type TypeDefns []TypeDefn

func (t TypeDefns) DftLen(name string) int {
	tdd := t.FindDefn(name)
	if tdd != nil {
		return tdd.DftLen
	}
	return -1
}

func (t TypeDefns) FindDefn(name string) *TypeDefn {
	for i,v := range t {
		if name == v.Name {
			return &t[i]
		}
	}
	return nil
}

func (t TypeDefns) GoType(name string) string {
	tdd := t.FindDefn(name)
	if tdd != nil {
		return tdd.Go
	}
	return ""
}

func (t TypeDefns) HtmlType(name string) string {
	tdd := t.FindDefn(name)
	if tdd != nil {
		return tdd.Html
	}
	return ""
}

func (t TypeDefns) SqlType(name string) string {
	tdd := t.FindDefn(name)
	if tdd != nil {
		return tdd.Sql
	}
	return ""
}

//============================================================================
//                        		Plugin Support
//============================================================================

// Plugin Support defines the basis for all the Database manager plugins. Each
// plugin must address all the definitions within this structure.  It also
// provides a consistent external interface used by the rest of the system.

type PluginData	struct {
	Name			string				// Name of database
	Types			*TypeDefns			// Type definitions useable for this database
	FlagsString 	func(name string) string
	ImportString 	func() string
	AddGo			bool			// Add "GO" after major sql statements
	CreateDB		bool
	NeedsUse		bool			// Need USE database; before doing anything
}

func (pd PluginData) FlagsStr(name string) string {
	return pd.FlagsString(name)
}

func (pd PluginData) ImportStr() string {
	return pd.ImportString()
}

//----------------------------------------------------------------------------
//					Global Plugin Support Functions
//----------------------------------------------------------------------------

// plugins provides the source of findong all registered functions. The index \
// into it is the database name as used in the JSON input (ie mariadb, mssql,
// mysql, postgresql, sqlite) Each plugin registers with this package at init()
// insuring that the package is available when needed.
var plugins		map[string]*PluginData
var mtxPlugins	sync.Mutex

// FindPlugin returns the Plugin interface for a name if possible. NIL is
// returned if it is not found.
func FindPlugin(name string) *PluginData {
	mtxPlugins.Lock()
	defer mtxPlugins.Unlock()
	if plugins == nil {
		return nil
	}
	if _, ok := plugins[name]; ok {
		return plugins[name]
	}
	return nil
}

// Register adds or replaces the given plugin in the plugins map.
func Register(pd *PluginData) error {
	mtxPlugins.Lock()
	defer mtxPlugins.Unlock()
	if plugins == nil {
		plugins = map[string]*PluginData{}
	}
	plugins[pd.Name] = pd
	return nil
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

