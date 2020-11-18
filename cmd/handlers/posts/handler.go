package postsHandler

import (
	"github.com/valyala/fasthttp"
)

type Handler interface {
	Create(ctx *fasthttp.RequestCtx)
	GetList(ctx *fasthttp.RequestCtx)
	Mark(ctx *fasthttp.RequestCtx)
	Delete(ctx *fasthttp.RequestCtx)
}

type handler struct {
	postService   postService
	postTransport postTransport
	errorWorker   errorWorker
}

func NewHandler(postService postService, postTransport postTransport, errorWorker errorWorker) Handler {
	return &handler{
		postService:   postService,
		postTransport: postTransport,
		errorWorker:   errorWorker,
	}
}

func (h *handler) Create(ctx *fasthttp.RequestCtx) {
	post, err := h.postTransport.CreateDecode(ctx)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	postReturn, err := h.postService.Create(post)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	err = h.postTransport.CreateEncode(postReturn, ctx)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}
}

func (h *handler) GetList(ctx *fasthttp.RequestCtx) {
	getPostListRequest, err := h.postTransport.GetListDecode(ctx)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	postReturn, err := h.postService.GetList(getPostListRequest)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	err = h.postTransport.GetListEncode(postReturn, ctx)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}
}

func (h *handler) Mark(ctx *fasthttp.RequestCtx) {
	markRequest, err := h.postTransport.SetMarkDecode(ctx)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	err = h.postService.SetMark(markRequest)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	err = h.postTransport.SetMarkEncode(ctx)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}
}

func (h *handler) Delete(ctx *fasthttp.RequestCtx) {
	request, err := h.postTransport.DeletePostDecode(ctx)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	err = h.postService.Delete(request)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	err = h.postTransport.DeletePostEncode(ctx)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}
}

