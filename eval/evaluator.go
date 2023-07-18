package eval

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/bigyihsuan/structlang/token"
	"github.com/bigyihsuan/structlang/trees/ast"
)

type Evaluator struct {
	Code    []ast.Stmt
	BaseEnv Env
}

func NewEvaluator(code []ast.Stmt) Evaluator {
	var e Evaluator
	e.Code = code
	e.BaseEnv = NewEnv()
	return e
}

func (e *Evaluator) Evaluate(currEnv *Env) error {
	var errs error
	for _, stmt := range e.Code {
		var err error = nil
		switch stmt := stmt.(type) {
		case ast.TypeDef:
			err = e.TypeDef(currEnv, stmt)
		case ast.VarDef:
			err = e.VarDef(currEnv, stmt)
		case ast.VarSet:
			err = e.VarSet(currEnv, stmt)
		default:
			fmt.Printf("unknown stmt: %T\n", stmt)
		}
		if err != nil {
			errs = errors.Join(errs, err)
		}

	}
	return errs
}

func (e *Evaluator) TypeDef(currEnv *Env, stmt ast.TypeDef) error {
	typename, _ := e.TypeName(currEnv, stmt.Type)
	structdef, _ := e.StructDef(currEnv, stmt.StructDef)

	currEnv.DefineType(typename.Name, structdef)
	return nil
}

func (e *Evaluator) StructDef(currEnv *Env, structDef ast.StructDef) (st StructType, err error) {
	st.Fields = make(map[string]TypeName)
	st.Vars = make([]TypeName, len(structDef.Vars))

	for _, structField := range structDef.Fields {
		fieldType, err := e.TypeName(currEnv, structField.Type)
		if err != nil {
			return st, err
		}
		for _, fieldName := range structField.Names {
			st.Fields[fieldName.Name] = fieldType
		}
	}
	for i, typeVar := range structDef.Vars {
		tn, err := e.TypeName(currEnv, typeVar)
		if err != nil {
			return st, err
		}
		st.Vars[i] = tn
	}

	return st, nil
}

func (e *Evaluator) TypeName(currEnv *Env, typename ast.Type) (TypeName, error) {
	var type_ TypeName
	name := typename.Name.Name
	vars := []TypeName{}
	for _, typeArg := range typename.Vars {
		arg, _ := e.TypeName(currEnv, typeArg)
		vars = append(vars, arg)
	}
	type_.Name = name
	type_.Vars = vars
	return type_, nil
}

func (e *Evaluator) VarDef(currEnv *Env, varDef ast.VarDef) error {
	lvalue, err := e.Lvalue(currEnv, varDef.Lvalue)
	if err != nil {
		return err
	}
	rvalue, err := e.Expr(currEnv, varDef.Rvalue)
	if err != nil {
		return err
	}
	currEnv.DefineVariable(lvalue.Name, rvalue)
	return nil
}

func (e *Evaluator) VarSet(currEnv *Env, varSet ast.VarSet) error {
	lvalue, err := e.Lvalue(currEnv, varSet.Lvalue)
	if err != nil {
		return err
	}
	rvalue, err := e.Expr(currEnv, varSet.Rvalue)
	if err != nil {
		return err
	}
	return currEnv.SetVariable(lvalue.Name, rvalue)
}

func (e *Evaluator) Lvalue(currEnv *Env, lvalue ast.Lvalue) (Identifier, error) {
	switch lvalue := lvalue.(type) {
	case ast.Ident:
		ident := NewIdentifier(lvalue)
		return ident, nil
	case ast.FieldAccess:
		base, err := e.Lvalue(currEnv, lvalue.Lvalue)
		if err != nil {
			return base, err
		}
		ident := base.NewAccess(lvalue.Field)
		return ident, nil
	default:
		fmt.Printf("unknown lvalue: %T\n", lvalue)
	}
	return Identifier{}, fmt.Errorf("unkown lvalue: %v", lvalue)
}

func (e *Evaluator) Expr(currEnv *Env, expr ast.Expr) (v Value, err error) {
	switch expr := expr.(type) {
	case ast.Literal:
		return e.Literal(currEnv, expr)
	case ast.Ident:
		ident, err := e.Lvalue(currEnv, expr)
		if err != nil {
			return v, err
		}
		val := currEnv.GetVariable(ident.Name)
		if val == nil {
			return v, fmt.Errorf("variable `%s` not defined", expr.Name)
		}
		return *val, nil
	case ast.StructLiteral:
		return e.StructLiteral(currEnv, expr)
	case ast.FieldAccess:
		return e.FieldAccess(currEnv, expr)
	case ast.PrefixExpr:
		return e.PrefixExpr(currEnv, expr)
	case ast.InfixExpr:
		return e.InfixExpr(currEnv, expr)
	default:
		fmt.Printf("unknown expr: %T\n", expr)
	}

	return v, nil
}

func (e *Evaluator) Literal(currEnv *Env, expr ast.Literal) (v Value, err error) {
	switch expr.Token.Type() {
	case token.INT:
		v, err := strconv.Atoi(expr.Token.Lexeme())
		return NewInt(v), err
	case token.FLOAT:
		v, err := strconv.ParseFloat(expr.Token.Lexeme(), 64)
		return NewFloat(v), err
	case token.TRUE:
		v, err := strconv.ParseBool(expr.Token.Lexeme())
		return NewBool(v), err
	case token.FALSE:
		v, err := strconv.ParseBool(expr.Token.Lexeme())
		return NewBool(v), err
	case token.STRING:
		return NewString(expr.Token.Lexeme()), nil
	case token.NIL:
		return NewNil(), nil
	default:
		return v, fmt.Errorf("unknown literal %s", expr.Token.Type().String())
	}
}

func (e *Evaluator) TypeVars(currEnv *Env, types []ast.Type) (tv []TypeName, err error) {
	for _, t := range types {
		typeName, err := e.TypeName(currEnv, t)
		if err != nil {
			return tv, err
		}
		tv = append(tv, typeName)
	}
	return tv, nil
}

func (e *Evaluator) StructLiteral(currEnv *Env, expr ast.StructLiteral) (v Value, err error) {
	// basic duck typing
	// check if all names+types in the struct literal match the ones in the type definition
	typename := expr.TypeName.Name.Name
	st := currEnv.GetType(typename)
	if st == nil {
		return v, fmt.Errorf("type not found: %s", typename)
	}
	structTemplate := (*st).Copy()

	typeVars, err := e.TypeVars(currEnv, expr.TypeName.Vars)
	if err != nil {
		return v, err
	}
	if len(typeVars) != len(structTemplate.Vars) {
		return v, fmt.Errorf("not enough type parameters: want %d, got %d", len(structTemplate.Vars), len(typeVars))
	}

	// overwrite template type variables with concrete types
	typeParams := make(map[string]TypeName)
	for idx, concreteType := range typeVars {
		typeVar := structTemplate.Vars[idx]
		for name, fieldType := range structTemplate.Fields {
			if fieldType.Name == typeVar.Name {
				structTemplate.Fields[name] = concreteType
				typeParams[fieldType.Name] = concreteType
			}
		}
	}

	fields := make(map[string]Value)
	for _, field := range expr.Fields {
		name := field.Name.Name
		value, err := e.Expr(currEnv, field.Value)
		if err != nil {
			return v, err
		}
		expFieldType, ok := structTemplate.Fields[name]
		if !ok {
			return v, fmt.Errorf("field `%s` not found in type `%s`", name, typename)
		} else if valType := value.TypeName().Name; valType != expFieldType.Name {
			return v, fmt.Errorf("unexpected type for field `%s`: got `%s`, want `%s`", name, valType, expFieldType.Name)
		}
		fields[name] = value
	}
	sv := NewStructValueFromType(structTemplate, typeParams, fields, typename)

	return sv, nil
}

func (e *Evaluator) FieldAccess(currEnv *Env, expr ast.FieldAccess) (v Value, err error) {
	var base Value
	switch l := expr.Lvalue.(type) {
	case ast.Ident:
		b := currEnv.GetVariable(l.Name)
		if b == nil {
			return v, fmt.Errorf("variable `%s` not defined", l.Name)
		}
		base = *b
	case ast.FieldAccess:
		b, err := e.FieldAccess(currEnv, l)
		if err != nil {
			return v, err
		}
		base = b
	}
	return base.Get(expr.Field.Name), nil
}

func (e *Evaluator) PrefixExpr(currEnv *Env, expr ast.PrefixExpr) (v Value, err error) {
	v, err = e.Expr(currEnv, expr.Right)
	if err != nil {
		return v, err
	}

	if neg, isNeg := v.(Neg); isNeg {
		switch expr.Op.Type() {
		case token.PLUS:
			return neg.Pos(), nil
		case token.MINUS:
			return neg.Neg(), nil
		}
	}

	return v, fmt.Errorf("invalid type `%T` for prefix op `%s`", v, expr.Op.Type())
}

func (e *Evaluator) InfixExpr(currEnv *Env, expr ast.InfixExpr) (v Value, err error) {
	left, err := e.Expr(currEnv, expr.Left)
	if err != nil {
		return left, err
	}
	right, err := e.Expr(currEnv, expr.Right)
	if err != nil {
		return right, err
	}

	lsum, isLsum := left.(Sum)
	rsum, isRsum := right.(Sum)
	if isLsum && isRsum {
		switch expr.Op.Type() {
		case token.PLUS:
			return lsum.Add(rsum), nil
		case token.MINUS:
			return lsum.Sub(rsum), nil
		}
	}

	lprod, isLprod := left.(Product)
	rprod, isRprod := right.(Product)
	if isLprod && isRprod {
		switch expr.Op.Type() {
		case token.STAR:
			return lprod.Mul(rprod), nil
		case token.SLASH:
			return lprod.Div(rprod), nil
		}
	}

	lcmp, isLcmp := left.(Cmp)
	rcmp, isRcmp := right.(Cmp)
	if isLcmp && isRcmp {
		switch expr.Op.Type() {
		case token.GT:
			return lcmp.Gt(rcmp), nil
		case token.GTEQ:
			return lcmp.GtEq(rcmp), nil
		case token.LT:
			return lcmp.Lt(rcmp), nil
		case token.LTEQ:
			return lcmp.LtEq(rcmp), nil
		case token.EQ:
			return lcmp.Eq(rcmp), nil
		}
	}
	return v, fmt.Errorf("invalid types `%T` and `%T` for infix op `%s`", left, right, expr.Op.Type())
}
