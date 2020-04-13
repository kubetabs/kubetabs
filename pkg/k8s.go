package pkg

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)


func KubeConfig() *kubernetes.Clientset {
	//在 kubeconfig 中使用当前上下文环境，config 获取支持 url 和 path 方式
	config, err := clientcmd.BuildConfigFromFlags("", "/tmp/admin.conf")
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func MetricsConfig() *metrics.Clientset {
	//在 kubeconfig 中使用当前上下文环境，config 获取支持 url 和 path 方式
	config, err := clientcmd.BuildConfigFromFlags("", "/tmp/admin.conf")
	if err != nil {
		panic(err.Error())
	}

	clientset, err := metrics.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return clientset
}
