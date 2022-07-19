package jack_compiler

type KeywordType uint8

const (
	CLASS KeywordType = iota + 1
	METHOD
	FUNCTION
	CONSTRUCTOR
	INT
	BOOLEAN
	CHAR
	VOID
	VAR
	STATIC
	FIELD
	LET
	DO
	IF
	ELSE
	WHILE
	RETURN
	TRUE
	FALSE
	NULL
	THIS
)

func (t KeywordType) String() string {
	strs := [...]string{
		"CLASS",
		"METHOD",
		"FUNCTION",
		"CONSTRUCTOR",
		"INT",
		"BOOLEAN",
		"CHAR",
		"VOID",
		"VAR",
		"STATIC",
		"FIELD",
		"LET",
		"DO",
		"IF",
		"ELSE",
		"WHILE",
		"RETURN",
		"TRUE",
		"FALSE",
		"NULL",
		"THIS",
	}

	return strs[t-1]
}
