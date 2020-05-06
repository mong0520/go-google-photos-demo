package models

type GeneralResponse struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Error   string      `json:"error_code"`
}

type SimpleAlbum struct {
	// CoverPhotoBaseUrl string
	Title string
	Url   string
}
