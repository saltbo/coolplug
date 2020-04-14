package rest

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/saltbo/coolplug/core"
	"github.com/saltbo/coolplug/model"
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
	router.GET("/plugins", r.list)
	router.POST("/plugins", r.install)
	router.DELETE("/plugins", r.uninstall)
}

func (r *PluginResource) list(c *gin.Context) {
	plugins := make([]model.Plugin, 0)
	r.engine.Database.Find(&plugins)
	c.JSON(http.StatusOK, plugins)
}

func (r *PluginResource) install(c *gin.Context) {
	name := c.PostForm("name")
	intro := c.PostForm("intro")

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

	mp := &model.Plugin{
		Name:     name,
		Intro:    intro,
		Thumb:    "",
		Status:   0,
		Filename: filename,
	}
	if err := r.engine.PluginInstall(mp); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("plugin install err: %s", err.Error()))
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully with fields name=%s.", file.Filename, name))
}

func (r *PluginResource) uninstall(c *gin.Context) {
	name := c.PostForm("name")

	mp := &model.Plugin{Name: name}
	if err := r.engine.Database.First(mp).Error; err != nil {
		//fmt.Errorf("plugin [%s] not exist", filename)
		return
	}

	if err := r.engine.PluginUninstall(mp.Filename); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("plugin install err: %s", err.Error()))
		return
	}

	//if err := os.Remove(mp.Filename); err != nil {
	//	c.Status(http.StatusBadRequest)
	//	return
	//}

	c.Status(http.StatusOK)
}
