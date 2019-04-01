package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/elazarl/goproxy"
)

func main() {
	//var Code string = `<script>alert('test')</script>`
	proxy := goproxy.NewProxyHttpServer()
	//proxy.OnRequest(goproxy.Not(goproxy.ReqHostMatches(regexp.MustCompile("(.*jdb247.*)|(.*umengcloud.*)|(.*openinstall.*)|(.*383014.*)")))).HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest(goproxy.Not(goproxy.UrlMatches(regexp.MustCompile("(ws.*)|(.*websocket)")))).HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			r.Header.Set("X-GoProxy", "yxorPoG-X")
			return r, nil
		})
	proxy.OnResponse().DoFunc(
		func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			if ctx.Req.Method != "GET" {
				return r
			}
			file := ctx.Req.Host + ctx.Req.URL.Path
			p, _ := filepath.Split(file)
			_, err := os.Stat(file)
			_, perr := os.Stat(p)
			if err == nil {
				return r
			} else if os.IsNotExist(perr) {
				os.MkdirAll(p, 0755)
			}

			log.Println("URL:", ctx.Req.URL, "PATH:", ctx.Req.URL.Path, "Accept:", ctx.Req.Header["Accept"])
			bs, _ := ioutil.ReadAll(r.Body)
			fp, err := os.Create(file)
			if err != nil {
				log.Println(err)
			} else {
				fp.Write(bs)
				fp.Close()
			}
			r.Body = ioutil.NopCloser(bytes.NewReader(bs))
			return r
		})
	log.Fatal(http.ListenAndServe(":8080", proxy))
}
