package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/mong0520/go-google-photos-demo/cache"
	"github.com/mong0520/go-google-photos-demo/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	Token_request_uri          string = "https://accounts.google.com/o/oauth2/auth"
	Response_type              string = "code"
	Client_id                  string = "58530366777-ipig2dojcbn8av0qjiifsh8prbbucejk.apps.googleusercontent.com"
	Redirect_uri               string = "http://localhost:5000/oauth2callback"
	Scope                      string = "https://www.googleapis.com/auth/photoslibrary"
	Access_type                string = "offline"
	oauthStateStringContextKey        = 987
	state                      string
	conf                       *oauth2.Config
	store                      sessions.CookieStore
)

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func AuthMiddleware(ctx *gin.Context) {
	session := ctx.Query("sessionID")
	fmt.Printf("session ID= %s\n", session)
	if session == "" {
		ctx.JSON(200, "please login first")
		return
	}

	cacheInst := *cache.GetCacheInstance()
	token := &oauth2.Token{}
	if val, ok := cacheInst[session]; ok {
		token = val
	} else {
		ctx.JSON(200, "Session ID can map to Access Token")
		return
	}

	ctx.Set("token", token)
	ctx.Set("conf", conf)
}

func GoogleAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Handle the exchange code to initiate a transport.
		session := sessions.Default(ctx)
		fmt.Printf("sessino state: %s, query state: %s", session.Get("state"), ctx.Query("state"))
		retrievedState := session.Get("state")
		if retrievedState != ctx.Query("state") {
			ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state: %s", retrievedState))
			return
		}

		tok, err := conf.Exchange(oauth2.NoContext, ctx.Query("code"))
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		client := conf.Client(oauth2.NoContext, tok)
		email, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		defer email.Body.Close()
		data, err := ioutil.ReadAll(email.Body)
		if err != nil {
			glog.Errorf("[Gin-OAuth] Could not read Body: %s", err)
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		var user models.User
		err = json.Unmarshal(data, &user)
		if err != nil {
			glog.Errorf("[Gin-OAuth] Unmarshal userinfo failed: %s", err)
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		// save userinfo, which could be used in Handlers
		ctx.Set("user", user)
		ctx.Set("conf", conf)
		ctx.Set("token", tok)
		ctx.Set("client", client)
	}
}

func LoginSetup(redirectURL string, scopes []string, secret []byte) {
	store = sessions.NewCookieStore(secret)
	store.Options(sessions.Options{
		MaxAge: int(24 * time.Hour * 30), //30 days
		Path:   "/",
	})

	conf = &oauth2.Config{
		ClientID:     os.Getenv("ClientID"),
		ClientSecret: os.Getenv("ClientSecret"),
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
	}
}

func LoginHandler(ctx *gin.Context) {
	state = randToken()
	session := sessions.Default(ctx)
	session.Set("state", state)
	session.Save()
	ctx.Writer.Write([]byte("<html><title>Golang Google</title> <body> <a href='" + GetLoginURL(state) + "'><button>Login with Google!</button> </a> </body></html>"))
}

func CallbackHandler(ctx *gin.Context) {
	cache := cache.GetCacheInstance()
	token := ctx.Value("token").(*oauth2.Token)
	session, _ := ctx.Cookie("goquestsession")
	cacheInst := *cache
	cacheInst[session] = token
	fmt.Println("cache = ", cacheInst[session])
	ctx.Redirect(http.StatusMovedPermanently, fmt.Sprintf("http://%s:5000/web", os.Getenv("HOSTNAME")))
}

func GetLoginURL(state string) string {
	return conf.AuthCodeURL(state)
}

func Session(name string) gin.HandlerFunc {
	return sessions.Sessions(name, store)
}
