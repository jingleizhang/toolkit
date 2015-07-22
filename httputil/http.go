package httputil

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	kUidLength = 32
)

var kWebUserAgent = []string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/38.0.2125.122 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/38.0.2125.122 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/600.1.25 (KHTML, like Gecko) Version/8.0 Safari/600.1.25",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10) AppleWebKit/600.1.25 (KHTML, like Gecko) Version/8.0 Safari/600.1.25",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/38.0.2125.111 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/38.0.2125.111 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_4) AppleWebKit/537.77.4 (KHTML, like Gecko) Version/7.0.5 Safari/537.77.4",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.57 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/38.0.2125.122 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/38.0.2125.101 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/38.0.2125.104 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/38.0.2125.111 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.65 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.120 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/38.0.2125.111 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/38.0.2125.122 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:33.0) Gecko/20100101 Firefox/33.0",
	"Mozilla/5.0 (Windows NT 6.1; rv:33.0) Gecko/20100101 Firefox/33.0"}

var kWapUserAgent = []string{
	"Mozilla/5.0 (iPad; U; CPU OS 3_2 like Mac OS X; en-us) AppleWebKit/531.21.10 (KHTML, like Gecko) Version/4.0.4 Mobile/7B334b Safari/531.21.10",
	"Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_0 like Mac OS X; en-us) AppleWebKit/532.9 (KHTML, like Gecko) Version/4.0.5 Mobile/8A293 Safari/6531.22.7"}

var kAlphabeta = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

var kDefaultTimeOut = 40

type HttpUtil struct {
	downLoadHosts []string
	strProxy      string
	cleaner       *HTMLCleaner
	client        *http.Client
	random        *rand.Rand
}

func NewHttpUtil(timeOutSeconds int) *HttpUtil {
	if timeOutSeconds <= 0 {
		timeOutSeconds = kDefaultTimeOut
	}

	httpUtil := &HttpUtil{
		cleaner: NewHTMLCleaner(),
		client:  newHttpClient(timeOutSeconds),
		random:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	if httpUtil.client != nil {
		return httpUtil
	}
	return nil
}

func NewProxyHttpUtil(proxy string, timeOutSeconds int) *HttpUtil {
	httpUtil := &HttpUtil{
		cleaner:  NewHTMLCleaner(),
		client:   newProxyHttpClient(proxy, timeOutSeconds),
		strProxy: proxy,
		random:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	if httpUtil.client != nil {
		return httpUtil
	}
	return nil
}

func GetHeaderAccept() string {
	return "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"
}

func GetHeaderAcceptEncode() string {
	return "gzip,deflate,sdch,application/json"
}

func GetHeaderUA() string {
	rand.NewSource(time.Now().UnixNano())
	index := rand.Intn(len(kWebUserAgent))
	return kWebUserAgent[index]
}

func GetHeaderContentType() string {
	return "application/x-www-form-urlencoded; param=value"
}

func newHttpClient(timeOutSeconds int) *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				timeout := time.Duration(timeOutSeconds) * time.Second
				deadline := time.Now().Add(timeout)
				c, err := net.DialTimeout(netw, addr, timeout)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
			DisableKeepAlives:     true,
			ResponseHeaderTimeout: time.Duration(timeOutSeconds) * time.Second,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		},
	}
	return client
}

func newProxyHttpClient(proxy string, timeOutSeconds int) *http.Client {
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		log.Println("error, got not valid proxy. ", proxy, err)
		return nil
	}
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				timeout := time.Duration(timeOutSeconds) * time.Second
				deadline := time.Now().Add(timeout)
				c, err := net.DialTimeout(netw, addr, timeout)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
			DisableKeepAlives:     true,
			ResponseHeaderTimeout: time.Duration(timeOutSeconds) * time.Second,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			Proxy:                 http.ProxyURL(proxyUrl),
		},
	}
	return client
}

func (httpUtil *HttpUtil) getUrlEncode(str string) string {
	return url.QueryEscape(str)
}

func (httpUtil *HttpUtil) getHeaderUA(needWapUserAgent bool) string {
	if needWapUserAgent {
		index := httpUtil.random.Intn(len(kWapUserAgent))
		return kWapUserAgent[index]
	}
	index := httpUtil.random.Intn(len(kWebUserAgent))
	return kWebUserAgent[index]
}

func (httpUtil *HttpUtil) getHeaderContentType() string {
	return "application/x-www-form-urlencoded; param=value"
}

func (httpUtil *HttpUtil) getRandomBaiduUID() string {
	var uid string
	for i := 0; i < kUidLength; i++ {
		uid = uid + kAlphabeta[httpUtil.random.Intn(len(kAlphabeta))]
	}
	uid = uid + ":SL=0:NR=50:FG=1"
	return uid
}

func (httpUtil *HttpUtil) GetStrProxy() string {
	return httpUtil.strProxy
}

func (httpUtil *HttpUtil) isGZHtml(resp *http.Response) bool {
	if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		return true
	}
	return false
}

func (httpUtil *HttpUtil) setHeader(req *http.Request, needWapUserAgent bool) {
	req.Header.Set("Accept", GetHeaderAccept())
	req.Header.Set("Accept-Encoding", "gzip,deflate,sdch")
	req.Header.Set("User-Agent", httpUtil.getHeaderUA(needWapUserAgent))
}

func (httpUtil *HttpUtil) unGzipHtml(html []byte) []byte {
	var b bytes.Buffer
	b.Write(html)

	ungz, err := gzip.NewReader(&b)

	defer ungz.Close()

	if err != nil {
		log.Println("gzip gzip.NewReader error|", err.Error(), "|len(html)|", len(html))
		return html
	}

	content, err := ioutil.ReadAll(ungz)
	if err != nil {
		log.Println("gzip ioutil.ReadAll got error|", err.Error(), "|len(html)|", len(html), "|len(content)|", len(content))
		return html
	}
	return content
}

func (httpUtil *HttpUtil) HttpGetByDetail(url string, cookies []*http.Cookie, needWapUserAgent, needClean bool, queryWords string) (status int, html string, respcookie []*http.Cookie, respinfo string) {
	status = 500
	respcookie = nil

	realurl := "<real_url>" + url + "</real_url>"
	resp := "<query>" + httpUtil.getUrlEncode(queryWords) + "</query>"
	resp += "<original_url>" + url + "</original_url>"

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("NewRequest got error|", err.Error(), "|url|", url)
		return status, "", nil, resp
	}

	if request == nil {
		log.Println("NewRequest got error|request is nil|url|", url)
		return status, "", nil, resp
	}

	httpUtil.setHeader(request, needWapUserAgent)

	if cookies != nil {
		for _, ck := range cookies {
			request.AddCookie(ck)
		}
	}

	response, err := httpUtil.client.Do(request)
	if err != nil {
		log.Println("Client.Do() got error|", err.Error())

		resp += realurl
		return status, "", nil, resp
	}

	defer func() {
		if response != nil && response.Body != nil {
			response.Body.Close()
		}
	}()

	status = response.StatusCode
	respcookie = response.Cookies()
	realurl = "<real_url>" + response.Request.URL.String() + "</real_url>"

	if status != 200 {
		resp += realurl
		return status, "", nil, resp
	}

	htmlByte, err := ioutil.ReadAll(response.Body)

	if err != nil {
		resp += realurl
		log.Println("httpUtil.http get ioutil.ReadAll got error|", err.Error())
		return status, "", nil, resp
	}

	if httpUtil.isGZHtml(response) {
		htmlByte = httpUtil.unGzipHtml(htmlByte)
	}

	resp += "<content_type>" + response.Header.Get("Content-Type") + "</content_type>"
	resp += realurl

	if needClean && len(htmlByte) > 10 {
		utf8Html := httpUtil.cleaner.ToUTF8(htmlByte)
		if utf8Html == nil {
			log.Println("httpUtil.cleaner.ToUTF8 got error|utf8Html is nil|url|", url)
			return status, "", respcookie, resp
		}
		cleanHtml := httpUtil.cleaner.CleanHTML(utf8Html)
		finialHtml := httpUtil.cleaner.RemoveNonUTF8Char(string(cleanHtml))
		return status, finialHtml, respcookie, resp
	}
	return status, string(htmlByte), respcookie, resp
}

func (httpUtil *HttpUtil) HTTPPostByDetail(hostUrl string, extHeaders map[string]string, postData map[string]string, reqCookie []*http.Cookie, needWapUserAgent, needClean bool, queryWords string) (status int, html string, respcookie []*http.Cookie, respinfo string) {
	status = 500
	respcookie = nil
	var resp string
	resp = "<query>" + httpUtil.getUrlEncode(queryWords) + "</query>"
	resp += "<real_url>POST:" + hostUrl + "?query=" + httpUtil.getUrlEncode(queryWords) + "</real_url>"
	resp += "<original_url>" + hostUrl + "</original_url>"

	params := url.Values{}
	for key, value := range postData {
		params.Set(key, value)
	}

	postDataStr := params.Encode()
	postDataBytes := []byte(postDataStr)

	reqest, err := http.NewRequest("POST", hostUrl, bytes.NewReader(postDataBytes))
	if err != nil {
		log.Println("NewRequest got error|", err.Error(), "|url|", hostUrl)
		return status, "", nil, resp
	}
	if reqest == nil {
		log.Println("NewRequest got error|request is nil|url|", hostUrl)
		return status, "", nil, resp
	}

	httpUtil.setHeader(reqest, needWapUserAgent)

	reqest.Header.Set("Content-Type", httpUtil.getHeaderContentType())

	for header, value := range extHeaders {
		reqest.Header.Set(header, value)
	}

	if reqCookie != nil {
		for _, ck := range reqCookie {
			reqest.AddCookie(ck)
		}
	}

	response, err := httpUtil.client.Do(reqest)

	if err != nil {
		log.Println("Client.Do got error|", err.Error())
		return status, "", respcookie, resp
	}

	defer func() {
		if response != nil && response.Body != nil {
			response.Body.Close()
		}
	}()

	status = response.StatusCode
	respcookie = response.Cookies()

	if response.StatusCode != 200 {
		return status, "", respcookie, resp
	}

	htmlByte, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Println("httpUtil.http post ioutil.ReadAll got error|", err.Error())
		return status, "", respcookie, resp
	}

	if httpUtil.isGZHtml(response) {
		htmlByte = httpUtil.unGzipHtml(htmlByte)
	}

	resp += "<content_type>" + response.Header.Get("Content-Type") + "</content_type>"

	if needClean && len(htmlByte) > 10 {
		utf8Html := httpUtil.cleaner.ToUTF8(htmlByte)
		if utf8Html == nil {
			log.Println("httpUtil.cleaner.ToUTF8 got error|", err.Error())
			return status, "", respcookie, resp
		}

		cleanHtml := httpUtil.cleaner.CleanHTML(utf8Html)
		finialHtml := httpUtil.cleaner.RemoveNonUTF8Char(string(cleanHtml))

		return status, finialHtml, respcookie, resp
	}
	return status, string(htmlByte), respcookie, resp
}
