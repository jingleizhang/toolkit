package httputil

import (
	"code.google.com/p/mahonia"
	"github.com/saintfish/chardet"
	"log"
	"strings"
	"unicode/utf8"
)

type HTMLCleaner struct {
	detector *chardet.Detector
	gbk      mahonia.Decoder
	big5     mahonia.Decoder
}

func NewHTMLCleaner() *HTMLCleaner {
	ret := HTMLCleaner{}
	ret.detector = chardet.NewHtmlDetector()
	ret.gbk = mahonia.NewDecoder("gb18030")
	ret.big5 = mahonia.NewDecoder("big5")
	return &ret
}

func (self *HTMLCleaner) DetectCharset(html []byte) []string {
	ret, err := self.detector.DetectAll(html)
	if err != nil {
		return nil
	}

	var result []string
	for _, c := range ret {
		result = append(result, strings.ToLower(c.Charset))
	}
	return result
}

func (self *HTMLCleaner) CleanHTML(src []byte) []byte {
	dst := []byte{}

	prev := byte(0)
	for i, ch := range src {
		if ch <= 32 || ch == 127 {
			if i > 0 && prev > 32 && prev != 127 {
				dst = append(dst, 32)
				prev = ch
			}
		} else {
			dst = append(dst, ch)
			prev = ch
		}
	}
	return dst
}

func (self *HTMLCleaner) isEmojiCharacter(b rune) bool {
	return len(string(b)) > 3
}

// filter non utf8 and emoji character

func (self *HTMLCleaner) RemoveEmojiChar(html string) string {
	v := make([]rune, 0, len(html)/2)
	for _, r := range html {
		if self.isEmojiCharacter(r) {
			continue
		}
		v = append(v, r)
	}
	return string(v)
}

func (self *HTMLCleaner) RemoveNonUTF8Char(html string) string {
	html = self.RemoveEmojiChar(html)
	if !utf8.ValidString(html) {
		// filter emoji
		v := make([]rune, 0, len(html))
		for i, r := range html {
			if r == utf8.RuneError {
				_, size := utf8.DecodeRuneInString(html[i:])
				if size == 1 {
					continue
				}
			}
			v = append(v, r)
		}

		html = string(v)
		return html
	}
	return html
}

func (self *HTMLCleaner) ToUTF8(html []byte) []byte {
	charsetSets := self.DetectCharset(html)
	detectedCharSets := "|"
	for _, charset := range charsetSets {
		// for log
		detectedCharSets += charset
		detectedCharSets += "|"

		if !strings.Contains(charset, "gb") && !strings.Contains(charset, "big") {
			charset = "utf-8"
		}
		if charset == "utf-8" || charset == "utf8" {
			return html
		} else if charset == "gb2312" || charset == "gb-2312" || charset == "gbk" || charset == "gb18030" || charset == "gb-18030" {
			ret, ok := self.gbk.ConvertStringOK(string(html))
			if ok {
				return []byte(ret)
			} else {
				continue
			}
		} else if charset == "big5" {
			ret, ok := self.big5.ConvertStringOK(string(html))
			if ok {
				return []byte(ret)
			} else {
				continue
			}
		}
	}
	log.Println("ERROR|chaset convert to utf8 failed|detectedCharSets:", detectedCharSets, "|charsetSets:", charsetSets)
	return nil
}
