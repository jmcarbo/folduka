// +build go1.8,gc,linux,!android go1.10,gc,darwin go1.14,gc,freebsd

// this file was generated by gomacro command: import _b "plugin"
// DO NOT EDIT! Any change will be lost when the file is re-generated

package imports

import (
	. "reflect"
	plugin "plugin"
)

// reflection: allow interpreted code to import "plugin"
func init() {
	Packages["plugin"] = Package{
	Name: "plugin",
	Binds: map[string]Value{
		"Open":	ValueOf(plugin.Open),
	}, Types: map[string]Type{
		"Plugin":	TypeOf((*plugin.Plugin)(nil)).Elem(),
		"Symbol":	TypeOf((*plugin.Symbol)(nil)).Elem(),
	}, 
	}
}
