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
