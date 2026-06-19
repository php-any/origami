package http

import (
	"sync"

	httpsrc "net/http"

	"github.com/php-any/origami/data"
)

var requestAttrBags sync.Map

func attachRequestAttrs(r *httpsrc.Request) {
	if r == nil {
		return
	}
	requestAttrBags.LoadOrStore(r, make(map[string]data.Value))
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

func requestAttrs(r *httpsrc.Request) map[string]data.Value {
	if r == nil {
		return nil
	}
	if v, ok := requestAttrBags.Load(r); ok {
		return v.(map[string]data.Value)
	}
	bag := make(map[string]data.Value)
	requestAttrBags.Store(r, bag)
	return bag
}

func beginRequest(r *httpsrc.Request) (*httpsrc.Request, data.ClassStmt) {
	attachRequestAttrs(r)
	return r, NewRequestClassFrom(r)
}
