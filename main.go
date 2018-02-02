//go:generate statik -src=./public

package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/apex/log"
	jlog "github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"
	"github.com/gorilla/pat"
	_ "github.com/kaihendry/longrunning/statik"
	"github.com/rakyll/statik/fs"
)

func init() {
	if os.Getenv("UP_STAGE") == "" {
		log.SetHandler(text.Default)
	} else {
		log.SetHandler(jlog.Default)
	}
}

func main() {
	app := pat.New()
	app.Post("/post", doSomething)

	statikFS, err := fs.New()
	if err != nil {
		log.WithError(err).Fatal("error compiling static resources")
	}

	app.PathPrefix("/").Handler(
		http.StripPrefix("/", http.FileServer(statikFS)))

	addr := ":" + os.Getenv("PORT")
	if err := http.ListenAndServe(addr, app); err != nil {
		log.WithError(err).Fatal("error listening")
	}
}

func doSomething(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var browserTime time.Time
	err := decoder.Decode(&browserTime)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	defer r.Body.Close()

	log.WithFields(log.Fields{
		"browserTime": browserTime,
	}).Info("input")

	// time.Sleep(20 * time.Second)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(time.Since(browserTime).String())
	return

}
