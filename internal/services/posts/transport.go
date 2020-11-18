package posts

import (
	"encoding/json"
	"errors"
	"github.com/Solar-2020/SolAr_Backend_2020/internal/models"
	"github.com/valyala/fasthttp"
	"time"
)

type Transport interface {
	CreateDecode(ctx *fasthttp.RequestCtx) (response models.InputPost, err error)
	CreateEncode(response models.Post, ctx *fasthttp.RequestCtx) (err error)

	GetListDecode(ctx *fasthttp.RequestCtx) (request models.GetPostListRequest, err error)
	GetListEncode(response []models.PostResult, ctx *fasthttp.RequestCtx) (err error)

	SetMarkDecode(ctx *fasthttp.RequestCtx) (request models.MarkPost, err error)
	SetMarkEncode(ctx *fasthttp.RequestCtx) (err error)

	DeletePostDecode(ctx *fasthttp.RequestCtx) (request models.DeletePostRequest, err error)
	DeletePostEncode(ctx *fasthttp.RequestCtx) (err error)
}

type transport struct {
}

func NewTransport() Transport {
	return &transport{}
}

func (t transport) CreateDecode(ctx *fasthttp.RequestCtx) (request models.InputPost, err error) {
	var inputPost models.InputPost
	err = json.Unmarshal(ctx.Request.Body(), &inputPost)
	if err != nil {
		return
	}
	inputPost.PublishDate = time.Now()

	userID, ok := ctx.UserValue("userID").(int)
	if ok {
		request = inputPost
		request.CreateBy = userID
		return
	}
	return request, errors.New("userID not found")
}

func (t transport) CreateEncode(response models.Post, ctx *fasthttp.RequestCtx) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}

func (t transport) GetListDecode(ctx *fasthttp.RequestCtx) (request models.GetPostListRequest, err error) {
	request.GroupID = ctx.QueryArgs().GetUintOrZero("groupID")
	request.Limit = ctx.QueryArgs().GetUintOrZero("limit")

	startFrom := string(ctx.QueryArgs().Peek("startFrom"))
	request.StartFrom, err = time.Parse(time.RFC3339, startFrom)
	if err != nil {
		return
	}


	if ctx.QueryArgs().Has("mark") {
		request.Mark.Defined = true
		request.Mark.Value = ctx.QueryArgs().GetBool("mark")
	}
	//request.UserID = 12
	userID, ok := ctx.UserValue("userID").(int)
	if ok {
		request.UserID = userID
		return
	}
	return
}

func (t transport) GetListEncode(response []models.PostResult, ctx *fasthttp.RequestCtx) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}


func (t transport) SetMarkDecode(ctx *fasthttp.RequestCtx) (request models.MarkPost, err error) {
	request.PostID = ctx.QueryArgs().GetUintOrZero("postId")
	request.GroupID = ctx.QueryArgs().GetUintOrZero("groupId")
	request.Mark = ctx.QueryArgs().GetBool("marked")

	userID, ok := ctx.UserValue("userID").(int)
	if ok {
		request.UserID = userID
		return
	}
	return
}

func (t transport) SetMarkEncode(ctx *fasthttp.RequestCtx) (err error) {
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	return
}

func (t transport) DeletePostDecode(ctx *fasthttp.RequestCtx) (request models.DeletePostRequest, err error) {
	request.PostID = ctx.QueryArgs().GetUintOrZero("postId")
	request.GroupID = ctx.QueryArgs().GetUintOrZero("groupId")

	userID, ok := ctx.UserValue("userID").(int)
	if ok {
		request.UserID = userID
		return
	}
	return
}

func (t transport) DeletePostEncode(ctx *fasthttp.RequestCtx) (err error) {
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	return
}