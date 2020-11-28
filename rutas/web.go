package rutas

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sam/controladores"
)

func Init() *gin.Engine {
	r := gin.Default()
	r.Use(CORSMiddleware())

	r.StaticFS("/files", http.Dir(os.Getenv("APP_FILES")))

	auth := r.Group("/v1")
	{
		auth.Any("/signin", controladores.Signin)
		auth.Any("/logout", controladores.Logout)
		auth.GET("/hello", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"time": time.Now(),
				"utc":  time.UTC.String(),
			})
		})
	}

	r.POST("/v1/keepin", controladores.Signin).Use(VerifyJWT())

	user := r.Group("/v1/user").Use(VerifyJWT())
	{
		user.GET("", controladores.AllUser)
		user.GET("/:id", controladores.FindUser)
		user.POST("", controladores.AddUser)
		user.PUT("/:id", controladores.UpdUser)
		user.DELETE("/:id", controladores.DelUser)
	}

	return r
}
