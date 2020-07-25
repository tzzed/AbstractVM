package evaluator

import (
	"avm/ast"
	"avm/ast/object"
	"avm/token"
	"github.com/golang-collections/collections/stack"
)

type Stack struct {
	Stack *stack.Stack
}

func New() *Stack {
	st := &Stack{}
	st.Stack = stack.New()
	return st
}

func (st *Stack) Eval(node ast.Node) object.Object {

	switch node := node.(type) {
	case *ast.Program:
		return st.evalStatements(node.Statements)
	case *ast.PushStatement:
		obj := st.evalPushStatement(node)
		return obj
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.IntValue}
	}

	return nil
}

func (st *Stack) evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, statement := range stmts {
		return st.Eval(statement)
	}

	return result
}

func (st *Stack) evalPushStatement(stmt *ast.PushStatement) object.Object {

	switch {
	case stmt.Name.String() == token.INT32:
		value := stmt.Value.(*ast.IntegerLiteral)
		st.Stack.Push(value.IntValue)
		return &object.Integer{Value: value.IntValue}
	case stmt.Name.String() == token.INT8:
		value := stmt.Value.(*ast.ByteLiteral)
		st.Stack.Push(value.ByteValue)
		return &object.Byte{Value: value.ByteValue}
	case stmt.Name.String() == token.INT16:
		value := stmt.Value.(*ast.ShortLiteral)
		st.Stack.Push(value.ShortValue)
		return &object.Short{Value: value.ShortValue}
	case stmt.Name.String() == token.FLOAT:
		value := stmt.Value.(*ast.FloatLiteral)
		st.Stack.Push(value.FloatValue)
		return &object.Float{Value: value.FloatValue}
	case stmt.Name.String() == token.DOUBLE:
		value := stmt.Value.(*ast.DoubleLiteral)
		st.Stack.Push(value.DoubleValue)
		return &object.Double{Value: value.DoubleValue}
	}

	return nil
}
