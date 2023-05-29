package args

import (
	"fmt"
	"os"
	"strings"
)

type TokenKind string

const (
	TokenName  = "Name"
	TokenLong  = "Long"
	TokenShort = "Short"
)

type Token struct {
	Kind  TokenKind
	Value string
}

func Scan(args []string) []Token {

	toks := []Token{}

	for _, arg := range args {
		toks = append(toks, scanArg(arg))
	}

	return toks
}

func scanArg(arg string) Token {

	if strings.HasPrefix(arg, "--") {
		return Token{Kind: TokenLong, Value: arg[2:]}
	} else if strings.HasPrefix(arg, "-") {
		return Token{Kind: TokenShort, Value: arg[1:]}
	}

	return Token{Kind: TokenName, Value: arg}
}

type ArgKind int

const (
	ArgBool = iota
	ArgString
)

type Arg struct {
	Kind      ArgKind
	Long      string
	Short     string
	Desc      string
	BoolVal   *bool
	StringVal *string
}

type Parser struct {
	flags []Arg
	names []string
}

func (p *Parser) GetFlag(long string, short string) (*Arg, bool) {
	for i := 0; i < len(p.flags); i++ {
		flag := &p.flags[i]
		if long != "" {
			if flag.Long == long {
				return flag, true
			}
		} else {
			if flag.Short == short {
				return flag, true
			}
		}
	}

	return nil, false
}

func (p *Parser) Parse(toks []Token) error {
	s := NewTokenStream(toks)
	p.names = []string{}

	for {
		t, ok := s.Next()
		if !ok {
			break
		}

		isLong := false

		switch t.Kind {

		case TokenLong:
			isLong = true
			fallthrough
		case TokenShort:

			var flag *Arg
			var ok bool

			if isLong {
				flag, ok = p.GetFlag(t.Value, "")
			} else {
				flag, ok = p.GetFlag("", t.Value)
			}

			if !ok {
				return fmt.Errorf("Invalid flag: '%s'", t.Value)
			}

			switch flag.Kind {
			case ArgBool:
				*flag.BoolVal = true
			case ArgString:
				nextTok, ok := s.Next()
				if !ok || nextTok.Kind != TokenName {
					return fmt.Errorf("Error parsing flag '%s'", t.Value)
				}
				*flag.StringVal = nextTok.Value
			}

		case TokenName:
			p.names = append(p.names, t.Value)
		}
	}

	return nil
}

var parser Parser

func AddBool(long string, short string, ptr *bool, desc string) {
	parser.flags = append(parser.flags, Arg{
		Kind:    ArgBool,
		Long:    long,
		Short:   short,
		BoolVal: ptr,
		Desc:    desc,
	})
}

func AddString(long string, short string, ptr *string, desc string) {
	parser.flags = append(parser.flags, Arg{
		Kind:      ArgString,
		Long:      long,
		Short:     short,
		StringVal: ptr,
		Desc:      desc,
	})
}

func Parse() ([]string, error) {
	return ParseAt(1)
}

func ParseAt(offset int) ([]string, error) {
	return ParseArgs(os.Args[offset:])
}

func ParseArgs(args []string) ([]string, error) {
	toks := Scan(args)

	err := parser.Parse(toks)
	if err != nil {
		return nil, err
	}

	return parser.names, nil
}
