package pkg

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"mvpidx/internal"
	"net/http"
	"strconv"
)

var Map = map[string][]int{
	"elf": {1, 2},
	"gnome": {3, 4},
}

var Enc internal.Encoder

var Usrs = map[int]Person {
	1: Person{Name: "Pietro", 		Age: 15},
	2: Person{Name: "Lucas", 		Age: 15},
	3: Person{Name: "Andrey", 		Age: 15},
	4: Person{Name: "Kerleston", 	Age: 15},
}

type Response struct {
	Status 	string 		`json:"Status"`
	Payload interface{} `json:"Payload"`
}

type ResumeS struct {
	Id   	string `json:"Id"`
	Name 	string `json:"Name"`
}
type Resume struct {
	Id   	int 	`json:"Id"`
	Name 	string 	`json:"Name"`
}
type Person struct {
	Name 	string 	`json:"Name"`
	Age 	int		`json:"Age"`
}

func decToken(token string) (string, string, error) {
	bts, err := b64.StdEncoding.DecodeString(token)
	if err != nil {
		return "","", err
	}

	var f interface{}

	err = json.Unmarshal(bts, &f)
	if err != nil {
		return "","", err
	}

	m := f.(map[string]interface{})
	nounce := m["nounce"].(string)
	user   := m["user"].(string)
	
	return nounce, user, nil
}

func ListP(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Token")
	w.Header().Set("Content-Type","application/json")
	respJsoned := Response{Status: "OK"}

	n, u, err := decToken(token)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("{\"err\":\"%s\"}", err.Error())))
		log.Printf("access uri=%s u=%s n=%s err=\"%s\"\n", r.URL.Path, u, n, err.Error())
		return
	}
	log.Printf("access uri=%s u=%s n=%s\n", r.URL.Path, u, n)

	accessList := Map[u]
	var respList []Resume
	for _,v := range accessList {
		respList = append(respList, Resume{Id: v, Name: Usrs[v].Name})
	}

	respJsoned.Payload = respList
	jsoned, err := json.Marshal(respJsoned)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("{\"err\":\"%s\"}", err.Error())))
		log.Printf("access uri=%s u=%s n=%s err=\"%s\"\n", r.URL.Path, u, n, err.Error())
		return
	}

	w.Write(jsoned)
}

func GetP(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Token")
	w.Header().Set("Content-Type","application/json")
	respJsoned := Response{Status: "OK"}

	n, u, err := decToken(token)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("{\"err\":\"%s\"}", err.Error())))
		log.Printf("access uri=%s u=%s n=%s err=\"%s\"\n", r.URL.Path, u, n, err.Error())
		return
	}
	log.Printf("access uri=%s u=%s n=%s\n", r.URL.Path, u, n)

	pid := r.URL.Query().Get("id")

	intPid, err := strconv.Atoi(pid)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("{\"err\":\"%s\"}", err.Error())))
		log.Printf("access uri=%s u=%s n=%s err=\"%s\"\n", r.URL.Path, u, n, err.Error())
		return
	}

	respJsoned.Payload = Usrs[intPid]
	jsoned, err := json.Marshal(respJsoned)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("{\"err\":\"%s\"}", err.Error())))
		log.Printf("access uri=%s u=%s n=%s err=\"%s\"\n", r.URL.Path, u, n, err.Error())
		return
	}

	w.Write(jsoned)
}

func ListPSecure(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Token")
	w.Header().Set("Content-Type","application/json")
	respJsoned := Response{Status: "OK"}

	n, u, err := decToken(token)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("{\"err\":\"%s\"}", err.Error())))
		log.Printf("access uri=%s u=%s n=%s err=\"%s\"\n", r.URL.Path, u, n, err.Error())
		return
	}
	log.Printf("access uri=%s u=%s n=%s\n", r.URL.Path, u, n)

	accessList := Map[u]
	var respList []ResumeS
	for _,v := range accessList {
		respList = append(respList, ResumeS{Id: Enc.Encode(n, v), Name: Usrs[v].Name})
	}

	respJsoned.Payload = respList
	jsoned, err := json.Marshal(respJsoned)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("{\"err\":\"%s\"}", err.Error())))
		log.Printf("access uri=%s u=%s n=%s err=\"%s\"\n", r.URL.Path, u, n, err.Error())
		return
	}

	w.Write(jsoned)
}

func GetPSecure(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Token")
	w.Header().Set("Content-Type","application/json")
	respJsoned := Response{Status: "OK"}

	n, u, err := decToken(token)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("{\"err\":\"%s\"}", err.Error())))
		log.Printf("access uri=%s u=%s n=%s err=\"%s\"\n", r.URL.Path, u, n, err.Error())
		return
	}


	pid := r.URL.Query().Get("id")

	nounce,id := Enc.Decode(pid)
	if nounce != n {
		w.Write([]byte("{\"err\":\"nounce do not match\"}"))
		log.Printf("access uri=%s u=%s n=%s err=\"%s\"\n", r.URL.Path, u, n, "nounce do not match")
		return
	}

	/*if err != nil {
		w.Write([]byte(fmt.Sprintf("{\"err\":\"%s\"}", err.Error())))
		log.Printf("access uri=%s u=%s n=%s err=\"%s\"\n", r.URL.Path, u, n, err.Error())
		return
	}*/

	respJsoned.Payload = Usrs[id]
	jsoned, err := json.Marshal(respJsoned)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("{\"err\":\"%s\"}", err.Error())))
		log.Printf("access uri=%s u=%s n=%s err=\"%s\"\n", r.URL.Path, u, n, err.Error())
		return
	}

	log.Printf("access uri=%s u=%s n=%s\n", r.URL.Path, u, n)
	w.Write(jsoned)
}