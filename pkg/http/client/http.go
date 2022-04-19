package client

import (
	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
	"telython/pkg/http"
	"time"
)

var client fasthttp.HostClient
var ParserPool fastjson.ParserPool
var addr string

func Init(_addr string, apiPath string) {
	if apiPath[0] != '/' {
		addr = "http://" + _addr + "/" + apiPath
	} else {
		addr = "http://" + _addr + apiPath
	}
	client = fasthttp.HostClient{
		Addr:                _addr,
		MaxIdleConnDuration: time.Minute,
		ReadTimeout:         30 * time.Second,
		WriteTimeout:        30 * time.Second,
	}
}

func GetError(value *fastjson.Value) *http.Error {
	if value == nil {
		return nil
	}
	if value.Exists("error") {
		return &http.Error{
			Code:    value.GetUint64("error"),
			Message: string(value.GetStringBytes("message")),
		}
	} else {
		return nil
	}
}

func Get(function string) (*fastjson.Value, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(addr + function)
	resp := fasthttp.AcquireResponse()
	err := client.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		return nil, err
	}
	p := ParserPool.Get()
	value, err := p.ParseBytes(resp.Body())
	ParserPool.Put(p)
	if err != nil {
		return nil, err
	}
	ReleaseResponse(resp)
	return value, nil
}
func Post(function string, json string) (*fastjson.Value, error) {
	req := fasthttp.AcquireRequest()
	req.SetBody([]byte(json))
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetRequestURI(addr + function)
	resp := fasthttp.AcquireResponse()
	err := client.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		return nil, err
	}
	p := ParserPool.Get()
	value, err := p.ParseBytes(resp.Body())
	ParserPool.Put(p)
	if err != nil {
		return nil, err
	}
	ReleaseResponse(resp)
	return value, nil
}
func Put(function string, json string) (*fastjson.Value, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(addr + function)
	req.Header.SetContentType("application/json")
	req.Header.SetMethodBytes([]byte("PUT"))
	req.SetBody([]byte(json))
	resp := fasthttp.AcquireResponse()
	err := client.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		return nil, err
	}
	p := ParserPool.Get()
	value, err := p.ParseBytes(resp.Body())
	ParserPool.Put(p)
	if err != nil {
		return nil, err
	}
	ReleaseResponse(resp)
	return value, nil
}
func Delete(function string, json string) (*fastjson.Value, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(addr + function)
	req.Header.SetContentType("application/json")
	req.Header.SetMethodBytes([]byte("DELETE"))
	req.SetBody([]byte(json))
	resp := fasthttp.AcquireResponse()
	err := client.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		return nil, err
	}
	p := ParserPool.Get()
	value, err := p.ParseBytes(resp.Body())
	ParserPool.Put(p)
	if err != nil {
		return nil, err
	}
	ReleaseResponse(resp)
	return value, nil
}

func ReleaseResponse(response *fasthttp.Response) {
	fasthttp.ReleaseResponse(response)
}
