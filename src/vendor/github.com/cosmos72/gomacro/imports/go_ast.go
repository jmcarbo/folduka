// this file was generated by gomacro command: import _b "go/ast"
// DO NOT EDIT! Any change will be lost when the file is re-generated

package imports

import (
	. "reflect"
	ast "go/ast"
	token "go/token"
)

// reflection: allow interpreted code to import "go/ast"
func init() {
	Packages["go/ast"] = Package{
	Name: "ast",
	Binds: map[string]Value{
		"Bad":	ValueOf(ast.Bad),
		"Con":	ValueOf(ast.Con),
		"FileExports":	ValueOf(ast.FileExports),
		"FilterDecl":	ValueOf(ast.FilterDecl),
		"FilterFile":	ValueOf(ast.FilterFile),
		"FilterFuncDuplicates":	ValueOf(ast.FilterFuncDuplicates),
		"FilterImportDuplicates":	ValueOf(ast.FilterImportDuplicates),
		"FilterPackage":	ValueOf(ast.FilterPackage),
		"FilterUnassociatedComments":	ValueOf(ast.FilterUnassociatedComments),
		"Fprint":	ValueOf(ast.Fprint),
		"Fun":	ValueOf(ast.Fun),
		"Inspect":	ValueOf(ast.Inspect),
		"IsExported":	ValueOf(ast.IsExported),
		"Lbl":	ValueOf(ast.Lbl),
		"MergePackageFiles":	ValueOf(ast.MergePackageFiles),
		"NewCommentMap":	ValueOf(ast.NewCommentMap),
		"NewIdent":	ValueOf(ast.NewIdent),
		"NewObj":	ValueOf(ast.NewObj),
		"NewPackage":	ValueOf(ast.NewPackage),
		"NewScope":	ValueOf(ast.NewScope),
		"NotNilFilter":	ValueOf(ast.NotNilFilter),
		"PackageExports":	ValueOf(ast.PackageExports),
		"Pkg":	ValueOf(ast.Pkg),
		"Print":	ValueOf(ast.Print),
		"RECV":	ValueOf(ast.RECV),
		"SEND":	ValueOf(ast.SEND),
		"SortImports":	ValueOf(ast.SortImports),
		"Typ":	ValueOf(ast.Typ),
		"Var":	ValueOf(ast.Var),
		"Walk":	ValueOf(ast.Walk),
	}, Types: map[string]Type{
		"ArrayType":	TypeOf((*ast.ArrayType)(nil)).Elem(),
		"AssignStmt":	TypeOf((*ast.AssignStmt)(nil)).Elem(),
		"BadDecl":	TypeOf((*ast.BadDecl)(nil)).Elem(),
		"BadExpr":	TypeOf((*ast.BadExpr)(nil)).Elem(),
		"BadStmt":	TypeOf((*ast.BadStmt)(nil)).Elem(),
		"BasicLit":	TypeOf((*ast.BasicLit)(nil)).Elem(),
		"BinaryExpr":	TypeOf((*ast.BinaryExpr)(nil)).Elem(),
		"BlockStmt":	TypeOf((*ast.BlockStmt)(nil)).Elem(),
		"BranchStmt":	TypeOf((*ast.BranchStmt)(nil)).Elem(),
		"CallExpr":	TypeOf((*ast.CallExpr)(nil)).Elem(),
		"CaseClause":	TypeOf((*ast.CaseClause)(nil)).Elem(),
		"ChanDir":	TypeOf((*ast.ChanDir)(nil)).Elem(),
		"ChanType":	TypeOf((*ast.ChanType)(nil)).Elem(),
		"CommClause":	TypeOf((*ast.CommClause)(nil)).Elem(),
		"Comment":	TypeOf((*ast.Comment)(nil)).Elem(),
		"CommentGroup":	TypeOf((*ast.CommentGroup)(nil)).Elem(),
		"CommentMap":	TypeOf((*ast.CommentMap)(nil)).Elem(),
		"CompositeLit":	TypeOf((*ast.CompositeLit)(nil)).Elem(),
		"Decl":	TypeOf((*ast.Decl)(nil)).Elem(),
		"DeclStmt":	TypeOf((*ast.DeclStmt)(nil)).Elem(),
		"DeferStmt":	TypeOf((*ast.DeferStmt)(nil)).Elem(),
		"Ellipsis":	TypeOf((*ast.Ellipsis)(nil)).Elem(),
		"EmptyStmt":	TypeOf((*ast.EmptyStmt)(nil)).Elem(),
		"Expr":	TypeOf((*ast.Expr)(nil)).Elem(),
		"ExprStmt":	TypeOf((*ast.ExprStmt)(nil)).Elem(),
		"Field":	TypeOf((*ast.Field)(nil)).Elem(),
		"FieldFilter":	TypeOf((*ast.FieldFilter)(nil)).Elem(),
		"FieldList":	TypeOf((*ast.FieldList)(nil)).Elem(),
		"File":	TypeOf((*ast.File)(nil)).Elem(),
		"Filter":	TypeOf((*ast.Filter)(nil)).Elem(),
		"ForStmt":	TypeOf((*ast.ForStmt)(nil)).Elem(),
		"FuncDecl":	TypeOf((*ast.FuncDecl)(nil)).Elem(),
		"FuncLit":	TypeOf((*ast.FuncLit)(nil)).Elem(),
		"FuncType":	TypeOf((*ast.FuncType)(nil)).Elem(),
		"GenDecl":	TypeOf((*ast.GenDecl)(nil)).Elem(),
		"GoStmt":	TypeOf((*ast.GoStmt)(nil)).Elem(),
		"Ident":	TypeOf((*ast.Ident)(nil)).Elem(),
		"IfStmt":	TypeOf((*ast.IfStmt)(nil)).Elem(),
		"ImportSpec":	TypeOf((*ast.ImportSpec)(nil)).Elem(),
		"Importer":	TypeOf((*ast.Importer)(nil)).Elem(),
		"IncDecStmt":	TypeOf((*ast.IncDecStmt)(nil)).Elem(),
		"IndexExpr":	TypeOf((*ast.IndexExpr)(nil)).Elem(),
		"InterfaceType":	TypeOf((*ast.InterfaceType)(nil)).Elem(),
		"KeyValueExpr":	TypeOf((*ast.KeyValueExpr)(nil)).Elem(),
		"LabeledStmt":	TypeOf((*ast.LabeledStmt)(nil)).Elem(),
		"MapType":	TypeOf((*ast.MapType)(nil)).Elem(),
		"MergeMode":	TypeOf((*ast.MergeMode)(nil)).Elem(),
		"Node":	TypeOf((*ast.Node)(nil)).Elem(),
		"ObjKind":	TypeOf((*ast.ObjKind)(nil)).Elem(),
		"Object":	TypeOf((*ast.Object)(nil)).Elem(),
		"Package":	TypeOf((*ast.Package)(nil)).Elem(),
		"ParenExpr":	TypeOf((*ast.ParenExpr)(nil)).Elem(),
		"RangeStmt":	TypeOf((*ast.RangeStmt)(nil)).Elem(),
		"ReturnStmt":	TypeOf((*ast.ReturnStmt)(nil)).Elem(),
		"Scope":	TypeOf((*ast.Scope)(nil)).Elem(),
		"SelectStmt":	TypeOf((*ast.SelectStmt)(nil)).Elem(),
		"SelectorExpr":	TypeOf((*ast.SelectorExpr)(nil)).Elem(),
		"SendStmt":	TypeOf((*ast.SendStmt)(nil)).Elem(),
		"SliceExpr":	TypeOf((*ast.SliceExpr)(nil)).Elem(),
		"Spec":	TypeOf((*ast.Spec)(nil)).Elem(),
		"StarExpr":	TypeOf((*ast.StarExpr)(nil)).Elem(),
		"Stmt":	TypeOf((*ast.Stmt)(nil)).Elem(),
		"StructType":	TypeOf((*ast.StructType)(nil)).Elem(),
		"SwitchStmt":	TypeOf((*ast.SwitchStmt)(nil)).Elem(),
		"TypeAssertExpr":	TypeOf((*ast.TypeAssertExpr)(nil)).Elem(),
		"TypeSpec":	TypeOf((*ast.TypeSpec)(nil)).Elem(),
		"TypeSwitchStmt":	TypeOf((*ast.TypeSwitchStmt)(nil)).Elem(),
		"UnaryExpr":	TypeOf((*ast.UnaryExpr)(nil)).Elem(),
		"ValueSpec":	TypeOf((*ast.ValueSpec)(nil)).Elem(),
		"Visitor":	TypeOf((*ast.Visitor)(nil)).Elem(),
	}, Proxies: map[string]Type{
		"Node":	TypeOf((*P_go_ast_Node)(nil)).Elem(),
		"Visitor":	TypeOf((*P_go_ast_Visitor)(nil)).Elem(),
	}, 
	}
}

// --------------- proxy for go/ast.Node ---------------
type P_go_ast_Node struct {
	Object	interface{}
	End_	func(interface{}) token.Pos
	Pos_	func(interface{}) token.Pos
}
func (P *P_go_ast_Node) End() token.Pos {
	return P.End_(P.Object)
}
func (P *P_go_ast_Node) Pos() token.Pos {
	return P.Pos_(P.Object)
}

// --------------- proxy for go/ast.Visitor ---------------
type P_go_ast_Visitor struct {
	Object	interface{}
	Visit_	func(_proxy_obj_ interface{}, node ast.Node) (w ast.Visitor)
}
func (P *P_go_ast_Visitor) Visit(node ast.Node) (w ast.Visitor) {
	return P.Visit_(P.Object, node)
}