package jack_compiler

import (
	"bytes"
	"encoding/json"
	"os"
	"strconv"
)

type ComplicationEngine struct {
	file *os.File
	body map[string]interface{}
}

func NewComplicationEngine(sourceName string) (*ComplicationEngine, error) {
	f, err := os.Create("testdata/" + sourceName + ".json")
	if err != nil {
		return nil, err
	}

	return &ComplicationEngine{
		file: f,
		body: map[string]interface{}{},
	}, nil
}

func (e *ComplicationEngine) Compile(t *Tokenizer) {
	class := map[string]interface{}{}
	e.body["class"] = class

	for t.HasMoreTokens() {
		switch t.CurrentToken() {
		case STATIC.String(), FIELD.String():
			e.CompileClassVarDec(t, class)
			t.Advance()
			continue

		case CONSTRUCTOR.String(), FUNCTION.String(), METHOD.String(), VOID.String():
			e.CompileSubroutine(t, class)
		}

		class[t.MakeTokenKey()] = t.CurrentToken()

		t.Advance()
	}
}

func (e *ComplicationEngine) CompileClassVarDec(t *Tokenizer, field map[string]interface{}) {
	target := map[string]interface{}{}
	field["classVarDec"] = target

	for t.HasMoreTokens() {
		target[t.MakeTokenKey()] = t.CurrentToken()

		if t.CurrentToken() == ";" {
			return
		}

		t.Advance()
	}
}

func (e *ComplicationEngine) CompileSubroutine(t *Tokenizer, field map[string]interface{}) {
	key := "subroutineDec"
	if _, ok := field[key]; ok {
		key = key + strconv.Itoa(t.currentLine)
	}

	target := map[string]interface{}{}
	field[key] = target

	for t.HasMoreTokens() {
		switch t.CurrentToken() {
		case "(":
			target[t.MakeTokenKey()] = t.CurrentToken()
			t.Advance()
			e.CompileParameterList(t, target)
			continue

		case "{":
			e.CompileSubroutineBody(t, target)
			return
		}

		target[t.MakeTokenKey()] = t.CurrentToken()

		t.Advance()
	}
}

func (e *ComplicationEngine) CompileSubroutineBody(t *Tokenizer, subroutineField map[string]interface{}) {
	target := map[string]interface{}{}
	subroutineField["subroutineBody"] = target

	for t.HasMoreTokens() {
		switch t.CurrentToken() {
		case "{":
			target[t.MakeTokenKey()] = t.CurrentToken()
			e.CompileStatements(t, target)
			target[t.MakeTokenKey()] = t.CurrentToken()
			t.Advance()
			return

		case VAR.String():
			e.CompileVarDec(t, target)
			t.Advance()
			continue
		}
	}
}

func (e *ComplicationEngine) CompileParameterList(t *Tokenizer, subroutineField map[string]interface{}) {
	target := map[string]interface{}{}
	subroutineField["parameterList"] = target

	for t.HasMoreTokens() {
		if t.CurrentToken() == ")" {
			return
		}

		target[t.MakeTokenKey()] = t.CurrentToken()
		t.Advance()
	}
}

func (e *ComplicationEngine) CompileVarDec(t *Tokenizer, subroutineField map[string]interface{}) {
	key := "varDec"
	if _, ok := subroutineField[key]; ok {
		key = key + strconv.Itoa(t.currentLine)
	}

	target := map[string]interface{}{}
	subroutineField[key] = target

	for t.HasMoreTokens() {
		target[t.MakeTokenKey()] = t.CurrentToken()

		tokenType, _ := t.TokenType()
		if target[tokenType.String()] == ";" {
			return
		}

		t.Advance()
	}
}

func (e *ComplicationEngine) CompileStatements(t *Tokenizer, field map[string]interface{}) {
	target := map[string]interface{}{}
	field["statements"] = target

	for t.HasMoreTokens() {
		switch t.CurrentToken() {
		case WHILE.String():
			e.CompileWhile(t, target)
		case IF.String():
		case LET.String():
			e.CompileLet(t, target)
		case DO.String():
			e.CompileDo(t, target)
		case RETURN.String():
			e.CompileReturn(t, target)
		case "}":
			return
		}

		t.Advance()
	}
}

func (e *ComplicationEngine) CompileDo(t *Tokenizer, field map[string]interface{}) {
	key := "doStatement"
	if _, ok := field[key]; ok {
		key = key + strconv.Itoa(t.currentLine)
	}

	target := map[string]interface{}{}
	field[key] = target

	for t.HasMoreTokens() {
		switch t.CurrentToken() {
		case "(":
			target[t.MakeTokenKey()] = t.CurrentToken()
			t.Advance()
			e.CompileExpressionList(t, target)
		}

		target[t.MakeTokenKey()] = t.CurrentToken()

		if t.CurrentToken() == ";" {
			return
		}

		t.Advance()
	}

}

func (e *ComplicationEngine) CompileLet(t *Tokenizer, statementsField map[string]interface{}) {
	key := "letStatement"
	if _, ok := statementsField[key]; ok {
		key = key + strconv.Itoa(t.currentLine)
	}

	target := map[string]interface{}{}
	statementsField[key] = target

	for t.HasMoreTokens() {
		switch t.CurrentToken() {
		case "[", "=":
			target[t.MakeTokenKey()] = t.CurrentToken()
			t.Advance()
			e.CompileExpression(t, target, false)
		}

		target[t.MakeTokenKey()] = t.CurrentToken()

		if t.CurrentToken() == ";" {
			return
		}

		t.Advance()
	}
}

func (e *ComplicationEngine) CompileWhile(t *Tokenizer, field map[string]interface{}) {
	key := "whileStatement"
	if _, ok := field[key]; ok {
		key = key + strconv.Itoa(t.currentLine)
	}

	target := map[string]interface{}{}
	field[key] = target

	for t.HasMoreTokens() {
		switch t.CurrentToken() {
		case "(":
			target[t.MakeTokenKey()] = t.CurrentToken()
			t.Advance()
			e.CompileExpression(t, target, true)

		case "{":
			target[t.MakeTokenKey()] = t.CurrentToken()
			t.Advance()
			e.CompileStatements(t, target)
		}

		target[t.MakeTokenKey()] = t.CurrentToken()

		if t.CurrentToken() == "}" {
			return
		}

		t.Advance()
	}
}

func (e *ComplicationEngine) CompileReturn(t *Tokenizer, field map[string]interface{}) {
	key := "returnStatement"
	if _, ok := field[key]; ok {
		key = key + strconv.Itoa(t.currentLine)
	}

	target := map[string]interface{}{}
	field[key] = target

	for t.HasMoreTokens() {
		if t.CurrentToken() != RETURN.String() {
			e.CompileExpression(t, target, false)
		}

		target[t.MakeTokenKey()] = t.CurrentToken()

		if t.CurrentToken() == ";" {
			return
		}

		t.Advance()
	}
}

func (e *ComplicationEngine) CompileIf() {
}

func (e *ComplicationEngine) CompileExpression(t *Tokenizer, field map[string]interface{}, isRecursion bool) {
	key := "expression"
	if _, ok := field[key]; ok {
		key = key + strconv.Itoa(t.currentLine)
	}

	target := map[string]interface{}{}
	field[key] = target

	for t.HasMoreTokens() {
		e.CompileTerm(t, target, isRecursion)

		switch t.CurrentToken() {
		case ";", ")", "]":
			return
		}

		t.Advance()
	}
}

func (e *ComplicationEngine) CompileTerm(t *Tokenizer, field map[string]interface{}, isRecursion bool) {
	key := "term"
	if _, ok := field[key]; ok {
		key = key + strconv.Itoa(t.currentLine)
	}

	target := map[string]interface{}{}
	field[key] = target

	for t.HasMoreTokens() {
		switch t.CurrentToken() {
		case ";":
			return
		case TRUE.String(), FALSE.String(), NULL.String(), THIS.String():
			target["KeywordConstant"] = t.CurrentToken()

		case "[", "(":
			target[t.MakeTokenKey()] = t.CurrentToken()
			t.Advance()
			e.CompileExpression(t, target, true)
		}

		if t.CurrentToken() == ")" {
			if !isRecursion {
				target[t.MakeTokenKey()] = t.CurrentToken()
				t.Advance()
			}
			return
		}

		if t.CurrentToken() == "]" {
			return
		}

		target[t.MakeTokenKey()] = t.CurrentToken()

		t.Advance()
	}
}

func (e *ComplicationEngine) CompileExpressionList(t *Tokenizer, field map[string]interface{}) {
	key := "expressionList"
	if _, ok := field[key]; ok {
		key = key + strconv.Itoa(t.currentLine)
	}

	target := map[string]interface{}{}
	field[key] = target

	for t.HasMoreTokens() {
		switch t.CurrentToken() {
		case ",":
			target[key] = t.CurrentToken()

		default:
			e.CompileExpression(t, target, true)
		}

		if t.CurrentToken() == ")" {
			return
		}

		t.Advance()
	}

}

func (e *ComplicationEngine) Write() error {
	var buf bytes.Buffer
	output, err := json.Marshal(e.body)
	if err != nil {
		return err
	}

	if err := json.Indent(&buf, output, "", "  "); err != nil {
		return err
	}

	_, err = e.file.WriteString(buf.String())
	return err
}

func (e *ComplicationEngine) Close() {
	e.file.Close()
}
