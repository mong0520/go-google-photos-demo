package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mong0520/go-google-photos-demo/handlers"
	"github.com/zalando/gin-oauth2/google"
	"golang.org/x/oauth2"
	"google.golang.org/api/photoslibrary/v1"
)

var redirectURL, credFile string
var tokenCache map[string]*oauth2.Token

type User struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
	Hd            string `json:"hd"`
}

func init() {
	bin := path.Base(os.Args[0])
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s
================
`, bin)
		flag.PrintDefaults()
	}
	flag.StringVar(&redirectURL, "redirect", "http://localhost:5000/auth/success", "URL to be redirected to after authorization.")
	flag.StringVar(&credFile, "cred-file", "./cred.json", "Credential JSON file")
}
func main() {
	flag.Parse()

	scopes := []string{
		"https://www.googleapis.com/auth/userinfo.email",
		photoslibrary.PhotoslibraryScope,
		// You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
	}
	secret := []byte("secret")
	sessionName := "goquestsession"

	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// init settings for google auth
	google.Setup(redirectURL, credFile, scopes, secret)
	router.Use(google.Session(sessionName))

	router.GET("/login", google.LoginHandler)
	router.GET("/albums", handlers.AlbumsHandler)

	// protected url group
	private := router.Group("/auth")
	private.Use(google.Auth())
	private.GET("/success", handlers.SuccessHandler)
	// private.GET("/api", func(ctx *gin.Context) {
	// 	token, err := ctx.Cookie("myphoto_cookie")
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	} else {
	// 		fmt.Println(token)
	// 	}
	// 	// ctx.JSON(200, gin.H{"message": "Hello from private for groups"})
	// })

	router.Run(":5000")
}

// func ShowUserInfo(ctx *gin.Context) {
// 	tokenStr, err := ctx.Cookie("myphoto_cookie")
// 	if err != nil {
// 		ctx.JSON(200, err)
// 	}
// 	token := &oauth2.Token{}
// 	err = json.Unmarshal([]byte(tokenStr), token)
// 	if err != nil {
// 		ctx.JSON(200, err)
// 	}
// 	conf := &oauth2.Config{
// 		ClientID:     os.Getenv("ClientID"),
// 		ClientSecret: os.Getenv("ClientSecret"),
// 		Scopes:       []string{photoslibrary.PhotoslibraryScope, "https://www.googleapis.com/auth/userinfo.email"},
// 	}
// 	client := conf.Client(oauth2.NoContext, token)
// 	email, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
// 	if err != nil {
// 		ctx.AbortWithError(http.StatusBadRequest, err)
// 		return
// 	}
// 	defer email.Body.Close()
// 	data, err := ioutil.ReadAll(email.Body)
// 	if err != nil {
// 		glog.Errorf("[Gin-OAuth] Could not read Body: %s", err)
// 		ctx.AbortWithError(http.StatusInternalServerError, err)
// 		return
// 	}

// 	var user User
// 	err = json.Unmarshal(data, &user)
// 	if err != nil {
// 		glog.Errorf("[Gin-OAuth] Unmarshal userinfo failed: %s", err)
// 		ctx.AbortWithError(http.StatusInternalServerError, err)
// 		return
// 	}
// 	ctx.JSON(200, user)
// }
