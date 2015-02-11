package rpc

import(
	"king/interfaces"
)

type outputstring map[string]interfaces.OutputString

type webSocketService struct {
	Methods outputstring
}

func (r *webSocketService) GetMethod(name string) interfaces.OutputString {
	if method, found := r.Methods[name]; found {
		return method
	}
	return nil
}

func (r *webSocketService) Exports(name string, method interfaces.OutputString) {
	if _, found := r.Methods[name]; found {
		return
	}

	r.Methods[name] = method
}


var WebSocketService = &webSocketService{outputstring{}}
