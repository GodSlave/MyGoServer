package template

const BGWebContent=`package backend

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"strconv"
	"net/http"
	"../bean"
	"github.com/GodSlave/MyGoServer/log"
	"encoding/json"
)

type {{.StructName}}WebData struct {
	sql *xorm.Engine
}

type {{.StructName}}Response struct {
	{{.StructName}}s   *[]bean.{{.StructName}}
	AllPage     int
	CurrentPage int
}

func (m *{{.StructName}}WebData) Init(router *gin.Engine) {
	router.GET("/{{.ModuleName}}/get{{.StructName}}s", m.get)
	router.POST("/{{.ModuleName}}/update{{.StructName}}", m.update)
	router.POST("/{{.ModuleName}}/delete{{.StructName}}", m.delete)
}

func (m *{{.StructName}}WebData) get(c *gin.Context) {
	page := c.DefaultQuery("page", "0")
	pagei, err := strconv.Atoi(page)
	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatus(http.StatusMethodNotAllowed)
		return
	}
	{{.StructName}}s := &[]bean.{{.StructName}}{}
	err = m.sql.Limit(20, 20*pagei).Find({{.StructName}}s)

	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	count, err := m.sql.Count(&bean.{{.StructName}}{})
	var allPage int
	if count%20 == 0 {
		allPage = int(count / 20)
	} else {
		allPage = int(count/20 + 1)
	}
	log.Info("%v", allPage)
	response := &{{.StructName}}Response{
		{{.StructName}}s:   {{.StructName}}s,
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
	{{.StructName}} := &bean.{{.StructName}}{
	}
	err := json.Unmarshal([]byte(content), {{.StructName}})
	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	if {{.StructName}}.Id > 0 {
		m.sql.Id({{.StructName}}.Id).Update({{.StructName}})
	} else {
		m.sql.Insert({{.StructName}})
	}
	c.Status(http.StatusOK)
}

func (m *{{.StructName}}WebData) delete(c *gin.Context) {
	content := c.DefaultQuery("id", "0")
	id, err := strconv.ParseInt(content,10,64)
	log.Info("%v", id)
	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	_, err = m.sql.Delete(&bean.{{.StructName}}{Id: id})
	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (m *{{.StructName}}WebData) SetEngine(sqlEngine *xorm.Engine) {
	m.sql=sqlEngine
}
`