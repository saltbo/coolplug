package rest

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/saltbo/coolplug/core"
)

type PluginResource struct {
	engine *core.Engine
}

func NewPluginResource(engine *core.Engine) *PluginResource {
	return &PluginResource{
		engine: engine,
	}
}

func (r *PluginResource) Register(router *gin.RouterGroup) {
	router.POST("/plugins", r.Upload)
}

func (r *PluginResource) Upload(c *gin.Context) {

	name := c.PostForm("name")
	email := c.PostForm("email")

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	filename := filepath.Base(file.Filename)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	if err := r.engine.Load(filename); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("load plugin file err: %s", err.Error()))
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully with fields name=%s and email=%s.", file.Filename, name, email))
}
