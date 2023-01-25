package main

import (
	"net/http"

	"aiwork.io/marketplace/internal"
	"aiwork.io/marketplace/routes/admin"
	"aiwork.io/marketplace/routes/anonymous"
	"aiwork.io/marketplace/routes/auth"
	"aiwork.io/marketplace/routes/client"
	"aiwork.io/marketplace/routes/tasks"
	"aiwork.io/marketplace/routes/users"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	xctx := internal.NewXContext()
	xctx.Init()

	server := internal.NewHttpServer()
	server.LoadHTMLGlob("public/*")
	server.GET("/verification", func(c *gin.Context) {
		c.HTML(http.StatusOK, "verification.html", gin.H{
			"state": c.Query("state"),
		})
	})

	// 3rd middlewares
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowHeaders = []string{"*"}
	server.Use(cors.New(config))

	anonymous.NewRouter(server.Group("/"), xctx)
	auth.NewRouter(server.Group("/auth"), xctx)
	admin.NewRouter(server.Group("/admin"), xctx)
	users.NewRouter(server.Group("/users"), xctx)
	tasks.NewRouter(server.Group("/tasks"), xctx)
	client.NewRouter(server.Group("/client"), xctx)

	server.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
