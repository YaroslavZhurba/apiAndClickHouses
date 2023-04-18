package server

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"limithttp/src/repository"
	"log"
)

// ApiServer
// contains server & IpRepository to count number of requests
// Future: add logger to print Info/Error msgs
type ApiServer struct {
	server       *fasthttp.Server
	ipRepository *repository.IpMaskRepository
}

func NewServer(maskSize, requestLimit uint32, cooldownSeconds float64) *ApiServer {
	apiServer := &ApiServer{
		server:       &fasthttp.Server{Name: "limit_server"},
		ipRepository: repository.NewIpRepository(maskSize, requestLimit, cooldownSeconds),
	}
	apiServer.server.Handler = apiServer.handlerFunc

	return apiServer
}

// TODO: add logger
func (s *ApiServer) Run(port uint) {
	if err := s.server.ListenAndServe(fmt.Sprintf(":%d", port)); err != nil {
		log.Println("error in ApiServer")
	}
}

func (s *ApiServer) handlerFunc(ctx *fasthttp.RequestCtx) {
	s.writeCors(ctx)
	switch string(ctx.Path()) {
	case "/query":
		s.processQuery(ctx)
	case "/clear":
		s.processClear(ctx)
	}
}

func (s *ApiServer) processQuery(ctx *fasthttp.RequestCtx) {
	ipInt, err := getIp(ctx)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	if !s.ipRepository.ProcessRequest(ipInt) { //check limit
		ctx.SetStatusCode(fasthttp.StatusTooManyRequests)
		return
	}

	ctx.SetBody([]byte("Hello")) //return static content if limit is okay
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (s *ApiServer) processClear(ctx *fasthttp.RequestCtx) {
	ipInt, err := getIp(ctx)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	s.ipRepository.ProcessClear(ipInt)

	ctx.SetBody([]byte("Hello")) //return static content if limit is okay
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (s *ApiServer) writeCors(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
