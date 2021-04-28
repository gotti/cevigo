package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

type Talker struct {
	*ole.IDispatch
}

type speakHandler struct {
	talker Talker
}

func (h speakHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := r.URL.Query().Get("cast")
	t := r.URL.Query().Get("text")
	fmt.Println("Cast:", c, "Text:", t)
	_, err := h.talker.PutProperty("Cast", c)
	if err != nil {
		w.WriteHeader(400)
		log.Fatal(err)
	}
	_, err = h.talker.GetProperty("Speak", t)
	if err != nil {
		w.WriteHeader(400)
		log.Fatal(err)
	}
}

type ttsHandler struct {
	talker Talker
}

func (h ttsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := r.URL.Query().Get("cast")
	t := r.URL.Query().Get("text")
	fmt.Println("Cast:", c, "Text:", t)
	_, err := h.talker.PutProperty("Cast", c)
	if err != nil {
		w.WriteHeader(400)
		log.Fatal(err)
	}
	_, err = h.talker.GetProperty("OutputWaveToFile", t, "./tmp.wav")
	if err != nil {
		w.WriteHeader(400)
		log.Fatal(err)
	}
	f, err := os.Open("./tmp.wav")
	if err != nil {
		w.WriteHeader(400)
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "audio/wav")
	w.WriteHeader(200)
	io.Copy(w, f)
}

func main() {
	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_DISABLE_OLE1DDE)
	talker, err := oleutil.CreateObject("CeVIO.Talk.RemoteService.Talker")
	if err != nil {
		log.Fatal("Initialization failed, Make sure you have installed CeVIO.")
		log.Fatal(err)
		os.Exit(1)
	}
	obj, err := talker.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		log.Fatal(err)
	}
	//speakはwindowsから直接音を出します．
	//speak?cast={キャスト名}&text={テキスト}
	http.Handle("/speak", speakHandler{Talker{obj}})
	//ttsは音声データをwavで応答します．
	//tts?cast={キャスト名}&text={テキスト}
	http.Handle("/tts", ttsHandler{Talker{obj}})
	err = http.ListenAndServe(":8085", nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
