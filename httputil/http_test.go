package httputil

import (
	"log"
	"testing"
)

func TestHttp302(t *testing.T) {
	//webpage with http 302
	url := "http://www.aibang.com/detail/beijing-station-jiaoweihutong-dao-shunyi/beijing/"

	httpUtil := NewHttpUtil(30)
	status, html, _, _ := httpUtil.HttpGetByDetail(url, nil, false, true, "")

	if status != 200 || len(html) < 50 {
		t.Error()
	}
}

func TestPost(t *testing.T) {
	//webpage with http 302
	url := "http://epub.sipo.gov.cn/patentoutline.action"

	postData := make(map[string]string)
	postData["pageSize"] = "20"
	postData["numWGSQ"] = "0"
	postData["numSYXX"] = "0"
	postData["numFMSQ"] = "0"
	postData["numFMGB"] = "0"
	postData["showType"] = "1"
	postData["strWord"] = "申请（专利权）人='%百度在线网络技术（北京）有限公司%'"

	httpUtil := NewHttpUtil(30)

	status, html, _, _ := httpUtil.HTTPPostByDetail(url, nil, postData, nil, false, true, "")

	if status != 200 || len(html) < 50 {
		t.Error()
	}
}

func TestHttpBaidu(t *testing.T) {
	//webpage with http 302
	url := "http://www.baidu.com/s?word=%E5%BC%A0%E7%A3%8A"

	httpUtil := NewHttpUtil(30)
	status, html, c, rsp := httpUtil.HttpGetByDetail(url, nil, false, true, "")
	if status != 200 || len(html) < 50 {
		t.Error()
	}
	log.Println("status, html'length, c, rsp:", status, len(html), c, rsp)
}

func TestUrlEncode(t *testing.T) {
	encodeUrl := GetUrlEncode("8pctgRBMALO/DT664deaIFGHU+GqHDoq")
	expectedEncodeUrl := "8pctgRBMALO%2FDT664deaIFGHU%2BGqHDoq"
	if encodeUrl != expectedEncodeUrl {
		t.Error("url encode not work!")
	}
}

func TestCrawlWapdata(t *testing.T) {
	url := "http://m.baidu.com/s?word=%E8%B4%B7%E6%AC%BE"
	httpUtil := NewHttpUtil(40)
	status, html, _, _ := httpUtil.HttpGetByDetail(url, nil, true, true, "")
	if status != 200 || len(html) < 50 {
		t.Error()
	}
	log.Println("status, html:", status, html)
}
