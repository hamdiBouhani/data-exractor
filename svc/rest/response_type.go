package rest

import (
	"sm-connector-be/common"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ResponseData(ctx *gin.Context, data interface{}) {
	rsp := gin.H{
		"success": true,
		"data":    data,
	}
	ctx.JSON(200, rsp)
}

// AbortWithError returns a JSON error using a gin.Context
func AbortWithError(
	g *gin.Context,
	status int,
	code string,
	message string,
	err error,
	request interface{},
) {
	data := gin.H{
		"code":    code,
		"message": message,
		"success": false,
		"data":    "",
	}

	if err != nil {
		data["error"] = err.Error()
	}
	logger := common.Log(g.Request.Context())
	logger.WithFields(logrus.Fields(data)).Errorf(`request data: %+v`, request)
	g.AbortWithStatusJSON(status, data)
}

func BindJsonErr(ctx *gin.Context, err error) {
	AbortWithError(ctx, 400, `P/R/0001`, `request json invalid`, err, nil)
}
