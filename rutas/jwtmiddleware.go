package rutas

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	m "github.com/sam/modelos"
)

func VerifyJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Token")
		branch := c.GetHeader("X-Branch")

		// Fix for testing purposes
		env := os.Getenv("APP_ENV")
		if env == "testing" {
			c.Set("USER", "1")
			c.Set("BRANCH", "1")
			c.Next()
			return
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, "Unauthorized | No autorizado")
			return
		}
		claim, err := m.ValidateTokenJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "Unauthorized | No autorizado  - "+err.Error())
			return
		}
		c.Set("USER", claim.User)
		c.Set("BRANCH", branch)
		c.Next()
	}
}
