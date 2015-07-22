package httputil

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestLinkExtract(t *testing.T) {
	text, err := ioutil.ReadFile("./debug_httputil")
	if err != nil {
		log.Println(err)
	}

	cleaner := NewHTMLCleaner()

	charset := cleaner.DetectCharset([]byte(text))

	ret := cleaner.ToUTF8(text)
	if ret != nil {
		err = ioutil.WriteFile("new.txt", ret, 0666)
		for _, c := range charset {
			log.Println("get charset:", c)
		}
	} else {
		t.Error("charset:", charset, ", error:", ret)
	}
}
