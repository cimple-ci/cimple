package messages

import "testing"

type myMessage struct {
}

func Test_Routes_Address_Reference(t *testing.T) {
	var resp interface{}
	message := &myMessage{}

	router := NewRouter()
	router.On(myMessage{}, func(m interface{}) {
		resp = m
	})

	router.Route(*message)

	if resp == nil {
		t.Error("message was not routed as expected")
	}
}

func Test_Routes_Pointer_Reference(t *testing.T) {
	var resp interface{}

	router := NewRouter()
	router.On(myMessage{}, func(m interface{}) {
		resp = m
	})

	message := &myMessage{}
	router.Route(message)

	if resp != message {
		t.Error("message was not routed as expected")
	}
}

func Test_Routes_Unknown(t *testing.T) {
	var resp interface{}
	message := &myMessage{}

	router := NewRouter()
	router.OnError(func(m interface{}) {
		resp = m
	})

	router.Route(message)

	if resp != message {
		t.Error("Message was not routed to OnError callback")
	}
}
