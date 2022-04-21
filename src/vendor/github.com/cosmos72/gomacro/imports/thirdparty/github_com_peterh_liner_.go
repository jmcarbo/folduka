// this file was generated by gomacro command: import "github.com/peterh/liner"
// DO NOT EDIT! Any change will be lost when the file is re-generated

package thirdparty

import (
	. "reflect"

	liner_ "github.com/peterh/liner"
)

// reflection: allow interpreted code to import "github.com/peterh/liner"
func init() {
	Packages["github.com/peterh/liner"] = Package{
		Name: "liner",
		Binds: map[string]Value{
			"ErrInternal":          ValueOf(&liner_.ErrInternal).Elem(),
			"ErrInvalidPrompt":     ValueOf(&liner_.ErrInvalidPrompt).Elem(),
			"ErrNotTerminalOutput": ValueOf(&liner_.ErrNotTerminalOutput).Elem(),
			"ErrPromptAborted":     ValueOf(&liner_.ErrPromptAborted).Elem(),
			"HistoryLimit":         ValueOf(liner_.HistoryLimit),
			"KillRingMax":          ValueOf(liner_.KillRingMax),
			"NewLiner":             ValueOf(liner_.NewLiner),
			"TabCircular":          ValueOf(liner_.TabCircular),
			"TabPrints":            ValueOf(liner_.TabPrints),
			"TerminalMode":         ValueOf(liner_.TerminalMode),
			"TerminalSupported":    ValueOf(liner_.TerminalSupported),
		}, Types: map[string]Type{
			"Completer":     TypeOf((*liner_.Completer)(nil)).Elem(),
			"ModeApplier":   TypeOf((*liner_.ModeApplier)(nil)).Elem(),
			"ShouldRestart": TypeOf((*liner_.ShouldRestart)(nil)).Elem(),
			"State":         TypeOf((*liner_.State)(nil)).Elem(),
			"TabStyle":      TypeOf((*liner_.TabStyle)(nil)).Elem(),
			"WordCompleter": TypeOf((*liner_.WordCompleter)(nil)).Elem(),
		}, Proxies: map[string]Type{
			"ModeApplier": TypeOf((*P_github_com_peterh_liner__ModeApplier)(nil)).Elem(),
		}, Untypeds: map[string]string{
			"HistoryLimit": "int:1000",
			"KillRingMax":  "int:60",
		},
	}
}

// --------------- proxy for github.com/peterh/liner.ModeApplier ---------------
type P_github_com_peterh_liner__ModeApplier struct {
	Object     interface{}
	ApplyMode_ func(interface{}) error
}

func (P *P_github_com_peterh_liner__ModeApplier) ApplyMode() error {
	return P.ApplyMode_(P.Object)
}
