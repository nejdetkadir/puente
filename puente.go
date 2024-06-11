package puente

import (
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"strings"
)

type (
	Context struct {
		routes  []Route
		onError func(err error) events.APIGatewayProxyResponse
	}
	Route struct {
		relatedGroup *Group
		path         string
		method       string
		handler      func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
	}
	Group struct {
		path    string
		context *Context
	}
	RouteBuilder interface {
		Get(path string, handler func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse)
		Post(path string, handler func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse)
		Put(path string, handler func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse)
		Patch(path string, handler func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse)
		Delete(path string, handler func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse)
		Group(path string) RouteBuilder
	}
	ContextType interface {
		RouteBuilder
		RegisteredRoutes() []Route
		ListenAPIGateway(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
		OnError(handler func(err error) events.APIGatewayProxyResponse)
		RouteMatcher(path string, method string) *Route
	}
)

const (
	HttpMethodGet    = "GET"
	HttpMethodPost   = "POST"
	HttpMethodPut    = "PUT"
	HttpMethodPatch  = "PATCH"
	HttpMethodDelete = "DELETE"
	ErrorNotFound    = "not_found"
)

func New() ContextType {
	return &Context{
		routes: make([]Route, 0),
		onError: func(err error) events.APIGatewayProxyResponse {
			var statusCode = 500

			if err.Error() == ErrorNotFound {
				statusCode = 404
			}

			return events.APIGatewayProxyResponse{
				StatusCode: statusCode,
				Body:       err.Error(),
			}
		},
	}
}

func (c *Context) Get(path string, handler func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse) {
	c.routes = append(c.routes, Route{path: path, method: HttpMethodGet, handler: handler})
}

func (c *Context) Post(path string, handler func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse) {
	c.routes = append(c.routes, Route{path: path, method: HttpMethodPost, handler: handler})
}

func (c *Context) Put(path string, handler func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse) {
	c.routes = append(c.routes, Route{path: path, method: HttpMethodPut, handler: handler})
}

func (c *Context) Patch(path string, handler func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse) {
	c.routes = append(c.routes, Route{path: path, method: HttpMethodPatch, handler: handler})
}

func (c *Context) Delete(path string, handler func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse) {
	c.routes = append(c.routes, Route{path: path, method: HttpMethodDelete, handler: handler})
}

func (c *Context) Group(path string) RouteBuilder {
	return &Group{
		path:    path,
		context: c,
	}
}

func (c *Context) RegisteredRoutes() []Route {
	return c.routes
}

func (c *Context) OnError(handler func(err error) events.APIGatewayProxyResponse) {
	c.onError = handler
}

func (c *Context) RouteMatcher(path string, method string) *Route {
	for _, route := range c.routes {
		if route.path == path && route.method == method {
			return &route
		}

		if route.method == method && matchDynamicRoute(route.path, path) {
			return &route
		}
	}

	return nil
}

func matchDynamicRoute(routePath, requestPath string) bool {
	routeParts := strings.Split(routePath, "/")
	requestParts := strings.Split(requestPath, "/")

	if len(routeParts) != len(requestParts) {
		return false
	}

	for i := range routeParts {
		if routeParts[i] != requestParts[i] && !strings.HasPrefix(routeParts[i], ":") {
			return false
		}
	}

	return true
}

func (c *Context) ListenAPIGateway(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var matchedRoute = c.RouteMatcher(request.Path, request.HTTPMethod)

	if matchedRoute == nil {
		return c.onError(errors.New(ErrorNotFound))
	}

	return matchedRoute.handler(request)
}

func (g *Group) Get(path string, handler func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse) {
	g.context.routes = append(g.context.routes, Route{path: g.path + path, method: HttpMethodGet, handler: handler})
}

func (g *Group) Post(path string, handler func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse) {
	g.context.routes = append(g.context.routes, Route{path: g.path + path, method: HttpMethodPost, handler: handler})
}

func (g *Group) Put(path string, handler func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse) {
	g.context.routes = append(g.context.routes, Route{path: g.path + path, method: HttpMethodPut, handler: handler})
}

func (g *Group) Patch(path string, handler func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse) {
	g.context.routes = append(g.context.routes, Route{path: g.path + path, method: HttpMethodPatch, handler: handler})
}

func (g *Group) Delete(path string, handler func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse) {
	g.context.routes = append(g.context.routes, Route{path: g.path + path, method: HttpMethodDelete, handler: handler})
}

func (g *Group) Group(path string) RouteBuilder {
	return &Group{
		path:    g.path + path,
		context: g.context,
	}
}
