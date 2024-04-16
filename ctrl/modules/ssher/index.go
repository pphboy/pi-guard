package ssher

import "github.com/gin-gonic/gin"

type Ssher interface {
	RegisterRoute(*gin.RouterGroup)
}

func NewSsher() Ssher {
	return nil
}
