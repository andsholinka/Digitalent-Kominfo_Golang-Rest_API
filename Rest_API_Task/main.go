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
	Panjang float32 `json:"panjang"`
	Lebar   float32 `json:"lebar"`
	Tinggi  float32 `json:"tinggi"`
}

type Hasil struct {
	JenisBangun string  `json:"Jenis Bangun Ruang"`
	Volume      float32 `json:"Volume"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/Hitung-Volume", Volume)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func Volume(w http.ResponseWriter, r *http.Request) {
	var hasilHitung []Hasil
	var sisi []Sisi
	if r.Method != "POST" {
		WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		WrapAPIError(w, r, "Can't read body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &sisi)
	if err != nil {
		WrapAPIError(w, r, "error unmarshal : "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, v := range sisi {
		hasilHitung = append(hasilHitung, Hasil{
			JenisBangun: v.jenisBangun(),
			Volume:      v.RumusVolume(),
		})
	}

	WrapAPIData(w, r, hasilHitung, http.StatusOK, "success")
}

func (s *Sisi) RumusVolume() float32 {
	return s.Panjang * s.Lebar * s.Tinggi
}

func (s *Sisi) jenisBangun() string {
	if s.Panjang == s.Lebar && s.Lebar == s.Tinggi {
		return "Kubus"
	} else {
		return "Balok"
	}

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
		log.Println(fmt.Sprintf("Can't wrap API error : %s", err))
	}
}

func WrapAPISuccess(w http.ResponseWriter, r *http.Request, message string, code int) {
	w.Header().Set("Content-Type", "appliction/json")
	w.WriteHeader(code)
	result, err := json.Marshal(map[string]interface{}{
		"code":   code,
		"status": message,
	})
	if err == nil {
		log.Println(message)
		w.Write(result)
	} else {
		log.Println(fmt.Sprintf("Can't wrap API success : %s", err))
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
		log.Println(fmt.Sprintf("Can't wrap API data : %s", err))
	}
}
