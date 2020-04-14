package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/saltbo/coolplug/core"
	"github.com/saltbo/coolplug/model"
	"github.com/saltbo/coolplug/rest/binding"
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
	router.DELETE("/plugins/:id", r.uninstall)
}

func (r *PluginResource) list(c *gin.Context) {
	plugins := make([]model.Plugin, 0)
	r.engine.Database.Find(&plugins)
	c.JSON(http.StatusOK, plugins)
}

func (r *PluginResource) install(c *gin.Context) {
	p := new(binding.Plugin)
	if err := c.ShouldBind(p); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	filename := fmt.Sprintf("./build/plugins/%s.so", strings.ToLower(p.Name))
	err := c.SaveUploadedFile(p.File, filename)
	if err != nil {
		c.String(http.StatusInternalServerError, "unknown error")
		return
	}

	mp := &model.Plugin{
		Name:     p.Name,
		Intro:    p.Intro,
		Thumb:    "",
		Status:   0,
		Filename: filename,
	}
	if err := r.engine.PluginInstall(mp); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("plugin install err: %s", err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

func (r *PluginResource) uninstall(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	mp := &model.Plugin{Model: gorm.Model{ID: uint(id)}}
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
