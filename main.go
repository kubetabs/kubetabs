package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"kubetabs/config"
	orm "kubetabs/database"
	"kubetabs/models"
	"kubetabs/router"
)

func main() {

	gin.SetMode(gin.DebugMode)

	log.Println(config.DatabaseConfig.Port)
	if config.ApplicationConfig.IsInit {
		if err := models.InitDb(); err != nil {
			log.Fatal("数据库初始化失败！")
		} else {
			config.SetApplicationIsInit()
		}
	}
	r := router.InitRouter()

	defer orm.Eloquent.Close()

	server := &http.Server{
		Addr:              config.ApplicationConfig.Host + ":" + config.ApplicationConfig.Port,
		Handler:           r,
	}

	go func() {
		err := server.ListenAndServe(); if err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %s\n", err)
		}
	}()

	// graceful shutdown
	c := make(chan os.Signal)

	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")

}
