package k8s

import (
	"github.com/gin-gonic/gin"
	"kubetabs/models"
	"kubetabs/pkg"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetDeploymentList(c *gin.Context) {
	var err error
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize = pkg.StrToInt(err, size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex = pkg.StrToInt(err, index)
	}
	namespace := c.Request.FormValue("namespace")

	clientset := pkg.KubeConfig()

	deployments, err := clientset.AppsV1().Deployments(namespace).List(metav1.ListOptions{})

	if err != nil {
		panic(err)
	}
	countItem := len(deployments.Items)
	var mp = make(map[string]interface{}, 3)
	pStart := (pageIndex - 1) * pageSize
	pEnd := pageIndex * pageSize
	if pStart >= countItem {
		mp["list"] = ""
	} else if pEnd >= countItem {
		mp["list"] = deployments.Items[pStart:countItem]
	} else {
		mp["list"] = deployments.Items[pStart:pEnd]
	}
	mp["count"] = countItem
	var res models.Response
	res.Data = mp

	c.JSON(http.StatusOK, res.ReturnOK())
}
