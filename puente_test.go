package puente

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	c := New()

	assert.NotNil(t, c)
	assert.Len(t, c.RegisteredRoutes(), 0)
}

func TestContext_Get(t *testing.T) {
	c := New()

	c.Get("/path", func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
		return events.APIGatewayProxyResponse{}
	})

	assert.Len(t, c.RegisteredRoutes(), 1)
	assert.Equal(t, "/path", c.RegisteredRoutes()[0].path)
	assert.Equal(t, HttpMethodGet, c.RegisteredRoutes()[0].method)
}

func TestContext_Post(t *testing.T) {
	c := New()

	c.Post("/path", func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
		return events.APIGatewayProxyResponse{}
	})

	assert.Len(t, c.RegisteredRoutes(), 1)
	assert.Equal(t, "/path", c.RegisteredRoutes()[0].path)
	assert.Equal(t, HttpMethodPost, c.RegisteredRoutes()[0].method)
}

func TestContext_Put(t *testing.T) {
	c := New()

	c.Put("/path", func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
		return events.APIGatewayProxyResponse{}
	})

	assert.Len(t, c.RegisteredRoutes(), 1)
	assert.Equal(t, "/path", c.RegisteredRoutes()[0].path)
	assert.Equal(t, HttpMethodPut, c.RegisteredRoutes()[0].method)
}

func TestContext_Patch(t *testing.T) {
	c := New()

	c.Patch("/path", func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
		return events.APIGatewayProxyResponse{}
	})

	assert.Len(t, c.RegisteredRoutes(), 1)
	assert.Equal(t, "/path", c.RegisteredRoutes()[0].path)
	assert.Equal(t, HttpMethodPatch, c.RegisteredRoutes()[0].method)
}

func TestContext_Delete(t *testing.T) {
	c := New()

	c.Delete("/path", func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
		return events.APIGatewayProxyResponse{}
	})

	assert.Len(t, c.RegisteredRoutes(), 1)
	assert.Equal(t, "/path", c.RegisteredRoutes()[0].path)
	assert.Equal(t, HttpMethodDelete, c.RegisteredRoutes()[0].method)
}

func TestContext_Group(t *testing.T) {
	c := New()

	group := c.Group("/group")
	group.Get("/path", func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
		return events.APIGatewayProxyResponse{}
	})

	assert.Len(t, c.RegisteredRoutes(), 1)
	assert.Equal(t, "/group/path", c.RegisteredRoutes()[0].path)
	assert.Equal(t, HttpMethodGet, c.RegisteredRoutes()[0].method)
}

func TestContext_ListenAPIGateway(t *testing.T) {
	c := New()

	c.Get("/path", func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
		return events.APIGatewayProxyResponse{
			Body:       "test response",
			StatusCode: 200,
		}
	})

	response := c.ListenAPIGateway(events.APIGatewayProxyRequest{
		Path:       "/path",
		HTTPMethod: HttpMethodGet,
	})

	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, "test response", response.Body)
}

func TestContext_ListenAPIGateway_NotFound(t *testing.T) {
	c := New()

	response := c.ListenAPIGateway(events.APIGatewayProxyRequest{})

	assert.Equal(t, 404, response.StatusCode)
	assert.Equal(t, ErrorNotFound, response.Body)
}

func TestContext_OnError(t *testing.T) {
	c := New()

	c.OnError(func(err error) events.APIGatewayProxyResponse {
		return events.APIGatewayProxyResponse{
			Body:       "custom error responder",
			StatusCode: 500,
		}
	})

	response := c.ListenAPIGateway(events.APIGatewayProxyRequest{})

	assert.Equal(t, 500, response.StatusCode)
	assert.Equal(t, "custom error responder", response.Body)
}

func TestContext_RouteMatcher(t *testing.T) {
	c := New()

	c.Get("/path", func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
		return events.APIGatewayProxyResponse{}
	})

	route := c.RouteMatcher("/path", HttpMethodGet)

	assert.NotNil(t, route)
	assert.Equal(t, "/path", route.path)
	assert.Equal(t, HttpMethodGet, route.method)
}

func TestContext_RouteMatcher_NotFound(t *testing.T) {
	c := New()

	route := c.RouteMatcher("/path", HttpMethodGet)

	assert.Nil(t, route)
}

func TestContext_RouteMatcher_With_Parameters(t *testing.T) {
	c := New()

	c.Get("/path/:id", func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
		return events.APIGatewayProxyResponse{}
	})

	route := c.RouteMatcher("/path/123", HttpMethodGet)

	assert.NotNil(t, route)
	assert.Equal(t, "/path/:id", route.path)
	assert.Equal(t, HttpMethodGet, route.method)
}

func TestContext_RouteMatcher_With_Parameters_NotFound(t *testing.T) {
	c := New()

	c.Get("/path/:id", func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
		return events.APIGatewayProxyResponse{}
	})

	route := c.RouteMatcher("/path", HttpMethodGet)

	assert.Nil(t, route)
}

func TestGroup_Get(t *testing.T) {
	c := New()

	group := c.Group("/group")
	group.Get("/path", func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
		return events.APIGatewayProxyResponse{}
	})

	assert.Len(t, c.RegisteredRoutes(), 1)
	assert.Equal(t, "/group/path", c.RegisteredRoutes()[0].path)
	assert.Equal(t, HttpMethodGet, c.RegisteredRoutes()[0].method)
}

func TestGroup_Post(t *testing.T) {
	c := New()

	group := c.Group("/group")
	group.Post("/path", func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
		return events.APIGatewayProxyResponse{}
	})

	assert.Len(t, c.RegisteredRoutes(), 1)
	assert.Equal(t, "/group/path", c.RegisteredRoutes()[0].path)
	assert.Equal(t, HttpMethodPost, c.RegisteredRoutes()[0].method)
}

func TestGroup_Put(t *testing.T) {
	c := New()

	group := c.Group("/group")
	group.Put("/path", func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
		return events.APIGatewayProxyResponse{}
	})

	assert.Len(t, c.RegisteredRoutes(), 1)
	assert.Equal(t, "/group/path", c.RegisteredRoutes()[0].path)
	assert.Equal(t, HttpMethodPut, c.RegisteredRoutes()[0].method)
}

func TestGroup_Patch(t *testing.T) {
	c := New()

	group := c.Group("/group")
	group.Patch("/path", func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
		return events.APIGatewayProxyResponse{}
	})

	assert.Len(t, c.RegisteredRoutes(), 1)
	assert.Equal(t, "/group/path", c.RegisteredRoutes()[0].path)
	assert.Equal(t, HttpMethodPatch, c.RegisteredRoutes()[0].method)
}

func TestGroup_Delete(t *testing.T) {
	c := New()

	group := c.Group("/group")
	group.Delete("/path", func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
		return events.APIGatewayProxyResponse{}
	})

	assert.Len(t, c.RegisteredRoutes(), 1)
	assert.Equal(t, "/group/path", c.RegisteredRoutes()[0].path)
	assert.Equal(t, HttpMethodDelete, c.RegisteredRoutes()[0].method)
}

func TestGroup_Group(t *testing.T) {
	c := New()

	group := c.Group("/group")
	subGroup := group.Group("/sub-group")
	subGroup.Get("/path", func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
		return events.APIGatewayProxyResponse{}
	})

	assert.Len(t, c.RegisteredRoutes(), 1)
	assert.Equal(t, "/group/sub-group/path", c.RegisteredRoutes()[0].path)
	assert.Equal(t, HttpMethodGet, c.RegisteredRoutes()[0].method)
}

func TestGroup_ListenAPIGateway(t *testing.T) {
	c := New()

	group := c.Group("/group")
	group.Get("/path", func(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
		return events.APIGatewayProxyResponse{
			Body:       "test response",
			StatusCode: 200,
		}
	})

	response := c.ListenAPIGateway(events.APIGatewayProxyRequest{
		Path:       "/group/path",
		HTTPMethod: HttpMethodGet,
	})

	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, "test response", response.Body)
}
