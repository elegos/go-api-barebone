package types

import (
	"net/http"
	"strings"
)

// Action the http handler
type Action http.HandlerFunc

// Route the route definition
type Route struct {
	SubRoutes map[string]Route // eventual sub-paths

	// RFC 7231
	Get     Action
	Head    Action
	Post    Action
	Put     Action
	Delete  Action
	Connect Action
	Options Action
	Trace   Action
	// RFC 5789
	Patch Action
}

// GetAction get the action for the given method/path
func (r *Route) GetAction(method string, path string) Action {
	splitPath := strings.Split(strings.Trim(path, "/"), "/")

	if len(splitPath) > 1 {
		if subRoute, ok := r.SubRoutes[splitPath[0]]; ok {
			return subRoute.GetAction(method, strings.Join(splitPath[1:], "/"))
		}

		return nil
	}

	subRoute, found := r.SubRoutes[splitPath[0]]
	if !found {
		return nil
	}

	switch strings.ToLower(method) {
	case "get":
		return subRoute.Get
	case "head":
		return subRoute.Head
	case "post":
		return subRoute.Post
	case "put":
		return subRoute.Put
	case "delete":
		return subRoute.Delete
	case "connect":
		return subRoute.Connect
	case "options":
		return subRoute.Options
	case "trace":
		return subRoute.Trace
	case "patch":
		return subRoute.Patch
	default:
		break
	}

	return nil
}
