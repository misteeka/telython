package main

import (
	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
	"time"
)

const ip = "127.0.0.1" // 192.168.1.237

var client fasthttp.HostClient
var parserPool fastjson.ParserPool

func init() {
	client = fasthttp.HostClient{
		Addr:                ip + ":8002",
		MaxIdleConnDuration: time.Minute,
		ReadTimeout:         30 * time.Second,
		WriteTimeout:        30 * time.Second,
	}
}

func Get(function string) (*fastjson.Value, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://127.0.0.1:8002/payments/" + function)
	resp := fasthttp.AcquireResponse()
	err := client.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		return nil, err
	}
	p := parserPool.Get()
	value, err := p.ParseBytes(resp.Body())
	parserPool.Put(p)
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
	req.SetRequestURI("http://127.0.0.1:8002/payments/" + function)
	resp := fasthttp.AcquireResponse()
	err := client.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		return nil, err
	}
	p := parserPool.Get()
	value, err := p.ParseBytes(resp.Body())
	parserPool.Put(p)
	if err != nil {
		return nil, err
	}
	ReleaseResponse(resp)
	return value, nil
}
func Put(function string, json string) (*fastjson.Value, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://127.0.0.1:8002/payments/" + function)
	req.Header.SetContentType("application/json")
	req.Header.SetMethodBytes([]byte("PUT"))
	req.SetBody([]byte(json))
	resp := fasthttp.AcquireResponse()
	err := client.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		return nil, err
	}
	p := parserPool.Get()
	value, err := p.ParseBytes(resp.Body())
	parserPool.Put(p)
	if err != nil {
		return nil, err
	}
	ReleaseResponse(resp)
	return value, nil
}
func Delete(function string, json string) (*fastjson.Value, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://127.0.0.1:8002/payments/" + function)
	req.Header.SetContentType("application/json")
	req.Header.SetMethodBytes([]byte("DELETE"))
	req.SetBody([]byte(json))
	resp := fasthttp.AcquireResponse()
	err := client.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		return nil, err
	}
	p := parserPool.Get()
	value, err := p.ParseBytes(resp.Body())
	parserPool.Put(p)
	if err != nil {
		return nil, err
	}
	ReleaseResponse(resp)
	return value, nil
}

func ReleaseResponse(response *fasthttp.Response) {
	fasthttp.ReleaseResponse(response)
}
