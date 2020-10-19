package cmd

import (
	"context"
	"fmt"
	"github.com/rendau/gl2dc/internal/adapters/rest"
	"github.com/rendau/gl2dc/internal/domain/core"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Execute() {
	var err error

	loadConf()

	cr := core.NewSt(
		viper.GetString("discord_webhook_url"),
		viper.GetString("graylog_link"),
	)

	httpApi := rest.New(
		viper.GetString("http_listen"),
		cr,
	)

	log.Println("Starting", "http_listen", viper.GetString("http_listen"))

	httpApi.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	var exitCode int

	select {
	case <-stop:
	case <-httpApi.Wait():
		exitCode = 1
	}

	log.Println("Shutting down...")

	restApiShutdownCtx, restApiShutdownCtxCancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer restApiShutdownCtxCancel()

	err = httpApi.Shutdown(restApiShutdownCtx)
	if err != nil {
		_ = fmt.Errorf("fail to shutdown http-api: %s", err.Error())
		exitCode = 1
	}

	os.Exit(exitCode)
}

func loadConf() {
	viper.SetDefault("http_listen", ":80")

	confFilePath := os.Getenv("CONF_PATH")
	if confFilePath == "" {
		confFilePath = "conf.yml"
	}
	viper.SetConfigFile(confFilePath)
	_ = viper.ReadInConfig()

	viper.AutomaticEnv() // read in environment variables that match
}
