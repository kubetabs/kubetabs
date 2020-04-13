package k8s

import (
	"github.com/gin-gonic/gin"
	"kubetabs/models"
	"kubetabs/pkg"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func GetNsList(c *gin.Context) {
	clientset := pkg.KubeConfig()
	ns, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	var mp = make(map[string]interface{}, 3)
	mp["list"] = ns.Items
	mp["count"] = len(ns.Items)
	var res models.Response
	res.Data = mp

	c.JSON(http.StatusOK, res.ReturnOK())
}