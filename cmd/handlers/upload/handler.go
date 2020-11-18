package uploadHandler

import (
	"github.com/valyala/fasthttp"
)

type Handler interface {
	File(ctx *fasthttp.RequestCtx)
	Photo(ctx *fasthttp.RequestCtx)
}

type handler struct {
	uploadService   uploadService
	uploadTransport uploadTransport
	errorWorker     errorWorker
}

func NewHandler(uploadService uploadService, uploadTransport uploadTransport, errorWorker errorWorker) Handler {
	return &handler{
		uploadService:   uploadService,
		uploadTransport: uploadTransport,
		errorWorker:     errorWorker,
	}
}

func (h *handler) File(ctx *fasthttp.RequestCtx) {
	request, err := h.uploadTransport.FileDecode(ctx)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	response, err := h.uploadService.File(request)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	err = h.uploadTransport.FileEncode(response, ctx)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}
}

func (h *handler) Photo(ctx *fasthttp.RequestCtx) {
	request, err := h.uploadTransport.PhotoDecode(ctx)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	response, err := h.uploadService.Photo(request)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	err = h.uploadTransport.PhotoEncode(response, ctx)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}
}

