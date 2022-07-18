package jack_compiler

import (
	"bufio"
	"errors"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Tokenizer struct {
	tokens      []string
	currentLine int
}

func NewTokenizer(f io.Reader) *Tokenizer {
	removeCommentOut := func(text string) string {
		isExistCommentOut1 := strings.Contains(text, "//")
		if isExistCommentOut1 {
			text = text[:strings.Index(text, "//")]

		}

		isExistCommentOut2 := strings.Contains(text, "/*")
		if isExistCommentOut2 {
			text = text[:strings.Index(text, "/*")]
		}

		return strings.TrimSpace(text)
	}

	extractTokens := func(word string) (tokens []string) {
		splitIntoTokens := func(text string, symbol string) (splitTexts []string) {
			splitTexts = strings.Split(text, symbol)

			if len(splitTexts) == 1 {
				return
			}

			var insertLocations []int
			for i := 1; i < len(splitTexts); i++ {
				insertLocations = append(insertLocations, i*2-1)
			}

			for _, location := range insertLocations {
				splitTexts = append(splitTexts[:location+1], splitTexts[location:]...)
				splitTexts[location] = symbol
			}

			var parseTexts []string
			for _, t := range splitTexts {
				if t != "" {
					parseTexts = append(parseTexts, t)
				}
			}

			splitTexts = parseTexts

			return
		}

		symbols := []string{"{", "}", "(", ")", "[", "]", ".", ",", ";", "+", "-", "*", "/", "&", "|", "<", ">", "=", "~"}
		for _, s := range symbols {
			if len(tokens) == 0 {
				tokens = splitIntoTokens(word, s)
				continue
			}

			var targetTokens []string
			for _, token := range tokens {
				targetTokens = append(targetTokens, splitIntoTokens(token, s)...)
			}
			tokens = targetTokens
		}

		return
	}

	extractStringConstant := func(text string) []string {
		var stringQuoteLocation []int
		targetText := text

		for {
			index := strings.Index(targetText, string('"'))
			if index == -1 {
				break
			}

			addIndex := 1
			if len(stringQuoteLocation) != 0 {
				addIndex = addIndex + stringQuoteLocation[len(stringQuoteLocation)-1]
			}

			stringQuoteLocation = append(stringQuoteLocation, index+addIndex)
			targetText = targetText[index+1:]
		}

		if len(stringQuoteLocation) == 0 {
			return strings.Split(text, " ")
		}

		stringConstants := make(map[string]string, len(stringQuoteLocation)%2)
		replaceText := text
		for len(stringQuoteLocation) != 0 {
			rand.Seed(time.Now().UnixNano())
			randVal := strconv.Itoa(rand.Intn(10000000)) + "_jack_compiler_rand_val"
			targetString := text[stringQuoteLocation[0]-1 : stringQuoteLocation[1]]
			stringConstants[randVal] = targetString
			replaceText = strings.Replace(replaceText, targetString, randVal, 1)
			stringQuoteLocation = stringQuoteLocation[2:]
		}

		splitTexts := strings.Split(replaceText, " ")
		for i, t := range splitTexts {
			for key := range stringConstants {
				if strings.Contains(t, key) {
					splitTexts[i] = strings.Replace(t, key, stringConstants[key], 1)
				}
			}
		}

		return splitTexts
	}

	tokenizer := new(Tokenizer)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lineText := scanner.Text()
		text := removeCommentOut(lineText)
		if text == "" {
			continue
		}

		for _, word := range extractStringConstant(text) {
			for _, token := range extractTokens(word) {
				tokenizer.tokens = append(tokenizer.tokens, token)
			}
		}
	}

	return tokenizer
}

func (t Tokenizer) HasMoreTokens() bool {
	return len(t.tokens) > t.currentLine
}

func (t *Tokenizer) Advance() {
	t.currentLine++
}

func (t Tokenizer) TokenType() (TokenType, error) {
	switch t.tokens[t.currentLine] {
	case "class", "method", "function", "constructor", "int", "boolean", "char", "void", "var", "static",
		"field", "let", "do", "if", "else", "while", "return", "true", "false", "null", "this":
		return KEYWORD, nil

	case "{", "}", "(", ")", "[", "]", ".", ",", ";", "+", "-", "*", "/", "&", "|", "<", ">", "=", "~":
		return SYMBOL, nil
	}

	currentToken := t.tokens[t.currentLine]
	if currentToken[:1] == string('"') && currentToken[len(currentToken)-1:] == string('"') {
		return STRING_CONST, nil
	}

	tokenNum, err := strconv.Atoi(currentToken)
	if err == nil && 0 <= tokenNum && tokenNum <= 32767 {
		return INT_CONST, nil
	}

	if _, err := strconv.Atoi(currentToken[:1]); err != nil {
		return IDENTIFIER, nil
	}

	return 0, errors.New("no applicable type for this token")
}

func (t Tokenizer) keyword() (KeywordType, error) {
	switch t.tokens[t.currentLine] {
	case "class":
		return CLASS, nil
	case "method":
		return METHOD, nil
	case "function":
		return FUNCTION, nil
	case "constructor":
		return CONSTRUCTOR, nil
	case "int":
		return INT, nil
	case "boolean":
		return BOOLEAN, nil
	case "char":
		return CHAR, nil
	case "void":
		return VOID, nil
	case "var":
		return VAR, nil
	case "static":
		return STATIC, nil
	case "field":
		return FIELD, nil
	case "let":
		return LET, nil
	case "do":
		return DO, nil
	case "if":
		return IF, nil
	case "else":
		return ELSE, nil
	case "while":
		return WHILE, nil
	case "return":
		return RETURN, nil
	case "true":
		return TRUE, nil
	case "false":
		return FALSE, nil
	case "null":
		return NULL, nil
	case "this":
		return THIS, nil
	default:
		return 0, errors.New("this keyword does not exist")
	}
}

func (t Tokenizer) Symbol() (string, error) {
	token, err := t.TokenType()
	if err != nil {
		return "", err
	}

	if token != SYMBOL {
		return "", errors.New("cannot call this function due to different token type")
	}

	return t.tokens[t.currentLine], nil
}

func (t Tokenizer) Identifier() (string, error) {
	token, err := t.TokenType()
	if err != nil {
		return "", err
	}

	if token != IDENTIFIER {
		return "", errors.New("cannot call this function due to different token type")
	}

	return t.tokens[t.currentLine], nil
}

func (t Tokenizer) IntVal() (int, error) {
	token, err := t.TokenType()
	if err != nil {
		return 0, err
	}

	if token != INT_CONST {
		return 0, errors.New("cannot call this function due to different token type")
	}

	intVal, err := strconv.Atoi(t.tokens[t.currentLine])
	if err != nil {
		return 0, err
	}

	return intVal, nil
}

func (t Tokenizer) StringVal() (string, error) {
	token, err := t.TokenType()
	if err != nil {
		return "", err
	}

	if token != STRING_CONST {
		return "", errors.New("cannot call this function due to different token type")
	}

	currentToken := t.tokens[t.currentLine]

	return currentToken[1 : len(currentToken)-2], nil
}
