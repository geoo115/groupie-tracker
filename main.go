package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

var artists []Artist
var locations Locations
var dates Dates
var relations Relations
var client *http.Client

type Artist struct {
	ID           int64    `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int64    `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    Locations
	Dates        Dates
	Relations    Relations
}

type Locations struct {
	Index []Location `json:"index"`
}

type Location struct {
	ID        int64    `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type Dates struct {
	Index []Date `json:"index"`
}

type Date struct {
	ID    int64    `json:"id"`
	Dates []string `json:"dates"`
}

type Relations struct {
	Index []Relation `json:"index"`
}

type Relation struct {
	ID             int64                  `json:"id"`
	DatesLocations map[string]interface{} `json:"datesLocations"`
}

func Getjson(url string, target interface{}) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API responded with status code: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func main() {
	client = &http.Client{Timeout: 10 * time.Second}

	Getjson("https://groupietrackers.herokuapp.com/api/artists", &artists)
	Getjson("https://groupietrackers.herokuapp.com/api/locations", &locations)
	Getjson("https://groupietrackers.herokuapp.com/api/dates", &dates)
	Getjson("https://groupietrackers.herokuapp.com/api/relation", &relations)
	AppendToStruct() // Associate locations, dates, and relations with each artist

	// Handler added to Read the Static files within the assets folder: .css, .js, img
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./assets"))))

	http.HandleFunc("/", dateHandler)
	http.HandleFunc("/loc", locHandle)
	http.HandleFunc("/dates", datesHandle)
	http.HandleFunc("/rel", relationsHandle)

	log.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func locHandle(w http.ResponseWriter, r *http.Request) {
	// Get the selected letter from the URL query parameter
	selectedLetter := r.FormValue("letter")

	// Filter the artists based on the selected letter
	var filteredArtists []Artist
	for _, artist := range artists {
		if strings.HasPrefix(strings.ToLower(artist.Name), strings.ToLower(selectedLetter)) {
			filteredArtists = append(filteredArtists, artist)
		}
	}

	// Prepare data for the template
	data := struct {
		Artists  []Artist
		Alphabet []string
	}{
		Artists:  filteredArtists,
		Alphabet: generateAlphabet(),
	}

	// Serve the HTML page with the artists
	tmpl := template.Must(template.ParseFiles("templates/locations.html"))
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func datesHandle(w http.ResponseWriter, r *http.Request) {
	// Get the selected letter from the URL query parameter
	selectedLetter := r.FormValue("letter")

	// Filter the artists based on the selected letter
	var filteredArtists []Artist
	for _, artist := range artists {
		if strings.HasPrefix(strings.ToLower(artist.Name), strings.ToLower(selectedLetter)) {
			filteredArtists = append(filteredArtists, artist)
		}
	}

	// Prepare data for the template
	data := struct {
		Artists  []Artist
		Alphabet []string
	}{
		Artists:  filteredArtists,
		Alphabet: generateAlphabet(),
	}

	// Serve the HTML page with the filtered artists
	tmpl := template.Must(template.ParseFiles("templates/dates.html"))
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func relationsHandle(w http.ResponseWriter, r *http.Request) {
	// Get the selected letter from the URL query parameter
	selectedLetter := r.FormValue("letter")

	// Filter the artists based on the selected letter
	var filteredArtists []Artist
	for _, artist := range artists {
		if strings.HasPrefix(strings.ToLower(artist.Name), strings.ToLower(selectedLetter)) {
			filteredArtists = append(filteredArtists, artist)
		}
	}

	// Prepare data for the template
	data := struct {
		Artists  []Artist
		Index    []Relation
		Alphabet []string
	}{
		Artists:  filteredArtists,
		Index:    relations.Index,
		Alphabet: generateAlphabet(),
	}

	// Serve the HTML page with the relations
	tmpl := template.Must(template.ParseFiles("templates/relations.html"))
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AppendToStruct() {
	for index := range locations.Index {
		artists[index].Locations.Index = append(artists[index].Locations.Index, locations.Index[index])
	}

	for index := range dates.Index {
		artists[index].Dates.Index = append(artists[index].Dates.Index, dates.Index[index])
	}

	for index := range relations.Index {
		artists[index].Relations.Index = append(artists[index].Relations.Index, relations.Index[index])
	}
}

func dateHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 Page Not Found", http.StatusNotFound)
		return
	}
	// Get the selected letter from the URL query parameter
	selectedLetter := r.FormValue("letter")

	// Filter the artists based on the selected letter
	var filteredArtists []Artist
	for _, artist := range artists {
		if strings.HasPrefix(strings.ToLower(artist.Name), strings.ToLower(selectedLetter)) {
			filteredArtists = append(filteredArtists, artist)
		}
	}

	// Prepare data for the template
	data := struct {
		Artists  []Artist
		Alphabet []string
	}{
		Artists:  filteredArtists,
		Alphabet: generateAlphabet(),
	}

	// Serve the HTML page with the filtered artists
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// generateAlphabet generates an alphabet array from A to Z
func generateAlphabet() []string {
	alphabet := make([]string, 26)
	for i := 0; i < 26; i++ {
		alphabet[i] = string('A' + i)
	}
	return alphabet
}
