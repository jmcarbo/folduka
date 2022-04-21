// this file was generated by gomacro command: import _b "compress/zlib"
// DO NOT EDIT! Any change will be lost when the file is re-generated

package imports

import (
	. "reflect"
	zlib "compress/zlib"
	io "io"
)

// reflection: allow interpreted code to import "compress/zlib"
func init() {
	Packages["compress/zlib"] = Package{
	Name: "zlib",
	Binds: map[string]Value{
		"BestCompression":	ValueOf(zlib.BestCompression),
		"BestSpeed":	ValueOf(zlib.BestSpeed),
		"DefaultCompression":	ValueOf(zlib.DefaultCompression),
		"ErrChecksum":	ValueOf(&zlib.ErrChecksum).Elem(),
		"ErrDictionary":	ValueOf(&zlib.ErrDictionary).Elem(),
		"ErrHeader":	ValueOf(&zlib.ErrHeader).Elem(),
		"HuffmanOnly":	ValueOf(zlib.HuffmanOnly),
		"NewReader":	ValueOf(zlib.NewReader),
		"NewReaderDict":	ValueOf(zlib.NewReaderDict),
		"NewWriter":	ValueOf(zlib.NewWriter),
		"NewWriterLevel":	ValueOf(zlib.NewWriterLevel),
		"NewWriterLevelDict":	ValueOf(zlib.NewWriterLevelDict),
		"NoCompression":	ValueOf(zlib.NoCompression),
	}, Types: map[string]Type{
		"Resetter":	TypeOf((*zlib.Resetter)(nil)).Elem(),
		"Writer":	TypeOf((*zlib.Writer)(nil)).Elem(),
	}, Proxies: map[string]Type{
		"Resetter":	TypeOf((*P_compress_zlib_Resetter)(nil)).Elem(),
	}, Untypeds: map[string]string{
		"BestCompression":	"int:9",
		"BestSpeed":	"int:1",
		"DefaultCompression":	"int:-1",
		"HuffmanOnly":	"int:-2",
		"NoCompression":	"int:0",
	}, 
	}
}

// --------------- proxy for compress/zlib.Resetter ---------------
type P_compress_zlib_Resetter struct {
	Object	interface{}
	Reset_	func(_proxy_obj_ interface{}, r io.Reader, dict []byte) error
}
func (P *P_compress_zlib_Resetter) Reset(r io.Reader, dict []byte) error {
	return P.Reset_(P.Object, r, dict)
}