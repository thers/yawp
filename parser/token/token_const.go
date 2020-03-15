package token

const (
	_ Token = iota

	ILLEGAL
	EOF
	COMMENT
	KEYWORD

	STRING
	BOOLEAN
	NULL
	NUMBER
	IDENTIFIER

	TEMPLATE_QUOTE // `

	AT   // @
	HASH // #

	PLUS           // +
	MINUS          // -
	MULTIPLY       // *
	SLASH          // /
	REMAINDER      // %
	EXPONENTIATION // **

	AND                  // &
	OR                   // |
	EXCLUSIVE_OR         // ^
	SHIFT_LEFT           // <<
	SHIFT_RIGHT          // >>
	UNSIGNED_SHIFT_RIGHT // >>>
	AND_NOT              // &^

	ADD_ASSIGN       // +=
	SUBTRACT_ASSIGN  // -=
	MULTIPLY_ASSIGN  // *=
	QUOTIENT_ASSIGN  // /=
	REMAINDER_ASSIGN // %=

	AND_ASSIGN                  // &=
	OR_ASSIGN                   // |=
	EXCLUSIVE_OR_ASSIGN         // ^=
	SHIFT_LEFT_ASSIGN           // <<=
	SHIFT_RIGHT_ASSIGN          // >>=
	UNSIGNED_SHIFT_RIGHT_ASSIGN // >>>=
	AND_NOT_ASSIGN              // &^=

	LOGICAL_AND // &&
	LOGICAL_OR  // ||
	INCREMENT   // ++
	DECREMENT   // --

	EQUAL        // ==
	STRICT_EQUAL // ===
	LESS         // <
	GREATER      // >
	ASSIGN       // =
	NOT          // !

	BITWISE_NOT // ~

	NOT_EQUAL        // !=
	STRICT_NOT_EQUAL // !==
	LESS_OR_EQUAL    // <=
	GREATER_OR_EQUAL // >=
	ARROW            // =>

	LEFT_PARENTHESIS  // (
	LEFT_BRACKET      // [
	LEFT_BRACE        // {
	COMMA             // ,
	PERIOD            // .
	OPTIONAL_CHAINING // ?.

	RIGHT_PARENTHESIS  // )
	RIGHT_BRACKET      // ]
	RIGHT_BRACE        // }
	SEMICOLON          // ;
	COLON              // :
	QUESTION_MARK      // ?
	NULLISH_COALESCING // ??
	DOTDOTDOT          // ...

	firstKeyword
	IF
	IN
	DO
	OF

	VAR
	LET
	CONST

	CLASS
	SUPER
	STATIC
	EXTENDS

	IMPORT
	AS
	EXPORT
	FROM

	// TS
	PUBLIC
	PROTECTED
	PRIVATE
	// /TS

	FOR
	NEW
	TRY

	THIS
	ELSE
	CASE
	VOID
	WITH

	WHILE
	BREAK
	CATCH
	THROW

	ASYNC
	AWAIT
	YIELD

	RETURN
	TYPEOF
	DELETE
	SWITCH

	DEFAULT
	FINALLY

	FUNCTION
	CONTINUE
	DEBUGGER

	INSTANCEOF
	lastKeyword
)

var token2string = [...]string{
	ILLEGAL:                     "ILLEGAL",
	EOF:                         "EOF",
	COMMENT:                     "COMMENT",
	KEYWORD:                     "KEYWORD",
	STRING:                      "STRING",
	BOOLEAN:                     "BOOLEAN",
	NULL:                        "NULL",
	NUMBER:                      "NUMBER",
	IDENTIFIER:                  "IDENTIFIER",
	TEMPLATE_QUOTE: 			 "`",
	PLUS:                        "+",
	MINUS:                       "-",
	MULTIPLY:                    "*",
	SLASH:                       "/",
	REMAINDER:                   "%",
	AND:                         "&",
	OR:                          "|",
	EXCLUSIVE_OR:                "^",
	SHIFT_LEFT:                  "<<",
	SHIFT_RIGHT:                 ">>",
	UNSIGNED_SHIFT_RIGHT:        ">>>",
	AND_NOT:                     "&^",
	ADD_ASSIGN:                  "+=",
	SUBTRACT_ASSIGN:             "-=",
	MULTIPLY_ASSIGN:             "*=",
	QUOTIENT_ASSIGN:             "/=",
	REMAINDER_ASSIGN:            "%=",
	AND_ASSIGN:                  "&=",
	OR_ASSIGN:                   "|=",
	EXCLUSIVE_OR_ASSIGN:         "^=",
	SHIFT_LEFT_ASSIGN:           "<<=",
	SHIFT_RIGHT_ASSIGN:          ">>=",
	UNSIGNED_SHIFT_RIGHT_ASSIGN: ">>>=",
	AND_NOT_ASSIGN:              "&^=",
	LOGICAL_AND:                 "&&",
	LOGICAL_OR:                  "||",
	INCREMENT:                   "++",
	DECREMENT:                   "--",
	EQUAL:                       "==",
	STRICT_EQUAL:                "===",
	LESS:                        "<",
	GREATER:                     ">",
	ASSIGN:                      "=",
	NOT:                         "!",
	BITWISE_NOT:                 "~",
	NOT_EQUAL:                   "!=",
	STRICT_NOT_EQUAL:            "!==",
	LESS_OR_EQUAL:               "<=",
	GREATER_OR_EQUAL:            ">=",
	ARROW:                       "=>",
	LEFT_PARENTHESIS:            "(",
	LEFT_BRACKET:                "[",
	LEFT_BRACE:                  "{",
	COMMA:                       ",",
	PERIOD:                      ".",
	OPTIONAL_CHAINING:           "?.",
	RIGHT_PARENTHESIS:           ")",
	RIGHT_BRACKET:               "]",
	RIGHT_BRACE:                 "}",
	SEMICOLON:                   ";",
	COLON:                       ":",
	QUESTION_MARK:               "?",
	NULLISH_COALESCING:          "??",
	EXPONENTIATION:              "**",
	DOTDOTDOT:                   "...",
	IF:                          "if",
	IN:                          "in",
	OF:                          "of",
	DO:                          "do",
	VAR:                         "var",
	LET:                         "let",
	CONST:                       "const",
	FOR:                         "for",
	NEW:                         "new",
	TRY:                         "try",
	THIS:                        "this",
	ELSE:                        "else",
	CASE:                        "case",
	VOID:                        "void",
	WITH:                        "with",
	WHILE:                       "while",
	BREAK:                       "break",
	CATCH:                       "catch",
	THROW:                       "throw",
	ASYNC:                       "async",
	AWAIT:                       "await",
	YIELD:                       "yield",
	RETURN:                      "return",
	TYPEOF:                      "typeof",
	DELETE:                      "delete",
	SWITCH:                      "switch",
	DEFAULT:                     "default",
	FINALLY:                     "finally",
	FUNCTION:                    "function",
	CONTINUE:                    "continue",
	DEBUGGER:                    "debugger",
	INSTANCEOF:                  "instanceof",
	CLASS:                       "class",
	SUPER:                       "super",
	STATIC:                      "static",
	EXTENDS:                     "extends",
	IMPORT:                      "import",
	AS:                          "as",
	EXPORT:                      "export",
}

var keywordTable = map[string]_keyword{
	"if": _keyword{
		token: IF,
	},
	"in": _keyword{
		token: IN,
	},
	"of": _keyword{
		token: OF,
	},
	"do": _keyword{
		token: DO,
	},
	"var": _keyword{
		token: VAR,
	},
	"for": _keyword{
		token: FOR,
	},
	"new": _keyword{
		token: NEW,
	},
	"try": _keyword{
		token: TRY,
	},
	"this": _keyword{
		token: THIS,
	},
	"else": _keyword{
		token: ELSE,
	},
	"case": _keyword{
		token: CASE,
	},
	"void": _keyword{
		token: VOID,
	},
	"with": _keyword{
		token: WITH,
	},
	"while": _keyword{
		token: WHILE,
	},
	"break": _keyword{
		token: BREAK,
	},
	"catch": _keyword{
		token: CATCH,
	},
	"throw": _keyword{
		token: THROW,
	},
	"async": _keyword{
		token: ASYNC,
	},
	"await": _keyword{
		token: AWAIT,
	},
	"yield": _keyword{
		token: YIELD,
	},
	"return": _keyword{
		token: RETURN,
	},
	"typeof": _keyword{
		token: TYPEOF,
	},
	"delete": _keyword{
		token: DELETE,
	},
	"switch": _keyword{
		token: SWITCH,
	},
	"default": _keyword{
		token: DEFAULT,
	},
	"finally": _keyword{
		token: FINALLY,
	},
	"function": _keyword{
		token: FUNCTION,
	},
	"continue": _keyword{
		token: CONTINUE,
	},
	"debugger": _keyword{
		token: DEBUGGER,
	},
	"instanceof": _keyword{
		token: INSTANCEOF,
	},
	"const": _keyword{
		token: CONST,
	},
	"class": _keyword{
		token: CLASS,
	},
	"extends": _keyword{
		token: EXTENDS,
	},
	"super": _keyword{
		token: SUPER,
	},
	"enum": _keyword{
		token:         KEYWORD,
		futureKeyword: true,
	},
	"export": _keyword{
		token: EXPORT,
	},
	"import": _keyword{
		token: IMPORT,
	},
	"as": _keyword{
		token: AS,
	},
	"from": _keyword{
		token: FROM,
	},
	"implements": _keyword{
		token:         KEYWORD,
		futureKeyword: true,
		strict:        true,
	},
	"interface": _keyword{
		token:         KEYWORD,
		futureKeyword: true,
		strict:        true,
	},
	"let": _keyword{
		token: LET,
	},
	"package": _keyword{
		token:         KEYWORD,
		futureKeyword: true,
		strict:        true,
	},
	"private": _keyword{
		token:         PRIVATE,
		futureKeyword: true,
		strict:        true,
	},
	"protected": _keyword{
		token:         PROTECTED,
		futureKeyword: true,
		strict:        true,
	},
	"public": _keyword{
		token:         PUBLIC,
		futureKeyword: true,
		strict:        true,
	},
	"static": _keyword{
		token: STATIC,
	},
}
