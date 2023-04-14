package main

import (
  "testing"
)

func TestEvalIntegerExpression(t *testing.T) {
  tests := []struct {
    input    string
    expected int64
  }{
    {"5", 5},
    {"10", 10},
    {"9000", 9000},
    {"-5", -5},
    {"-10", -10},
    {"5 + 5 + 5 + 5 - 10", 10},
    {"2 * 2 * 2 * 2 * 2", 32},
    {"-50 + 100 + -50", 0},
    {"5 * 2 + 10", 20},
    {"5 + 2 * 10", 25},
    {"20 + 2 * -10", 0},
    {"50 / 2 * 2 + 10", 60},
    {"2 * (5 + 10)", 30},
    {"3 * 3 * 3 + 10", 37},
    {"3 * (3 * 3) + 10", 37},
    {"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
  }

  for _, tt := range tests {
    testIntegerObject(t, testEval(tt.input), tt.expected)
  }
}

func TestEvalBooleanExpression(t *testing.T) {
  tests := []struct {
    input    string
    expected bool
  }{
    {"true", true},
    {"false", false},
  }

  for _, tt := range tests {
    testBooleanObject(t, testEval(tt.input), tt.expected)
  }
}

func TestBangOperator(t *testing.T) {
  tests := []struct {
    input    string
    expected bool
  }{
    {"!true", false},
    {"!false", true},
    {"!5", false},
    {"!!true", true},
    {"!!false", false},
    {"!!5", true},
    {"true", true},
    {"false", false},
    {"1 < 2", true},
    {"1 > 2", false},
    {"1 < 1", false},
    {"1 > 1", false},
    {"1 == 1", true},
    {"1 != 1", false},
    {"1 == 2", false},
    {"1 != 2", true},
    {"true == true", true},
    {"false == false", true},
    {"true == false", false},
    {"true != false", true},
    {"false != true", true},
    {"(1 < 2) == true", true},
    {"(1 < 2) == false", false},
    {"(1 > 2) == true", false},
    {"(1 > 2) == false", true},
  }

  for _, tt := range tests {
    testBooleanObject(t, testEval(tt.input), tt.expected)
  }
}

func TestIfElseExpressions(t *testing.T) {
  tests := []struct {
    input    string
    expected interface{}
  }{
    {"if (true) { 10 }", 10},
    {"if (false) { 10 }", nil},
    {"if (1) { 10 }", 10},
    {"if (1 < 2) { 10 }", 10},
    {"if (1 > 2) { 10 }", nil},
    {"if (1 > 2) { 10 } else { 20 }", 20},
    {"if (1 < 2) { 10 } else { 20 }", 10},
  }

  for _, tt := range tests {
    evaluated := testEval(tt.input)

    integer, ok := tt.expected.(int)

    if ok {
      testIntegerObject(t, evaluated, int64(integer))
    } else {
      testNullObject(t, evaluated)
    }
  }
}

func testEval(input string) Object {
  return Eval(NewParser(NewLexer(input)).Parse())
}

func testIntegerObject(t *testing.T, obj Object, expected int64) bool {
  result, ok := obj.(*Integer)

  if !ok {
    t.Errorf("Object is not an Integer, got=%T (%+v)", obj, obj)
    return false
  }

  if result.Value != expected {
    t.Errorf(
      "Object has wrong value, got=%d, want=%d",
      result.Value,
      expected,
    )
    return false
  }

  return true
}

func testBooleanObject(t *testing.T, obj Object, expected bool) bool {
  result, ok := obj.(*Boolean)

  if !ok {
    t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
    return false
  }

  if result.Value != expected {
    t.Errorf(
      "object has wrong value. got=%t, want=%t",
      result.Value,
      expected,
    )
    return false
  }

  return true
}

func testNullObject(t *testing.T, obj Object) bool {
  if obj != NULL_LIT {
    t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
    return false
  }

  return true
}
