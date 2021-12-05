package olewrapper

import (
	"github.com/go-ole/go-ole"
)

type Talker struct {
	Obj    *ole.IUnknown
	Talker *ole.IDispatch
}

func NewTalker(obj *ole.IUnknown, talker *ole.IDispatch) Talker {
	return Talker{obj, talker}
}

func (t *Talker) Release() {
	t.Talker.Release()
	t.Obj.Release()
}
