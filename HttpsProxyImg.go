package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"git.code4.in/mobilegameserver/config"
	"git.code4.in/mobilegameserver/logging"
	"git.code4.in/mobilegameserver/unibase"
	"github.com/elazarl/goproxy"
)

func main() {
	config.SetConfig("logfilename", "/tmp/httpimageserver.log")
	unibase.InitServerLogger("HM")
	//var Code string = `<script>alert('test')</script>`
	proxy := goproxy.NewProxyHttpServer()
	//proxy.OnRequest(goproxy.Not(goproxy.ReqHostMatches(regexp.MustCompile("(.*jdb247.*)|(.*umengcloud.*)|(.*openinstall.*)|(.*383014.*)")))).HandleConnect(goproxy.AlwaysMitm)
	//proxy.OnRequest(goproxy.Not(goproxy.UrlMatches(regexp.MustCompile("(ws.*)|(.*websocket)")))).HandleConnect(goproxy.AlwaysMitm)
	//proxy.OnRequest(goproxy.Not(goproxy.UrlMatches(regexp.MustCompile("(ws.*)|(.*websocket)|(.*openinstall.*)|(.*jdb247.*)")))).HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest(goproxy.Not(goproxy.UrlMatches(regexp.MustCompile("(api.*)|(ws.*)|(.*websocket)|(.*openinstall.*)|(.*jdb247.*)|(.*sxxqsw.*)|(.*fungaming.*)")))).HandleConnect(goproxy.AlwaysMitm)
	//proxy.OnRequest(goproxy.UrlMatches(regexp.MustCompile("(.*png)|(.*jpg)|(.*jpeg)|(.*mp3)"))).HandleConnect(goproxy.AlwaysMitm)
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
			ext := filepath.Ext(file)
			if ext == ".mp3" || ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
				p, _ := filepath.Split(file)
				_, err := os.Stat(file)
				_, perr := os.Stat(p)
				if err == nil {
					return r
				} else if os.IsNotExist(perr) {
					os.MkdirAll(p, 0755)
				}

				logging.Debug("下载:%s,%s", file, ctx.Req.URL.String())
				bs, _ := ioutil.ReadAll(r.Body)
				fp, err := os.Create(file)
				if err != nil {
					logging.Debug(err.Error())
				} else {
					fp.Write(bs)
					fp.Close()
				}
				r.Body = ioutil.NopCloser(bytes.NewReader(bs))
			}
			return r
		})
	log.Fatal(http.ListenAndServe(":8081", proxy))
	logging.Final()
}
