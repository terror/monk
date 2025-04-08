package main

import (
  "fmt"
)

var (
  NULL_LIT  = &Null{}
  TRUE_LIT  = &Boolean{Value: true}
  FALSE_LIT = &Boolean{Value: false}
)

func Eval(node Node, env *Environment) Object {
  switch node := node.(type) {
  case *Program:
    return evalProgram(node, env)
  case *BlockStatement:
    return evalBlockStatement(node, env)
  case *BooleanExpression:
    return nativeBoolToBooleanObject(node.Value)
  case *ExpressionStatement:
    return Eval(node.Expression, env)
  case *IfExpression:
    return evalIfExpression(node, env)
  case *InfixExpression:
    left := Eval(node.Left, env)
    if isError(left) {
      return left
    }
    right := Eval(node.Right, env)
    if isError(right) {
      return right
    }
    return evalInfixExpression(node.Operator, left, right)
  case *IntegerLiteral:
    return &Integer{Value: node.Value}
  case *PrefixExpression:
    right := Eval(node.Right, env)
    if isError(right) {
      return right
    }
    return evalPrefixExpression(node.Operator, right)
  case *Identifier:
    return evalIdentifier(node, env)
  case *LetStatement:
    val := Eval(node.Value, env)
    if isError(val) {
      return val
    }
    env.Set(node.Name.Value, val)
    return NULL_LIT
  case *ReturnStatement:
    val := Eval(node.ReturnValue, env)
    if isError(val) {
      return val
    }
    return &ReturnValue{Value: val}
  case *FunctionLiteral:
    params := node.Parameters
    body := node.Body
    return &Function{Parameters: params, Body: body, Env: env}
  case *CallExpression:
    function := Eval(node.Function, env)
    if isError(function) {
      return function
    }
    args := evalExpressions(node.Arguments, env)
    if len(args) == 1 && isError(args[0]) {
      return args[0]
    }
    return applyFunction(function, args)
  }

  return nil
}

func evalProgram(program *Program, env *Environment) Object {
  var result Object

  for _, statement := range program.Statements {
    result = Eval(statement, env)

    switch result := result.(type) {
    case *ReturnValue:
      return result.Value
    case *Error:
      return result
    }
  }

  return result
}

func evalBlockStatement(block *BlockStatement, env *Environment) Object {
  var result Object

  for _, statement := range block.Statements {
    result = Eval(statement, env)

    if result != nil {
      rt := result.Type()
      if rt == RETURN_VALUE_OBJ || rt == ERROR_OBJ {
        return result
      }
    }
  }

  return result
}

func evalStatements(statements []Statement, env *Environment) Object {
  var result Object

  for _, statement := range statements {
    result = Eval(statement, env)

    if returnValue, ok := result.(*ReturnValue); ok {
      return returnValue.Value
    }
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
    return newError("unknown operator: %s%s", operator, right.Type())
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
    return newError("unknown operator: -%s", right.Type())
  }

  value := right.(*Integer).Value

  return &Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right Object) Object {
  switch {
  case left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ:
    return evalIntegerInfixExpression(operator, left, right)
  case operator == "==":
    return nativeBoolToBooleanObject(left == right)
  case operator == "!=":
    return nativeBoolToBooleanObject(left != right)
  case left.Type() != right.Type():
    return newError("type mismatch: %s %s %s",
      left.Type(), operator, right.Type())
  default:
    return newError("unknown operator: %s %s %s",
      left.Type(), operator, right.Type())
  }
}

func evalIntegerInfixExpression(operator string, left, right Object) Object {
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
    return newError("unknown operator: %s %s %s",
      left.Type(), operator, right.Type())
  }
}

func evalIfExpression(ie *IfExpression, env *Environment) Object {
  condition := Eval(ie.Condition, env)
  if isError(condition) {
    return condition
  }

  if isTruthy(condition) {
    return Eval(ie.Consequence, env)
  } else if ie.Alternative != nil {
    return Eval(ie.Alternative, env)
  } else {
    return NULL_LIT
  }
}

func evalIdentifier(node *Identifier, env *Environment) Object {
  val, ok := env.Get(node.Value)
  if !ok {
    return newError("identifier not found: %s", node.Value)
  }
  return val
}

func evalExpressions(exps []Expression, env *Environment) []Object {
  var result []Object

  for _, e := range exps {
    evaluated := Eval(e, env)
    if isError(evaluated) {
      return []Object{evaluated}
    }
    result = append(result, evaluated)
  }

  return result
}

func applyFunction(fn Object, args []Object) Object {
  function, ok := fn.(*Function)
  if !ok {
    return newError("not a function: %s", fn.Type())
  }

  extendedEnv := extendFunctionEnv(function, args)
  evaluated := Eval(function.Body, extendedEnv)
  return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(fn *Function, args []Object) *Environment {
  env := NewEnclosedEnvironment(fn.Env)

  for paramIdx, param := range fn.Parameters {
    env.Set(param.Value, args[paramIdx])
  }

  return env
}

func unwrapReturnValue(obj Object) Object {
  if returnValue, ok := obj.(*ReturnValue); ok {
    return returnValue.Value
  }

  return obj
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

func newError(format string, a ...interface{}) *Error {
  return &Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj Object) bool {
  if obj != nil {
    return obj.Type() == ERROR_OBJ
  }
  return false
}
