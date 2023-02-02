package v1

import (
	"net/http"

	"github.com/go-gin-example/pkg/app"
	"github.com/go-gin-example/pkg/logging"
	articel_servcie "github.com/go-gin-example/service/article_service"

	"github.com/beego/beego/v2/core/validation"
	"github.com/gin-gonic/gin"
	"github.com/go-gin-example/models"
	"github.com/go-gin-example/pkg/e"
	"github.com/go-gin-example/pkg/setting"
	"github.com/go-gin-example/pkg/util"
	"github.com/unknwon/com"
)

//获取单个文章
func GetArticle(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	var data interface{}
	if valid.HasErrors() {
		app.MakeErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, data)
		return
	}

	articleService := articel_servcie.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, data)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, data)
		return
	}

	article, err := articleService.Get()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, article)
}

//获取多个文章
func GetArticles(c *gin.Context) {
	appG := app.Gin{c}
	data := make(map[string]interface{})
	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()

		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	var tagId int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()

		valid.Min(tagId, 1, "tagId").Message("标签ID必须大于0")
	}

	if valid.HasErrors() {
		app.MakeErrors(valid.Errors)

		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	articleService := articel_servcie.Article{
		TagID: tagId,
		State: state,

		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}

	list, err := articleService.GetAll()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}
	count, err := articleService.Count()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_COUNT_ARTICLE_FAIL, nil)
		return
	}

	data["list"] = list
	data["total"] = count

	appG.Response(http.StatusOK, e.SUCCESS, data)
}

//新增文章
func AddArticle(c *gin.Context) {
	tagId := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	createdBy := c.Query("created_by")
	state := com.StrTo(c.Query("state")).MustInt()
	converImageUrl := c.Query("conver_image_url")

	valid := validation.Validation{}
	valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(desc, "desc").Message("简述不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	valid.Required(converImageUrl, "converImageUrl").Message("图片路径不能为空")
	valid.MaxSize(converImageUrl, 255, "converImageUrl").Message("图片路径长度最长为255字符")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		for _, err := range valid.Errors {
			logging.Info("err key:%s, err message:%s", err.Key, err.Message)
		}

		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": make(map[string]interface{}),
		})
	}

	if !models.ExitsTagById(tagId) {
		code = e.ERROR_NOT_EXIST_TAG
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": make(map[string]interface{}),
		})
	}

	data := make(map[string]interface{})
	data["tag_id"] = tagId
	data["title"] = title
	data["desc"] = desc
	data["content"] = content
	data["created_by"] = createdBy
	data["state"] = state
	data["conver_image_url"] = converImageUrl

	models.AddArticle(data)
	code = e.SUCCESS

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}

//修改文章
func UpdateArticle(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}

	id := com.StrTo(c.Param("id")).MustInt()
	tagId := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	modifiedBy := c.Query("modified_by")
	converImageUrl := c.Query("conver_image_url")

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Min(id, 1, "id").Message("ID必须大于0")
	valid.MaxSize(title, 100, "title").Message("标题最长为100字符")
	valid.MaxSize(converImageUrl, 255, "converImageUrl").Message("图片路径最长为255字符")
	valid.MaxSize(desc, 255, "desc").Message("简述最长为255字符")
	valid.MaxSize(content, 65535, "content").Message("内容最长为65535字符")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")

	code := e.INVALID_PARAMS
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			logging.Info("err key: %s, err message: %s", err.Key, err.Message)
		}

		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": make(map[string]interface{}),
		})
	}
	articleService := articel_servcie.Article{ID: id}
	if _, err := articleService.ExistByID(); err != nil {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
	}

	if !models.ExitsTagById(tagId) {
		code = e.ERROR_NOT_EXIST_TAG
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": make(map[string]interface{}),
		})
	}

	data := make(map[string]interface{})
	if tagId > 0 {
		data["tag_id"] = tagId
	}
	if title != "" {
		data["title"] = title
	}
	if desc != "" {
		data["desc"] = desc
	}
	if content != "" {
		data["content"] = content
	}
	if converImageUrl != "" {
		data["conver_image_url"] = converImageUrl
	}
	data["modified_by"] = modifiedBy

	models.EditArticle(id, data)
	code = e.SUCCESS

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}

//删除文章
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			logging.Info("err key: %s, err message: %s", err.Key, err.Message)
		}

		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": make(map[string]string),
		})
	}

	articelService := articel_servcie.Article{ID: id}
	if _, err := articelService.ExistByID(); err != nil {
		code = e.ERROR_NOT_EXIST_ARTICLE
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": make(map[string]string),
		})
	}

	models.DeleteArticle(id)
	code = e.SUCCESS

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}
