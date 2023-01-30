package util

import (
	"github.com/gin-gonic/gin"
	"github.com/go-gin-example/pkg/setting"
	"github.com/unknwon/com"
)

func GetPage(ctx *gin.Context) int {
	result := 0
	page, _ := com.StrTo(ctx.Query("page")).Int()
	if page > 0 {
		result = (page - 1) * setting.AppSetting.PageSize
	}

	return result
}
