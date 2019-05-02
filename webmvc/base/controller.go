package base

import ()

// HttpResponse defines the response data structure
type HttpResponse struct {
	Body       string
	StatusCode int
}

// ControllerInterface defines the methods allowed for each controller
type ControllerInterface interface {
	Get(subpath string, queries map[string]string) *HttpResponse
	Post(subpath string, queries map[string]string) *HttpResponse
	Put(subpath string, queries map[string]string) *HttpResponse
	Delete(subpath string, queries map[string]string) *HttpResponse
}

// Controller is the base class that implements the ControllerInterface
// and acts as the template for all inherited sub-controllers
type Controller struct{}

func (c *Controller) Get(subpath string, queries map[string]string) *HttpResponse {
	return &HttpResponse{
		Body:       "GET Method Not Allowed",
		StatusCode: 405,
	}
}

func (c *Controller) Post(subpath string, queries map[string]string) *HttpResponse {
	return &HttpResponse{
		Body:       "POST Method Not Allowed",
		StatusCode: 405,
	}
}

func (c *Controller) Put(subpath string, queries map[string]string) *HttpResponse {
	return &HttpResponse{
		Body:       "PUT Method Not Allowed",
		StatusCode: 405,
	}
}

func (c *Controller) Delete(subpath string, queries map[string]string) *HttpResponse {
	return &HttpResponse{
		Body:       "DELETE Method Not Allowed",
		StatusCode: 405,
	}
}
