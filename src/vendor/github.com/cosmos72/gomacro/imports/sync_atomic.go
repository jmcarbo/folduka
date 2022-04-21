// this file was generated by gomacro command: import _b "sync/atomic"
// DO NOT EDIT! Any change will be lost when the file is re-generated

package imports

import (
	. "reflect"
	atomic "sync/atomic"
)

// reflection: allow interpreted code to import "sync/atomic"
func init() {
	Packages["sync/atomic"] = Package{
	Name: "atomic",
	Binds: map[string]Value{
		"AddInt32":	ValueOf(atomic.AddInt32),
		"AddInt64":	ValueOf(atomic.AddInt64),
		"AddUint32":	ValueOf(atomic.AddUint32),
		"AddUint64":	ValueOf(atomic.AddUint64),
		"AddUintptr":	ValueOf(atomic.AddUintptr),
		"CompareAndSwapInt32":	ValueOf(atomic.CompareAndSwapInt32),
		"CompareAndSwapInt64":	ValueOf(atomic.CompareAndSwapInt64),
		"CompareAndSwapPointer":	ValueOf(atomic.CompareAndSwapPointer),
		"CompareAndSwapUint32":	ValueOf(atomic.CompareAndSwapUint32),
		"CompareAndSwapUint64":	ValueOf(atomic.CompareAndSwapUint64),
		"CompareAndSwapUintptr":	ValueOf(atomic.CompareAndSwapUintptr),
		"LoadInt32":	ValueOf(atomic.LoadInt32),
		"LoadInt64":	ValueOf(atomic.LoadInt64),
		"LoadPointer":	ValueOf(atomic.LoadPointer),
		"LoadUint32":	ValueOf(atomic.LoadUint32),
		"LoadUint64":	ValueOf(atomic.LoadUint64),
		"LoadUintptr":	ValueOf(atomic.LoadUintptr),
		"StoreInt32":	ValueOf(atomic.StoreInt32),
		"StoreInt64":	ValueOf(atomic.StoreInt64),
		"StorePointer":	ValueOf(atomic.StorePointer),
		"StoreUint32":	ValueOf(atomic.StoreUint32),
		"StoreUint64":	ValueOf(atomic.StoreUint64),
		"StoreUintptr":	ValueOf(atomic.StoreUintptr),
		"SwapInt32":	ValueOf(atomic.SwapInt32),
		"SwapInt64":	ValueOf(atomic.SwapInt64),
		"SwapPointer":	ValueOf(atomic.SwapPointer),
		"SwapUint32":	ValueOf(atomic.SwapUint32),
		"SwapUint64":	ValueOf(atomic.SwapUint64),
		"SwapUintptr":	ValueOf(atomic.SwapUintptr),
	}, Types: map[string]Type{
		"Value":	TypeOf((*atomic.Value)(nil)).Elem(),
	}, 
	}
}
