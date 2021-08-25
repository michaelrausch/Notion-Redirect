package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

type ConfigFile struct {
	Notionurl string
}

func main() {
	r := mux.NewRouter()

	r.PathPrefix("/").
		Methods("GET").
		HandlerFunc(RedirectHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	conf, err := readConf("config.yaml")

	// There was an error loading or parsing config.yaml
	// Return an internal server error
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal Server Error")
		return
	}

	path := conf.Notionurl + r.URL.Path

	http.Redirect(w, r, path, http.StatusPermanentRedirect)
}

func readConf(filename string) (*ConfigFile, error) {
	buf, err := ioutil.ReadFile(filename)

	// Error reading the file
	if err != nil {
		return nil, err
	}

	c := &ConfigFile{}
	err = yaml.Unmarshal(buf, c)

	// Error parsing the YAML file
	if err != nil {
		return nil, err
	}

	return c, nil
}
