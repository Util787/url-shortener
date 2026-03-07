package rest

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())

	router.GET("/", func(c *gin.Context) {
		// Serve the SPA entry point. Try multiple locations so it works both
		// in local dev (where index.html is in repo root) and in the Docker image
		// (we copy it next to the binary to /usr/local/bin/index.html).
		candidates := []string{"./index.html", "/usr/local/bin/index.html", "./src/index.html"}
		var idxPath string
		for _, p := range candidates {
			if _, err := os.Stat(p); err == nil {
				idxPath = p
				break
			}
		}
		if idxPath == "" {
			c.String(http.StatusNotFound, "index.html not found on server")
			return
		}
		c.File(idxPath)
	})

	// Prevent requests for favicon from being routed to the short-url redirect handler.
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	router.GET("/config.js", func(c *gin.Context) {
		apiBase := strings.TrimRight(h.redirectBaseURL, "/")
		if apiBase == "" {
			scheme := "http"
			if c.Request.TLS != nil {
				scheme = "https"
			}
			apiBase = scheme + "://" + c.Request.Host
		}

		c.Header("Content-Type", "application/javascript; charset=utf-8")
		c.String(http.StatusOK, "window.__API_BASE_URL__ = %q;", apiBase)
	})

	v1 := router.Group("/")

	v1.Use(newBasicMiddleware(h.log))

	{
		shortener := v1.Group("/")
		{
			shortener.GET("/:short_url_id", h.RedirectURL)
			shortener.GET("/random", h.GetRandomURL)
			shortener.POST("/save", h.SaveURL)
			shortener.POST("/delete", h.DeleteURL)
		}
	}
	return router
}
