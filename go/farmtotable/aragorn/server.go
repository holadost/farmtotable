package aragorn

import "github.com/gin-gonic/gin"

func RunServer() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})
	r.Run()
}
