package k8s

import (
	"github.com/gin-gonic/gin"
	"kubetabs/models"
	"kubetabs/pkg"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strings"
)

func GetPodsNoHttp(hostname string) int {
	clientset := pkg.KubeConfig()

	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + hostname,
	})
	if err != nil {
		panic(err)
	}

	return len(pods.Items)
}

func GetPodList(c *gin.Context) {
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
	hostname := c.Request.FormValue("hostname")
	podnameForm := strings.TrimSpace(c.Request.FormValue("podname"))

	clientset := pkg.KubeConfig()
	//clientset.AppsV1beta1().Deployments()

	var podItem []v1.Pod
	if hostname == "" {
		pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		podItem = pods.Items
	} else {
		pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{
			FieldSelector: "spec.nodeName=" + hostname,
		})
		if err != nil {
			panic(err.Error())
		}
		podItem = pods.Items
	}

	var podItemFilter []v1.Pod
	if podnameForm != "" {
		for _, pod := range podItem {
			if strings.Contains(pod.ObjectMeta.Name, podnameForm) {
				podItemFilter =  append(podItemFilter, pod)
			}
		}
	} else {
		podItemFilter = podItem
	}

	countFilter := len(podItemFilter)
	var mp = make(map[string]interface{}, 3)
	pStart := (pageIndex - 1) * pageSize
	pEnd := pageIndex * pageSize
	if pStart >= countFilter {
		mp["list"] = ""
	} else if pEnd >= countFilter {
		mp["list"] = podItemFilter[pStart:countFilter]
	} else {
		mp["list"] = podItemFilter[pStart:pEnd]
	}
	mp["count"] = len(podItemFilter)
	var res models.Response
	res.Data = mp

	c.JSON(http.StatusOK, res.ReturnOK())
}

func GetPodMetrics(c *gin.Context) {
	namespace := c.Request.FormValue("namespace")
	podName := c.Request.FormValue("podname")

	clientset := pkg.MetricsConfig()

	podMetric, err := clientset.MetricsV1beta1().PodMetricses(namespace).Get(podName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	var mp = make(map[string]interface{}, 3)
	mp["list"] = podMetric

	var res models.Response
	res.Data = mp

	c.JSON(http.StatusOK, res.ReturnOK())
}
