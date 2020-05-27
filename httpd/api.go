package httpd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/viper"

	"go.uber.org/zap"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

var FirstCode string
var Cli *http.Client
var Token *oauth2.Token
var microsoftGraphAPI map[string]string

var MicrosoftGraphOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/auth/callback",
	ClientID:     os.Getenv("CID"),
	ClientSecret: os.Getenv("CS"),
	Scopes:       []string{"files.readwrite", "mail.read"},
	Endpoint:     microsoft.AzureADEndpoint(""),
}

//初始化api
func InitMGApi() {
	microsoftGraphAPI = map[string]string{
		"my profile":   "https://graph.microsoft.com/v1.0/me/",
		"my drive all": "https://graph.microsoft.com/v1.0/me/drive/root/children",
		"users":        "https://graph.microsoft.com/v1.0/users",
		"message":      "https://graph.microsoft.com/v1.0/me/messages",
	}
}

//登录
func MicrosoftLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	InitMGApi()
	u := MicrosoftGraphOauthConfig.AuthCodeURL(viper.GetString("auth.state"))
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)

}

//获取code
func GetCode(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	if r.FormValue("state") != viper.GetString("auth.state") {
		Loges.Info("state is error")
		http.Redirect(w, r, "/auth/login", http.StatusTemporaryRedirect)
		return
	}
	FirstCode = r.URL.Query().Get("code")
	fmt.Println("=============")
	fmt.Println("===  ", FirstCode)
	fmt.Println("=============")
	http.Redirect(w, r, "/me", http.StatusTemporaryRedirect)

}

//通过code获取token
func SetCli(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	code := FirstCode
	if len(code) <= 1 {
		http.Redirect(w, r, "/auth/login", http.StatusTemporaryRedirect)
		return
	}
	Token, err := MicrosoftGraphOauthConfig.Exchange(ctx, code)
	if err != nil {
		Loges.Error("get token  is err:", zap.Error(err))
		return
	}
	Cli = MicrosoftGraphOauthConfig.Client(ctx, Token)

}

//get 请求函数标准
func GetMic(url string) []byte {

	if Cli == nil {
		Loges.Error("Cli  is err,: empty")
		return nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Loges.Error("new request is err:", zap.Error(err))
	}
	resp, err := Cli.Do(req)
	if err != nil {
		Loges.Error("do request is err:", zap.Error(err))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Loges.Error("ioutil ReadAll   is err:", zap.Error(err))
	}
	return body

}

//Post 请求函数标准
func PostMic(url string) []byte {

	if Cli == nil {
		Loges.Error("Cli  is err,: empty")
		return nil
	}

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		Loges.Error("new request is err:", zap.Error(err))
	}

	resp, err := Cli.Do(req)
	if err != nil {
		Loges.Error("do request is err:", zap.Error(err))
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		Loges.Error("ioutil ReadAll   is err:", zap.Error(err))
	}

	return body

}

//url := "https://login.microsoftonline.com/common/oauth2/v2.0/token"

//显示组下的用户
func ShowUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	if Cli == nil {
		SetCli(w, r)
	}

	url := microsoftGraphAPI["users"]
	body := GetMic(url)
	w.Write(body)

}

//当前用户的配置文件
func ShowProfile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if Cli == nil {
		SetCli(w, r)
	}

	url := microsoftGraphAPI["my profile"]
	body := GetMic(url)
	w.Write(body)

}
func ShowDrive(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if Cli == nil {
		SetCli(w, r)
	}

	url := microsoftGraphAPI["my drive all"]
	body := GetMic(url)
	w.Write(body)

}

func ShowMsg(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if Cli == nil {
		SetCli(w, r)
	}

	url := microsoftGraphAPI["message"]
	body := GetMic(url)
	w.Write(body)

}
