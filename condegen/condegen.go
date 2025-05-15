package codegen

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
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
