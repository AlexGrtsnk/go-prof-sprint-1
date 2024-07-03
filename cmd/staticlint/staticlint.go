package main

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	checks := map[string]bool{
		"SA5000": true,
		"SA6000": true,
		"SA9004": true,
		"SA4006": true,
		"ST2021": true,
	}
	var mychecks []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		// добавляем в массив нужные проверки
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}
	mychecks = append(mychecks, ErrCheckAnalyzer, printf.Analyzer, shadow.Analyzer, structtag.Analyzer)
	multichecker.Main(
		mychecks...,
	)
}

var ErrCheckAnalyzer = &analysis.Analyzer{
	Name: "errcheck",
	Doc:  "check for unchecked errors",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	expr := func(x *ast.ExprStmt) {
		// проверяем, что выражение представляет собой вызов функции,
		// у которой возвращаемая ошибка никак не обрабатывается
		if call, ok := x.X.(*ast.CallExpr); ok {
			if isReturnError(pass, call) {
				pass.Reportf(x.Pos(), "expression returns unchecked error")
			}
		}
	}
	expr1 := func(x *ast.ExprStmt) {
		// проверяем, что выражение представляет собой вызов функции main,
		// у которой есть вызов os.exit()
		if call, ok := x.X.(*ast.CallExpr); ok {
			if isOsExit(pass, call) {
				pass.Reportf(x.Pos(), "no os.exit allowed in main")
			}
		}
	}
	tuplefunc := func(x *ast.AssignStmt) {
		// рассматриваем присваивание, при котором
		// вместо получения ошибок используется '_'
		// a, b, _ := tuplefunc()
		// проверяем, что это вызов функции
		if call, ok := x.Rhs[0].(*ast.CallExpr); ok {
			results := resultErrors(pass, call)
			for i := 0; i < len(x.Lhs); i++ {
				// перебираем все идентификаторы слева от присваивания
				if id, ok := x.Lhs[i].(*ast.Ident); ok && id.Name == "_" && results[i] {
					pass.Reportf(id.NamePos, "assignment with unchecked error")
				}
			}
		}
	}
	errfunc := func(x *ast.AssignStmt) {
		// множественное присваивание: a, _ := b, myfunc()
		// ищем ситуацию, когда функция справа возвращает ошибку,
		// а соответствующий идентификатор слева равен '_'
		for i := 0; i < len(x.Lhs); i++ {
			if id, ok := x.Lhs[i].(*ast.Ident); ok {
				// вызов функции справа
				if call, ok := x.Rhs[i].(*ast.CallExpr); ok {
					if id.Name == "_" && isReturnError(pass, call) {
						pass.Reportf(id.NamePos, "assignment with unchecked error")
					}
				}
			}
		}
	}

	for _, file := range pass.Files {
		// функцией ast.Inspect проходим по всем узлам AST
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.ExprStmt: // выражение
				expr(x)
				expr1(x)
			case *ast.AssignStmt: // оператор присваивания
				// справа одно выражение x,y := myfunc()
				if len(x.Rhs) == 1 {
					tuplefunc(x)
				} else {
					// справа несколько выражений x,y := z,myfunc()
					errfunc(x)
				}
			}
			return true
		})
	}
	return nil, nil
}

var errorType = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

func isErrorType(t types.Type) bool {
	return types.Implements(t, errorType)
}

// resultErrors возвращает булев массив со значениями true,
// если тип i-го возвращаемого значения соответствует ошибке.
func resultErrors(pass *analysis.Pass, call *ast.CallExpr) []bool {
	switch t := pass.TypesInfo.Types[call].Type.(type) {
	case *types.Named: // возвращается значение
		return []bool{isErrorType(t)}
	case *types.Pointer: // возвращается указатель
		return []bool{isErrorType(t)}
	case *types.Tuple: // возвращается несколько значений
		s := make([]bool, t.Len())
		for i := 0; i < t.Len(); i++ {
			switch mt := t.At(i).Type().(type) {
			case *types.Named:
				s[i] = isErrorType(mt)
			case *types.Pointer:
				s[i] = isErrorType(mt)
			}
		}
		return s
	}
	return []bool{false}
}

// isReturnError возвращает true, если среди возвращаемых значений есть ошибка.
func isReturnError(pass *analysis.Pass, call *ast.CallExpr) bool {
	for _, isError := range resultErrors(pass, call) {
		if isError {
			return true
		}
	}
	return false
}

// resultExits возвращает булев массив со значениями true,
// если тип i-го возвращаемого значения соответствует exit.
func resultExits(pass *analysis.Pass, call *ast.CallExpr) []bool {
	switch t := pass.TypesInfo.Types[call].Type.(type) {
	case *types.Named: // возвращается значение
		return []bool{isExitType(t)}
	case *types.Pointer: // возвращается указатель
		return []bool{isExitType(t)}
	case *types.Tuple: // возвращается несколько значений
		s := make([]bool, t.Len())
		for i := 0; i < t.Len(); i++ {
			switch mt := t.At(i).Type().(type) {
			case *types.Named:
				s[i] = isExitType(mt)
			case *types.Pointer:
				s[i] = isExitType(mt)
			}
		}
		return s
	}
	return []bool{false}
}

var exitType = types.Universe.Lookup("exit").Type().Underlying().(*types.Interface)

func isExitType(t types.Type) bool {
	return types.Implements(t, exitType)
}

// isOsExit возвращает true, если среди функции есть вызов os.exit()
func isOsExit(pass *analysis.Pass, call *ast.CallExpr) bool {
	for _, isExit := range resultExits(pass, call) {
		if isExit {
			return true
		}
	}
	return false
}
