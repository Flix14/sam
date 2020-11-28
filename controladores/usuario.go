package controladores

import (
	"strconv"

	"github.com/gin-gonic/gin"
	m "github.com/sam/modelos"
)

type UserReq struct {
	ID        int    `json:"id"`
	User      string `json:"username"`
	Password  string `json:"password"`
	Password2 string `json:"password2"`
}

func AllUser(c *gin.Context) {
	db := m.DB
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	id := c.DefaultQuery("id", "")
	username := c.DefaultQuery("username", "")
	name := c.DefaultQuery("name", "")
	// sortBy := c.DefaultQuery("sortBy", "id")
	// order := c.DefaultQuery("order", "desc")

	var users []m.User
	var count int64

	if id != "" {
		db = db.Where("id = ?", id)
	}
	if username != "" {
		db = db.Where("username LIKE ?", username+"%")
	}
	if name != "" {
		db = db.Where("name LIKE ?", name+"%")
	}

	db.Scopes(m.Pagination(page, limit)).Find(&users)
	db.Model(m.User{}).Count(&count)
	paginator := m.Paginator{
		Limit:       limit,
		Page:        page,
		TotalRecord: count,
		Records:     users,
	}
	c.JSON(200, paginator)
}

func FindUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user m.User
	err := user.Find(id)
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, user)
}

func AddUser(c *gin.Context) {
	var user m.User
	var userReq UserReq
	err := c.BindJSON(&userReq)
	user.ID = userReq.ID
	user.Username = userReq.User
	user.Password = userReq.Password

	// TODO validate two passwords equal

	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	newuser, err := user.Add()
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(201, newuser)
}

func UpdUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user m.User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}

	// TODO validate two passwords equal

	user.ID = id
	newuser, err := user.Update()
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, newuser)
}

func DelUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user m.User
	user.ID = id
	err := user.Remove()
	if err != nil {
		c.JSON(500, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "ok"})
}
