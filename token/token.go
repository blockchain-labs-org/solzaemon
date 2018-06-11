package token

type Token int

const (
	ILLEGAL Token = iota

	IDENT  // contract
	INT    // 10
	STRING // "string"

	ADD // +
	SUB // -
	MUL // *
	POW // **
	QUO // /
	REM // %

	AND     // &
	OR      // |
	XOR     // ^
	SHL     // <<
	SHR     // >>
	AND_NOT // &^

	ASSIGN // =
	EQ     // ==

	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	PERIOD // .

	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	SEMICOLON // ;
	COLON     // :
)
