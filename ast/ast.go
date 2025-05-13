package ast

import (
	"bytes"
	"js-compiler/token"
	"strings"
)

type Node interface {
	TokenLiteral() string // usado apenas para debugging e tests
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

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

// Program eh o root de AST produzida pelo parser
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

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(rs.TokenLiteral() + " ; ")
	return out.String()
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

func (i *Identifier) String() string { return i.Value }

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
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

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
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
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
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

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
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

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
	Name       *Identifier
}

func (fl *FunctionLiteral) expressionNode()      {}
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

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, p := range ce.Arguments {
		args = append(args, p.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}

	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ""))
	out.WriteString("}")
	return out.String()
}

// // Node represents a node in the AST
// type Node interface {
// 	TokenLiteral() string
// 	String() string
// }
//
// // Statement represents a statement in the program
// type Statement interface {
// 	Node
// 	statementNode()
// }
//
// // Expression represents an expression in the program
// type Expression interface {
// 	Node
// 	expressionNode()
// }
//
// // Program represents the entire program
// type Program struct {
// 	Statements []Statement
// }
//
// func (p *Program) TokenLiteral() string {
// 	if len(p.Statements) > 0 {
// 		return p.Statements[0].TokenLiteral()
// 	}
// 	return ""
// }
//
// func (p *Program) String() string {
// 	var out bytes.Buffer
//
// 	for _, s := range p.Statements {
// 		out.WriteString(s.String())
// 	}
//
// 	return out.String()
// }
//
// // VarStatement represents a variable declaration
// type VarStatement struct {
// 	Token token.Token // the VAR token
// 	Name  *Identifier
// 	Value Expression
// }
//
// func (vs *VarStatement) statementNode()       {}
// func (vs *VarStatement) TokenLiteral() string { return vs.Token.Literal }
// func (vs *VarStatement) String() string {
// 	var out bytes.Buffer
//
// 	out.WriteString(vs.TokenLiteral() + " ")
// 	out.WriteString(vs.Name.String())
// 	out.WriteString(" = ")
//
// 	if vs.Value != nil {
// 		out.WriteString(vs.Value.String())
// 	}
//
// 	out.WriteString(";")
//
// 	return out.String()
// }
//
// // ReturnStatement represents a return statement
// type ReturnStatement struct {
// 	Token       token.Token // the RETURN token
// 	ReturnValue Expression
// }
//
// func (rs *ReturnStatement) statementNode()       {}
// func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
// func (rs *ReturnStatement) String() string {
// 	var out bytes.Buffer
//
// 	out.WriteString(rs.TokenLiteral() + " ")
//
// 	if rs.ReturnValue != nil {
// 		out.WriteString(rs.ReturnValue.String())
// 	}
//
// 	out.WriteString(";")
//
// 	return out.String()
// }
//
// // ExpressionStatement represents an expression statement
// type ExpressionStatement struct {
// 	Token      token.Token // the first token of the expression
// 	Expression Expression
// }
//
// func (es *ExpressionStatement) statementNode()       {}
// func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
// func (es *ExpressionStatement) String() string {
// 	if es.Expression != nil {
// 		return es.Expression.String()
// 	}
// 	return ""
// }
//
// // BlockStatement represents a block of statements
// type BlockStatement struct {
// 	Token      token.Token // the { token
// 	Statements []Statement
// }
//
// func (bs *BlockStatement) statementNode()       {}
// func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
// func (bs *BlockStatement) String() string {
// 	var out bytes.Buffer
//
// 	for _, s := range bs.Statements {
// 		out.WriteString(s.String())
// 	}
//
// 	return out.String()
// }
//
// // IfStatement represents an if statement
// type IfExpression struct {
// 	Token       token.Token
// 	Condition   Expression
// 	Consequence *BlockStatement
// 	Alternative *BlockStatement
// }
//
// func (ie *IfExpression) expressionNode()      {}
// func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
// func (ie *IfExpression) String() string {
// 	var out bytes.Buffer
//
// 	out.WriteString("if")
// 	out.WriteString(ie.Condition.String())
// 	out.WriteString(" ")
// 	out.WriteString(ie.Consequence.String())
//
// 	if ie.Alternative != nil {
// 		out.WriteString("else ")
// 		out.WriteString(ie.Alternative.String())
// 	}
// 	return out.String()
// }
//
// // WhileStatement represents a while loop
// type WhileStatement struct {
// 	Token     token.Token // the WHILE token
// 	Condition Expression
// 	Body      *BlockStatement
// }
//
// func (ws *WhileStatement) statementNode()       {}
// func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Literal }
// func (ws *WhileStatement) String() string {
// 	var out bytes.Buffer
//
// 	out.WriteString("while")
// 	out.WriteString("(")
// 	out.WriteString(ws.Condition.String())
// 	out.WriteString(")")
// 	out.WriteString(" ")
// 	out.WriteString(ws.Body.String())
//
// 	return out.String()
// }
//
// // PrintStatement represents a print statement
// // TODO: adapt to js console.log
// type ConsoleLog struct {
// 	Token token.Token // the PRINT token
// 	Value Expression
// }
//
// func (ps *ConsoleLog) statementNode()       {}
// func (ps *ConsoleLog) TokenLiteral() string { return ps.Token.Literal }
// func (ps *ConsoleLog) String() string {
// 	var out bytes.Buffer
//
// 	out.WriteString(ps.TokenLiteral() + " ")
// 	out.WriteString(ps.Value.String())
// 	out.WriteString(";")
//
// 	return out.String()
// }
//
// // Identifier represents an identifier
// type Identifier struct {
// 	Token token.Token // the IDENT token
// 	Value string
// }
//
// func (i *Identifier) expressionNode()      {}
// func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
// func (i *Identifier) String() string       { return i.Value }
//
// // IntegerLiteral represents an integer literal
// type IntegerLiteral struct {
// 	Token token.Token // the INT token
// 	Value int64
// }
//
// func (il *IntegerLiteral) expressionNode()      {}
// func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
// func (il *IntegerLiteral) String() string       { return il.Token.Literal }
//
// // StringLiteral represents a string literal
// type StringLiteral struct {
// 	Token token.Token // the STRING token
// 	Value string
// }
//
// func (sl *StringLiteral) expressionNode()      {}
// func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
// func (sl *StringLiteral) String() string       { return "\"" + sl.Value + "\"" }
//
// // BooleanLiteral represents a boolean literal
// type BooleanLiteral struct {
// 	Token token.Token // the true/false token
// 	Value bool
// }
//
// func (bl *BooleanLiteral) expressionNode()      {}
// func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
// func (bl *BooleanLiteral) String() string       { return bl.Token.Literal }
//
// // PrefixExpression represents a prefix operator expression
// type PrefixExpression struct {
// 	Token    token.Token // the prefix token, e.g. !
// 	Operator string
// 	Right    Expression
// }
//
// func (pe *PrefixExpression) expressionNode()      {}
// func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
// func (pe *PrefixExpression) String() string {
// 	var out bytes.Buffer
//
// 	out.WriteString("(")
// 	out.WriteString(pe.Operator)
// 	out.WriteString(pe.Right.String())
// 	out.WriteString(")")
//
// 	return out.String()
// }
//
// // InfixExpression represents an infix operator expression
// type InfixExpression struct {
// 	Token    token.Token // the operator token, e.g. +
// 	Left     Expression
// 	Operator string
// 	Right    Expression
// }
//
// func (ie *InfixExpression) expressionNode()      {}
// func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
// func (ie *InfixExpression) String() string {
// 	var out bytes.Buffer
//
// 	out.WriteString("(")
// 	out.WriteString(ie.Left.String())
// 	out.WriteString(" " + ie.Operator + " ")
// 	out.WriteString(ie.Right.String())
// 	out.WriteString(")")
//
// 	return out.String()
// }
//
// // FunctionLiteral represents a function literal
// type FunctionLiteral struct {
// 	Token      token.Token // the 'fn' token
// 	Parameters []*Identifier
// 	Body       *BlockStatement
// 	Name       string
// }
//
// func (fl *FunctionLiteral) expressionNode()      {}
// func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
// func (fl *FunctionLiteral) String() string {
// 	var out bytes.Buffer
//
// 	params := []string{}
// 	for _, p := range fl.Parameters {
// 		params = append(params, p.String())
// 	}
//
// 	out.WriteString(fl.TokenLiteral())
// 	if fl.Name != "" {
// 		out.WriteString(" " + fl.Name)
// 	}
// 	out.WriteString("(")
// 	out.WriteString(strings.Join(params, ", "))
// 	out.WriteString(") ")
// 	out.WriteString(fl.Body.String())
//
// 	return out.String()
// }
//
// // CallExpression represents a function call
// type CallExpression struct {
// 	Token     token.Token // the '(' token
// 	Function  Expression  // Identifier or FunctionLiteral
// 	Arguments []Expression
// }
//
// func (ce *CallExpression) expressionNode()      {}
// func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
// func (ce *CallExpression) String() string {
// 	var out bytes.Buffer
//
// 	args := []string{}
// 	for _, a := range ce.Arguments {
// 		args = append(args, a.String())
// 	}
//
// 	out.WriteString(ce.Function.String())
// 	out.WriteString("(")
// 	out.WriteString(strings.Join(args, ", "))
// 	out.WriteString(")")
//
// 	return out.String()
// }
//
// // Print the AST for debugging
// func ConsoleLogAST(node Node, indent string) {
// 	switch node := node.(type) {
// 	case *Program:
// 		fmt.Printf("%sProgram\n", indent)
// 		for _, stmt := range node.Statements {
// 			ConsoleLogAST(stmt, indent+"  ")
// 		}
// 	case *VarStatement:
// 		fmt.Printf("%sVarStatement: %s = \n", indent, node.Name.Value)
// 		ConsoleLogAST(node.Value, indent+"  ")
// 	case *ReturnStatement:
// 		fmt.Printf("%sReturnStatement:\n", indent)
// 		ConsoleLogAST(node.ReturnValue, indent+"  ")
// 	case *ExpressionStatement:
// 		fmt.Printf("%sExpressionStatement:\n", indent)
// 		ConsoleLogAST(node.Expression, indent+"  ")
// 	case *BlockStatement:
// 		fmt.Printf("%sBlockStatement:\n", indent)
// 		for _, stmt := range node.Statements {
// 			ConsoleLogAST(stmt, indent+"  ")
// 		}
// 	case *IfExpression:
// 		fmt.Printf("%sIfStatement:\n", indent)
// 		fmt.Printf("%sCondition:\n", indent)
// 		ConsoleLogAST(node.Condition, indent+"  ")
// 		fmt.Printf("%sConsequence:\n", indent)
// 		ConsoleLogAST(node.Consequence, indent+"  ")
// 		if node.Alternative != nil {
// 			fmt.Printf("%sAlternative:\n", indent)
// 			ConsoleLogAST(node.Alternative, indent+"  ")
// 		}
// 	case *WhileStatement:
// 		fmt.Printf("%sWhileStatement:\n", indent)
// 		fmt.Printf("%sCondition:\n", indent)
// 		ConsoleLogAST(node.Condition, indent+"  ")
// 		fmt.Printf("%sBody:\n", indent)
// 		ConsoleLogAST(node.Body, indent+"  ")
// 	case *ConsoleLog:
// 		fmt.Printf("%sConsoleLog:\n", indent)
// 		ConsoleLogAST(node.Value, indent+"  ")
// 	case *Identifier:
// 		fmt.Printf("%sIdentifier: %s\n", indent, node.Value)
// 	case *IntegerLiteral:
// 		fmt.Printf("%sIntegerLiteral: %d\n", indent, node.Value)
// 	case *StringLiteral:
// 		fmt.Printf("%sStringLiteral: %s\n", indent, node.Value)
// 	case *BooleanLiteral:
// 		fmt.Printf("%sBooleanLiteral: %t\n", indent, node.Value)
// 	case *PrefixExpression:
// 		fmt.Printf("%sPrefixExpression: %s\n", indent, node.Operator)
// 		ConsoleLogAST(node.Right, indent+"  ")
// 	case *InfixExpression:
// 		fmt.Printf("%sInfixExpression: %s\n", indent, node.Operator)
// 		fmt.Printf("%sLeft:\n", indent)
// 		ConsoleLogAST(node.Left, indent+"  ")
// 		fmt.Printf("%sRight:\n", indent)
// 		ConsoleLogAST(node.Right, indent+"  ")
// 	case *FunctionLiteral:
// 		fmt.Printf("%sFunctionLiteral: %s\n", indent, node.Name)
// 		fmt.Printf("%sParameters:\n", indent)
// 		for _, param := range node.Parameters {
// 			ConsoleLogAST(param, indent+"  ")
// 		}
// 		fmt.Printf("%sBody:\n", indent)
// 		ConsoleLogAST(node.Body, indent+"  ")
// 	case *CallExpression:
// 		fmt.Printf("%sCallExpression:\n", indent)
// 		fmt.Printf("%sFunction:\n", indent)
// 		ConsoleLogAST(node.Function, indent+"  ")
// 		fmt.Printf("%sArguments:\n", indent)
// 		for _, arg := range node.Arguments {
// 			ConsoleLogAST(arg, indent+"  ")
// 		}
// 	}
// }
//
// type LetStatement struct {
// 	Token token.Token
// 	Name  *Identifier
// 	Value Expression
// }
//
// func (ls *LetStatement) statementNode()       {}
// func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
// func (ls *LetStatement) String() string {
// 	var out bytes.Buffer
//
// 	out.WriteString(ls.TokenLiteral() + " ")
// 	out.WriteString(ls.Name.String())
// 	out.WriteString(" = ")
//
// 	if ls.Value != nil {
// 		out.WriteString(ls.Value.String())
// 	}
// 	out.WriteString(";")
// 	return out.String()
// }
//
// // AssignmentExpression represents an assignment
// type AssignmentExpression struct {
// 	Token token.Token // the = token
// 	Name  *Identifier
// 	Value Expression
// }
//
// func (ae *AssignmentExpression) expressionNode()      {}
// func (ae *AssignmentExpression) TokenLiteral() string { return ae.Token.Literal }
// func (ae *AssignmentExpression) String() string {
// 	var out bytes.Buffer
// 	out.WriteString(ae.Name.String())
// 	out.WriteString(" = ")
// 	out.WriteString(ae.Value.String())
// 	return out.String()
// }
//
// // EmptyExpression represents an empty expression
// type EmptyExpression struct {
// 	Token token.Token
// }
//
// func (ee *EmptyExpression) expressionNode()      {}
// func (ee *EmptyExpression) TokenLiteral() string { return ee.Token.Literal }
// func (ee *EmptyExpression) String() string       { return "" }
//
// type ArrayLiteral struct {
// 	Token    token.Token
// 	Elements []Expression
// }
//
// func (al *ArrayLiteral) expressionNode()      {}
// func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
// func (al *ArrayLiteral) String() string {
// 	var out bytes.Buffer
//
// 	elements := []string{}
//
// 	for _, el := range al.Elements {
// 		elements = append(elements, el.String())
// 	}
//
// 	out.WriteString("[")
// 	out.WriteString(strings.Join(elements, ", "))
// 	out.WriteString("]")
// 	return out.String()
// }
//
// type IndexExpression struct {
// 	Token token.Token
// 	Left  Expression
// 	Index Expression
// }
//
// func (ie *IndexExpression) expressionNode()      {}
// func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
// func (ie *IndexExpression) String() string {
// 	var out bytes.Buffer
//
// 	out.WriteString("(")
// 	out.WriteString(ie.Left.String())
// 	out.WriteString("[")
// 	out.WriteString(ie.Index.String())
// 	out.WriteString("])")
// 	return out.String()
// }
//
// type HashLiteral struct {
// 	Token token.Token
// 	Pairs map[Expression]Expression
// }
//
// func (hl *HashLiteral) expressionNode()      {}
// func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
// func (hl *HashLiteral) String() string {
// 	var out bytes.Buffer
//
// 	pairs := []string{}
// 	for key, value := range hl.Pairs {
// 		pairs = append(pairs, key.String()+":"+value.String())
// 	}
//
// 	out.WriteString("{")
// 	out.WriteString(strings.Join(pairs, ""))
// 	out.WriteString("}")
// 	return out.String()
// }
