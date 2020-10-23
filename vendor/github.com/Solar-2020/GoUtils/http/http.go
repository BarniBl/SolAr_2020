package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

type Message interface{}

type Service interface {
	Address() string
}

type MessageEncodeUrl interface {
	Encode() (urlQuery string, err error)
}

type MessageDecodeBody interface {
	Decode([]byte) (err error)
}

type ServiceEndpoint struct {
	Service Service
	Endpoint string
	Method string
	ContentType string
}

func (e *ServiceEndpoint) Send(message Message, response Message) (err error) {
	var httpResponse *http.Response

	switch e.Method {
	case "GET":
		httpResponse, err = e.sendGet(message)
	case "POST", "PUT", "DELETE":
		httpResponse, err = e.sendWithBody(message)
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	if httpResponse.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad response: %d", httpResponse.StatusCode)
		return
	}
	defer httpResponse.Body.Close()

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	if r, ok := response.(MessageDecodeBody); ok {
		err = r.Decode(body)
	} else {
		err = json.Unmarshal(body, response)
	}
	if err != nil {
		fmt.Println("cannot unmarshal request")
	}
	return
}

func (e *ServiceEndpoint) sendWithBody(request Message) (response *http.Response, err error) {
	body, err := json.Marshal(request)
	if err != nil {
		return
	}

	var queryParams string
	if urlEncode, ok := request.(MessageEncodeUrl); ok {
		queryParams, err = urlEncode.Encode()
		if err != nil {
			return
		}
	}

	client := http.Client{}
	req, err := http.NewRequest(e.Method, e.getFullAddress(queryParams), bytes.NewReader(body))
	if err != nil {
		return
	}
	if e.ContentType == "" {
		e.ContentType = "application/json"
	}
	req.Header.Set("Content-Type", e.ContentType)

	response, err = client.Do(req)
	return
}

func (e *ServiceEndpoint) sendGet(message Message) (response *http.Response, err error) {
	response, err = http.Get(e.getFullAddress())
	return
}

func (e *ServiceEndpoint) getFullAddress(params ...string) string {
	parts := strings.SplitN(e.Service.Address(), "://", 2)
	params_ := append(make([]string, 0), parts[1], e.Endpoint)
	params_ = append(params_, params...)
	return fmt.Sprintf("%s://%s", parts[0], path.Join(params_...))
	//return fmt.Sprintf("%s%s", e.Service.Address(), e.Endpoint)
}

func DecodeDefault(ctx *fasthttp.RequestCtx) (request interface{}, err error) {
	err = json.Unmarshal(ctx.Request.Body(), &request)
	return
}

func EncodeDefault(response interface{}, ctx *fasthttp.RequestCtx) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}