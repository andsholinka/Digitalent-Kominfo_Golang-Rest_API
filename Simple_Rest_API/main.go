package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Sisi struct {
	JenisBangun string `json:"jenis_bangun"`
	Panjang     int    `json:"panjang"`
	Lebar       int    `json:"lebar"`
}

type Hasil struct {
	JenisBangun string `json:"jenis_bangun"`
	Luas        int    `json:"luas"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/hitung-luas", Luas)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func Luas(w http.ResponseWriter, r *http.Request) {

	var hasilHitung []Hasil
	var sisi []Sisi
	if r.Method != "POST" {
		WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		WrapAPIError(w, r, "can't read body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &sisi)
	if err != nil {
		WrapAPIError(w, r, "error unmarshal : "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, v := range sisi {
		hasilHitung = append(hasilHitung, Hasil{
			JenisBangun: v.JenisBangun,
			Luas:        v.RumusLuas(),
		})
	}

	WrapAPIData(w, r, hasilHitung, http.StatusOK, "success")
}

func (s *Sisi) RumusLuas() int {
	return s.Panjang * s.Lebar
}

func WrapAPIError(w http.ResponseWriter, r *http.Request, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	result, err := json.Marshal(map[string]interface{}{
		"code":          code,
		"error_type":    http.StatusText(code),
		"error_details": message,
	})
	if err == nil {
		w.Write(result)
	} else {
		log.Println(fmt.Sprintf("can't wrap API error : %s", err))
	}
}

func WrapAPISuccess(w http.ResponseWriter, r *http.Request, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	result, err := json.Marshal(map[string]interface{}{
		"code":   code,
		"status": message,
	})
	if err == nil {
		log.Println(message)
		w.Write(result)
	} else {
		log.Println(fmt.Sprintf("can't wrap API success : %s", err))
	}
}

func WrapAPIData(w http.ResponseWriter, r *http.Request, data interface{}, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	result, err := json.Marshal(map[string]interface{}{
		"code":   code,
		"status": message,
		"data":   data,
	})
	if err == nil {
		log.Println(message)
		w.Write(result)
	} else {
		log.Println(fmt.Sprintf("can't wrap API data : %s", err))
	}
}
