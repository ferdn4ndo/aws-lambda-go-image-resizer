package handlers

import (
	"fmt"
	"net/http"
)

type HealthHandler struct{}

func (handler *HealthHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	_, _ = fmt.Fprintf(writer, "I'm fine")
}
