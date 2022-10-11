package apimanagement

import (
	"strconv"
	"strings"
	"time"
)

var (
	st        ArtistsSimpleApi
	dataReady bool
)

func WaitForReady() {
	for !dataReady {
		time.Sleep(time.Millisecond * 20)
	}
}

// referesh data from API every n minute, must be invocate in new Goroutine !!
func ReadEssentialAPI(min time.Duration) {
	for {
		dataReady = false
		ConfigApi()
		GetAllArtistsSimpleApi()
		GetSearchBarData()
		dataReady = true
		time.Sleep(min * time.Minute)
	}
}

func GetAllArtistsSimpleApi() *ArtistsSimpleApi {
	getApi(api.Artists, &st.Artists)
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
	for _, each := range st.Artists {
		if each.ArtistName == name {
			return each.Id
		}
	}
	return -1
}

// return a new artists colection with a int slice
func GetNewSliceByIdArtistsSimpleApi(table []int) ArtistsSimpleApi {
	var newArtists ArtistsSimpleApi
	for _, each := range st.Artists {
		for _, id := range table {
			if each.Id == id {
				newArtists.Artists = append(newArtists.Artists, Artist{each.Id, each.Img, each.ArtistName})
			}
		}
	}
	return newArtists
}

// get all id match the search
func GetIdSearch(s string) []int {
	const key = 1
	const value = 0
	searchTable := strings.Split(s, " - ")
	switch searchTable[key] {
	case "FirstAlbum_Date":
		return st.Bar.FirstAlbum_Date[searchTable[value]]
	case "Artist_BandName":
		return st.Bar.Artist_BandName[searchTable[value]]
	case "Creation_Date":
		return st.Bar.Creation_Date[searchTable[value]]
	case "Locations":
		return st.Bar.Locations[searchTable[value]]
	case "Members":
		return st.Bar.Members[searchTable[value]]
	}
	return []int{}
}

// super, il n'y a plus d'erreur !!!!!!!!!!!!!!!!!!!!!!!!! David

// load data for sear bar
func GetSearchBarData() {
	barTemp := BarData{}
	var allArtistData []BarArtistsData
	var allLocations map[string][]Locations
	getApi(api.Artists, &allArtistData)
	getApi(api.Locations, &allLocations)
	barTemp.FirstAlbum_Date = make(map[string][]int)
	barTemp.Artist_BandName = make(map[string][]int)
	barTemp.Creation_Date = make(map[string][]int)
	barTemp.Locations = make(map[string][]int)
	barTemp.Members = make(map[string][]int)
	for _, data := range allArtistData {
		barTemp.Artist_BandName[data.ArtistName] = append(barTemp.Artist_BandName[data.ArtistName], data.Id)
		for _, members := range data.Members {
			barTemp.Members[members] = append(barTemp.Members[members], data.Id)
		}
		Cdate := strconv.Itoa(data.CreationDate)
		barTemp.Creation_Date[Cdate] = append(barTemp.Creation_Date[Cdate], data.Id)
	}
	for _, data := range allLocations["index"] {
		for _, loca := range data.Locations {
			barTemp.Locations[loca] = append(barTemp.Locations[loca], data.Id)
		}
	}
	st.Bar = barTemp
}
