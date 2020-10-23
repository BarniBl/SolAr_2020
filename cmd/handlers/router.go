package handlers

import (
	"fmt"
	httputils "github.com/Solar-2020/GoUtils/http"
	postsHandler "github.com/Solar-2020/SolAr_Backend_2020/cmd/handlers/posts"
	uploadHandler "github.com/Solar-2020/SolAr_Backend_2020/cmd/handlers/upload"
	"github.com/Solar-2020/SolAr_Backend_2020/internal/errorWorker"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"runtime/debug"
)

func NewFastHttpRouter(posts postsHandler.Handler, upload uploadHandler.Handler, middleware httputils.Middleware) *fasthttprouter.Router {
	router := fasthttprouter.New()

	//router.Handle("GET", "/health", check)

	router.PanicHandler = panicHandler
	router.Handle("GET", "/health", middleware.Log(httputils.HealthCheckHandler))

	router.Handle("POST", "/api/posts/post", middleware.CORS(posts.Create))
	router.Handle("GET", "/api/posts/posts", middleware.CORS(posts.GetList))

	router.Handle("POST", "/api/upload/photo", middleware.CORS(upload.Photo))
	router.Handle("POST", "/api/upload/file", middleware.CORS(upload.File))

	return router
}

func panicHandler(ctx *fasthttp.RequestCtx, err interface{}) {
	fmt.Printf("Request falied with panic: %s, error: %v\nTrace:\n", string(ctx.Request.RequestURI()), err)
	fmt.Println(string(debug.Stack()))
	errorWorker.NewErrorWorker().ServeFatalError(ctx)
}
