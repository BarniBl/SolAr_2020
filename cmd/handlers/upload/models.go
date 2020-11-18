package uploadHandler

import (
	"github.com/Solar-2020/SolAr_Backend_2020/internal/models"
	"github.com/valyala/fasthttp"
)

type uploadService interface {
	File(request models.WriteFile) (response models.File, err error)
	Photo(request models.WritePhoto) (response models.Photo, err error)
}

type uploadTransport interface {
	FileDecode(ctx *fasthttp.RequestCtx) (request models.WriteFile, err error)
	FileEncode(response models.File, ctx *fasthttp.RequestCtx) (err error)

	PhotoDecode(ctx *fasthttp.RequestCtx) (request models.WritePhoto, err error)
	PhotoEncode(response models.Photo, ctx *fasthttp.RequestCtx) (err error)
}

type errorWorker interface {
	ServeJSONError(ctx *fasthttp.RequestCtx, serveError error)
}
