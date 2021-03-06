package controllers

import (
	"log"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/majimoe/server/util"
)

func showStack(err interface{}) {
	stack := make([]byte, 4096)
	runtime.Stack(stack, true)

	log.Printf("PANIC: %s\n%s", err, stack)
}

func Recovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			showStack(err)
			util.Render.HTML(c.Writer, http.StatusInternalServerError, "error/500", nil)
			c.Abort(0)
		}
	}()

	c.Next()
}
