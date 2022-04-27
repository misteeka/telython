package client

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
	"telython/pkg/http"
	"time"
)

var ParserPool fastjson.ParserPool

type Client struct {
	driver *fasthttp.HostClient
	addr   string
}

func New(_addr string, apiPath string) *Client {
	var addr string
	if apiPath[0] != '/' {
		addr = "https://" + _addr + "/" + apiPath
	} else {
		addr = "https://" + _addr + apiPath
	}
	client := &fasthttp.HostClient{
		Addr:                _addr,
		MaxIdleConnDuration: time.Minute,
		ReadTimeout:         30 * time.Second,
		WriteTimeout:        30 * time.Second,
		IsTLS:               true,
	}
	return &Client{
		driver: client,
		addr:   addr,
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

func (client *Client) Get(function string) (*fastjson.Value, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(client.addr + function)
	resp := fasthttp.AcquireResponse()
	err := client.driver.Do(req, resp)
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
func (client *Client) Post(function string, json string) (*fastjson.Value, error) {
	req := fasthttp.AcquireRequest()
	req.SetBody([]byte(json))
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetRequestURI(client.addr + function)
	resp := fasthttp.AcquireResponse()
	err := client.driver.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(resp.Body()))
	p := ParserPool.Get()
	value, err := p.ParseBytes(resp.Body())
	ParserPool.Put(p)
	if err != nil {
		return nil, err
	}
	ReleaseResponse(resp)
	return value, nil
}
func (client *Client) Put(function string, json string) (*fastjson.Value, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(client.addr + function)
	req.Header.SetContentType("application/json")
	req.Header.SetMethodBytes([]byte("PUT"))
	req.SetBody([]byte(json))
	resp := fasthttp.AcquireResponse()
	err := client.driver.Do(req, resp)
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
func (client *Client) Delete(function string, json string) (*fastjson.Value, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(client.addr + function)
	req.Header.SetContentType("application/json")
	req.Header.SetMethodBytes([]byte("DELETE"))
	req.SetBody([]byte(json))
	resp := fasthttp.AcquireResponse()
	err := client.driver.Do(req, resp)
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
