package main //import "onedrive_demo"

import (
	"fmt"
	"microsoftGraph_demo/httpd"
	"net/http"

	"github.com/julienschmidt/httprouter"

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

	router := httprouter.New()
	//	router.Handle("GET","/index",http.FileServer(http.Dir("template/")))
	router.GET("/auth/callback", httpd.GetCode)
	router.GET("/auth/login", httpd.MicrosoftLogin)
	router.GET("/me", httpd.ShowUser)
	router.GET("/profile", httpd.ShowProfile)

	httpd.Loges.Fatal("services :", zap.Any("", http.ListenAndServe(viper.GetString("server.hostPort"), router)))

}
