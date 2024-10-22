package pkg

import (
	"log"
	"mvpidx/internal"
	"net/http"
)

func ServerStart() {
	host := "0.0.0.0:8080"
	serverMux := http.NewServeMux()
	server := &http.Server{
		Addr: 		host,
		Handler: 	serverMux,
	}

	serverMux.HandleFunc("/person/list/", ListP)
	serverMux.HandleFunc("/person/get/", GetP)
	serverMux.HandleFunc("/person/list/secure/", ListPSecure)
	serverMux.HandleFunc("/person/get/secure/", GetPSecure)

	log.Printf("Listening %s\n",host)
	Enc = *internal.Init("0123456789abcdef")
	err := server.ListenAndServe()
	if err != nil {
		log.Println(err.Error())
	}
}