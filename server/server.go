package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"statistical-analysis/service"

	"github.com/gorilla/mux"
)

var (
	ErrUidNotFound = errors.New("Could not find uid")
	ErrStats       = errors.New("Could not get stats")
)

func Start(port string) {
	router := mux.NewRouter()

	router.HandleFunc("/{uid}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uid, ok := vars["uid"]
		if !ok {
			fmt.Println(ErrUidNotFound)
			response := fmt.Sprintf(`{"error": "%s"}`, ErrUidNotFound.Error())
			w.Write([]byte(response))
			return
		}

		fmt.Println(uid)
		data, err := service.CalculateStats(uid)
		if err != nil {
			fmt.Println(err)
			response := fmt.Sprintf(`{"error": "%s"}`, err.Error())
			w.Write([]byte(response))
			return
		}
		dataResp, err := json.Marshal(data)
		if err != nil {
			response := fmt.Sprintf(`{"error": "%s"}`, err.Error())
			w.Write([]byte(response))
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write(dataResp)
		return
	})

	log.Fatalln(http.ListenAndServe(":"+port, router))
}
