package controladores

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	m "github.com/sam/modelos"
)

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func Signin(c *gin.Context) {
	var creds Credentials
	user := new(m.User)

	err := c.BindJSON(&creds)
	if err != nil {
		c.JSON(500, gin.H{"msg": err})
		return
	}

	if creds.Password == "" || len(creds.Password) < 4 {
		c.JSON(400, gin.H{"msg": "invalid credentials | credenciales invalidas"})
		return
	}

	user, err = user.VerifyCredentials(creds.Username, creds.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": err.Error()})
		return
	}

	token, exp, err := user.GenerateTokenJWT()
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: exp,
	})
	c.JSON(200, gin.H{"user": user, "token": token})
}

func Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Unix(0, 0),
	})
}
