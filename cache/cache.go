package cache

import "golang.org/x/oauth2"

var Cache *map[string]*oauth2.Token

func GetCacheInstance() *map[string]*oauth2.Token {
	if Cache == nil {
		c := map[string]*oauth2.Token{
			"init": nil,
		}
		Cache = &c
		return &c
	}

	return Cache
}
