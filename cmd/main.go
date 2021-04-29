package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

type component struct {
	Id    string
	Name  string
	Value uint
}

type components struct {
	Length int
}

type Talker struct {
	*ole.IDispatch
	parameters *map[string]*parameter
}

type speakHandler struct {
	talker Talker
}

type parameter struct {
	Name  string
	Value interface{}
}

func (param *parameter) set(value interface{}) {
	param.Value = value
}

func (talker *Talker) setIntParameter(param string, value string) error {
	i, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if (*(talker.parameters))[param].Value.(int) != i {
		talker.PutProperty("ToneScale", i)
		print((*(talker.parameters))[param].Name)
		_, err = talker.PutProperty((*(talker.parameters))[param].Name, i)
		if err != nil {
			log.Fatal(err)
			return err
		}
		(*(talker.parameters))[param].set(i)
	}
	return nil
}

func (talker *Talker) setParameters(url *url.Values) error {
	cast := url.Get("cast")
	text := url.Get("text")
	fmt.Println("Cast:", cast, "Text:", text)
	_, err := talker.PutProperty("Cast", cast)
	if err != nil {
		log.Fatal(err)
		return err
	}
	speed := url.Get("speed")
	if speed != "" {
		err := talker.setIntParameter("speed", speed)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	tone := url.Get("tone")
	if tone != "" {
		err := talker.setIntParameter("tone", tone)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	alpha := url.Get("alpha")
	if alpha != "" {
		err := talker.setIntParameter("alpha", alpha)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	toneScale := url.Get("toneScale")
	if toneScale != "" {
		err := talker.setIntParameter("toneScale", toneScale)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	genki := url.Get("genki")
	hutu := url.Get("hutu")
	ikari := url.Get("ikari")
	kanasimi := url.Get("kanasimi")
	if genki != "" || hutu != "" || ikari != "" || kanasimi != "" {
		o, err := talker.GetProperty("Components")
		if err != nil {
			log.Fatal(err)
		}
		a := o.ToIDispatch()
		if genki != "" {
			b, _ := a.GetProperty("ByName", "元気")
			d := b.ToIDispatch()
			i, _ := strconv.Atoi(genki)
			_, err = d.PutProperty("Value", i)
		}
		if hutu != "" {
			b, _ := a.GetProperty("ByName", "普通")
			d := b.ToIDispatch()
			i, _ := strconv.Atoi(hutu)
			_, err = d.PutProperty("Value", i)
		}
		if genki != "" {
			b, _ := a.GetProperty("ByName", "怒り")
			d := b.ToIDispatch()
			i, _ := strconv.Atoi(ikari)
			_, err = d.PutProperty("Value", i)
		}
		if genki != "" {
			b, _ := a.GetProperty("ByName", "哀しみ")
			d := b.ToIDispatch()
			i, _ := strconv.Atoi(kanasimi)
			_, err = d.PutProperty("Value", i)
		}
	}
	return err
}

func (h speakHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()
	h.talker.setParameters(&url)
	t := url.Get("text")
	_, err := h.talker.GetProperty("Speak", t)
	if err != nil {
		w.WriteHeader(400)
		log.Fatal(err)
	}
}

type ttsHandler struct {
	talker Talker
}

func (h ttsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()
	h.talker.setParameters(&url)
	t := url.Get("text")
	_, err := h.talker.GetProperty("OutputWaveToFile", t, "./tmp.wav")
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
	var parameters = map[string]*parameter{
		"cast":      {"Cast", "さとうささら"},
		"speed":     {"Speed", 50},
		"tone":      {"Tone", 50},
		"alpha":     {"Alpha", 50},
		"toneScale": {"ToneScale", 50},
		"genki":     {"元気", 100},
		"hutu":      {"普通", 0},
		"ikari":     {"怒り", 0},
		"kanasimi":  {"悲しみ", 0},
	}
	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_DISABLE_OLE1DDE)
	talker, err := oleutil.CreateObject("CeVIO.Talk.RemoteService.Talker")
	fmt.Println(talker.VTable().QueryInterface)
	if err != nil {
		log.Fatal("Initialization failed, Make sure you have installed CeVIO.")
		log.Fatal(err)
		os.Exit(1)
	}
	obj, err := talker.QueryInterface(ole.IID_IDispatch)
	//fmt.Println(obj.GetSingleIDOfName("ToneScale"))
	_, err = obj.PutProperty("Cast", "さとうささら")
	if err != nil {
		log.Fatal(err)
	}
	//speakはwindowsから直接音を出します．
	//speak?cast={キャスト名}&text={テキスト}
	handler := Talker{obj, &parameters}
	http.Handle("/speak", speakHandler{handler})
	//ttsは音声データをwavで応答します．
	//tts?cast={キャスト名}&text={テキスト}
	http.Handle("/tts", ttsHandler{handler})
	err = http.ListenAndServe(":8085", nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
