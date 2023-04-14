package main

var (
  NULL_LIT  = &Null{}
  TRUE_LIT  = &Boolean{Value: true}
  FALSE_LIT = &Boolean{Value: false}
)

func Eval(node Node) Object {
  switch node := node.(type) {
  case *BooleanExpression:
    return nativeBoolToBooleanObject(node.Value)
  case *ExpressionStatement:
    return Eval(node.Expression)
  case *IntegerLiteral:
    return &Integer{Value: node.Value}
  case *PrefixExpression:
    return evalPrefixExpression(node.Operator, Eval(node.Right))
  case *Program:
    return evalStatements(node.Statements)
  }

  return nil
}

func evalStatements(statements []Statement) Object {
  var result Object

  for _, statement := range statements {
    result = Eval(statement)
  }

  return result
}

func nativeBoolToBooleanObject(input bool) *Boolean {
  if input {
    return TRUE_LIT
  } else {
    return FALSE_LIT
  }
}

func evalPrefixExpression(operator string, right Object) Object {
  switch operator {
  case "!":
    return evalBangOperatorExpression(right)
  case "-":
    return evalMinusPrefixOperatorExpression(right)
  default:
    return NULL_LIT
  }
}

func evalBangOperatorExpression(right Object) Object {
  switch right {
  case TRUE_LIT:
    return FALSE_LIT
  case FALSE_LIT:
    return TRUE_LIT
  case NULL_LIT:
    return TRUE_LIT
  default:
    return FALSE_LIT
  }
}

func evalMinusPrefixOperatorExpression(right Object) Object {
  if right.Type() != INTEGER_OBJ {
    return NULL_LIT
  }

  value := right.(*Integer).Value

  return &Integer{Value: -value}
}
