package main

var (
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
