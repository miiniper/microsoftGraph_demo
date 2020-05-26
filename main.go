package main //import "onedrive_demo"

import (
	"fmt"
	"microsoftGraph_demo/httpd"
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"

	"github.com/spf13/viper"
)

func main() {
	fmt.Println("Server starting ... ")
	viper.SetConfigName("conf")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		//fmt.Printf("Config file changed: %s", e.Name)
		httpd.Loges.Info("Config file changed: ", zap.Any("", e.Name))
	})

	service, err := httpd.New(viper.GetString("server.hostport"))
	if err != nil {
		panic(err)
	}
	err = service.Start()
	if err != nil {
		panic(err)
	}
	defer service.Close()

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGKILL)
	<-terminate

}
