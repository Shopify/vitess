package asthelpergen

import (
	"go/types"
	"log"

	"github.com/dave/jennifer/jen"
)

type equalsGen struct {
	todo    []types.Type
	methods []jen.Code
	scope   *types.Scope
}

var _ generator = (*equalsGen)(nil)

func newEqualsGen(scope *types.Scope) *equalsGen {
	return &equalsGen{
		scope: scope,
	}
}

func (e *equalsGen) visitStruct(t types.Type, stroct *types.Struct) error {
	return nil
}

func (e *equalsGen) visitSlice(t types.Type, slice *types.Slice) error {
	return nil
}

func (e *equalsGen) visitInterface(t types.Type, iface *types.Interface) error {
	e.todo = append(e.todo, t)
	return nil
}

func (e *equalsGen) createFile(pkgName string) (string, *jen.File) {
	out := jen.NewFile(pkgName)
	out.HeaderComment(licenseFileHeader)
	out.HeaderComment("Code generated by ASTHelperGen. DO NOT EDIT.")

	alreadyDone := map[string]bool{}
	for len(e.todo) > 0 {
		t := e.todo[0]
		underlying := t.Underlying()
		typeName := printableTypeName(t)
		e.todo = e.todo[1:]

		if alreadyDone[typeName] {
			continue
		}

		if e.tryInterface(underlying, t) ||
			e.tryStruct(underlying, t) ||
			e.trySlice(underlying, t) ||
			e.tryPtr(underlying, t) {
			alreadyDone[typeName] = true
			continue
		}

		log.Printf("don't know how to handle %s %T", typeName, underlying)
	}

	for _, method := range e.methods {
		out.Add(method)
	}

	return "equals.go", out
}

func (e *equalsGen) tryInterface(underlying, t types.Type) bool {
	iface, ok := underlying.(*types.Interface)
	if !ok {
		return false
	}

	err := e.makeInterfaceEqualsMethod(t, iface)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return true
}

const equalsName = "Equals"

func (e *equalsGen) makeInterfaceEqualsMethod(t types.Type, iface *types.Interface) error {

	/*
		func EqualsAST(inA, inB AST) bool {
			if inA == inB {
				return true
			}
			if inA == nil || inB8 == nil {
				return false
			}
			switch a := inA.(type) {
			case *SubImpl:
				b, ok := inB.(*SubImpl)
				if !ok {
					return false
				}
				return EqualsSubImpl(a, b)
			}
			return false
		}
	*/
	stmts := []jen.Code{
		jen.If(jen.Id("inA == nil").Op("&&").Id("inB == nil")).Block(jen.Return(jen.True())),
		jen.If(jen.Id("inA == nil").Op("||").Id("inB == nil")).Block(jen.Return(jen.False())),
	}

	var cases []jen.Code
	_ = findImplementations(e.scope, iface, func(t types.Type) error {
		if _, ok := t.Underlying().(*types.Interface); ok {
			return nil
		}
		typeString := types.TypeString(t, noQualifier)
		caseBlock := jen.Case(jen.Id(typeString)).Block(
			jen.Id("b, ok := inB.").Call(jen.Id(typeString)),
			jen.If(jen.Id("!ok")).Block(jen.Return(jen.False())),
			jen.Return(e.compareValueType(t, jen.Id("a"), jen.Id("b"), true)),
		)
		cases = append(cases, caseBlock)
		return nil
	})

	cases = append(cases,
		jen.Default().Block(
			jen.Comment("this should never happen"),
			jen.Return(jen.False()),
		))

	stmts = append(stmts, jen.Switch(jen.Id("a := inA.(type)").Block(
		cases...,
	)))

	stmts = append(stmts, jen.Return(jen.False()))

	typeString := types.TypeString(t, noQualifier)
	funcName := equalsName + printableTypeName(t)
	funcDecl := jen.Func().Id(funcName).Call(jen.List(jen.Id("inA"), jen.Id("inB")).Id(typeString)).Bool().Block(stmts...)
	e.addFunc(funcName, funcDecl)

	return nil
}

func (e *equalsGen) compareValueType(t types.Type, a, b *jen.Statement, eq bool) *jen.Statement {
	switch t.Underlying().(type) {
	case *types.Basic:
		if eq {
			return a.Op("==").Add(b)
		}
		return a.Op("!=").Add(b)
	}

	e.todo = append(e.todo, t)
	var neg = "!"
	if eq {
		neg = ""
	}
	return jen.Id(neg+equalsName+printableTypeName(t)).Call(a, b)
}

func (e *equalsGen) addFunc(name string, code jen.Code) {
	e.methods = append(e.methods, jen.Comment(name+" does deep equals."), code)
}

func (e *equalsGen) tryStruct(underlying, t types.Type) bool {
	stroct, ok := underlying.(*types.Struct)
	if !ok {
		return false
	}

	err := e.makeStructEqualsMethod(t, stroct)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return true
}

func (e *equalsGen) makeStructEqualsMethod(t types.Type, stroct *types.Struct) error {
	/*
		func EqualsRefOfRefContainer(inA RefContainer, inB RefContainer) bool {
			return EqualsRefOfLeaf(inA.ASTImplementationType, inB.ASTImplementationType) &&
				EqualsAST(inA.ASTType, inB.ASTType) && inA.NotASTType == inB.NotASTType
		}

	*/

	typeString := types.TypeString(t, noQualifier)
	funcName := equalsName + printableTypeName(t)
	funcDecl := jen.Func().Id(funcName).Call(jen.List(jen.Id("a"), jen.Id("b")).Id(typeString)).Bool().
		Block(jen.Return(e.compareAllStructFields(stroct)))
	e.addFunc(funcName, funcDecl)

	return nil
}

func (e *equalsGen) compareAllStructFields(stroct *types.Struct) jen.Code {
	var basicsPred []*jen.Statement
	var others []*jen.Statement
	for i := 0; i < stroct.NumFields(); i++ {
		field := stroct.Field(i)
		if field.Type().Underlying().String() == "interface{}" {
			// we can safely ignore this, we do not want ast to contain interface{} types.
			continue
		}
		fieldA := jen.Id("a").Dot(field.Name())
		fieldB := jen.Id("b").Dot(field.Name())
		pred := e.compareValueType(field.Type(), fieldA, fieldB, true)
		if _, ok := field.Type().(*types.Basic); ok {
			basicsPred = append(basicsPred, pred)
			continue
		}
		others = append(others, pred)
	}

	var ret *jen.Statement
	for _, pred := range basicsPred {
		if ret == nil {
			ret = pred
		} else {
			ret = ret.Op("&&").Line().Add(pred)
		}
	}

	for _, pred := range others {
		if ret == nil {
			ret = pred
		} else {
			ret = ret.Op("&&").Line().Add(pred)
		}
	}

	if ret == nil {
		return jen.True()
	}
	return ret
}

func (e *equalsGen) tryPtr(underlying, t types.Type) bool {
	ptr, ok := underlying.(*types.Pointer)
	if !ok {
		return false
	}

	if strct, isStruct := ptr.Elem().Underlying().(*types.Struct); isStruct {
		e.makePtrToStructCloneMethod(t, strct)
		return true
	}

	return false
}

func (e *equalsGen) makePtrToStructCloneMethod(t types.Type, strct *types.Struct) {
	typeString := types.TypeString(t, noQualifier)
	funcName := equalsName + printableTypeName(t)

	//func EqualsRefOfType(a,b  *Type) *Type
	funcDeclaration := jen.Func().Id(funcName).Call(jen.Id("a"), jen.Id("b").Id(typeString)).Bool()
	stmts := []jen.Code{
		jen.If(jen.Id("a == b")).Block(jen.Return(jen.True())),
		jen.If(jen.Id("a == nil").Op("||").Id("b == nil")).Block(jen.Return(jen.False())),
		jen.Return(e.compareAllStructFields(strct)),
	}

	e.methods = append(e.methods, funcDeclaration.Block(stmts...))
}

func (e *equalsGen) trySlice(underlying, t types.Type) bool {
	slice, ok := underlying.(*types.Slice)
	if !ok {
		return false
	}

	err := e.makeSliceEqualsMethod(t, slice)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return true
}

func (e *equalsGen) makeSliceEqualsMethod(t types.Type, slice *types.Slice) error {
	/*
		func EqualsSliceOfRefOfLeaf(a, b []*Leaf) bool {
			if len(a) != len(b) {
				return false
			}
			for i := 0; i < len(a); i++ {
				if !EqualsRefOfLeaf(a[i], b[i]) {
					return false
				}
			}
			return false
		}
	*/

	stmts := []jen.Code{jen.If(jen.Id("len(a) != len(b)")).Block(jen.Return(jen.False())),
		jen.For(jen.Id("i := 0; i < len(a); i++")).Block(
			jen.If(e.compareValueType(slice.Elem(), jen.Id("a[i]"), jen.Id("b[i]"), false)).Block(jen.Return(jen.False()))),
		jen.Return(jen.True()),
	}

	typeString := types.TypeString(t, noQualifier)
	funcName := equalsName + printableTypeName(t)
	funcDecl := jen.Func().Id(funcName).Call(jen.List(jen.Id("a"), jen.Id("b")).Id(typeString)).Bool().Block(stmts...)
	e.addFunc(funcName, funcDecl)
	return nil
}
