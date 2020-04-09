package rest

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/saltbo/coolplug/core"
)

type Resource interface {
	Register(router *gin.RouterGroup)
}

func Boot(e *core.Engine) {
	resources := []Resource{
		NewPluginResource(e),
	}

	apiRouter := e.Router.Group("/api")
	for _, resource := range resources {
		resource.Register(apiRouter)
	}

	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err := e.Router.Run(); err != nil {
		log.Fatalln(err)
	}
}
