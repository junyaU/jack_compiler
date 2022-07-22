package jack_compiler

import (
	"os"
)

type ComplicationEngine struct {
	file *os.File
}

func NewComplicationEngine(sourceName string) (*ComplicationEngine, error) {
	f, err := os.Create("testdata/" + sourceName + ".xml")
	if err != nil {
		return nil, err
	}

	return &ComplicationEngine{
		file: f,
	}, nil
}

func (e *ComplicationEngine) CompileClass() error {
	return nil
}

func (e *ComplicationEngine) CompileClassVarDec() error {
	return nil
}

func (e *ComplicationEngine) CompileSubroutine() error {
	return nil
}

func (e *ComplicationEngine) CompileParameterList() error {
	return nil
}

func (e *ComplicationEngine) CompileVarDec() error {
	return nil
}

func (e *ComplicationEngine) CompileStatements() error {
	return nil
}

func (e *ComplicationEngine) CompileDo() error {
	return nil
}

func (e *ComplicationEngine) CompileLet() error {
	return nil
}

func (e *ComplicationEngine) CompileWhile() error {
	return nil
}

func (e *ComplicationEngine) CompileReturn() error {
	return nil
}

func (e *ComplicationEngine) CompileIf() error {
	return nil
}

func (e *ComplicationEngine) CompileExpression() error {
	return nil
}

func (e *ComplicationEngine) CompileTerm() error {
	return nil
}

func (e *ComplicationEngine) CompileExpressionList() error {
	return nil
}

func (e *ComplicationEngine) Close() {
	e.file.Close()
}
