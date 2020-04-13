package tools

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"kubetabs/models"
	"kubetabs/models/tools"
	"kubetabs/pkg"
	"kubetabs/utils"
	"net/http"
	"text/template"
)

func Preview(c *gin.Context) {
	table := tools.SysTables{}
	id, err := utils.StringToInt64(c.Param("tableId"))
	pkg.AssertErr(err, "", -1)
	table.TableId = id
	t1, err := template.ParseFiles("template/model.go.template")
	pkg.AssertErr(err, "", -1)
	t2, err := template.ParseFiles("template/api.go.template")
	pkg.AssertErr(err, "", -1)
	tab, _ := table.Get()
	var b1 bytes.Buffer
	err = t1.Execute(&b1, tab)
	var b2 bytes.Buffer
	err = t2.Execute(&b2, tab)

	mp := make(map[string]interface{})
	mp["template/model.go.template"] = b1.String()
	mp["template/api.go.template"] = b2.String()
	var res models.Response
	res.Data = mp

	c.JSON(http.StatusOK, res.ReturnOK())
}
