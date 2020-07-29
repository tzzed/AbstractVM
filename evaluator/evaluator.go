package evaluator

import (
	"avm/ast"
	"avm/token"
	"errors"
	"fmt"
)

type Node struct {
	v    Value
	next *Node
}

type Stack struct {
	head *Node
	size int
}

func New() *Stack {
	s := &Stack{nil, 0}
	return s
}

func (s *Stack) Size() int {
	return s.size
}

func (s *Stack) Push(v Value) {
	s.head = &Node{v: v, next: s.head}
	s.size++
}

func (s *Stack) Pop() (Value, error) {
	if s.size == 0 {
		return Value{}, errors.New("error stack is empty")
	}

	v := s.head.v
	s.head = s.head.next
	s.size--
	return v, nil

}

func (s *Stack) Eval(node ast.Node) (Value, error) {

	switch node := node.(type) {
	case *ast.Program:
		return s.evalStatements(node.Statements)
	case *ast.PushStatement:
		return s.evalPushStatement(node)
	case *ast.AddStatement:
		return s.evalAdd()
	case *ast.IntegerLiteral:
		return Value{V: node.IntValue}, nil
	case *ast.AssertStatement:
		return s.evalAssert(node)
	default:
		return Value{}, fmt.Errorf("unknown instruction %v", node)
	}

}

func (s *Stack) evalStatements(stmts []ast.Statement) (Value, error) {

	for _, statement := range stmts {
		return s.Eval(statement)
	}

	return Value{}, errors.New("no statement")
}

func convertAstToValue(n string, expr ast.Expression) (Value, error) {
	switch {
	case n == token.INT32:
		value := expr.(*ast.IntegerLiteral)
		v := NewInt32Value(value.IntValue)
		return v, nil
	case n == token.INT8:
		value := expr.(*ast.ByteLiteral)
		v := NewInt8Value(value.ByteValue)
		return v, nil
	case n == token.INT16:
		value := expr.(*ast.ShortLiteral)
		v := NewInt16Value(value.ShortValue)
		return v, nil
	case n == token.FLOAT:
		value := expr.(*ast.FloatLiteral)
		v := NewFloatValue(value.FloatValue)
		return v, nil
	case n == token.DOUBLE:
		value := expr.(*ast.DoubleLiteral)
		v := NewDoubleValue(value.DoubleValue)
		return v, nil
	}

	return Value{}, fmt.Errorf("bad statement %s", n)
}

func (s *Stack) evalPushStatement(stmt *ast.PushStatement) (Value, error) {

	v, err := convertAstToValue(stmt.Name.String(), stmt.Value)
	if err != nil {
		return Value{}, err
	}

	s.Push(v)
	return v, nil
}

func (s *Stack) Print() {
	tmp := s.head
	for tmp != nil {
		fmt.Println(tmp.v)
		tmp = tmp.next
	}
	fmt.Println()
}

func (s *Stack) evalAdd() (Value, error) {
	a, err := s.Pop()
	if err != nil {
		return Value{}, errors.New("empty stack")
	}

	b, err := s.Pop()
	if err != nil {
		return Value{}, errors.New("not enough elem in stack")
	}

	switch GetBiggerType(a, b) {
	case CharValue:
		ca, err := a.ConvertToChar()
		if err != nil {
			return a, err
		}

		cb, err := b.ConvertToChar()
		if err != nil {
			return b, err
		}

		v := NewInt8Value(ca + cb)
		s.Push(v)
		return v, nil
	case ShortValue:
		sa, err := a.ConvertToShort()
		if err != nil {
			return Value{}, err
		}
		sb, err := b.ConvertToShort()
		if err != nil {
			return Value{}, err
		}

		v := NewInt16Value(sa + sb)
		s.Push(v)
		return v, nil
	case IntegerValue:
		ia, err := a.ConvertToInteger()
		if err != nil {
			return a, err
		}

		ib, err := b.ConvertToInteger()
		if err != nil {
			return b, err
		}

		v := NewInt32Value(ia + ib)
		s.Push(v)
		return v, nil
	case FloatValue:
		fa, err := a.ConvertToFloat()
		if err != nil {
			fmt.Println(fa)
			return a, err
		}

		fb, err := b.ConvertToFloat()
		if err != nil {
			fmt.Println(fb)
			return b, err
		}

		v := NewFloatValue(fa + fb)
		s.Push(v)
		return v, nil
	case DoubleValue:
		da, err := a.ConvertToDouble()
		if err != nil {
			return a, err
		}

		db, err := b.ConvertToDouble()
		if err != nil {
			return b, err
		}

		v := NewDoubleValue(da + db)
		s.Push(v)
		return v, nil
	}

	return Value{}, fmt.Errorf("unsupported type %s or %s", a.Type, b.Type)
}

func (s *Stack) evalAssert(stmt *ast.AssertStatement) (Value, error) {
	res := s.head
	v, err := convertAstToValue(stmt.Name.String(), stmt.Value)
	if err != nil {
		return Value{}, err
	}

	if res.v.V == v.V && res.v.Type == v.Type {
		return v, nil
	}

	return v, fmt.Errorf("expected %s(%v) stack contains  %s(%v)", v.Type, v.V, res.v.Type, res.v.V )
}
