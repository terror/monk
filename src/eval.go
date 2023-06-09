package main

var (
  NULL_LIT  = &Null{}
  TRUE_LIT  = &Boolean{Value: true}
  FALSE_LIT = &Boolean{Value: false}
)

func Eval(node Node) Object {
  switch node := node.(type) {
  case *BlockStatement:
    return evalStatements(node.Statements)
  case *BooleanExpression:
    return nativeBoolToBooleanObject(node.Value)
  case *ExpressionStatement:
    return Eval(node.Expression)
  case *IfExpression:
    return evalIfExpression(node)
  case *InfixExpression:
    left := Eval(node.Left)
    right := Eval(node.Right)
    return evalInfixExpression(node.Operator, left, right)
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

func evalInfixExpression(operator string,
  left, right Object,
) Object {
  switch {
  case left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ:
    return evalIntegerInfixExpression(operator, left, right)
  case operator == "==":
    return nativeBoolToBooleanObject(left == right)
  case operator == "!=":
    return nativeBoolToBooleanObject(left != right)
  default:
    return NULL_LIT
  }
}

func evalIntegerInfixExpression(operator string,
  left, right Object,
) Object {
  leftVal := left.(*Integer).Value
  rightVal := right.(*Integer).Value
  switch operator {
  case "+":
    return &Integer{Value: leftVal + rightVal}
  case "-":
    return &Integer{Value: leftVal - rightVal}
  case "*":
    return &Integer{Value: leftVal * rightVal}
  case "/":
    return &Integer{Value: leftVal / rightVal}
  case "==":
    return nativeBoolToBooleanObject(leftVal == rightVal)
  case "!=":
    return nativeBoolToBooleanObject(leftVal != rightVal)
  case "<":
    return nativeBoolToBooleanObject(leftVal < rightVal)
  case ">":
    return nativeBoolToBooleanObject(leftVal > rightVal)
  default:
    return NULL_LIT
  }
}

func evalIfExpression(ie *IfExpression) Object {
  condition := Eval(ie.Condition)

  if isTruthy(condition) {
    return Eval(ie.Consequence)
  } else if ie.Alternative != nil {
    return Eval(ie.Alternative)
  } else {
    return NULL_LIT
  }
}

func isTruthy(obj Object) bool {
  switch obj {
  case NULL_LIT:
    return false
  case TRUE_LIT:
    return true
  case FALSE_LIT:
    return false
  default:
    return true
  }
}
