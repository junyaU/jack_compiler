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
	f, err := os.Create(sourceName + ".json")
	if err != nil {
		return nil, err
	}

	return &ComplicationEngine{
		file: f,
		body: map[string]interface{}{},
	}, nil
}

func (e *ComplicationEngine) Compile(t *Tokenizer) {
	field := map[string]interface{}{}
	e.body["class"] = field

	for t.HasMoreTokens() {
		switch t.CurrentToken() {
		case STATIC.String(), FIELD.String():
			e.compileClassVarDec(t, field)
			t.Advance()
			continue

		case CONSTRUCTOR.String(), FUNCTION.String(), METHOD.String(), VOID.String():
			e.compileSubroutine(t, field)
			continue

		default:
			field[t.MakeTokenKey()] = t.CurrentToken()
		}

		t.Advance()
	}
}

func (e *ComplicationEngine) compileClassVarDec(t *Tokenizer, parent map[string]interface{}) {
	field := e.registerStructureField("classVarDec", t.currentLine, parent)

	for t.HasMoreTokens() {
		field[t.MakeTokenKey()] = t.CurrentToken()

		if t.CurrentToken() == ";" {
			return
		}

		t.Advance()
	}
}

func (e *ComplicationEngine) compileSubroutine(t *Tokenizer, parent map[string]interface{}) {
	field := e.registerStructureField("subroutineDec", t.currentLine, parent)

	for t.HasMoreTokens() {
		switch t.CurrentToken() {
		case "(":
			field[t.MakeTokenKey()] = t.CurrentToken()
			t.Advance()
			e.compileParameterList(t, field)
			continue

		case "{":
			e.compileSubroutineBody(t, field)
			return
		}

		field[t.MakeTokenKey()] = t.CurrentToken()

		t.Advance()
	}
}

func (e *ComplicationEngine) compileSubroutineBody(t *Tokenizer, parent map[string]interface{}) {
	field := e.registerStructureField("subroutineBody", t.currentLine, parent)

	for t.HasMoreTokens() {
		switch t.CurrentToken() {
		case "{":
			field[t.MakeTokenKey()] = t.CurrentToken()
			e.compileStatements(t, field)
			field[t.MakeTokenKey()] = t.CurrentToken()
			t.Advance()
			return

		case VAR.String():
			e.compileVarDec(t, field)
			t.Advance()
			continue
		}
	}
}

func (e *ComplicationEngine) compileParameterList(t *Tokenizer, parent map[string]interface{}) {
	field := e.registerStructureField("parameterList", t.currentLine, parent)

	for t.HasMoreTokens() {
		if t.CurrentToken() == ")" {
			return
		}

		field[t.MakeTokenKey()] = t.CurrentToken()
		t.Advance()
	}
}

func (e *ComplicationEngine) compileVarDec(t *Tokenizer, parent map[string]interface{}) {
	field := e.registerStructureField("varDec", t.currentLine, parent)

	for t.HasMoreTokens() {
		field[t.MakeTokenKey()] = t.CurrentToken()

		tokenType, _ := t.TokenType()
		if field[tokenType.String()] == ";" {
			return
		}

		t.Advance()
	}
}

func (e *ComplicationEngine) compileStatements(t *Tokenizer, parent map[string]interface{}) {
	field := e.registerStructureField("statements", t.currentLine, parent)

	for t.HasMoreTokens() {
		switch t.CurrentToken() {
		case WHILE.String():
			e.compileWhile(t, field)
		case IF.String():
			e.compileIf(t, field)
		case LET.String():
			e.compileLet(t, field)
		case DO.String():
			e.compileDo(t, field)
		case RETURN.String():
			e.compileReturn(t, field)
		case "}":
			return
		}

		t.Advance()
	}
}

func (e *ComplicationEngine) compileDo(t *Tokenizer, parent map[string]interface{}) {
	field := e.registerStructureField("doStatement", t.currentLine, parent)

	for t.HasMoreTokens() {
		switch t.CurrentToken() {
		case "(":
			field[t.MakeTokenKey()] = t.CurrentToken()
			t.Advance()
			e.compileExpressionList(t, field)
		}

		field[t.MakeTokenKey()] = t.CurrentToken()

		if t.CurrentToken() == ";" {
			return
		}

		t.Advance()
	}

}

func (e *ComplicationEngine) compileLet(t *Tokenizer, parent map[string]interface{}) {
	field := e.registerStructureField("letStatement", t.currentLine, parent)

	for t.HasMoreTokens() {
		switch t.CurrentToken() {
		case "[", "=":
			field[t.MakeTokenKey()] = t.CurrentToken()
			t.Advance()
			e.compileExpression(t, field, false)
		}

		field[t.MakeTokenKey()] = t.CurrentToken()

		if t.CurrentToken() == ";" {
			return
		}

		t.Advance()
	}
}

func (e *ComplicationEngine) compileWhile(t *Tokenizer, parent map[string]interface{}) {
	field := e.registerStructureField("whileStatement", t.currentLine, parent)

	for t.HasMoreTokens() {
		switch t.CurrentToken() {
		case "(":
			field[t.MakeTokenKey()] = t.CurrentToken()
			t.Advance()
			e.compileExpression(t, field, true)

		case "{":
			field[t.MakeTokenKey()] = t.CurrentToken()
			t.Advance()
			e.compileStatements(t, field)
		}

		field[t.MakeTokenKey()] = t.CurrentToken()

		if t.CurrentToken() == "}" {
			return
		}

		t.Advance()
	}
}

func (e *ComplicationEngine) compileReturn(t *Tokenizer, parent map[string]interface{}) {
	field := e.registerStructureField("returnStatement", t.currentLine, parent)

	for t.HasMoreTokens() {
		if t.CurrentToken() != RETURN.String() {
			e.compileExpression(t, field, false)
		}

		field[t.MakeTokenKey()] = t.CurrentToken()

		if t.CurrentToken() == ";" {
			return
		}

		t.Advance()
	}
}

func (e *ComplicationEngine) compileIf(t *Tokenizer, parent map[string]interface{}) {
	field := e.registerStructureField("ifStatement", t.currentLine, parent)

	for t.HasMoreTokens() {
		switch t.CurrentToken() {
		case "(":
			field[t.MakeTokenKey()] = t.CurrentToken()
			t.Advance()
			e.compileExpression(t, field, true)

		case "{":
			field[t.MakeTokenKey()] = t.CurrentToken()
			t.Advance()
			e.compileStatements(t, field)
		}

		if t.CurrentToken() == "}" {
			field[t.MakeTokenKey()] = t.CurrentToken()

			if t.NextToken() == ELSE.String() {
				t.Advance()
				continue
			}

			return
		}

		field[t.MakeTokenKey()] = t.CurrentToken()

		t.Advance()
	}
}

func (e *ComplicationEngine) compileExpression(t *Tokenizer, parent map[string]interface{}, isRecursion bool) {
	field := e.registerStructureField("expression", t.currentLine, parent)

	for t.HasMoreTokens() {
		e.compileTerm(t, field, isRecursion)

		switch t.CurrentToken() {
		case ";", ")", "]":
			return
		}

		t.Advance()
	}
}

func (e *ComplicationEngine) compileTerm(t *Tokenizer, parent map[string]interface{}, isRecursion bool) {
	field := e.registerStructureField("term", t.currentLine, parent)

	for t.HasMoreTokens() {
		switch t.CurrentToken() {
		case ";":
			return

		case "[", "(":
			field[t.MakeTokenKey()] = t.CurrentToken()
			t.Advance()
			e.compileExpression(t, field, true)
		}

		if t.CurrentToken() == ")" {
			if !isRecursion {
				field[t.MakeTokenKey()] = t.CurrentToken()
				t.Advance()
			}
			return
		}

		if t.CurrentToken() == "]" {
			return
		}

		field[t.MakeTokenKey()] = t.CurrentToken()

		t.Advance()
	}
}

func (e *ComplicationEngine) compileExpressionList(t *Tokenizer, parent map[string]interface{}) {
	field := e.registerStructureField("expressionList", t.currentLine, parent)

	for t.HasMoreTokens() {
		switch t.CurrentToken() {
		case ",":
			field[t.MakeTokenKey()] = t.CurrentToken()

		default:
			e.compileExpression(t, field, true)
		}

		if t.CurrentToken() == ")" {
			return
		}

		t.Advance()
	}
}

func (e *ComplicationEngine) registerStructureField(fieldName string, identifier int, parent map[string]interface{}) map[string]interface{} {
	if _, ok := parent[fieldName]; ok {
		fieldName = fieldName + strconv.Itoa(identifier)
	}

	field := map[string]interface{}{}
	parent[fieldName] = field

	return field
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
