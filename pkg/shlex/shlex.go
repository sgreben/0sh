// Modfied for github.com/sgreben/0sh. Original LICENSE below.

/*
Copyright 2012 Google Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package shlex

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// TokenType is a top-level token classification: A word, space, comment, unknown.
type TokenType int

// RuneTokenClass is the type of a UTF-8 character classification: A quote, space, escape.
type RuneTokenClass int

// the internal state used by the lexer state machine
type lexerState int

// Token is a (type, Value) pair representing a lexographical token.
type Token struct {
	Type  TokenType
	Value string
}

// Equal reports whether tokens a, and b, are equal.
// Two tokens are equal if both their types and values are equal. A nil token can
// never be equal to another token.
func (a *Token) Equal(b *Token) bool {
	if a == nil || b == nil {
		return false
	}
	if a.Type != b.Type {
		return false
	}
	return a.Value == b.Value
}

// Named classes of UTF-8 runes
const (
	spaceRunes            = " \t\r\n"
	semicolonRunes        = ";"
	escapingQuoteRunes    = `"`
	nonEscapingQuoteRunes = "'"
	escapeRunes           = `\`
	commentRunes          = "#"
)

// Classes of rune token
const (
	unknownRuneClass RuneTokenClass = iota
	spaceRuneClass
	escapingQuoteRuneClass
	nonEscapingQuoteRuneClass
	escapeRuneClass
	commentRuneClass
	semicolonRuneClass
	eofRuneClass
)

// Classes of lexographic token
const (
	TokenTypeUnknown TokenType = iota
	TokenTypeWord
	TokenTypeWordThenNewline
	TokenTypeNewline
	TokenTypeWordThenSemicolon
	TokenTypeSemicolon
	TokenTypeComment
)

// Lexer state machine states
const (
	LexerStateStart           lexerState = iota // no runes have been seen
	LexerStateInWord                            // processing regular runes in a word
	LexerStateEscaping                          // we have just consumed an escape rune; the next rune is literal
	LexerStateEscapingQuoted                    // we have just consumed an escape rune within a quoted string
	LexerStateQuotingEscaping                   // we are within a quoted string that supports escaping ("...")
	LexerStateQuoting                           // we are within a string that does not support escaping ('...')
	LexerStateComment                           // we are within a comment (everything following an unquoted or unescaped #
)

// TokenClassifier is used for classifying rune characters.
type TokenClassifier map[rune]RuneTokenClass

// AddRuneClass adds a rune class
func (t TokenClassifier) AddRuneClass(runes string, Type RuneTokenClass) {
	for _, runeChar := range runes {
		t[runeChar] = Type
	}
}

// newDefaultClassifier creates a new classifier for ASCII characters.
func newDefaultClassifier() TokenClassifier {
	t := TokenClassifier{}
	t.AddRuneClass(spaceRunes, spaceRuneClass)
	t.AddRuneClass(escapingQuoteRunes, escapingQuoteRuneClass)
	t.AddRuneClass(nonEscapingQuoteRunes, nonEscapingQuoteRuneClass)
	t.AddRuneClass(escapeRunes, escapeRuneClass)
	t.AddRuneClass(commentRunes, commentRuneClass)
	t.AddRuneClass(semicolonRunes, semicolonRuneClass)
	return t
}

// ClassifyRune classifiees a rune
func (t TokenClassifier) ClassifyRune(runeVal rune) RuneTokenClass {
	return t[runeVal]
}

// Tokenizer turns an input stream into a sequence of typed tokens
type Tokenizer struct {
	input      bufio.Reader
	classifier TokenClassifier
}

// NewTokenizer creates a new tokenizer from an input stream.
func NewTokenizer(r io.Reader) *Tokenizer {
	input := bufio.NewReader(r)
	classifier := newDefaultClassifier()
	return &Tokenizer{
		input:      *input,
		classifier: classifier}
}

// scanStream scans the stream for the next token using the internal state machine.
// It will panic if it encounters a rune which it does not know how to handle.
func (t *Tokenizer) scanStream(substWord func(string) string) (*Token, error) {
	state := LexerStateStart
	var inSingleQuotes bool
	var Type TokenType
	var Value []rune
	var nextRune rune
	var nextRuneType RuneTokenClass
	var err error
	unquoteSingle := func(s string) (string, error) {
		return strconv.Unquote("`" + s + "`")
	}
	unquoteDouble := func(s string) (string, error) {
		return strconv.Unquote(`"` + s + `"`)
	}
	unquote := func(word string) (string, error) {
		if inSingleQuotes {
			return unquoteSingle(word)
		}
		return unquoteDouble(word)
	}
	substWordIfShould := func(word string) string {
		if inSingleQuotes {
			return word
		}
		return substWord(word)
	}
	postProcessWord := func(word string) (string, error) {
		word, err := unquote(word)
		if err != nil {
			return "", err
		}
		return substWordIfShould(word), nil
	}
	for {
		nextRune, _, err = t.input.ReadRune()
		nextRuneType = t.classifier.ClassifyRune(nextRune)

		if err == io.EOF {
			nextRuneType = eofRuneClass
			err = nil
		} else if err != nil {
			return nil, err
		}

		switch state {
		case LexerStateStart: // no runes read yet
			{
				switch nextRuneType {
				case eofRuneClass:
					{
						return nil, io.EOF
					}
				case spaceRuneClass:
					switch nextRune {
					case '\n':
						token := &Token{
							Type:  TokenTypeNewline,
							Value: "\n"}
						return token, nil
					}
				case semicolonRuneClass:
					token := &Token{
						Type:  TokenTypeSemicolon,
						Value: ";"}
					return token, nil
				case escapingQuoteRuneClass:
					{
						Type = TokenTypeWord
						state = LexerStateQuotingEscaping
					}
				case nonEscapingQuoteRuneClass:
					{
						Type = TokenTypeWord
						state = LexerStateQuoting
					}
				case escapeRuneClass:
					{
						Type = TokenTypeWord
						state = LexerStateEscaping
					}
				case commentRuneClass:
					{
						Type = TokenTypeComment
						state = LexerStateComment
					}
				default:
					{
						Type = TokenTypeWord
						Value = append(Value, nextRune)
						state = LexerStateInWord
					}
				}
			}
		case LexerStateInWord: // in a regular word
			{
				switch nextRuneType {
				case eofRuneClass:
					{
						value, err := postProcessWord(string(Value))
						if err != nil {
							return nil, err
						}
						token := &Token{
							Type:  Type,
							Value: value,
						}
						return token, err
					}
				case spaceRuneClass:
					{
						value, err := postProcessWord(string(Value))
						if err != nil {
							return nil, err
						}
						switch nextRune {
						case '\n':
							token := &Token{
								Type:  TokenTypeWordThenNewline,
								Value: value,
							}
							return token, err
						}
						token := &Token{
							Type:  Type,
							Value: value,
						}
						return token, err
					}
				case semicolonRuneClass:
					{
						value, err := postProcessWord(string(Value))
						if err != nil {
							return nil, err
						}
						token := &Token{
							Type:  TokenTypeWordThenSemicolon,
							Value: value}
						return token, err
					}
				case escapingQuoteRuneClass:
					{
						state = LexerStateQuotingEscaping
						inSingleQuotes = false
					}
				case nonEscapingQuoteRuneClass:
					{
						state = LexerStateQuoting
						inSingleQuotes = true
					}
				case escapeRuneClass:
					{
						state = LexerStateEscaping
					}
				default:
					{
						Value = append(Value, nextRune)
					}
				}
			}
		case LexerStateEscaping: // the rune after an escape character
			{
				switch nextRuneType {
				case eofRuneClass:
					{
						err = fmt.Errorf("EOF found after escape character")
						token := &Token{
							Type:  Type,
							Value: string(Value)}
						return token, err
					}
				default:
					{
						state = LexerStateInWord
						switch nextRune {
						case '"', '\'', ';':
							state = LexerStateInWord
							Value = append(Value, nextRune)
						default:
							state = LexerStateInWord
							Value = append(Value, '\\', nextRune)
						}
					}
				}
			}
		case LexerStateEscapingQuoted: // the next rune after an escape character, in double quotes
			{
				switch nextRuneType {
				case eofRuneClass:
					{
						err = fmt.Errorf("EOF found after escape character")
						token := &Token{
							Type:  Type,
							Value: string(Value)}
						return token, err
					}
				default:
					{
						switch nextRune {
						case '"', '\'':
							state = LexerStateQuotingEscaping
							Value = append(Value, nextRune)
						default:
							state = LexerStateQuotingEscaping
							Value = append(Value, '\\', nextRune)
						}
					}
				}
			}
		case LexerStateQuotingEscaping: // in escaping double quotes
			{
				switch nextRuneType {
				case eofRuneClass:
					{
						err = fmt.Errorf("EOF found when expecting closing quote")
						token := &Token{
							Type:  Type,
							Value: string(Value)}
						return token, err
					}
				case escapingQuoteRuneClass:
					{
						state = LexerStateInWord
					}
				case escapeRuneClass:
					{
						state = LexerStateEscapingQuoted
					}
				default:
					{
						Value = append(Value, nextRune)
					}
				}
			}
		case LexerStateQuoting: // in non-escaping single quotes
			{
				switch nextRuneType {
				case eofRuneClass:
					{
						err = fmt.Errorf("EOF found when expecting closing quote")
						token := &Token{
							Type:  Type,
							Value: string(Value)}
						return token, err
					}
				case nonEscapingQuoteRuneClass:
					{
						state = LexerStateInWord
					}
				default:
					{
						Value = append(Value, nextRune)
					}
				}
			}
		case LexerStateComment: // in a comment
			{
				switch nextRuneType {
				case eofRuneClass:
					{
						token := &Token{
							Type:  Type,
							Value: string(Value)}
						return token, err
					}
				case spaceRuneClass:
					{
						if nextRune == '\n' {
							state = LexerStateStart
							token := &Token{
								Type:  Type,
								Value: string(Value)}
							return token, err
						}
						Value = append(Value, nextRune)
					}
				default:
					{
						Value = append(Value, nextRune)
					}
				}
			}
		default:
			{
				return nil, fmt.Errorf("Unexpected state: %v", state)
			}
		}
	}
}

// Next returns the next token in the stream.
func (t *Tokenizer) Next(substWord func(string) string) (*Token, error) {
	return t.scanStream(substWord)
}
