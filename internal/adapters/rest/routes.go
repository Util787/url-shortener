package rest

import "github.com/gin-gonic/gin"

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())

	v1 := router.Group("/api/v1")
	v1.Use(newBasicMiddleware(h.log))

	{
		shortener := v1.Group("/")
		{
			shortener.GET("/redirect", h.RedirectURL)
			shortener.GET("/random", h.GetRandomURL)
			shortener.POST("/save", h.SaveURL)
			shortener.POST("/delete", h.DeleteURL)
		}
	}
	return router
}
