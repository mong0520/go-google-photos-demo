package main

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mong0520/go-google-photos-demo/handlers"
	photoslibrary "google.golang.org/api/photoslibrary/v1"
)

const (
	serviceName = "googlephotos-uploader-go-api"
)

var albumService *photoslibrary.AlbumsService

func run() {
	router := gin.New()
	godotenv.Load()
	router.Use(cors.Default())

	// router.Static("/web", "./web")
	// router.GET("/login", handlers.LoginHandler)
	// router.GET("/login2", handlers.LoginHandler2)
	// router.GET("/oauth2callback", handlers.CallbackHander2)
	router.GET("/albums", handlers.AlbumsHandler)
	router.GET("/healthcheck", handlers.HealthCheckHandler)

	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	router.Run(addr)
}

func main() {
	run()
}

// package main

// import (
// 	"flag"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"path"

// 	"github.com/gin-gonic/gin"
// 	"github.com/zalando/gin-oauth2/google"
// )

// var redirectURL, credFile string

// func init() {
// 	bin := path.Base(os.Args[0])
// 	flag.Usage = func() {
// 		fmt.Fprintf(os.Stderr, `
// Usage of %s
// ================
// `, bin)
// 		flag.PrintDefaults()
// 	}
// 	flag.StringVar(&redirectURL, "redirect", "http://127.0.0.1:8081/auth/api", "URL to be redirected to after authorization.")
// 	flag.StringVar(&credFile, "cred-file", "./cred.json", "Credential JSON file")
// }
// func main() {
// 	flag.Parse()

// 	scopes := []string{
// 		"https://www.googleapis.com/auth/userinfo.email",
// 		// You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
// 	}
// 	secret := []byte("secret")
// 	sessionName := "goquestsession"

// 	router := gin.Default()
// 	// init settings for google auth
// 	google.Setup(redirectURL, credFile, scopes, secret)
// 	router.Use(google.Session(sessionName))

// 	router.GET("/login", google.LoginHandler)

// 	// protected url group
// 	private := router.Group("/auth")
// 	private.Use(google.Auth())
// 	private.GET("/", UserInfoHandler)
// 	private.GET("/api", func(ctx *gin.Context) {
// 		ctx.JSON(200, gin.H{"message": "Hello from private for groups"})
// 	})

// 	router.Run("127.0.0.1:8081")
// }

// func UserInfoHandler(ctx *gin.Context) {
// 	ctx.JSON(http.StatusOK, gin.H{"Hello": "from private", "user": ctx.MustGet("user").(google.User)})
// }
