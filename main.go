package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mong0520/go-google-photos-demo/handlers"
	"google.golang.org/api/photoslibrary/v1"
)

var redirectURL, credFile string

func init() {
	godotenv.Load()
	bin := path.Base(os.Args[0])
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s
================
`, bin)
		flag.PrintDefaults()
	}
	flag.StringVar(&redirectURL, "redirect", fmt.Sprintf("http://%s:5000/callback", os.Getenv("HOSTNAME")), "URL to be redirected to after authorization.")
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
	handlers.LoginSetup(redirectURL, scopes, secret)
	router.Use(handlers.Session(sessionName))

	router.GET("/login", handlers.LoginHandler)
	router.GET("/albums", handlers.AuthMiddleware, handlers.AlbumsHandler)
	router.Static("/web", "./web")

	// protected url group
	private := router.Group("/callback")
	private.Use(handlers.GoogleAuthMiddleware())
	private.GET("/", handlers.CallbackHandler)

	router.Run(":5000")
}
