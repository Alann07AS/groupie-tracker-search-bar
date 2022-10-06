package apimanagement

import (
	"strconv"
)

var st []ArtistsSimpleApi
var stRefresh int

func GetAllArtistsSimpleApi() *[]ArtistsSimpleApi {
	if st != nil && stRefresh < 10 {
		stRefresh++
		return &st
	}
	stRefresh = 0
	getApi(api.Artists, &st)
	return &st
}

// get mores info by id searh
func GetAllInfoArtists(id int) AllInfoArtist {
	var artistInfo AllInfoArtist
	sId := strconv.Itoa(id)
	getApi(api.Relation+"/"+sId, &artistInfo.Relation)
	getApi(api.Artists+"/"+sId, &artistInfo)
	return artistInfo
}

// Get element by group Name, return Id api
func GetAllArtistsSimpleApiByName(name string) int {
	for _, each := range st {
		if each.ArtistName == name {
			return each.Id
		}
	}
	return -1
}

// super, il n'y a plus d'erreur !!!!!!!!!!!!!!!!!!!!!!!!! David
