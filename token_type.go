package jack_compiler

type TokenType uint8

const (
	KEYWORD TokenType = iota + 1
	SYMBOL
	IDENTIFIER
	INT_CONST
	STRING_CONST
)

func (t TokenType) String() string {
	strs := [...]string{"KEYWORD", "SYMBOL", "IDENTIFIER", "INT_CONST", "STRING_CONST"}
	return strs[t-1]
}
