package pieql

type Token int

const (
	//Special Tokens
	ILLEGAL Token = iota
	EOF
	WS

	// Miscellaneous
	COMMA

	//Literals
	literal_beg
	IDENT
	literal_end

	// Operators
	operator_beg
	MUL
	operator_end

	// Keywords
	keyword_beg
	SELECT
	FROM
	keyword_end
)
