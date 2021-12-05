package main

import (
	"flag"
	"fmt"
	"github.com/gotti/cevigo/pkg/cevioai"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Talker struct {
	cevioai.ITalker2V40
}

type speakHandler struct {
	talker Talker
}

var defaultParameters = map[string]int{
	"Volume":    50,
	"Speed":     50,
	"Tone":      50,
	"ToneScale": 50,
}

func (talker *Talker) setParameters(url *url.Values) error {
	cast := url.Get("cast")
	text := url.Get("text")
	fmt.Println("Cast:", cast, "Text:", text)
	speed := url.Get("speed")
	volume := url.Get("volume")
	talker.SetCast("さとうささら")
	if volume != "" {
		v, err := strconv.Atoi(volume)
		if err != nil {
			return err
		}
		err = talker.SetSpeed(v)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	if speed != "" {
		v, err := strconv.Atoi(speed)
		if err != nil {
			return err
		}
		err = talker.SetSpeed(v)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	tone := url.Get("tone")
	if tone != "" {
		v, err := strconv.Atoi(tone)
		if err != nil {
			return err
		}
		err = talker.SetTone(v)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	alpha := url.Get("alpha")
	if alpha != "" {
		v, err := strconv.Atoi(alpha)
		if err != nil {
			return err
		}
		err = talker.SetTone(v)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	toneScale := url.Get("toneScale")
	if toneScale != "" {
		v, err := strconv.Atoi(toneScale)
		if err != nil {
			return err
		}
		err = talker.SetTone(v)
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
		o, err := talker.GetComponents()
		if err != nil {
			log.Fatal(err)
		}
		if genki != "" {
			b, err := o.ByName("元気")
			if err != nil {
				return err
			}
			i, _ := strconv.Atoi(genki)
			err = b.SetValue(i)
			if err != nil {
				return err
			}
		}
		if hutu != "" {
			b, _ := o.ByName("普通")
			if err != nil {
				return err
			}
			i, _ := strconv.Atoi(hutu)
			err = b.SetValue(i)
			if err != nil {
				return err
			}
		}
		if ikari != "" {
			b, _ := o.ByName("怒り")
			if err != nil {
				return err
			}
			i, _ := strconv.Atoi(ikari)
			err = b.SetValue(i)
			if err != nil {
				return err
			}
		}
		if kanasimi != "" {
			b, _ := o.ByName("悲しみ")
			if err != nil {
				return err
			}
			i, _ := strconv.Atoi(kanasimi)
			err = b.SetValue(i)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func splitMultiple(text string, seps []string) (ret []string) {
	var list, next []string
	list = append(list, text)
	for _, sep := range seps {
		for _, t := range list {
			next = append(next, strings.Split(t, sep)...)
		}
		list = next
		next = []string{}
	}
	var splittedText []string
	for _, txt := range list {
		t := []rune(txt)
		if len(t) == 0 {
			continue
		}
		l := 200
		for i := 0; i < len(t); i += l {
			var o []rune
			if i+l <= len(t) {
				o = t[i:(i + l)]
			} else {
				o = t[i:]
			}
			splittedText = append(splittedText, string(o))
		}
	}
	for i := 0; i < len(splittedText); i++ {
		var concated string
		var j int
		for j = i; j < i+5; j++ {
			if j >= len(splittedText) {
				break
			}
			if len([]rune(concated+splittedText[j])) >= 200 {
				break
			}
		}
		for _, s := range splittedText[i:j] {
			concated += s + "。"
		}
		ret = append(ret, concated)
		i = j
	}
	return ret
}

func (h speakHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()
	h.talker.setParameters(&url)
	t := url.Get("text")
	texts := splitMultiple(t, []string{"。", " ", "　", "\n"})
	fmt.Println("texts:", texts)
	for _, t := range texts {
		state, err := h.talker.Speak(t)
		if err != nil {
			w.WriteHeader(400)
			log.Fatal(err)
		}
		err = state.Wait()
		if err != nil {
			w.WriteHeader(400)
			log.Fatal(err)
		}
	}
	w.WriteHeader(200)
}

type ttsHandler struct {
	talker Talker
}

func (h ttsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func main() {
	apidiff := flag.String("api", "cevio", "cevio, or cevioai")
	flag.Parse()
	var apiname string
	if *apidiff == "cevio" {
		apiname = "CeVIO.Talk.RemoteService.Talker"
	} else if *apidiff == "cevioai" {
		apiname = "CeVIO.Talk.RemoteService2.Talker2"
	} else {
		println("set cevio or cevioai to --api")
		os.Exit(1)
	}
	talker := cevioai.NewITalker2V40(apiname)
	talker.SetCast("さとうささら")
	fmt.Printf("connected to %s", apiname)
	//speakはwindowsから直接音を出します．
	//speak?cast={キャスト名}&text={テキスト}
	handler := Talker{talker}
	http.Handle("/speak", speakHandler{handler})
	//ttsは音声データをwavで応答します．
	//tts?cast={キャスト名}&text={テキスト}
	http.Handle("/tts", ttsHandler{handler})
	err := http.ListenAndServe(":8085", nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
