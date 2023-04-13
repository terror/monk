package main

import (
  "bytes"
  "strings"
)

type Node interface {
  TokenLiteral() string
  String() string
}

type Statement interface {
  Node
  statementNode()
}

type Expression interface {
  Node
  expressionNode()
}

type Program struct {
  Statements []Statement
}

func (p *Program) TokenLiteral() string {
  if len(p.Statements) > 0 {
    return p.Statements[0].TokenLiteral()
  } else {
    return ""
  }
}

func (p *Program) String() string {
  var out bytes.Buffer

  for _, s := range p.Statements {
    out.WriteString(s.String())
  }

  return out.String()
}

type Identifier struct {
  Token Token
  Value string
}

func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

func (i *Identifier) expressionNode() {}

func (i *Identifier) String() string {
  return i.Value
}

type IntegerLiteral struct {
  Token Token
  Value int64
}

func (i *IntegerLiteral) TokenLiteral() string { return i.Token.Literal }

func (i *IntegerLiteral) expressionNode() {}

func (i *IntegerLiteral) String() string {
  return i.Token.Literal
}

type FunctionLiteral struct {
  Token      Token
  Parameters []*Identifier
  Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }

func (fl *FunctionLiteral) String() string {
  var out bytes.Buffer

  params := []string{}

  for _, p := range fl.Parameters {
    params = append(params, p.String())
  }

  out.WriteString(fl.TokenLiteral())
  out.WriteString("(")
  out.WriteString(strings.Join(params, ", "))
  out.WriteString(") ")
  out.WriteString(fl.Body.String())

  return out.String()
}

type LetStatement struct {
  Token Token
  Name  *Identifier
  Value Expression
}

func (ls *LetStatement) TokenLiteral() string {
  return ls.Token.Literal
}

func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) String() string {
  var out bytes.Buffer

  out.WriteString(ls.TokenLiteral() + " ")
  out.WriteString(ls.Name.String())
  out.WriteString(" = ")

  if ls.Value != nil {
    out.WriteString(ls.Value.String())
  }

  out.WriteString(";")

  return out.String()
}

type ReturnStatement struct {
  Token       Token
  ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string {
  return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
  var out bytes.Buffer

  out.WriteString(rs.TokenLiteral() + " ")

  if rs.ReturnValue != nil {
    out.WriteString(rs.ReturnValue.String())
  }

  out.WriteString(";")

  return out.String()
}

type ExpressionStatement struct {
  Token      Token
  Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) TokenLiteral() string {
  return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
  if es.Expression != nil {
    return es.Expression.String()
  } else {
    return ""
  }
}

type BlockStatement struct {
  Token      Token
  Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

func (bs *BlockStatement) String() string {
  var out bytes.Buffer

  for _, s := range bs.Statements {
    out.WriteString(s.String())
  }

  return out.String()
}

type PrefixExpression struct {
  Token    Token
  Operator string
  Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

func (pe *PrefixExpression) String() string {
  var out bytes.Buffer

  out.WriteString("(")
  out.WriteString(pe.Operator)
  out.WriteString(pe.Right.String())
  out.WriteString(")")

  return out.String()
}

type InfixExpression struct {
  Left     Expression
  Token    Token
  Operator string
  Right    Expression
}

func (ie *InfixExpression) expressionNode() {}

func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }

func (ie *InfixExpression) String() string {
  var out bytes.Buffer

  out.WriteString("(")
  out.WriteString(ie.Left.String())
  out.WriteString(" " + ie.Operator + " ")
  out.WriteString(ie.Right.String())
  out.WriteString(")")

  return out.String()
}

type Boolean struct {
  Token Token
  Value bool
}

func (b *Boolean) expressionNode() {}

func (b *Boolean) TokenLiteral() string { return b.Token.Literal }

func (b *Boolean) String() string { return b.Token.Literal }

type IfExpression struct {
  Token       Token
  Condition   Expression
  Consequence *BlockStatement
  Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }

func (ie *IfExpression) String() string {
  var out bytes.Buffer

  out.WriteString("if")
  out.WriteString(ie.Condition.String())
  out.WriteString(" ")
  out.WriteString(ie.Consequence.String())

  if ie.Alternative != nil {
    out.WriteString("else ")
    out.WriteString(ie.Alternative.String())
  }

  return out.String()
}
