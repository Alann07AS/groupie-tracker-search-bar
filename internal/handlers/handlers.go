package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"gt-alann/config"
	apimanagement "gt-alann/internal/apiManagement"
	"gt-alann/internal/serverManagement"
)

var appConfig *config.Config

func ConfigHandle() {
	appConfig = config.ConfigLoad()
}

// _____________HANDLER______________//
func HomeHandle(w http.ResponseWriter, r *http.Request) {
	renderTmpl(w, "allArtists", nil)
}

func AllArtistsHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 ERROR, page not found.", http.StatusNotFound)
		return
	}
	r.ParseForm()
	valueSearch := r.FormValue("search")
	if valueSearch == "" {
		renderTmpl(w, "allArtists", apimanagement.GetAllArtistsSimpleApi())
	} else {
		if list := apimanagement.GetIdSearch(valueSearch); list != nil {
			renderTmpl(w, "allArtists", apimanagement.GetNewSliceByIdArtistsSimpleApi(list))
		} else {
			fmt.Fprint(w, "bad search")
		}
	}
}

func ArtistHandle(w http.ResponseWriter, r *http.Request) {
	id := 0
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "method error:", http.StatusInternalServerError)
	}
	val := r.FormValue("group")
	id = apimanagement.GetAllArtistsSimpleApiByName(val)
	if val == "" {
		http.Error(w, "Bad method, not method get", http.StatusBadRequest)
		return
	}
	renderTmpl(w, "artist", apimanagement.GetAllInfoArtists(id))
}

func AdminHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	serverManagement.ServeurAction(r.FormValue("bt"))
	renderTmpl(w, "admin", nil)
}

//_____________HANDLER______________//

func renderTmpl(w http.ResponseWriter, tmplName string, data any) {
	templateCache := appConfig.TemplateCache

	tmpl, ok := templateCache[tmplName+".page.tmpl"]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := make(map[string]*template.Template)
	pages, err := filepath.Glob("templates/*.page.tmpl")
	if err != nil {
		return cache, err
	}
	for _, eachPage := range pages {
		name := filepath.Base(eachPage)
		tmpl := template.Must(template.ParseFiles(eachPage))

		layouts, err := filepath.Glob("templates/layouts/*.layout.tmpl")
		if err != nil {
			return cache, err
		}

		if len(layouts) > 0 {
			tmpl.ParseGlob("templates/layouts/*.layout.tmpl")
		}
		cache[name] = tmpl
	}
	fmt.Println(cache)
	return cache, nil
}
