package scanner

import (
	"unicode"

	"github.com/blockchain-labs-org/solzaemon/token"
)

var tokMap = map[rune]token.Token{
	'(': token.LPAREN,
	')': token.RPAREN,
	'{': token.LBRACE,
	'}': token.RBRACE,
	'[': token.LBRACK,
	']': token.RBRACK,
	';': token.SEMICOLON,
	'.': token.PERIOD,
	'^': token.XOR,
}

// Scanner is lexical scanner
type Scanner struct {
	src []rune
	pos int
}

func NewScanner(src []rune) *Scanner {
	return &Scanner{src: src}
}

func (s *Scanner) peek() rune {
	return s.src[s.pos]
}

func (s *Scanner) next() rune {
	if s.pos >= len(s.src) {
		return 0
	}

	ret := s.src[s.pos]
	s.pos++
	return ret
}

func (s *Scanner) backup() {
	s.pos--
}

func (s *Scanner) skipBlank() {
	for isBlank(s.peek()) {
		s.next()
	}
}

func (s *Scanner) Scan() (tok token.Token, lit string) {
	if s.pos >= len(s.src) {
		return 0, ""
	}
	s.skipBlank()
	switch ch := s.next(); {
	case ch == '"':
		s.backup()
		tok = token.STRING
		lit = s.scanString()
		return
	case isLetter(ch):
		s.backup()
		tok = token.IDENT
		lit = s.scanIdent()
		return
	case isDigit(ch):
		s.backup()
		tok = token.INT
		lit = s.scanNumber()
		return
	case ch == '=':
		if s.next() == '=' {
			tok = token.EQ
			lit = "=="
		} else {
			s.backup()
			tok = token.ASSIGN
			lit = string(ch)
		}
	case ch == '*':
		if s.next() == '*' {
			tok = token.POW
			lit = "**"
		} else {
			s.backup()
			tok = token.MUL
			lit = string(ch)
		}
	default:
		tk, ok := tokMap[ch]
		if ok {
			tok = tk
			lit = string(ch)
		} else {
			tok = token.ILLEGAL
			lit = string(ch)
		}
	}
	return
}

func (s *Scanner) scanUntilSemicolon() string {
	var ret []rune
	started := false
done:
	for {
		switch ch := s.next(); {
		case ch == ';':
			ret = append(ret, ch)
			if started {
				break done
			}
			started = true
		default:
			ret = append(ret, ch)
		}
	}
	return string(ret)

}

func (s *Scanner) scanIdent() string {
	var ret []rune
	for {
		ch := s.next()
		if !isLetter(ch) && !isDigit(ch) {
			s.backup()
			break
		}
		ret = append(ret, ch)
	}
	return string(ret)
}

func (s *Scanner) scanString() string {
	var ret []rune
	started := false
done:
	for {
		switch ch := s.next(); {
		case ch == '"':
			ret = append(ret, ch)
			if started {
				break done
			}
			started = true
		default:
			ret = append(ret, ch)
		}
	}
	return string(ret)
}

func (s *Scanner) scanNumber() string {
	var ret []rune
done:
	for {
		switch ch := s.next(); {
		case isDigit(ch):
			ret = append(ret, ch)
		default:
			s.backup()
			break done
		}
	}
	return string(ret)
}

func isBlank(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n'
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}
