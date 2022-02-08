package cevioai

import (
	"github.com/gotti/cevigo/pkg/olewrapper"
	"log"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

type ITalker2V40 struct {
	*olewrapper.Talker
}

func NewITalker2V40(apiname string) ITalker2V40 {
	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_DISABLE_OLE1DDE)
	obj, err := oleutil.CreateObject(apiname)
	if err != nil {
		log.Fatalf("Initialization failed, Make sure you have installed %s", apiname)
		log.Fatal(err)
	}
	talker, err := obj.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		log.Fatal(err)
	}
	s := olewrapper.NewTalker(obj, talker)
	t := ITalker2V40{&s}
	return t
}

func (t *ITalker2V40) putIntProp(key string, value int) error {
	_, err := t.Talker.Talker.PutProperty(key, value)
	if err != nil {
		return err
	}
	return nil
}

func (t *ITalker2V40) getIntProp(key string) (int, error) {
	val, err := t.Talker.Talker.GetProperty(key)
	if err != nil {
		return 0, err
	}
	i := int(val.Val)
	return i, nil
}

func (t *ITalker2V40) SetVolume(val int) error {
	return t.putIntProp("Volume", val)
}
func (t *ITalker2V40) GetVolume() (int, error) {
	return t.getIntProp("Volume")
}
func (t *ITalker2V40) SetSpeed(val int) error {
	return t.putIntProp("Speed", val)
}
func (t *ITalker2V40) GetSpeed() (int, error) {
	return t.getIntProp("Speed")
}
func (t *ITalker2V40) SetTone(val int) error {
	return t.putIntProp("Tone", val)
}
func (t *ITalker2V40) GetTone() (int, error) {
	return t.getIntProp("Tone")
}
func (t *ITalker2V40) SetToneScale(val int) error {
	return t.putIntProp("ToneScale", val)
}
func (t *ITalker2V40) GetToneScale() (int, error) {
	return t.getIntProp("ToneScale")
}
func (t *ITalker2V40) SetAlpha(val int) error {
	return t.putIntProp("Alpha", val)
}
func (t *ITalker2V40) GetAlpha() (int, error) {
	return t.getIntProp("Alpha")
}

func (t *ITalker2V40) GetComponents() (*ITalkerComponentArray2, error) {
	v, err := t.Talker.Talker.GetProperty("Components")
	if err != nil {
		return nil, err
	}
	ar := NewITalkerComponentArray2(v.ToIDispatch())
	return &ar, nil
}

func (t *ITalker2V40) SetCast(val string) (string, error) {
	v, err := t.Talker.Talker.PutProperty("Cast", val)
	if err != nil {
		return "", err
	}
	return v.ToString(), nil
}

func (t *ITalker2V40) GetCast() (string, error) {
	v, err := t.Talker.Talker.GetProperty("Cast")
	if err != nil {
		return "", err
	}
	return v.ToString(), nil
}

func (t *ITalker2V40) GetAvailableCasts() (*IStringArray, error) {
	v, err := t.Talker.Talker.GetProperty("AvailableCasts")
	if err != nil {
		return nil, err
	}
	d := IStringArray{dispatch: v.ToIDispatch()}
	return &d, nil
}

func (t *ITalker2V40) Speak(value string) (*ISpeakingState2, error) {
	s, err := t.Talker.Talker.GetProperty("Speak", value)
	if err != nil {
		return nil, err
	}
	d := ISpeakingState2{dispatch: s.ToIDispatch()}
	return &d, nil
}

func (t *ITalker2V40) OutputWaveToFile(text string, path string) (bool, error) {
	v, err := t.Talker.Talker.GetProperty("OutputWaveToFile", text, path)
	if err != nil {
		return false, err
	}
	d := v.Value().(bool)
	return d, nil
}

type IStringArray struct {
	dispatch *ole.IDispatch
}

func (s IStringArray) ToGoArray() (ret []string, err error) {
	l, err := s.GetLength()
	if err != nil {
		return make([]string, 0), err
	}
	for i := 0; i < l; i++ {
		d, err := s.GetAt(i)
		if err != nil {
			return make([]string, 0), err
		}
		ret = append(ret, d)
	}
	return ret, nil
}

func (s IStringArray) GetLength() (int, error) {
	v, err := s.dispatch.GetProperty("Length")
	if err != nil {
		return 0, err
	}
	r := int(v.Val)
	return r, nil
}

func (s IStringArray) GetAt(val int) (string, error) {
	v, err := s.dispatch.CallMethod("At", val)
	if err != nil {
		return "", err
	}
	r := v.Value().(string)
	return r, nil
}
