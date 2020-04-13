package k8s

import (
	"github.com/gin-gonic/gin"
	"kubetabs/models"
	"kubetabs/pkg"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetScList(c *gin.Context) {
	var err error
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize = pkg.StrToInt(err, size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex = pkg.StrToInt(err, index)
	}

	clientset := pkg.KubeConfig()
	scs, err := clientset.StorageV1().StorageClasses().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	countSc := len(scs.Items)
	var mp = make(map[string]interface{}, 3)
	pStart := (pageIndex - 1) * pageSize
	pEnd := pageIndex * pageSize
	if pStart >= countSc {
		mp["list"] = ""
	} else if pEnd >= countSc {
		mp["list"] = scs.Items[pStart:countSc]
	} else {
		mp["list"] = scs.Items[pStart:pEnd]
	}
	mp["count"] = countSc
	mp["pageIndex"] = pageIndex
	mp["pageSize"] = pageSize

	var res models.Response
	res.Data = mp

	c.JSON(http.StatusOK, res.ReturnOK())
}
