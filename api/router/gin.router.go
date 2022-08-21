package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type __router struct{}

var (
	ginDispatcher = gin.Default()
)

func NewGinRouter() Router {
	return &__router{}
}

func (r __router) GET(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	ginDispatcher.GET(uri, gin.WrapF(f))
}

func (r __router) POST(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	ginDispatcher.POST(uri, gin.WrapF(f))
}

func (r __router) SERVE(port string) {
	fmt.Printf(" GIN Http Server is running on port %v:", port)
	http.ListenAndServe(port, ginDispatcher)
}
