package solver

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"github.com/coreos/goproxy"
)

var (
	_spider = newSpider()
)

type spider struct {
	proxy *goproxy.ProxyHttpServer
}

func Run(port string) {
	_spider.Init()
	_spider.Run(port)
}

func Close() {
	db.Close()
}

func newSpider() *spider {
	sp := &spider{}
	sp.proxy = goproxy.NewProxyHttpServer()
	sp.proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	return sp
}

func (s *spider) Run(port string) {
	log.Println("proxy server at port:" + port)
	log.Fatal(http.ListenAndServe(":"+port, s.proxy))
}

func (s *spider) Init() {
	requestHandleFunc := func(request *http.Request, ctx *goproxy.ProxyCtx) (req *http.Request, resp *http.Response) {
		req = request
		if ctx.Req.URL.Host == `abc.com` {
			resp = new(http.Response)
			resp.StatusCode = 200
			resp.Header = make(http.Header)
			resp.Header.Add("Content-Disposition", "attachment; filename=ca.crt")
			resp.Header.Add("Content-Type", "application/octet-stream")
			resp.Body = ioutil.NopCloser(bytes.NewReader(goproxy.CA_CERT))

		} else if strings.Contains(ctx.Req.URL.Host, "ads") {
			resp = new(http.Response)
			resp.Header = make(http.Header)
			resp.StatusCode = 200
			resp.Body = ioutil.NopCloser(bytes.NewReader([]byte("")))
		}
		return
	}
	responseHandleFunc := func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if resp == nil {
			return resp
		}
		log.Println(ctx.Req.URL.Host + ctx.Req.URL.Path)

		if ctx.Req.URL.Path == "/api/1.0/lplist_matches.json" || ctx.Req.URL.Path == "/api/1.0/lpcreate_match.json" || ctx.Req.URL.Path == "/api/1.0/lpmatch_detail.json" {
			//send letterpress match data to webserver
			bs, _ := ioutil.ReadAll(resp.Body)
			// println(string(bs))
			go setMatch(bs)
			if ctx.Req.URL.Path == "/api/1.0/lpcreate_match.json" {
				go checkVowalsMatch(bs)
			}
			resp.Body = ioutil.NopCloser(bytes.NewReader(bs))
		} else if ctx.Req.URL.Path == "/api/1.0/lp_check_word.json" {
			bs, _ := ioutil.ReadAll(resp.Body)
			if strings.Contains(string(bs), "\"found\":false") {
				go deleteWordDb(ctx.Req.URL.RawQuery)
				exec.Command("adb", "shell", "input", "tap", "50", "50").Run() // tap clear
				go func() {
					time.Sleep(300 * time.Millisecond)
					exec.Command("adb", "shell", "input", "tap", "500", "1000").Run() // tap OK
				}()
			}
			resp.Body = ioutil.NopCloser(bytes.NewReader(bs))
		} else if ctx.Req.URL.Path == "/api/1.0/lpsubmit_turn.json" {
			bs, _ := ioutil.ReadAll(resp.Body)

			if strings.Contains(string(bs), "\"matchStatus\":1,") || strings.Contains(string(bs), "\"matchStatus\":3,") { // status==1 ongoing match,  3 new created match
				go func() {
					exec.Command("adb", "shell", "input", "tap", "50", "50").Run() // tap go back to match list
				}()
			}
			if !strings.Contains(string(bs), "You must enter a valid credentials") {
				go markPlayedWord(bs)
			}

			resp.Body = ioutil.NopCloser(bytes.NewReader(bs))
		} else if ctx.Req.URL.Path == "/api/1.0/lpendmatch_inturn.json" {
			bs, _ := ioutil.ReadAll(resp.Body)
			println(string(bs))
			go func() {
				time.Sleep(500 * time.Millisecond)
				exec.Command("adb", "shell", "input", "tap", "500", "1200").Run() // tap REMOVE GAME
			}()
			resp.Body = ioutil.NopCloser(bytes.NewReader(bs))
		} else if ctx.Req.URL.Path == "/api/1.0/lpload_stats.json" {
			bs, _ := ioutil.ReadAll(resp.Body)
			println(string(bs))
			resp.Body = ioutil.NopCloser(bytes.NewReader(bs))
		} else if strings.Contains(ctx.Req.URL.Host, "ads") {
			bs, _ := ioutil.ReadAll(resp.Body)
			println(string(bs))
			resp.Body = ioutil.NopCloser(bytes.NewReader([]byte("")))
		}
		return resp
	}
	s.proxy.OnResponse().DoFunc(responseHandleFunc)
	s.proxy.OnRequest().DoFunc(requestHandleFunc)
}

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

// Parse query string
func parseURLquery(query string) (m map[string][]string, mk []string, err error) {
	m = make(map[string][]string)
	mk = make([]string, 0)
	for query != "" {
		key := query
		if i := strings.IndexAny(key, "&;"); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = ""
		}
		if key == "" {
			continue
		}
		value := ""
		if i := strings.Index(key, "="); i >= 0 {
			key, value = key[:i], key[i+1:]
		}
		key, err1 := url.QueryUnescape(key)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		value, err1 = url.QueryUnescape(value)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		m[key] = append(m[key], value)
		mk = append(mk, key)
	}
	return
}

// Encode the values
func encodeURLquery(m map[string][]string, mk []string) string {
	var buf bytes.Buffer
	for _, k := range mk {
		vs := m[k]
		prefix := url.QueryEscape(k) + "="
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(url.QueryEscape(v))
		}
	}
	return buf.String()
}

func orPanic(err error) {
	if err != nil {
		panic(err)
	}
}
