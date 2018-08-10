package User

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"strconv"
	"net/http"
	"MyServer/bean"
	"github.com/GodSlave/MyGoServer/log"
	"encoding/json"
)

type BaseUserWebData struct {
	sql *xorm.Engine
}

type BaseUserResponse struct {
	BaseUsers   *[]bean.BaseUser
	AllPage     int
	CurrentPage int
}

func (m *BaseUserWebData) Init(router *gin.Engine) {
	router.GET("/User/getBaseUsers", m.get)
	router.POST("/User/updateBaseUser", m.update)
	router.POST("/User/deleteBaseUser", m.delete)
}

func (m *BaseUserWebData) get(c *gin.Context) {
	page := c.DefaultQuery("page", "0")
	pagei, err := strconv.Atoi(page)
	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatus(http.StatusMethodNotAllowed)
		return
	}
	BaseUsers := &[]bean.BaseUser{}
	err = m.sql.Limit(20, 20*pagei).Find(BaseUsers)

	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	count, err := m.sql.Count(&bean.BaseUser{})
	var allPage int
	if count%20 == 0 {
		allPage = int(count / 20)
	} else {
		allPage = int(count/20 + 1)
	}
	log.Info("%v", allPage)
	response := &BaseUserResponse{
		BaseUsers:   BaseUsers,
		AllPage:     allPage,
		CurrentPage: pagei,
	}
	data, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
	}
	log.Info(string(data))
	c.String(http.StatusOK, string(data))
	c.Done()
}

func (m *BaseUserWebData) update(c *gin.Context) {
	content := c.PostForm("key")
	BaseUser := &bean.BaseUser{
	}
	err := json.Unmarshal([]byte(content), BaseUser)
	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	if BaseUser.Id > 0 {
		m.sql.Id(BaseUser.Id).Update(BaseUser)
	} else {
		m.sql.Insert(BaseUser)
	}
	c.Status(http.StatusOK)
}

func (m *BaseUserWebData) delete(c *gin.Context) {
	content := c.DefaultQuery("id", "0")
	id, err := strconv.ParseInt(content,10,64)
	log.Info("%v", id)
	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	_, err = m.sql.Delete(&bean.BaseUser{Id: id})
	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
