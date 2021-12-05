package cevioai

import (
	"math/rand"
	"testing"
)

func TestSpeed(t *testing.T) {
	talker := NewITalker2V40("CeVIO.Talk.RemoteService2.Talker2")
	defer talker.Release()
	talker.SetCast("さとうささら")
	r := rand.Int() % 100
	err := talker.SetSpeed(r)
	if err != nil {
		t.Errorf("%e", err)
	}
	g, err := talker.GetSpeed()
	if g != r {
		t.Errorf("get %d, expected %d", g, r)
	}
}
func TestVolume(t *testing.T) {
	talker := NewITalker2V40("CeVIO.Talk.RemoteService2.Talker2")
	talker.SetCast("さとうささら")
	defer talker.Release()
	r := rand.Int() % 100
	err := talker.SetVolume(r)
	if err != nil {
		t.Errorf("%e", err)
	}
	g, err := talker.GetVolume()
	if g != r {
		t.Errorf("get %d, expected %d", g, r)
	}
}

func TestCasts(t *testing.T) {
	talker := NewITalker2V40("CeVIO.Talk.RemoteService2.Talker2")
	talker.SetCast("さとうささら")
	defer talker.Release()
	a, err := talker.GetAvailableCasts()
	if err != nil {
		t.Errorf("%e", err)
	}
	b, err := a.GetAt(0)
	if err != nil {
		t.Errorf("%e", err)
	}
	c, err := a.GetLength()
	if err != nil {
		t.Errorf("%e", err)
	}
	if b != "さとうささら" || c != 1 {
		t.Errorf("got %s, %d, expected %s, %d", b, c, "さとうささら", 1)
	}
}

func TestGetComponents(t *testing.T) {
	talker := NewITalker2V40("CeVIO.Talk.RemoteService2.Talker2")
	talker.SetCast("さとうささら")
	defer talker.Release()
	c, err := talker.GetComponents()
	if err != nil {
		t.Errorf("%e", err)
	}
	l, err := c.GetLength()
	if err != nil {
		t.Errorf("%e", err)
	}
	if l != 4 {
		t.Errorf("got %d, expected %d", l, 4)
	}
	na, err := c.GetAt(0)
	n, err := na.GetName()
	if err != nil {
		t.Errorf("%e", err)
	}
	if n != "元気" {
		t.Errorf("got %s, expected %s", n, "元気")
	}
	r := rand.Int() % 100
	err = na.SetValue(r)
	if err != nil {
		t.Errorf("%e", err)
	}
	v, err := na.GetValue()
	if err != nil {
		t.Errorf("%e", err)
	}
	if v != r {
		t.Errorf("%d", v)
	}
}
