package arn

import (
	"github.com/gojp/kana"
	"github.com/ikawaha/kagome/tokenizer"
)

var japaneseTokenizer = tokenizer.New()

// JapaneseToken represents a single token in a sentence.
type JapaneseToken struct {
	Original string
	Hiragana string
	Katakana string
	Romaji   string
}

// NeedsFurigana tells you whether furigana are needed or not.
func (token *JapaneseToken) NeedsFurigana() bool {
	return !kana.IsHiragana(token.Original) && !kana.IsKatakana(token.Original) && !kana.IsLatin(token.Original)
}

// TokenizeJapanese splits the given sentence into tokens.
func TokenizeJapanese(japanese string) []*JapaneseToken {
	var tokens []*JapaneseToken

	for _, token := range japaneseTokenizer.Tokenize(japanese) {
		// Ignore start and end of sentence tokens
		if token.Class == tokenizer.DUMMY {
			continue
		}

		features := token.Features()
		hiragana := ""
		katakana := ""
		romaji := ""

		if len(features) >= 9 {
			romaji = kana.KanaToRomaji(features[8])
			hiragana = kana.RomajiToHiragana(romaji)
			katakana = features[8]
		}

		tokens = append(tokens, &JapaneseToken{
			Original: token.Surface,
			Hiragana: hiragana,
			Katakana: katakana,
			Romaji:   romaji,
		})
	}

	return tokens
}
