package http

import (
	"sync"

	httpsrc "net/http"

	"github.com/php-any/origami/data"
)

var requestAttrBags sync.Map

type requestAttrBag struct {
	sync.RWMutex
	values map[string]data.Value
}

func attachRequestAttrs(r *httpsrc.Request) {
	if r == nil {
		return
	}
	requestAttrBags.LoadOrStore(r, &requestAttrBag{values: make(map[string]data.Value)})
}

var requestFormatterSlots sync.Map

func attachRequestFormatter(r *httpsrc.Request, slot *formatHandlerSlot) {
	if r == nil || slot == nil {
		return
	}
	requestFormatterSlots.Store(r, slot)
}

func requestFormatterFor(r *httpsrc.Request) *formatHandlerSlot {
	if r == nil {
		return nil
	}
	if v, ok := requestFormatterSlots.Load(r); ok {
		return v.(*formatHandlerSlot)
	}
	return nil
}

func detachRequestAttrs(r *httpsrc.Request) {
	if r != nil {
		requestAttrBags.Delete(r)
		requestFormatterSlots.Delete(r)
	}
}

func requestAttrs(r *httpsrc.Request) *requestAttrBag {
	if r == nil {
		return nil
	}
	if v, ok := requestAttrBags.Load(r); ok {
		return v.(*requestAttrBag)
	}
	bag := &requestAttrBag{values: make(map[string]data.Value)}
	actual, _ := requestAttrBags.LoadOrStore(r, bag)
	if actual != nil {
		return actual.(*requestAttrBag)
	}
	return bag
}

func beginRequest(r *httpsrc.Request) (*httpsrc.Request, data.ClassStmt) {
	attachRequestAttrs(r)
	return r, NewRequestClassFrom(r)
}
