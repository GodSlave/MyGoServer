package template

const BGWebContent=`package {{.StructName}}

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"strconv"
	"net/http"
	"MyServer/bean"
	"github.com/GodSlave/MyGoServer/log"
	"encoding/json"
)

type {{.StructName}}WebData struct {
	sql *xorm.Engine
}

type {{.StructName}}Response struct {
	{{.StructName}}s   *[]bean.{{.StructName}}Item
	AllPage     int
	CurrentPage int
}

func (m *{{.StructName}}WebData) Init(router *gin.Engine) {
	router.GET("/{{.StructName}}/get", m.get)
	router.POST("/{{.StructName}}/update", m.update)
	router.POST("/{{.StructName}}/delete", m.delete)
}

func (m *{{.StructName}}WebData) get(c *gin.Context) {
	page := c.DefaultQuery("page", "0")
	pagei, err := strconv.Atoi(page)
	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatus(http.StatusMethodNotAllowed)
		return
	}
	{{.StructName}}Items := &[]bean.{{.StructName}}Item{}
	err = m.sql.Limit(20, 20*pagei).Find({{.StructName}}Items)

	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	count, err := m.sql.Count(&bean.{{.StructName}}Items{})
	var allPage int
	if count%20 == 0 {
		allPage = int(count / 20)
	} else {
		allPage = int(count/20 + 1)
	}
	log.Info("%v", allPage)
	response := &{{.StructName}}Response{
		{{.StructName}}Items:   {{.StructName}}Items,
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

func (m *{{.StructName}}WebData) update(c *gin.Context) {
	content := c.PostForm("key")
	{{.StructName}}Item := &bean.{{.StructName}}Item{
	}
	err := json.Unmarshal([]byte(content), {{.StructName}}Item)
	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	if {{.StructName}}Item.Id > 0 {
		m.sql.Id({{.StructName}}Item.Id).Update({{.StructName}}Item)
	} else {
		m.sql.Insert({{.StructName}}Item)
	}
	c.Status(http.StatusOK)
}

func (m *{{.StructName}}WebData) delete(c *gin.Context) {
	content := c.DefaultQuery("id", "0")
	id, err := strconv.Atoi(content)
	log.Info("%v", id)
	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	_, err = m.sql.Delete(&bean.{{.StructName}}Item{Id: int32(id)})
	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
`