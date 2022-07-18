package jack_compiler

import (
	"bufio"
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
