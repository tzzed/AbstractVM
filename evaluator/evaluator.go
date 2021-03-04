package evaluator

import (
	"avm/ast"
	"avm/token"
	"errors"
	"fmt"
	"math"
)

type Node struct {
	v    Value
	next *Node
}

type Stack struct {
	head *Node
	size int
}

func NewStack() *Stack {
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
	if s.IsEmpty() {
		return Value{}, errors.New("error: pop on empty stack")
	}

	v := s.head.v
	s.head = s.head.next
	s.size--
	return v, nil
}

func (s *Stack) Clear() {
	if s.IsEmpty() {
		return
	}
	temp := s.head
	for temp != nil {
		_, _ = s.Pop()
		temp = temp.next
	}
}

func (s *Stack) IsEmpty() bool {
	return s.size == 0
}

func (s *Stack) Dup() {
	if s.IsEmpty() {
		return
	}

	dup := s.head.v
	s.Push(dup)
}

// Dump display each value on the stack.
func (s *Stack) Dump() (Value, error) {
	tmp := s.head
	for tmp != nil {
		fmt.Println(tmp.v)
		tmp = tmp.next
	}
	fmt.Println()

	return Value{}, nil
}

func (s *Stack) Swap() error {
	if s.size < 2 {
		return fmt.Errorf("error: Swap require stack size greater than 2: got %d", s.size)
	}

	a, _ := s.Pop()
	b, _ := s.Pop()
	s.Push(a)
	s.Push(b)
	return nil
}

func (s *Stack) InsertAt(p int, v Value) error {
	if p < 0 || p > s.size {
		return fmt.Errorf("index %d out of range", p)
	}

	st := NewStack()
	temp := s.head
	i := 0
	for temp != nil {
		if i == p {
			s.Push(v)
			break
		}

		va, err := s.Pop()
		if err != nil {
			return err
		}

		st.Push(va)
		temp = temp.next
		i++
	}

	temp = st.head
	for temp != nil {
		va, err := st.Pop()
		if err != nil {
			return err
		}

		s.Push(va)
		temp = temp.next
	}

	return nil
}

func (s *Stack) Peek(index int) (Value, error) {
	//stack size is between 0-15
	if index < 0 || index > 15 {
		return Value{}, fmt.Errorf("index %d out of range", index)
	}

	tmp := s.head
	i := 0
	for tmp != nil {
		if i == index {
			return tmp.v, nil
		}

		i++
		tmp = tmp.next
	}

	return Value{}, fmt.Errorf("no value at index %d", index)
}

func (s *Stack) Eval(node ast.Node) (Value, error) {
	switch n := node.(type) {
	case *ast.Program:
		return s.evalStatements(n.Statements)
	case *ast.PushStatement:
		return s.evalPushStatement(n)
	case *ast.AddStatement:
		return s.evalAdd()
	case *ast.AssertStatement:
		return s.evalAssert(n)
	case *ast.MulStatement:
		return s.evalMul()
	case *ast.DivStatement:
		return s.evalDiv()
	case *ast.ModStatement:
		return s.evalMod()
	case *ast.DumpStatement:
		return s.Dump()
	case *ast.PopStatement:
		return s.Pop()
	case *ast.ExpressionStatement:
		return s.Eval(n.Expression)
	case *ast.IntegerLiteral:
		return Value{V: n.IntValue, Type: IntegerValue}, nil
	case *ast.InfixExpression:
		left, err := s.Eval(n.Left)
		if err != nil {
			return Value{}, err
		}
		right, err := s.Eval(n.Right)
		if err != nil {
			return Value{}, err
		}
		return evalInfixExpression(n.Operator, left, right)
	default:
		return Value{}, fmt.Errorf("unknown instruction ")
	}

}

func evalIntegerInfixExpression(op string, left, right Value) (Value, error) {

	leftVal, err := left.ConvertToInteger()
	if err != nil {
		return left, err
	}
	rightVal, err := right.ConvertToInteger()
	if err != nil {
		return right, err
	}
	switch op {
	case token.PLUS:
		return NewInt32Value(leftVal + rightVal), nil
	case token.MINUS:
		return NewInt32Value(leftVal - rightVal), nil

	case token.ASTERISK:
		return NewInt32Value(leftVal * rightVal), nil
	case token.SLASH:
		return NewInt32Value(leftVal / rightVal), nil
	default:
		return Value{}, fmt.Errorf("no infix evaluator for %q\n", op)
	}
}

func evalInfixExpression(op string, left, right Value) (Value, error) {
	return evalIntegerInfixExpression(op, left, right)
}

func (s *Stack) evalStatements(stmts []ast.Statement) (Value, error) { return s.Eval(stmts[0]) }

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
	fmt.Printf("Here: %T\n", stmt.Value)
	v := Value{}
	var err error
	switch stmt.Value.(type) {
	case *ast.InfixExpression:
		v, err = s.Eval(stmt.Value)
		if err != nil {
			return Value{}, err
		}
		break
	default:
		v, err = convertAstToValue(stmt.Name.String(), stmt.Value)
		if err != nil {
			return Value{}, err
		}

	}


	s.Push(v)
	return v, nil
}

func (s *Stack) evalAdd() (Value, error) {
	if s.size < 2 {
		return Value{}, fmt.Errorf("stack size must be greater than 2: got %d", s.size)
	}

	a, _ := s.Pop()
	b, _ := s.Pop()
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
			return a, err
		}

		fb, err := b.ConvertToFloat()
		if err != nil {
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
	fmt.Printf("Here: %T\n", stmt.Value)
	v := Value{}
	var err error
	switch stmt.Value.(type) {
	case *ast.InfixExpression:
		v, err = s.Eval(stmt.Value)
		if err != nil {
			return Value{}, err
		}
		break
	default:
		v, err = convertAstToValue(stmt.Name.String(), stmt.Value)
		if err != nil {
			return Value{}, err
		}

	}

	fmt.Printf("left:: %v  \n", v)
	if s.IsEmpty() {
		return Value{}, errors.New("cannot check value empty stack")
	}

	res := s.head

	if res.v.V == v.V && res.v.Type == v.Type {
		return v, nil
	}

	return v, fmt.Errorf("expected %s(%v) stack contains  %s(%v)", v.Type, v.V, res.v.Type, res.v.V)
}

func (s *Stack) evalMod() (Value, error) {
	a, err := s.Pop()
	if err != nil {
		return Value{}, err
	}

	b, err := s.Pop()
	if err != nil {
		return Value{}, err
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

		// if the divisor by 0 returns error
		if cb == 0 {
			return Value{}, errors.New("error: integer divide by zero")
		}

		v := NewInt8Value(ca % cb)
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

		if sb == 0 {
			return Value{}, errors.New("error: integer divide by zero")
		}
		v := NewInt16Value(int16(sa % sb))
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

		if b.V.(int32) == 0 {
			return Value{}, errors.New("error: integer divide by zero")
		}

		v := NewInt32Value(ia % ib)
		s.Push(v)
		return v, nil
	case FloatValue:
		fa, err := a.ConvertToFloat()
		if err != nil {
			return a, err
		}

		fb, err := b.ConvertToFloat()
		if err != nil {
			return b, err
		}

		if b.V.(float32) == 0 {
			return Value{}, errors.New("error: integer divide by zero")
		}

		f := math.Mod(float64(fa), float64(fb))
		f = math.Round(f * math.Pow10(2)) / math.Pow10(2)
		v := NewFloatValue(float32(f))
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

		if b.V.(float64) == 0 || a.V.(float64) == 0 {
			return Value{}, errors.New("error: integer divide by zero")
		}
		f := math.Mod(da, db)
		v := NewDoubleValue(f)
		s.Push(v)
		return v, nil
	}

	return Value{}, nil
}

func (s *Stack) evalDiv() (Value, error) {
	a, err := s.Pop()
	if err != nil {
		return Value{}, err
	}

	b, err := s.Pop()
	if err != nil {
		return Value{}, err
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

		// if the divisor by 0 returns error
		if cb == 0 {
			return Value{}, errors.New("error: integer divide by zero")
		}

		v := NewInt8Value(ca / cb)
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

		if sb == 0 {
			return Value{}, errors.New("error: integer divide by zero")
		}
		v := NewInt16Value(int16(sa / sb))
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

		if b.V.(int32) == 0 {
			return Value{}, errors.New("error: integer divide by zero")
		}

		v := NewInt32Value(ia / ib)
		s.Push(v)
		return v, nil
	case FloatValue:
		fa, err := a.ConvertToFloat()
		if err != nil {
			return a, err
		}

		fb, err := b.ConvertToFloat()
		if err != nil {
			return b, err
		}

		if b.V.(float32) == 0 {
			return Value{}, errors.New("error: integer divide by zero")
		}

		v := NewFloatValue(fa / fb)
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

		if b.V.(float64) == 0 || a.V.(float64) == 0 {
			return Value{}, errors.New("error: integer divide by zero")
		}
		v := NewDoubleValue(da / db)
		s.Push(v)
		return v, nil
	}

	return Value{}, nil
}

func (s *Stack) evalMul() (Value, error) {

	a, err := s.Pop()
	if err != nil {
		return Value{}, err
	}

	b, err := s.Pop()
	if err != nil {
		return Value{}, err
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

		v := NewInt8Value(ca * cb)
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

		v := NewInt16Value(sa * sb)
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

		v := NewInt32Value(ia * ib)
		s.Push(v)
		return v, nil
	case FloatValue:
		fa, err := a.ConvertToFloat()
		if err != nil {
			return a, err
		}

		fb, err := b.ConvertToFloat()
		if err != nil {
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

		v := NewDoubleValue(da * db)
		s.Push(v)
		return v, nil
	}

	return Value{}, nil
}
