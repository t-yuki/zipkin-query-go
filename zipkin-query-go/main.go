package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/t-yuki/zipkin-go/models"
	"github.com/t-yuki/zipkin-go/storage"
)

var flagListen = flag.String("http", ":8081", "http listen addr")

func main() {
	log.Print("start zipkin-query-go")
	stor, err := storage.Open()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/api/v1/spans", &StoreSpansHandler{stor})

	ln, err := net.Listen("tcp", *flagListen)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		err := http.Serve(ln, nil)
		if err != nil {
			log.Print(err)
		}
		stor.Close()
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Print("stop zipkin-query-go")
	ln.Close()
}

type StoreSpansHandler struct {
	stor storage.Storage
}

func (h StoreSpansHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Print("HandleStoreSpans")
	var req []*models.Span
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "json.Decode: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := h.stor.StoreSpans(req)
	if err != nil {
		http.Error(w, "storage.StoreSpans: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	return
}
