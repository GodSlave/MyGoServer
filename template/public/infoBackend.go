package info

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"strconv"
	"net/http"
	"MyServer/bean"
	"github.com/GodSlave/MyGoServer/log"
	"encoding/json"
)

type infoWebData struct {
	sql *xorm.Engine
}

type infoResponse struct {
	infos   *[]bean.infoItem
	AllPage     int
	CurrentPage int
}

func (m *infoWebData) Init(router *gin.Engine) {
	router.GET("/info/get", m.get)
	router.POST("/info/update", m.update)
	router.POST("/info/delete", m.delete)
}

func (m *infoWebData) get(c *gin.Context) {
	page := c.DefaultQuery("page", "0")
	pagei, err := strconv.Atoi(page)
	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatus(http.StatusMethodNotAllowed)
		return
	}
	infoItems := &[]bean.infoItem{}
	err = m.sql.Limit(20, 20*pagei).Find(infoItems)

	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	count, err := m.sql.Count(&bean.infoItems{})
	var allPage int
	if count%20 == 0 {
		allPage = int(count / 20)
	} else {
		allPage = int(count/20 + 1)
	}
	log.Info("%v", allPage)
	response := &infoResponse{
		infoItems:   infoItems,
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

func (m *infoWebData) update(c *gin.Context) {
	content := c.PostForm("key")
	infoItem := &bean.infoItem{
	}
	err := json.Unmarshal([]byte(content), infoItem)
	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	if infoItem.Id > 0 {
		m.sql.Id(infoItem.Id).Update(infoItem)
	} else {
		m.sql.Insert(infoItem)
	}
	c.Status(http.StatusOK)
}

func (m *infoWebData) delete(c *gin.Context) {
	content := c.DefaultQuery("id", "0")
	id, err := strconv.Atoi(content)
	log.Info("%v", id)
	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	_, err = m.sql.Delete(&bean.infoItem{Id: int32(id)})
	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
