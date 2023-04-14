package main

func Eval(node Node) Object {
  switch node := node.(type) {
  case *BooleanExpression:
    return &Boolean{Value: node.Value}
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
