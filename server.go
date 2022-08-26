package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Data struct {
	Websites []string `json:"websites"`
}

var u Data
var status_map = make(map[string]string)

func main() {
	fmt.Println("Server started listening...")
	go Status_update()
	http.HandleFunc("/", StatusHandler)
	http.ListenAndServe("localhost:8080", nil)
}

func Status_update() {
	for {
		for _, v := range u.Websites {
			resp, err := http.Get(v)
			if err == nil && resp.StatusCode >= 200 && resp.StatusCode <= 299 {
				status_map[v] = "UP"
			} else {
				status_map[v] = "DOWN"
			}
		}
		time.Sleep(60 * time.Second)
	}
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	data_acc := Data{}

	switch r.Method {
	case "POST":
		err := json.NewDecoder(r.Body).Decode(&data_acc)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		var f int
		for _, v := range data_acc.Websites {
			f = 0
			for _, a := range u.Websites {
				if v == a {
					f = 1
				}
			}
			if f == 0 {
				u.Websites = append(u.Websites, v)
			}
		}
		fmt.Fprint(w, "\n")
		fmt.Fprintf(w, "Data got: %s\n", data_acc.Websites)
		fmt.Fprintf(w, "Post request Completed\n")

	case "GET":
		url := r.URL.Query().Get("url")
		fmt.Fprint(w, "\n")
		if url != "" {
			if _, ok := status_map[url]; ok {
				fmt.Fprintf(w, "%s: %s", url[7:], status_map[url])
			} else {
				fmt.Fprint(w, "Website URL not present in the list")
			}
			return
		} else {
			fmt.Fprint(w, "Printing status of all the urls present in the map\n")
			for k, v := range status_map {
				fmt.Fprintf(w, "%s: %s\n", k[7:], v)
			}
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
