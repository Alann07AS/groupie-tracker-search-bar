package apimanagement

import (
	"fmt"
	"log"
	"net/http"
	"sort"
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
		log.Println("APIdata Load...(refresh:", dataReady, ")")
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

// get all artists in range creation and first album dates
func GetAllArtistInFilters(r *http.Request) []int {
	var list1, list2, list3, finalList []int
	r.ParseForm()
	fMin, _ := strconv.Atoi(r.FormValue("minDateAlbum"))
	fMax, _ := strconv.Atoi(r.FormValue("maxDateAlbum"))
	aMin, _ := strconv.Atoi(r.FormValue("minDateCreation"))
	aMax, _ := strconv.Atoi(r.FormValue("maxDateCreation"))

	for i, each := range st.Bar.FirstAlbum_YearDate {
		if ii, _ := strconv.Atoi(i); ii <= fMax && ii >= fMin {
			list1 = append(list1, each...)
		}
	}
	for i, each := range st.Bar.Creation_Date {
		if ii, _ := strconv.Atoi(i); ii <= aMax && ii >= aMin {
			list2 = append(list2, each...)
		}
	}
	finalList = CompareList(list1, list2)
	for i, each := range st.Bar.NbMembers {
		if r.Form.Has(i) {
			fmt.Println(each)
			list3 = append(list3, each...)
		}
	}
	if len(list3) != 0 {
		finalList = CompareList(finalList, list3)
	}
	if loca := r.FormValue("locationsFilter"); loca != "" {
		fmt.Println(st.Bar.Locations[loca])
		fmt.Println(finalList)
		finalList = CompareList(finalList, st.Bar.Locations[loca])
	}
	return finalList
}

func CompareList(list1, list2 []int) []int {
	result := []int{}
	for _, each1 := range list1 {
		for _, each2 := range list2 {
			if each1 == each2 {
				result = append(result, each1)
			}
		}
	}
	return result
}

// get all id match the search
func GetIdSearch(s string) []int {
	const key = 1
	const value = 0
	searchTable := strings.Split(s, " - ")
	if len(searchTable) >= 2 {
		switch searchTable[key] {
		case "FirstAlbum_Date":
			return st.Bar.FirstAlbum_FullDate[searchTable[value]]
		case "Artist_BandName":
			return st.Bar.Artist_BandName[searchTable[value]]
		case "Creation_Date":
			return st.Bar.Creation_Date[searchTable[value]]
		case "Locations":
			return st.Bar.Locations[searchTable[value]]
		case "Members":
			return st.Bar.Members[searchTable[value]]
		}
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
	barTemp.FirstAlbum_FullDate = make(map[string][]int)
	barTemp.FirstAlbum_YearDate = make(map[string][]int)
	barTemp.Artist_BandName = make(map[string][]int)
	barTemp.Creation_Date = make(map[string][]int)
	barTemp.Locations = make(map[string][]int)
	barTemp.Members = make(map[string][]int)
	barTemp.NbMembers = make(map[string][]int)
	for _, data := range allArtistData {
		barTemp.Artist_BandName[data.ArtistName] = append(barTemp.Artist_BandName[data.ArtistName], data.Id)
		for _, members := range data.Members {
			barTemp.Members[members] = append(barTemp.Members[members], data.Id)
		}
		indexNbmembers := strconv.Itoa(len(data.Members))
		barTemp.NbMembers[indexNbmembers] = append(barTemp.NbMembers[indexNbmembers], data.Id)
		Cdate := strconv.Itoa(data.CreationDate)
		barTemp.Creation_Date[Cdate] = append(barTemp.Creation_Date[Cdate], data.Id)
		barTemp.FirstAlbum_FullDate[data.FirstAlbum] = append(barTemp.FirstAlbum_FullDate[data.FirstAlbum], data.Id)
	}
	for _, data := range allLocations["index"] {
		for _, loca := range data.Locations {
			barTemp.Locations[loca] = append(barTemp.Locations[loca], data.Id)
		}
	}
	FirstAlbumKeys := []string{}
	for i := range barTemp.FirstAlbum_FullDate {
		FirstAlbumKeys = append(FirstAlbumKeys, string(i[6:]))
	}
	sort.Strings(FirstAlbumKeys)
	for _, each := range FirstAlbumKeys {
		barTemp.FirstAlbum_YearDate[each] = append(barTemp.FirstAlbum_YearDate[each], []int{}...)
	}
	for i, each := range barTemp.FirstAlbum_FullDate {
		barTemp.FirstAlbum_YearDate[i[6:]] = append(barTemp.FirstAlbum_YearDate[i[6:]], each...)
	}
	st.Bar = barTemp
}
