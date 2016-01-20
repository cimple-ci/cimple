package messages

import "reflect"

func NewRouter() *Router {
	router := &Router{
		maps: make(map[interface{}]func(interface{})),
	}
	return router
}

type Router struct {
	maps map[interface{}]func(interface{})
}

func (r *Router) On(typ interface{}, callback func(interface{})) {
	r.maps[reflect.TypeOf(typ)] = callback
}

func (r *Router) Route(msg interface{}) {
	if val, ok := r.maps[reflect.TypeOf(msg)]; ok {
		val(msg)
	}
}
