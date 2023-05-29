package args

type TokenStream struct {
	toks []Token
	cur  int
}

func NewTokenStream(toks []Token) *TokenStream {
	return &TokenStream{toks: toks, cur: 0}
}

func (s *TokenStream) Next() (*Token, bool) {
	if s.cur >= len(s.toks) {
		return nil, false
	}

	s.cur += 1
	return &s.toks[s.cur-1], true
}

func (s *TokenStream) Peek() (*Token, bool) {
	if s.cur >= len(s.toks) {
		return nil, false
	}

	return &s.toks[s.cur], true
}
