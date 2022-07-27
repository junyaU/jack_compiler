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
		tokenType, _ := t.TokenType()
		if tokenType == KEYWORD {
			keyword, _ := t.Keyword()
			switch keyword {
			case STATIC, FIELD:
				e.CompileClassVarDec(t, class)
				continue
			case CONSTRUCTOR, FUNCTION, METHOD, VOID:
				e.CompileSubroutine(t, class)
			}
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

	var isLoadParameter bool
	for t.HasMoreTokens() {
		target[t.MakeTokenKey()] = t.CurrentToken()

		if t.CurrentToken() == "}" {
			return
		}

		t.Advance()

		if t.CurrentToken() == "(" && !isLoadParameter {
			e.CompileParameterList(t, target)
			isLoadParameter = true
		}

		if t.CurrentToken() == "{" {
			e.CompileSubroutineBody(t, target)
		}
	}
}

func (e *ComplicationEngine) CompileSubroutineBody(t *Tokenizer, subroutineField map[string]interface{}) {
	target := map[string]interface{}{}
	subroutineField["subroutineBody"] = target

	for t.HasMoreTokens() {
		if t.CurrentToken() == "}" {
			return
		}

		tokenType, _ := t.TokenType()
		keyword, _ := t.Keyword()
		if tokenType == KEYWORD && keyword == VAR {
			e.CompileVarDec(t, target)
			t.Advance()
			continue
		}

		e.CompileStatements(t, target)

		if t.CurrentToken() == "}" {
			return
		}
		t.Advance()
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

func (e *ComplicationEngine) CompileStatements(t *Tokenizer, subroutineBodyField map[string]interface{}) {
	target := map[string]interface{}{}
	subroutineBodyField["statements"] = target

	for t.HasMoreTokens() {
		if t.CurrentToken() == "}" {
			return
		}

		keyword, _ := t.Keyword()
		switch keyword {
		case WHILE:
		case IF:
		case LET:
			e.CompileLet(t, target)
		case DO:
		case RETURN:
		}

		t.Advance()
	}
}

func (e *ComplicationEngine) CompileDo() error {
	return nil
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
			target[t.MakeTokenKey()] = t.CurrentToken()
			fallthrough

		case ";":
			return

		default:
			target[t.MakeTokenKey()] = t.CurrentToken()
		}

		t.Advance()
	}
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
		case "true", "false", "null", "this":
			target["KeywordConstant"] = t.CurrentToken()

		case "[", "(":
			target[t.MakeTokenKey()] = t.CurrentToken()
			t.Advance()
			e.CompileExpression(t, target, true)
			fallthrough

		case ")":
			if !isRecursion {
				target[t.MakeTokenKey()] = t.CurrentToken()
				t.Advance()
			}
			return

		case "]":
			return
		}

		target[t.MakeTokenKey()] = t.CurrentToken()

		t.Advance()
	}
}

func (e *ComplicationEngine) CompileExpressionList() error {
	return nil
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
