package codegen

import (
	"fmt"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"js-compiler/ast"
	"os"
	"regexp"
	"strings"
)

type Generator struct {
	module        *ir.Module
	context       *Context
	stringCounter int
	blockCount    int
	consoleLog    *ir.Func
}

type Context struct {
	currentFunction *ir.Func
	blocks          map[string]*ir.Block
	namedValues     map[string]value.Value
	stringConstants map[string]*ir.Global
	parent          *Context
}

type Value struct {
	Value value.Value
	Type  types.Type
}

func New() *Generator {
	module := ir.NewModule()
	context := newContext(nil)

	printFunc := declarePrintf(module)
	declareRuntime(module)

	return &Generator{
		module:        module,
		context:       context,
		stringCounter: 0,
		blockCount:    0,
		consoleLog:    printFunc,
	}
}

func (g *Generator) Generate(program *ast.Program) (*ir.Module, error) {
	// First pass: process function declarations
	fmt.Fprintf(os.Stderr, "DEBUG: Starting Generate method\n")
	fmt.Fprintf(os.Stderr, "DEBUG: Total statements: %d\n", len(program.Statements))

	// Process function declarations first
	for i, stmt := range program.Statements {
		if varStmt, ok := stmt.(*ast.LetStatement); ok {
			if _, ok := varStmt.Value.(*ast.FunctionLiteral); ok {
				_, err := g.generateStatement(stmt)
				if err != nil {
					return nil, fmt.Errorf("error processing function declaration %d: %w", i, err)
				}
			}
		} else if exprStmt, ok := stmt.(*ast.ExpressionStatement); ok {
			if funcLit, ok := exprStmt.Expression.(*ast.FunctionLiteral); ok {
				_, err := g.generateExpression(funcLit)
				if err != nil {
					return nil, fmt.Errorf("error processing function declaration %d: %w", i, err)
				}
			}
		}
	}

	// Create main function
	fmt.Fprintf(os.Stderr, "DEBUG: Creating main function\n")
	mainFunc := g.module.NewFunc("main", types.NewFunc(types.I32))
	mainBlock := mainFunc.NewBlock("entry")
	g.context.currentFunction = mainFunc
	g.context.blocks["entry"] = mainBlock

	fmt.Fprintf(os.Stderr, "DEBUG: Processing statements in main function\n")
	// Process statements in main function
	for i, stmt := range program.Statements {
		if varStmt, ok := stmt.(*ast.LetStatement); ok {
			if _, ok := varStmt.Value.(*ast.FunctionLiteral); ok {
				continue
			}
		} else if exprStmt, ok := stmt.(*ast.ExpressionStatement); ok {
			if _, ok := exprStmt.Expression.(*ast.FunctionLiteral); ok {
				continue
			}
		}

		fmt.Fprintf(os.Stderr, "DEBUG: Generating statement %d: %T\n", i, stmt)
		_, err := g.generateStatement(stmt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Failed to generate statement %d: %v\n", i, err)
			return nil, fmt.Errorf("error processing statement %d: %w", i, err)
		}
	}

	// Ensure main function has a return
	fmt.Fprintf(os.Stderr, "DEBUG: Checking main function terminator\n")
	if mainBlock.Term == nil {
		fmt.Fprintf(os.Stderr, "DEBUG: Adding default return to main function\n")
		mainBlock.NewRet(constant.NewInt(types.I32, 0))
	}

	// Ensure all functions have a return
	fmt.Fprintf(os.Stderr, "DEBUG: Checking terminator for all functions\n")
	for _, fn := range g.module.Funcs {
		for _, block := range fn.Blocks {
			if block.Term == nil {
				if fn.Sig.RetType.Equal(types.Void) {
					block.NewRet(nil)
				} else {
					block.NewRet(constant.NewInt(types.I32, 0))
				}
			}
		}
	}

	fmt.Fprintf(os.Stderr, "DEBUG: Total functions in module: %d\n", len(g.module.Funcs))
	for _, fn := range g.module.Funcs {
		fmt.Fprintf(os.Stderr, "DEBUG: Function %s\n", fn.Name())
	}

	return g.module, nil
}

func CompileToLLVM(program *ast.Program) (string, error) {
	generator := New()
	module, err := generator.Generate(program)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	_, err = module.WriteTo(&buf)
	if err != nil {
		return "", err
	}

	ir := buf.String()
	ir = fixFunctionDeclarations(ir)
	return ir, nil

}

func (g *Generator) debugExpression(expr ast.Expression) {
	if expr == nil {
		fmt.Fprintf(os.Stderr, "NIL EXPRESSION DETECTED\n")
		if g.context != nil && g.context.currentFunction != nil {
			fmt.Fprintf(os.Stderr, "Current function: %s\n", g.context.currentFunction.Name())
		}
	}
}

// Lookup looks up a value in the context
func (c *Context) Lookup(name string) (value.Value, bool) {
	if val, ok := c.namedValues[name]; ok {
		return val, true
	}
	if c.parent != nil {
		return c.parent.Lookup(name)
	}
	return nil, false
}

func (g *Generator) generateStatement(stmt ast.Statement) (value.Value, error) {
	switch stmt := stmt.(type) {
	// case *ast.LetStatement:
	// 	return g.generateVarStatement(stmt)
	// case *ast.ReturnStatement:
	// 	return g.generateReturnStatement(stmt)
	// case *ast.ExpressionStatement:
	// 	if stmt.Expression == nil {
	// 		return nil, fmt.Errorf("nil expression in expression statement")
	// 	}
	// 	return g.generateExpression(stmt.Expression)
	// case *ast.BlockStatement:
	// 	return g.generateBlockStatement(stmt)
	// case *ast.IfExpression:
	// 	return g.generateIfStatement(stmt)
	// case *ast.WhileStatement:
	// 	return g.generateWhileStatement(stmt)
	// case *ast.PrintStatement:
	// 	return g.generatePrintStatement(stmt)
	default:
		return nil, fmt.Errorf("unknown statement type: %T", stmt)
	}
}

func (g *Generator) generateExpression(expr ast.Expression) (value.Value, error) {
	g.debugExpression(expr)
	if expr == nil {
		return nil, fmt.Errorf("nil expression")
	}

	switch expr := expr.(type) {
	case *ast.Identifier:
		return g.generateIdentifier(expr)
	case *ast.IntegerLiteral:
		return g.generateIntegerLiteral(expr)
	// case *ast.StringLiteral:
	// 	return g.generateStringLiteral(expr)
	// case *ast.PrefixExpression:
	// 	return g.generatePrefixExpression(expr)
	// case *ast.InfixExpression:
	// 	return g.generateInfixExpression(expr)
	// case *ast.FunctionLiteral:
	// 	return g.generateFunctionLiteral(expr)
	// case *ast.CallExpression:
	// 	return g.generateCallExpression(expr)
	// case *ast.Boolean:
	// 	return g.generateBooleanLiteral(expr)
	// case *ast.AssignmentExpression:
	// 	return g.generateAssignmentExpression(expr)
	// case *ast.EmptyExpression:
	// 	return constant.NewInt(types.I32, 0), nil
	default:
		return nil, fmt.Errorf("unknown expression type: %T", expr)
	}
}

func (g *Generator) generateIdentifier(expr *ast.Identifier) (value.Value, error) {
	if val, ok := g.context.Lookup(expr.Value); ok {
		if val == nil {
			fmt.Fprintf(os.Stderr, "WARNING: variable %s found in context but has nil value\n", expr.Value)
			return constant.NewInt(types.I32, 0), nil
		}

		if g.context.currentFunction == nil || len(g.context.currentFunction.Blocks) == 0 {
			fmt.Fprintf(os.Stderr, "WARNING: No current function or blocks when generating identifier\n")
			return constant.NewInt(types.I32, 0), nil
		}

		currentBlock := g.context.currentFunction.Blocks[len(g.context.currentFunction.Blocks)-1]
		if pointerType, ok := val.Type().(*types.PointerType); ok {
			if pointerType.ElemType.Equal(types.I32) {
				return constant.NewInt(types.I32, 0), nil
			}

			if _, ok := pointerType.ElemType.(*types.FuncType); ok {
				return val, nil
			}

			if pointerType.ElemType.Equal(types.I32) ||
				pointerType.ElemType.Equal(types.I64) ||
				pointerType.ElemType.Equal(types.I8) ||
				(pointerType.ElemType.String() == "*i8") {
				return currentBlock.NewLoad(pointerType.ElemType, val), nil
			}
		}

	}

	for _, fn := range g.module.Funcs {
		if fn.Name() == expr.Value {
			return fn, nil
		}
	}

	return constant.NewInt(types.I32, 0), nil
}

func (g *Generator) generateIntegerLiteral(expr *ast.IntegerLiteral) (value.Value, error) {
	return constant.NewInt(types.I32, int64(expr.Value)), nil
}

func (g *Generator) generateStringLiteral(expr *ast.StringLiteral) (value.Value, error) {
	return g.getStringConstant(expr.Value), nil
}

func (g *Generator) getStringConstant(str string) value.Value {
	if global, ok := g.context.stringConstants[str]; ok {
		zero := constant.NewInt(types.I64, 0)
		strType := types.NewArray(uint64(len(str)+1), types.I8)
		return constant.NewGetElementPtr(strType, global, zero, zero)
	}

	processedStr := strings.ReplaceAll(str, "\\n", "\n")
	strType := types.NewArray(uint64(len(processedStr)+1), types.I8)
	strConst := g.module.NewGlobalDef(fmt.Sprintf(".str.%d", g.stringCounter),
		constant.NewCharArrayFromString(processedStr+"\x00"))
	g.stringCounter++

	zero := constant.NewInt(types.I64, 0)
	strPtr := constant.NewGetElementPtr(strType, strConst, zero, zero)
	g.context.stringConstants[str] = strConst
	return strPtr
}

func (g *Generator) generateBooleanLiteral(expr *ast.Boolean) (value.Value, error) {
	if expr.Value {
		return constant.NewInt(types.I1, 1), nil
	}
	return constant.NewInt(types.I1, 0), nil
}

func newContext(parent *Context) *Context {
	return &Context{
		currentFunction: nil,
		blocks:          make(map[string]*ir.Block),
		namedValues:     make(map[string]value.Value),
		stringConstants: make(map[string]*ir.Global),
		parent:          parent,
	}
}

func declarePrintf(module *ir.Module) *ir.Func {
	printfType := types.NewFunc(types.I32, types.NewPointer(types.I8))
	fn := module.NewFunc("printf", printfType)
	fn.Linkage = enum.LinkageExternal
	fn.Sig.Variadic = true
	fn.CallingConv = enum.CallingConvC

	if len(fn.Params) > 0 {
		fn.Params[0].SetName("format")
	}

	return fn
}

func declareRuntime(module *ir.Module) {
	declareExternalFunction(module, "malloc", types.NewPointer(types.I8), types.I64)
	declareExternalFunction(module, "free", types.Void, types.NewPointer(types.I8))
	declareExternalFunction(module, "strlen", types.I64, types.NewPointer(types.I8))
	declareExternalFunction(module, "strcpy", types.NewPointer(types.I8),
		types.NewPointer(types.I8), types.NewPointer(types.I8))
	declareExternalFunction(module, "abs", types.I32, types.I32)
	declareExternalFunction(module, "pow", types.Double, types.Double, types.Double)
	declareExternalFunction(module, "exit", types.Void, types.I32)
}

func getStandardFunctionDeclarations() string {
	return `
; Standard function declarations
declare i32 @printf(i8*, ...)
declare i8* @malloc(i64)
declare void @free(i8*)
declare i64 @strlen(i8*)
declare i8* @strcpy(i8*, i8*)
declare i32 @abs(i32)
declare double @pow(double, double)
declare void @exit(i32)

; Standard format strings
@.fmt.int = private constant [4 x i8] c"%d\0A\00"
@.fmt.str = private constant [3 x i8] c"%s\00"
`
}

func declareExternalFunction(module *ir.Module, name string, retType types.Type, paramTypes ...types.Type) *ir.Func {
	funcType := types.NewFunc(retType, paramTypes...)
	fn := module.NewFunc(name, funcType)
	for i := range paramTypes {
		if i < len(fn.Params) {
			fn.Params[i].SetName(fmt.Sprintf("param%d", i))
		}
	}

	return fn
}

func fixFunctionDeclarations(ir string) string {
	fmt.Fprintf(os.Stderr, "DEBUG: Initial IR before adding standard declarations:\n%s\n", ir)

	// Remove any existing function declarations
	ir = regexp.MustCompile(`declare\s+(?:external\s+)?(?:ccc\s+)?[^@]+@(?:printf|malloc|free|strlen|strcpy|abs|pow|exit)\(.*\)`).ReplaceAllString(ir, "")

	// Add standard function declarations at the top
	ir = getStandardFunctionDeclarations() + "\n" + ir

	fmt.Fprintf(os.Stderr, "DEBUG: IR after adding standard declarations:\n%s\n", ir)

	// Fix main function signature
	mainPattern := regexp.MustCompile(`define\s+[^@]+\s+@main\s*\([^)]*\)\s*(?:#[0-9]+)?\s*{`)
	ir = mainPattern.ReplaceAllStringFunc(ir, func(s string) string {
		fmt.Fprintf(os.Stderr, "DEBUG: Main function found: %s\n", s)
		return "define i32 @main() {"
	})

	ir = regexp.MustCompile(`call\s+i32\s*\([^)]*\)\s*@printf`).ReplaceAllString(ir, "call i32 @printf(i8*, ...)")

	//TEST:
	ir = regexp.MustCompile(`\(([^)]*)\)\s*\(\.\.\.\)`).ReplaceAllString(ir, "($1, ...)")
	// now catch printf specifically and insert the full variadic signature
	ir = regexp.MustCompile(`call\s+i32\s*\([^)]*\)\s*@printf`).ReplaceAllString(ir, "call i32 (i8*, ...) @printf")

	// Fix getelementptr instructions for printf calls
	ir = regexp.MustCompile(`getelementptr\s+\[(\d+)\s+x\s+i8\],\s*i8\*\s*getelementptr`).ReplaceAllString(ir, "getelementptr [$1 x i8], [$1 x i8]*")
	ir = regexp.MustCompile(`i8\*\s*getelementptr\(\[(\d+)\s+x\s+i8\],\s*i8\*\s*getelementptr`).ReplaceAllString(ir, "i8* getelementptr([$1 x i8], [$1 x i8]*")

	// Convert array types to pointer types in printf calls (catch the common pattern)
	// before our existing hack, insert something like:
	re := regexp.MustCompile(`getelementptr\s+\[(\d+)\s+x\s+i8\],\s*\[\d+\s+x\s+i8\]\*\s*(@\.str\.\d+),\s*i64\s+0,\s*i64\s+0`)
	ir = re.ReplaceAllString(ir,
		`getelementptr inbounds ([$1 x i8], [$1 x i8]* $2, i64 0, i64 0)`)

	badFuncPattern := regexp.MustCompile(`define\s+i32\s+\(i32\)\s+@(\w+)`)
	matches := badFuncPattern.FindAllString(ir, -1)
	for _, m := range matches {
		fmt.Fprintf(os.Stderr, "ERROR: Invalid function definition: %s\n", m)
	}

	// Fix string constant declarations
	ir = regexp.MustCompile(`@\.str\.[0-9]+ = .*c"([^"]*)".*`).ReplaceAllStringFunc(ir, func(s string) string {
		matches := regexp.MustCompile(`c"([^"]*)"`).FindStringSubmatch(s)
		if len(matches) > 1 {
			content := strings.ReplaceAll(matches[1], "\\n", "\n")
			content = strings.ReplaceAll(content, "\\00", "\x00")
			return strings.Replace(s, matches[1], content, 1)
		}
		return s
	})

	//INFO: defines the main function correctly
	ir = regexp.MustCompile(`define\s+i32\s+\(i32\)\s+@(\w+)\s*\(\)\s*\{`).
		ReplaceAllString(ir, "define i32 @$1(i32 %n) {")
	ir = regexp.MustCompile(`declare\s+i8\*\s+\(i64\)\s+@malloc\(\)`).ReplaceAllString(ir, "declare i8* @malloc(i64)")
	ir = regexp.MustCompile(`declare\s+void\s+\(i8\*\)\s+@free\(\)`).ReplaceAllString(ir, "declare void @free(i8*)")
	ir = regexp.MustCompile(`declare\s+i64\s+\(i8\*\)\s+@strlen\(\)`).ReplaceAllString(ir, "declare i64 @strlen(i8*)")
	ir = regexp.MustCompile(`declare\s+i8\*\s+\(i8\*,\s*i8\*\)\s+@strcpy\(\)`).ReplaceAllString(ir, "declare i8* @strcpy(i8*, i8*)")
	ir = regexp.MustCompile(`declare\s+i32\s+\(i32\)\s+@abs\(\)`).ReplaceAllString(ir, "declare i32 @abs(i32)")
	ir = regexp.MustCompile(`declare\s+double\s+\(double,\s*double\)\s+@pow\(\)`).ReplaceAllString(ir, "declare double @pow(double, double)")
	ir = regexp.MustCompile(`declare\s+void\s+\(i32\)\s+@exit\(\)`).ReplaceAllString(ir, "declare void @exit(i32)")

	fmt.Fprintf(os.Stderr, "DEBUG: Final IR after all fixes:\n%s\n", ir)

	// Final verification
	if !strings.Contains(ir, "define i32 @main()") {
		fmt.Fprintf(os.Stderr, "ERROR: Main function could not be generated\n")
	} else {
		fmt.Fprintf(os.Stderr, "DEBUG: Main function successfully generated\n")
	}

	return ir
}
