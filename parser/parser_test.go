package parser_test

import (
	"fmt"
	"testing"

	"monkey/ast"
	"monkey/lexer"
	"monkey/parser"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// help method

func testProgramStatementCount(t *testing.T, program *ast.Program, count int) {
	require.Equal(t, count, len(program.Statements), "program.Statements count wrong!")
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) {
	require.Equal(t, "let", stmt.TokenLiteral())

	require.IsType(t, new(ast.LetStatement), stmt)

	letStmt, ok := stmt.(*ast.LetStatement)
	require.Equal(t, true, ok, "not *ast.LetStatement")

	require.Equal(t, name, letStmt.Name.Value)
	require.Equal(t, name, letStmt.Name.TokenLiteral())
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) {
	switch v := expected.(type) {
	case int:
		testIntegerLiteral(t, exp, int64(v))
	case int64:
		testIntegerLiteral(t, exp, v)
	case string:
		testIdentifier(t, exp, v)
	case bool:
		testBooleanLiteral(t, exp, v)
	default:
		require.Failf(t, "type of exp not handled.", "got=%T", exp)
	}
}

func testInfixExpression(
	t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{},
) {
	infix, ok := exp.(*ast.InfixExpression)
	require.True(t, ok, "not *ast.InfixExpression type")
	testLiteralExpression(t, infix.Left, left)
	require.Equal(t, operator, infix.Operator)
	testLiteralExpression(t, infix.Right, right)
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, value int64) {
	integ, ok := exp.(*ast.IntegerLiteral)
	require.True(t, ok, "not *ast.IntegerLiteral type")
	require.Equal(t, value, integ.Value)
	require.Equal(t, fmt.Sprintf("%d", value), integ.TokenLiteral())
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) {
	ident, ok := exp.(*ast.Identifier)
	require.True(t, ok, "not *ast.Identifier type")
	require.Equal(t, value, ident.Value)
	require.Equal(t, value, ident.TokenLiteral())
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) {
	bo, ok := exp.(*ast.Boolean)
	require.True(t, ok, "exp not *ast.Boolean type")
	require.Equal(t, value, bo.Value)
	require.Equal(t, fmt.Sprintf("%t", value), bo.TokenLiteral())
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	// check errors
	errs := p.Errors()
	for _, err := range errs {
		assert.Fail(t, err)
	}
	require.Equal(t, 0, len(errs))
}

// Test
func TestLetStatement(t *testing.T) {
	input := `
        let x = 5;
        let y = 10;
        let foobar = 838383; 
    `

	p := parser.New(lexer.New(input))
	program := p.ParseProgram()

	checkParserErrors(t, p)
	require.NotNil(t, program, "ParseProgram() returend nil")
	testProgramStatementCount(t, program, 3)

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		testLetStatement(t, stmt, tt.expectedIdentifier)
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		p := parser.New(lexer.New(tt.input))
		program := p.ParseProgram()
		checkParserErrors(t, p)

		testProgramStatementCount(t, program, 1)

		testLetStatement(t, program.Statements[0], tt.expectedIdentifier)

		letStmt, ok := program.Statements[0].(*ast.LetStatement)
		require.True(t, ok, "program.Statements[0] is not ast.LetStatement")
		testLiteralExpression(t, letStmt.Value, tt.expectedValue)
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
        return 5;
        return 10;
        return 838383; 
    `

	p := parser.New(lexer.New(input))
	program := p.ParseProgram()

	checkParserErrors(t, p)
	require.NotNil(t, program, "ParseProgram() returend nil")
	testProgramStatementCount(t, program, 3)

	for _, stmt := range program.Statements {
		assert.IsType(t, new(ast.ReturnStatement), stmt)
		assert.Equal(t, "return", stmt.TokenLiteral())
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		p := parser.New(lexer.New(tt.input))
		program := p.ParseProgram()
		checkParserErrors(t, p)

		testProgramStatementCount(t, program, 1)

		stmt, ok := program.Statements[0].(*ast.ReturnStatement)
		require.True(t, ok, "program.Statements[0] is not ast.ReturnStatement")

		require.Equal(t, "return", stmt.TokenLiteral())
		testLiteralExpression(t, stmt.ReturnValue, tt.expectedValue)
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	p := parser.New(lexer.New(input))

	program := p.ParseProgram()

	checkParserErrors(t, p)

	testProgramStatementCount(t, program, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement")

	ident, ok := stmt.Expression.(*ast.Identifier)
	require.True(t, ok, "expression not *ast.Identifier")

	assert.Equal(t, "foobar", ident.Value)
	assert.Equal(t, "foobar", ident.TokenLiteral())
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	p := parser.New(lexer.New(input))

	program := p.ParseProgram()

	checkParserErrors(t, p)
	testProgramStatementCount(t, program, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement")

	testIntegerLiteral(t, stmt.Expression, 5)
}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue interface{}
	}{
		{"!5;", "!", 5},
		{"-15", "-", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		p := parser.New(lexer.New(tt.input))
		program := p.ParseProgram()
		checkParserErrors(t, p)

		testProgramStatementCount(t, program, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement")

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		require.True(t, ok, "expression not *ast.PrefixExpression")

		require.Equal(t, tt.operator, exp.Operator)

		testLiteralExpression(t, exp.Right, tt.integerValue)
	}
}

func TestParsingInfixExpression(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range tests {
		p := parser.New(lexer.New(tt.input))
		program := p.ParseProgram()
		checkParserErrors(t, p)

		testProgramStatementCount(t, program, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement")

		testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func TestOperatorPrecedneceExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}

	for i, tt := range tests {
		p := parser.New(lexer.New(tt.input))
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()

		require.Equal(t, tt.expected, actual, "case index %d", i)
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		p := parser.New(lexer.New(tt.input))
		program := p.ParseProgram()
		checkParserErrors(t, p)

		testProgramStatementCount(t, program, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement")

		boolean, ok := stmt.Expression.(*ast.Boolean)
		require.True(t, ok, "exp not *ast.Boolean type")
		require.Equal(t, tt.expectedBoolean, boolean.Value)
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)

	testProgramStatementCount(t, program, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement")

	exp, ok := stmt.Expression.(*ast.IfExpression)
	require.True(t, ok, "exp not *ast.IfExpression type")

	testInfixExpression(t, exp.Condition, "x", "<", "y")

	require.Equal(t, 1, len(exp.Consequence.Statements), "consequence is not 1 statements")

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "consequence.Statements[0] is not ast.ExpressionStatement")
	testIdentifier(t, consequence.Expression, "x")
	require.Nil(t, exp.Alternative, "exp.Alternative.Statements was not nil")
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)

	testProgramStatementCount(t, program, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement")

	exp, ok := stmt.Expression.(*ast.IfExpression)
	require.True(t, ok, "exp not *ast.IfExpression type")

	testInfixExpression(t, exp.Condition, "x", "<", "y")

	// consequence
	require.Equal(t, 1, len(exp.Consequence.Statements), "consequence.Statements count wrong.")
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "consequence.Statements[0] is not ast.ExpressionStatement")
	testIdentifier(t, consequence.Expression, "x")

	// alternative
	require.Equal(t, 1, len(exp.Alternative.Statements), "exp.Alternative.Statements count wrong.")
	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "alternative.Statements[0] is not ast.ExpressionStatement")
	testIdentifier(t, alternative.Expression, "y")
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)

	testProgramStatementCount(t, program, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement")

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	require.True(t, ok, "exp not *ast.FunctionIteral type")
	require.Equal(t, 2, len(function.Parameters), "function literal parameters wrong!")

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	require.Equal(t, 1, len(function.Body.Statements), "function.Body.Statements count wrong!")

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "function.Body not *ast.ExpressionStatement type")

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		p := parser.New(lexer.New(tt.input))
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement")

		function, ok := stmt.Expression.(*ast.FunctionLiteral)
		require.True(t, ok, "exp not *ast.FunctionIteral type")
		require.Equal(t, len(tt.expectedParams), len(function.Parameters), "function literal parameters wrong!")
		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"
	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)

	testProgramStatementCount(t, program, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement.")

	call, ok := stmt.Expression.(*ast.CallExpression)
	require.True(t, ok, "exp not *ast.CallExpression type")
	testIdentifier(t, call.Function, "add")
	require.Equal(t, 3, len(call.Arguments), "call.Arguments count wrong.")

	testLiteralExpression(t, call.Arguments[0], 1)
	testInfixExpression(t, call.Arguments[1], 2, "*", 3)
	testInfixExpression(t, call.Arguments[2], 4, "+", 5)
}
