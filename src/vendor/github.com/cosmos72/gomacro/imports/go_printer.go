// this file was generated by gomacro command: import _b "go/printer"
// DO NOT EDIT! Any change will be lost when the file is re-generated

package imports

import (
	. "reflect"
	printer "go/printer"
)

// reflection: allow interpreted code to import "go/printer"
func init() {
	Packages["go/printer"] = Package{
	Name: "printer",
	Binds: map[string]Value{
		"Fprint":	ValueOf(printer.Fprint),
		"RawFormat":	ValueOf(printer.RawFormat),
		"SourcePos":	ValueOf(printer.SourcePos),
		"TabIndent":	ValueOf(printer.TabIndent),
		"UseSpaces":	ValueOf(printer.UseSpaces),
	}, Types: map[string]Type{
		"CommentedNode":	TypeOf((*printer.CommentedNode)(nil)).Elem(),
		"Config":	TypeOf((*printer.Config)(nil)).Elem(),
		"Mode":	TypeOf((*printer.Mode)(nil)).Elem(),
	}, 
	}
}