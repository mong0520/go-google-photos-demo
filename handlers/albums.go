package handlers

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mong0520/go-google-photos-demo/models"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	photoslibrary "google.golang.org/api/photoslibrary/v1"
)

var albumService *photoslibrary.AlbumsService
var defaultPageSize = int64(50)

func GetAlbums(c *gin.Context) {
	// those context value was set in middleware
	token := c.MustGet("token").(*oauth2.Token)
	conf := c.MustGet("conf").(*oauth2.Config)
	httpClient := conf.Client(context.Background(), token)
	service := &photoslibrary.Service{}

	service, err := photoslibrary.New(httpClient)
	if err != nil {
		c.JSON(200, err)
	}

	albumService = photoslibrary.NewAlbumsService(service)
	albums, err := listAlbums()
	c.JSON(200, albums)
}

func listAlbums() (albums []*models.SimpleAlbum, err error) {
	albumList := albumService.List()
	ret, err := albumList.PageSize(defaultPageSize).Do()
	albumList.Do()
	if err != nil {
		return albums, err
	}
	for _, album := range ret.Albums {
		fmt.Println(album.Title, album.ProductUrl)
		simpleAlbum := &models.SimpleAlbum{
			Title:             album.Title,
			CoverPhotoBaseUrl: album.CoverPhotoBaseUrl,
			Url:               album.ProductUrl,
		}
		albums = append(albums, simpleAlbum)
	}
	for {
		nextPageToken := ret.NextPageToken
		if nextPageToken == "" {
			break
		}
		ret, err = albumList.PageToken(nextPageToken).PageSize(defaultPageSize).Do()
		if err != nil {
			log.Fatal(err)
			return albums, err
		}
		for _, album := range ret.Albums {
			fmt.Println(album.Title, album.ProductUrl)
			simpleAlbum := &models.SimpleAlbum{
				Title:             album.Title,
				CoverPhotoBaseUrl: album.CoverPhotoBaseUrl,
				Url:               album.ProductUrl,
			}
			albums = append(albums, simpleAlbum)
		}
	}

	return albums, nil
}
