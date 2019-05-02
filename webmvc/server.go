package webmvc

import (
	"fmt"
	"net/http"
	"strings"
	"webmvc/base"
)

// NewServer creates a new server instance and is the entry serving point
// of the framework
type NewServer struct {
	Routes *base.Routes
}

func buildQueries(rawQueries string) map[string]string {
	requestQueries := strings.Split(rawQueries, "&")

	queries := make(map[string]string)
	if requestQueries[0] != "" && len(requestQueries) > 0 {
		for _, query := range requestQueries {
			queryStreams := strings.Split(query, "=")
			for i := 0; i < len(queryStreams); i += 2 {
				queries[queryStreams[i]] = queryStreams[i+1]
			}
		}
	}
	return queries
}

func (s *NewServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestPath := r.URL.Path
	queries := buildQueries(r.URL.RawQuery)

	if s.Routes.Exists(requestPath) {
		controller, subpath := s.Routes.GetController(requestPath)
		var response *base.HttpResponse
		switch r.Method {
		case http.MethodPost:
			response = controller.Post(subpath, queries)
		case http.MethodPut:
			response = controller.Put(subpath, queries)
		case http.MethodDelete:
			response = controller.Delete(subpath, queries)
		default: // include http.MethodGet
			response = controller.Get(subpath, queries)
		}

		if response.StatusCode == 200 {
			fmt.Fprint(w, response.Body)
		} else {
			http.Error(w, response.Body, response.StatusCode)
		}
	} else {
		http.Error(w, "Not Found", 404)
	}
}

func CreateNewServer() *NewServer {
	routes := base.CreateNewRouter()
	server := &NewServer{routes}
	return server
}
