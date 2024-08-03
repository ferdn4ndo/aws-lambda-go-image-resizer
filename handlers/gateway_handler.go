package handlers

import (
	"path/filepath"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/core"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
)

type GatewayHandler struct {
	initialized       bool
	muxLambda         *gorillamux.GorillaMuxAdapter
	healthHandler     *HealthHandler
	resizeCropHandler *ResizeCropHandler
}

func (gateway *GatewayHandler) ServeHTTP(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request.Path = filepath.Clean(request.Path)
	gateway.initRouter()

	return gateway.muxLambda.Proxy(request)
}

func (gateway *GatewayHandler) initRouter() {
	if !gateway.initialized {
		router := mux.NewRouter().StrictSlash(false)
		gateway.healthHandler = new(HealthHandler)
		gateway.resizeCropHandler = new(ResizeCropHandler)

		router.HandleFunc("/health", gateway.healthHandler.ServeHTTP)
		router.HandleFunc("/{optional}", gateway.resizeCropHandler.ServeHTTP).Methods("GET")

		gateway.muxLambda = gorillamux.New(router)
		gateway.initialized = true
	}
}
