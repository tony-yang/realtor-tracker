package base

import (
	"strings"
)

// Routes create the route controller mapping and routes the traffic to
// the corresponding controller for further processing
type Routes struct {
	routes map[string]ControllerInterface
}

func CreateNewRouter() *Routes {
	return &Routes{
		map[string]ControllerInterface{},
	}
}

func (r *Routes) Exists(path string) bool {
	resources := strings.Split(path, "/")[1:]
	controllerPath := "/" + resources[0]

	if _, ok := r.routes[controllerPath]; ok {
		return true
	}
	return false
}

func (r *Routes) GetController(path string) (ControllerInterface, string) {
	resources := strings.Split(path, "/")[1:]
	controllerPath := "/" + resources[0]
	if len(resources) > 1 {
		subpath := strings.Join(resources[1:], "/")
		return r.routes[controllerPath], subpath
	} else {
		return r.routes[controllerPath], ""
	}
}

func (r *Routes) RegisterRoute(path string, controller ControllerInterface) {
	resources := strings.Split(path, "/")[1:]
	controllerPath := "/" + resources[0]
	r.routes[controllerPath] = controller
}
