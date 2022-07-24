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

func (e *ComplicationEngine) CompileClass(t *Tokenizer) {
	class := map[string]interface{}{}
	e.body["class"] = class

	var i int
	for t.HasMoreTokens() {
		i++
		tokenType, _ := t.TokenType()
		tokenTypeVal := tokenType.String()

		if tokenType == KEYWORD {
			keyword, _ := t.Keyword()
			switch keyword {
			case STATIC, FIELD:
				e.CompileClassVarDec(t)
				continue
			case CONSTRUCTOR, FUNCTION, METHOD, VOID:
				e.CompileSubroutine(t, i)
			}
		}

		if _, ok := class[tokenType.String()]; ok {
			tokenTypeVal = tokenTypeVal + "_" + strconv.Itoa(i)
		}

		class[tokenTypeVal] = t.CurrentToken()

		t.Advance()
	}
}

func (e *ComplicationEngine) CompileClassVarDec(t *Tokenizer) {
	target := map[string]interface{}{}
	classElement := e.body["class"].(map[string]interface{})
	classElement["classVarDec"] = target

	var i int
	for t.HasMoreTokens() {
		i++

		tokenType, _ := t.TokenType()
		tokenTypeVal := tokenType.String()
		if _, ok := target[tokenType.String()]; ok {
			tokenTypeVal = tokenTypeVal + "_" + strconv.Itoa(i)
		}

		target[tokenTypeVal] = t.CurrentToken()
		if t.CurrentToken() == ";" {
			return
		}

		t.Advance()
	}
}

func (e *ComplicationEngine) CompileSubroutine(t *Tokenizer, i int) {
	target := map[string]interface{}{}
	classElement := e.body["class"].(map[string]interface{})

	key := "subroutineDec"
	if _, ok := classElement[key]; ok {
		key = key + strconv.Itoa(i)
	}

	classElement[key] = target

	var index int
	var isLoadParameter bool
	for t.HasMoreTokens() {
		index++

		tokenType, _ := t.TokenType()
		tokenTypeVal := tokenType.String()
		if _, ok := target[tokenTypeVal]; ok {
			tokenTypeVal = tokenTypeVal + "_" + strconv.Itoa(i)
		}
		target[tokenTypeVal] = t.CurrentToken()

		if target[tokenTypeVal] == "}" {
			return
		}

		t.Advance()

		if target[tokenTypeVal] == "(" && !isLoadParameter {
			e.CompileParameterList(t, target)
			isLoadParameter = true
		}

		if target[tokenTypeVal] == "{" {
			e.CompileSubroutineBody(t, target)
		}
	}
}

func (e *ComplicationEngine) CompileSubroutineBody(t *Tokenizer, subroutineField map[string]interface{}) {
	target := map[string]interface{}{}
	subroutineField["subroutineBody"] = target

	var index int
	for t.HasMoreTokens() {
		if t.CurrentToken() == "}" {
			return
		}

		index++

		tokenType, _ := t.TokenType()
		keyword, _ := t.Keyword()
		if tokenType == KEYWORD && keyword == VAR {
			e.CompileVarDec(t, target, index)
			t.Advance()
			continue
		}

		e.CompileStatements(t, target)

		t.Advance()
	}
}

func (e *ComplicationEngine) CompileParameterList(t *Tokenizer, subroutineField map[string]interface{}) {
	target := map[string]interface{}{}
	subroutineField["parameterList"] = target

	var index int
	for t.HasMoreTokens() {
		if t.CurrentToken() == ")" {
			return
		}

		index++
		tokenType, _ := t.TokenType()
		tokenTypeVal := tokenType.String()
		if _, ok := target[tokenType.String()]; ok {
			tokenTypeVal = tokenTypeVal + "_" + strconv.Itoa(index)
		}
		target[tokenTypeVal] = t.CurrentToken()

		t.Advance()
	}
}

func (e *ComplicationEngine) CompileVarDec(t *Tokenizer, subroutineField map[string]interface{}, i int) {
	key := "varDec"
	if _, ok := subroutineField[key]; ok {
		key = key + strconv.Itoa(i)
	}
	target := map[string]interface{}{}
	subroutineField[key] = target

	var index int
	for t.HasMoreTokens() {
		index++

		tokenType, _ := t.TokenType()
		tokenTypeVal := tokenType.String()
		if _, ok := target[tokenType.String()]; ok {
			tokenTypeVal = tokenTypeVal + "_" + strconv.Itoa(index)
		}
		target[tokenTypeVal] = t.CurrentToken()

		if target[tokenTypeVal] == ";" {
			return
		}

		t.Advance()
	}
}

func (e *ComplicationEngine) CompileStatements(t *Tokenizer, subroutineBodyField map[string]interface{}) {
	target := map[string]interface{}{}
	subroutineBodyField["statements"] = target

	var index int
	for t.HasMoreTokens() {
		index++

		if t.CurrentToken() == "}" {
			return
		}

		keyword, _ := t.Keyword()
		switch keyword {
		case WHILE:
		case IF:
		case LET:
			e.CompileLet(t, target, index)
		case DO:
		case RETURN:
		}

		t.Advance()
	}
}

func (e *ComplicationEngine) CompileDo() error {
	return nil
}

func (e *ComplicationEngine) CompileLet(t *Tokenizer, statementsField map[string]interface{}, i int) {
	key := "letStatement"
	if _, ok := statementsField[key]; ok {
		key = key + strconv.Itoa(i)
	}

	target := map[string]interface{}{}
	statementsField[key] = target

	var index int
	for t.HasMoreTokens() {
		index++
		tokenType, _ := t.TokenType()
		tokenTypeVal := tokenType.String()
		if _, ok := target[tokenType.String()]; ok {
			tokenTypeVal = tokenTypeVal + "_" + strconv.Itoa(index)
		}
		target[tokenTypeVal] = t.CurrentToken()

		switch tokenType {
		case SYMBOL:
			switch t.CurrentToken() {
			case "[", "=":
				//e.CompileExpression()
			}
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

func (e *ComplicationEngine) CompileExpression(t *Tokenizer, field map[string]interface{}, i int) {
	key := "expression"
	if _, ok := field[key]; ok {
		key = key + strconv.Itoa(i)
	}

	target := map[string]interface{}{}
	field[key] = target

	var index int
	for t.HasMoreTokens() {
		if t.CurrentToken() == ";" {
			return
		}

		e.CompileTerm(t, target, index)

		t.Advance()
	}
}

func (e *ComplicationEngine) CompileTerm(t *Tokenizer, field map[string]interface{}, i int) {
	key := "term"
	if _, ok := field[key]; ok {
		key = key + strconv.Itoa(i)
	}

	target := map[string]interface{}{}
	field[key] = target

	for t.HasMoreTokens() {

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
