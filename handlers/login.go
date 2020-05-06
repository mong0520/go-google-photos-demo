package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mong0520/go-google-photos-demo/cache"
	"golang.org/x/oauth2"
)

var (
	Token_request_uri          string = "https://accounts.google.com/o/oauth2/auth"
	Response_type              string = "code"
	Client_id                  string = "58530366777-ipig2dojcbn8av0qjiifsh8prbbucejk.apps.googleusercontent.com"
	Redirect_uri               string = "http://localhost:5000/oauth2callback"
	Scope                      string = "https://www.googleapis.com/auth/photoslibrary"
	Access_type                string = "offline"
	oauthStateStringContextKey        = 987
)

func SuccessHandler(ctx *gin.Context) {
	cache := cache.GetCacheInstance()
	fmt.Println("setting token cache")
	token := ctx.Value("token").(*oauth2.Token)
	session, _ := ctx.Cookie("goquestsession")
	fmt.Printf("session = %s\n", session)
	cacheInst := *cache
	cacheInst[session] = token
	fmt.Println("cache = ", cacheInst[session])
	ctx.JSON(http.StatusOK, "loging success")
	// ctx.Redirect(http.StatusMovedPermanently, "http://api.nt1.me:5000/albums")
}

// func LoginHandler2(c *gin.Context) {
// 	fmt.Println("test")
// 	godotenv.Load()
// 	_, err := user.Current()
// 	if err != nil {
// 		panic(err)
// 	}
// 	// ask the user to authenticate on google in the browser
// 	conf := &oauth2.Config{
// 		ClientID:     os.Getenv("ClientID"),
// 		ClientSecret: os.Getenv("ClientSecret"),
// 		Scopes:       []string{photoslibrary.PhotoslibraryScope},
// 		Endpoint: oauth2.Endpoint{
// 			AuthURL:  google.Endpoint.AuthURL,
// 			TokenURL: google.Endpoint.TokenURL,
// 		},
// 		RedirectURL: Redirect_uri,
// 	}
// 	redirectUrl, _ := getRedirectUrl(conf)
// 	// oauth2ns.AuthenticateUser()

// 	fmt.Println(redirectUrl)
// 	c.Redirect(http.StatusMovedPermanently, redirectUrl)

// 	// client, err := authenticateUser(config)
// }

// func getRedirectUrl(oauthConfig *oauth2.Config) (string, error) {
// 	// validate params
// 	if oauthConfig == nil {
// 		return "", stacktrace.NewError("oauthConfig can't be nil")
// 	}

// 	// add transport for self-signed certificate to context
// 	tr := &http.Transport{
// 		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
// 	}
// 	sslcli := &http.Client{Transport: tr}
// 	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, sslcli)

// 	// Some random string, random for each request
// 	oauthStateString := rndm.String(8)
// 	ctx = context.WithValue(ctx, oauthStateStringContextKey, oauthStateString)
// 	urlString := oauthConfig.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline)

// 	// if IP != "127.0.0.1" {
// 	// 	urlString = fmt.Sprintf("%s&device_id=%s&device_name=%s", urlString, "MyMacbok", "MyMacbok")
// 	// }
// 	return urlString, nil
// }

// // func authenticateUser(oauthConfig *oauth2.Config, options ...AuthenticateUserOption) (*AuthorizedClient, error) {

// // 	clientChan, stopHTTPServerChan, cancelAuthentication := startHTTPServer(ctx, oauthConfig)
// // 	log.Println(color.CyanString("You will now be taken to your browser for authentication or open the url below in a browser."))
// // 	log.Println(color.CyanString(urlString))
// // 	log.Println(color.CyanString("If you are opening the url manually on a different machine you will need to curl the result url on this machine manually."))
// // 	time.Sleep(1000 * time.Millisecond)
// // 	err := open.Run(urlString)
// // 	if err != nil {
// // 		log.Println(color.RedString("Failed to open browser, you MUST do the manual process."))
// // 	}
// // 	time.Sleep(600 * time.Millisecond)

// // 	// shutdown the server after timeout
// // 	go func() {
// // 		log.Printf("Authentication will be cancelled in %s seconds", strconv.Itoa(authTimeout))
// // 		time.Sleep(authTimeout * time.Second)
// // 		stopHTTPServerChan <- struct{}{}
// // 	}()

// // 	select {
// // 	// wait for client on clientChan
// // 	case client := <-clientChan:
// // 		// After the callbackHandler returns a client, it's time to shutdown the server gracefully
// // 		stopHTTPServerChan <- struct{}{}
// // 		return client, nil

// // 		// if authentication process is cancelled first return an error
// // 	case <-cancelAuthentication:
// // 		return nil, fmt.Errorf("authentication timed out and was cancelled")
// // 	}
// // }
