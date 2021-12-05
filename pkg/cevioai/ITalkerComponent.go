package cevioai

import (
	"fmt"

	"github.com/go-ole/go-ole"
)

type ITalkerComponentArray2 struct {
	dispatch *ole.IDispatch
}

func NewITalkerComponentArray2(d *ole.IDispatch) ITalkerComponentArray2 {
	return ITalkerComponentArray2{dispatch: d}
}

func (t *ITalkerComponentArray2) GetLength() (int, error) {
	v, err := t.dispatch.GetProperty("Length")
	if err != nil {
		return 0, err
	}
	i := int(v.Val)
	return i, nil
}

func (t *ITalkerComponentArray2) GetAt(val int) (*ITalkerComponent2, error) {
	v, err := t.dispatch.GetProperty("At", val)
	if err != nil {
		return nil, err
	}
	c := ITalkerComponent2{dispatch: v.ToIDispatch()}
	return &c, nil
}

func (t *ITalkerComponentArray2) ByName(val string) (*ITalkerComponent2, error) {
	v, err := t.dispatch.GetProperty("ByName", val)
	if err != nil {
		return nil, err
	}
	c := ITalkerComponent2{dispatch: v.ToIDispatch()}
	return &c, nil
}

type ITalkerComponent2 struct {
	dispatch *ole.IDispatch
}

func (t *ITalkerComponent2) GetId() (string, error) {
	v, err := t.dispatch.GetProperty("Id")
	if err != nil {
		return "", err
	}
	i, ok := v.Value().(string)
	if !ok {
		return "", fmt.Errorf("Failed to interface assertion in getting Id of emotion")
	}
	return i, nil
}

func (t *ITalkerComponent2) GetName() (string, error) {
	v, err := t.dispatch.GetProperty("Name")
	if err != nil {
		return "", err
	}
	i, ok := v.Value().(string)
	if !ok {
		return "", fmt.Errorf("Failed to interface assertion in getting name of emotion")
	}
	return i, nil
}

func (t *ITalkerComponent2) SetValue(val int) error {
	_, err := t.dispatch.PutProperty("Value", val)
	if err != nil {
		return err
	}
	return nil
}

func (t *ITalkerComponent2) GetValue() (int, error) {
	v, err := t.dispatch.GetProperty("Value")
	if err != nil {
		return 0, err
	}
	i := int(v.Val)
	return i, nil
}
