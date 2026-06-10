// Package pattern provides a lexer, parser, and validator for the pattern language.
package pattern

import (
	"fmt"
	"strings"

	"github.com/Kyanite/noise/internal/logging"
)

// Node represents a node in the AST
type Node interface {
	node()
	String() string
	Pos() Position
}

// Expression represents an expression in the pattern language
type Expression interface {
	Node
	expressionNode()
}

// Statement represents a statement in the pattern language
type Statement interface {
	Node
	statementNode()
}

// Program represents a complete program
type Program struct {
	Statements []Statement
}

func (p *Program) node() {}
func (p *Program) String() string {
	var out strings.Builder
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

func (p *Program) Pos() Position {
	if len(p.Statements) > 0 {
		return p.Statements[0].Pos()
	}
	return Position{}
}

// PatternStatement represents a pattern definition
type PatternStatement struct {
	Name       string
	Expression Expression
	Position   Position
}

func (ps *PatternStatement) node()          {}
func (ps *PatternStatement) statementNode() {}
func (ps *PatternStatement) String() string {
	return ps.Name + " = " + ps.Expression.String() + ";"
}
func (ps *PatternStatement) Pos() Position { return ps.Position }

// SequenceExpression represents a sequence of values
type SequenceExpression struct {
	Values   []Expression
	Position Position
}

func (se *SequenceExpression) node()           {}
func (se *SequenceExpression) expressionNode() {}
func (se *SequenceExpression) String() string {
	var values []string
	for _, v := range se.Values {
		values = append(values, v.String())
	}
	return "[" + strings.Join(values, ", ") + "]"
}
func (se *SequenceExpression) Pos() Position { return se.Position }

// LiteralExpression represents a literal value
type LiteralExpression struct {
	Value    PatternValue
	Position Position
}

func (le *LiteralExpression) node()           {}
func (le *LiteralExpression) expressionNode() {}
func (le *LiteralExpression) String() string {
	switch v := le.Value.(type) {
	case NumberValue:
		return fmt.Sprintf("%v", v.Value)
	case StringValue:
		return fmt.Sprintf(`"%s"`, v.Value)
	case SampleValue:
		return fmt.Sprintf("<%s>", v.Name)
	case NoteValue:
		return fmt.Sprintf("<%s%s%d>", v.Note, v.Accident, v.Octave)
	case RestValue:
		return "~"
	default:
		return "?"
	}
}
func (le *LiteralExpression) Pos() Position { return le.Position }

// FunctionCallExpression represents a function call
type FunctionCallExpression struct {
	Name      string
	Arguments []Expression
	Position  Position
}

func (fce *FunctionCallExpression) node()           {}
func (fce *FunctionCallExpression) expressionNode() {}
func (fce *FunctionCallExpression) String() string {
	var args []string
	for _, arg := range fce.Arguments {
		args = append(args, arg.String())
	}
	return fce.Name + "(" + strings.Join(args, ", ") + ")"
}
func (fce *FunctionCallExpression) Pos() Position { return fce.Position }

// BinaryExpression represents a binary operation
type BinaryExpression struct {
	Left     Expression
	Operator TokenType
	Right    Expression
	Position Position
}

func (be *BinaryExpression) node()           {}
func (be *BinaryExpression) expressionNode() {}
func (be *BinaryExpression) String() string {
	return "(" + be.Left.String() + " " + be.Operator.String() + " " + be.Right.String() + ")"
}
func (be *BinaryExpression) Pos() Position { return be.Position }

// UnaryExpression represents a unary operation
type UnaryExpression struct {
	Operator TokenType
	Right    Expression
	Position Position
}

func (ue *UnaryExpression) node()           {}
func (ue *UnaryExpression) expressionNode() {}
func (ue *UnaryExpression) String() string {
	return "(" + ue.Operator.String() + ue.Right.String() + ")"
}
func (ue *UnaryExpression) Pos() Position { return ue.Position }

// GroupExpression represents a grouped expression
type GroupExpression struct {
	Expression Expression
	Position   Position
}

func (ge *GroupExpression) node()           {}
func (ge *GroupExpression) expressionNode() {}
func (ge *GroupExpression) String() string {
	return "(" + ge.Expression.String() + ")"
}
func (ge *GroupExpression) Pos() Position { return ge.Position }

// ParameterExpression represents a parameter assignment
type ParameterExpression struct {
	Name     string
	Value    Expression
	Position Position
}

func (pe *ParameterExpression) node()           {}
func (pe *ParameterExpression) expressionNode() {}
func (pe *ParameterExpression) String() string {
	return pe.Name + ":" + pe.Value.String()
}
func (pe *ParameterExpression) Pos() Position { return pe.Position }

// ModifierExpression represents a pattern modifier
type ModifierExpression struct {
	Base      Expression
	Modifier  string
	Arguments []Expression
	Position  Position
}

func (me *ModifierExpression) node()           {}
func (me *ModifierExpression) expressionNode() {}
func (me *ModifierExpression) String() string {
	var args []string
	for _, arg := range me.Arguments {
		args = append(args, arg.String())
	}
	if len(args) > 0 {
		return me.Base.String() + "." + me.Modifier + "(" + strings.Join(args, ", ") + ")"
	}
	return me.Base.String() + "." + me.Modifier
}
func (me *ModifierExpression) Pos() Position { return me.Position }

// TimeExpression represents a time-related expression
type TimeExpression struct {
	Value    Expression
	TimeUnit string
	Position Position
}

func (te *TimeExpression) node()           {}
func (te *TimeExpression) expressionNode() {}
func (te *TimeExpression) String() string {
	return te.Value.String() + te.TimeUnit
}
func (te *TimeExpression) Pos() Position { return te.Position }

// RangeExpression represents a range expression
type RangeExpression struct {
	Start    Expression
	End      Expression
	Step     Expression // optional
	Position Position
}

func (re *RangeExpression) node()           {}
func (re *RangeExpression) expressionNode() {}
func (re *RangeExpression) String() string {
	if re.Step != nil {
		return re.Start.String() + ".." + re.End.String() + ":" + re.Step.String()
	}
	return re.Start.String() + ".." + re.End.String()
}
func (re *RangeExpression) Pos() Position { return re.Position }

// ListExpression represents a list of values
type ListExpression struct {
	Elements []Expression
	Position Position
}

func (le *ListExpression) node()           {}
func (le *ListExpression) expressionNode() {}
func (le *ListExpression) String() string {
	var elements []string
	for _, el := range le.Elements {
		elements = append(elements, el.String())
	}
	return "[" + strings.Join(elements, ", ") + "]"
}
func (le *ListExpression) Pos() Position { return le.Position }

// StructExpression represents a struct/record expression
type StructExpression struct {
	Fields   []*ParameterExpression
	Position Position
}

func (se *StructExpression) node()           {}
func (se *StructExpression) expressionNode() {}
func (se *StructExpression) String() string {
	var fields []string
	for _, field := range se.Fields {
		fields = append(fields, field.String())
	}
	return "{" + strings.Join(fields, ", ") + "}"
}
func (se *StructExpression) Pos() Position { return se.Position }

// IndexExpression represents an index access
type IndexExpression struct {
	Base     Expression
	Index    Expression
	Position Position
}

func (ie *IndexExpression) node()           {}
func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) String() string {
	return ie.Base.String() + "[" + ie.Index.String() + "]"
}
func (ie *IndexExpression) Pos() Position { return ie.Position }

// PropertyExpression represents a property access
type PropertyExpression struct {
	Base     Expression
	Property string
	Position Position
}

func (pe *PropertyExpression) node()           {}
func (pe *PropertyExpression) expressionNode() {}
func (pe *PropertyExpression) String() string {
	return pe.Base.String() + "." + pe.Property
}
func (pe *PropertyExpression) Pos() Position { return pe.Position }

// ASTVisitor defines the interface for visiting AST nodes
type ASTVisitor interface {
	Visit(node Node) interface{}
}

// ASTWalker implements the visitor pattern for traversing the AST
type ASTWalker struct {
	visitor ASTVisitor
}

func NewASTWalker(visitor ASTVisitor) *ASTWalker {
	return &ASTWalker{visitor: visitor}
}

func (w *ASTWalker) Walk(node Node) interface{} {
	if node == nil {
		return nil
	}
	return w.Visit(node)
}

func (w *ASTWalker) Visit(node Node) interface{} {
	switch n := node.(type) {
	case *Program:
		return w.visitProgram(n)
	case *PatternStatement:
		return w.visitPatternStatement(n)
	case *SequenceExpression:
		return w.visitSequenceExpression(n)
	case *LiteralExpression:
		return w.visitLiteralExpression(n)
	case *FunctionCallExpression:
		return w.visitFunctionCallExpression(n)
	case *BinaryExpression:
		return w.visitBinaryExpression(n)
	case *UnaryExpression:
		return w.visitUnaryExpression(n)
	case *GroupExpression:
		return w.visitGroupExpression(n)
	case *ParameterExpression:
		return w.visitParameterExpression(n)
	case *ModifierExpression:
		return w.visitModifierExpression(n)
	case *TimeExpression:
		return w.visitTimeExpression(n)
	case *RangeExpression:
		return w.visitRangeExpression(n)
	case *ListExpression:
		return w.visitListExpression(n)
	case *StructExpression:
		return w.visitStructExpression(n)
	case *IndexExpression:
		return w.visitIndexExpression(n)
	case *PropertyExpression:
		return w.visitPropertyExpression(n)
	default:
		return w.visitor.Visit(node)
	}
}

func (w *ASTWalker) visitProgram(node *Program) interface{} {
	for _, stmt := range node.Statements {
		w.Walk(stmt)
	}
	return w.visitor.Visit(node)
}

func (w *ASTWalker) visitPatternStatement(node *PatternStatement) interface{} {
	w.Walk(node.Expression)
	return w.visitor.Visit(node)
}

func (w *ASTWalker) visitSequenceExpression(node *SequenceExpression) interface{} {
	for _, value := range node.Values {
		w.Walk(value)
	}
	return w.visitor.Visit(node)
}

func (w *ASTWalker) visitLiteralExpression(node *LiteralExpression) interface{} {
	return w.visitor.Visit(node)
}

func (w *ASTWalker) visitFunctionCallExpression(node *FunctionCallExpression) interface{} {
	for _, arg := range node.Arguments {
		w.Walk(arg)
	}
	return w.visitor.Visit(node)
}

func (w *ASTWalker) visitBinaryExpression(node *BinaryExpression) interface{} {
	w.Walk(node.Left)
	w.Walk(node.Right)
	return w.visitor.Visit(node)
}

func (w *ASTWalker) visitUnaryExpression(node *UnaryExpression) interface{} {
	w.Walk(node.Right)
	return w.visitor.Visit(node)
}

func (w *ASTWalker) visitGroupExpression(node *GroupExpression) interface{} {
	w.Walk(node.Expression)
	return w.visitor.Visit(node)
}

func (w *ASTWalker) visitParameterExpression(node *ParameterExpression) interface{} {
	w.Walk(node.Value)
	return w.visitor.Visit(node)
}

func (w *ASTWalker) visitModifierExpression(node *ModifierExpression) interface{} {
	w.Walk(node.Base)
	for _, arg := range node.Arguments {
		w.Walk(arg)
	}
	return w.visitor.Visit(node)
}

func (w *ASTWalker) visitTimeExpression(node *TimeExpression) interface{} {
	w.Walk(node.Value)
	return w.visitor.Visit(node)
}

func (w *ASTWalker) visitRangeExpression(node *RangeExpression) interface{} {
	w.Walk(node.Start)
	w.Walk(node.End)
	if node.Step != nil {
		w.Walk(node.Step)
	}
	return w.visitor.Visit(node)
}

func (w *ASTWalker) visitListExpression(node *ListExpression) interface{} {
	for _, element := range node.Elements {
		w.Walk(element)
	}
	return w.visitor.Visit(node)
}

func (w *ASTWalker) visitStructExpression(node *StructExpression) interface{} {
	for _, field := range node.Fields {
		w.Walk(field)
	}
	return w.visitor.Visit(node)
}

func (w *ASTWalker) visitIndexExpression(node *IndexExpression) interface{} {
	w.Walk(node.Base)
	w.Walk(node.Index)
	return w.visitor.Visit(node)
}

func (w *ASTWalker) visitPropertyExpression(node *PropertyExpression) interface{} {
	w.Walk(node.Base)
	return w.visitor.Visit(node)
}

// ASTPrinter is a visitor that prints the AST
type ASTPrinter struct {
	indent int
}

func NewASTPrinter() *ASTPrinter {
	return &ASTPrinter{indent: 0}
}

func (p *ASTPrinter) Visit(node Node) interface{} {
	if !logging.DebugEnabled() {
		return nil
	}
	logging.Debugf("%s%T", p.indentPrefix(), node)
	p.indent++
	defer func() { p.indent-- }()
	return nil
}

func (p *ASTPrinter) indentPrefix() string {
	if p.indent <= 0 {
		return ""
	}
	builder := strings.Builder{}
	for i := 0; i < p.indent; i++ {
		builder.WriteString("  ")
	}
	return builder.String()
}

// PrintAST prints the AST to stdout
func PrintAST(node Node) {
	if !logging.DebugEnabled() {
		return
	}
	printer := NewASTPrinter()
	walker := NewASTWalker(printer)
	walker.Walk(node)
}
