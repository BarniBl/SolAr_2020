package upload

import (
	"encoding/json"
	"errors"
	"github.com/Solar-2020/SolAr_Backend_2020/internal/models"
	"github.com/valyala/fasthttp"
)

type Transport interface {
	FileDecode(ctx *fasthttp.RequestCtx) (request models.WriteFile, err error)
	FileEncode(response models.File, ctx *fasthttp.RequestCtx) (err error)

	PhotoDecode(ctx *fasthttp.RequestCtx) (request models.WritePhoto, err error)
	PhotoEncode(response models.Photo, ctx *fasthttp.RequestCtx) (err error)
}

type transport struct {
}

func NewTransport() Transport {
	return &transport{}
}

func (t transport) FileDecode(ctx *fasthttp.RequestCtx) (request models.WriteFile, err error) {
	//body := ctx.FormValue("body")
	file, err := ctx.FormFile("file")
	if err != nil {
		return
	}

	//err = json.Unmarshal(body, &request)
	//if err != nil {
	//	return
	//}

	request.Name = file.Filename
	request.File = file
	userID, ok := ctx.UserValue("userID").(int)
	if ok {
		request.UserID = userID
		return
	}
	return request, errors.New("userID not found")
}

func (t transport) FileEncode(response models.File, ctx *fasthttp.RequestCtx) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}

func (t transport) PhotoDecode(ctx *fasthttp.RequestCtx) (request models.WritePhoto, err error) {
	//body := ctx.FormValue("body")
	file, err := ctx.FormFile("file")
	if err != nil {
		return
	}

	//err = json.Unmarshal(body, &request)
	//if err != nil {
	//	return
	//}
	request.Name = file.Filename

	request.File = file
	userID, ok := ctx.UserValue("userID").(int)
	if ok {
		request.UserID = userID
		return
	}
	return request, errors.New("userID not found")
}

func (t transport) PhotoEncode(response models.Photo, ctx *fasthttp.RequestCtx) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}
