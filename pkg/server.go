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

	// init encryption function passing the key
	s := Sec{enc: *internal.Init("0123456789abcdef")}

	// IDOR vulnerable routes
	serverMux.HandleFunc("/person/list/", ListP)
	serverMux.HandleFunc("/person/get/", GetP)
	// Secured routes
	serverMux.HandleFunc("/person/list/secure/", s.ListPSecure)
	serverMux.HandleFunc("/person/get/secure/", s.GetPSecure)

	log.Printf("Listening %s\n",host)
	err := server.ListenAndServe()
	if err != nil {
		log.Println(err.Error())
	}
}