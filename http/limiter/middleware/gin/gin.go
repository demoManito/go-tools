package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/demoManito/go-tools/http/limiter"
)

// NewMiddleware gin 限流中间件
func NewMiddleware(limter *limiter.Limiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := limter.Wait(ctx); err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		ctx.Next()
	}
}
