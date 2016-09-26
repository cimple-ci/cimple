package messages

import "reflect"

func NewRouter() *Router {
	router := &Router{
		maps: make(map[interface{}]func(interface{})),
	}
	return router
}

type Router struct {
	maps    map[interface{}]func(interface{})
	onError func(interface{})
}

func (r *Router) On(typ interface{}, callback func(interface{})) {
	r.maps[reflect.TypeOf(typ)] = callback
}

func (r *Router) OnError(callback func(interface{})) {
	r.onError = callback
}

func (r *Router) Route(msg interface{}) {
	typ := reflect.TypeOf(msg)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	val, ok := r.maps[typ]
	if ok {
		val(msg)
	} else {
		if r.onError != nil {
			r.onError(msg)
		}
	}
}
