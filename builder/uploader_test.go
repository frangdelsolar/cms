package builder_test

import (
	"net/http"
	"testing"

	"github.com/frangdelsolar/cms/builder"
	th "github.com/frangdelsolar/cms/builder/test_helpers"
	"github.com/stretchr/testify/assert"
)

const testFilePath = "test_data/test.json"

func TestUploaderGetsCreated(t *testing.T) {
	e, err := th.GetDefaultEngine()
	assert.NoError(t, err, "GetDefaultEngine should not return an error")

	t.Log("Testing Upload App is registered")
	app, err := e.Admin.GetApp("upload")
	assert.NoError(t, err, "GetApp should not return an error")
	assert.NotNil(t, app, "GetApp should return a non-nil App")

	handler := func(w http.ResponseWriter, r *http.Request) {}

	t.Log("Testing Upload routes are registered")
	expectedRoutes := []builder.RouteHandler{
		builder.NewRouteHandler("/file", handler, "file-new", true),
		builder.NewRouteHandler("/file/{id}/delete", handler, "file-delete", true),
		builder.NewRouteHandler("/static/{path:.*}", handler, "file-static", true),
	}

	routes := e.Server.GetRoutes()
	for _, expectedRoute := range expectedRoutes {
		found := false
		for _, route := range routes {
			if route.Route == expectedRoute.Route {
				assert.Equal(t, expectedRoute.Name, route.Name, "Route name should be the same")
				assert.Equal(t, expectedRoute.RequiresAuth, route.RequiresAuth, "Route requires auth should be the same")
				found = true
			}
		}

		assert.True(t, found, "Expected route not found: %s", expectedRoute.Route)
	}
}

func TestAnonymousCanUploadAllowed(t *testing.T) {
	e, err := th.GetDefaultEngine()
	assert.NoError(t, err, "GetDefaultEngine should not return an error")

	// Create a helper request to get the detail
	request, _, _ := th.NewRequestWithFile(
		http.MethodPost,
		"",
		testFilePath,
		true,
		nil,
		nil,
	)

	var result builder.Upload

	response, err := th.ExecuteApiCall(
		t,
		e.Engine.GetUploadPostHandler(e.Engine.Config.UploaderConfig),
		request,
		&result,
	)

	t.Log(response.Message)
	assert.NoError(t, err, "ApiDetail should not return an error")
	assert.NotNil(t, response, "ApiDetail should return a non-nil response")
	assert.Equal(t, response.Success, true, "Success should be true")

	assert.NotNil(t, result.FileName, "FileName should be something", result.FileName)
	assert.NotNil(t, result.FilePath, "FilePath should be something", result.FilePath)
	assert.NotNil(t, result.Url, "Url should be something", result.Url)

	// clean up
	builder.DeleteFile(result.FilePath)
}
func TestAnonymousCanNotUploadForbidden(t *testing.T) {
	e, err := th.GetDefaultEngine()
	assert.NoError(t, err, "GetDefaultEngine should not return an error")

	// Create a helper request to get the detail
	request, _, _ := th.NewRequestWithFile(
		http.MethodPost,
		"",
		testFilePath,
		false,
		nil,
		nil,
	)

	var result builder.Upload

	response, err := th.ExecuteApiCall(
		t,
		e.Engine.GetUploadPostHandler(e.Engine.Config.UploaderConfig),
		request,
		&result,
	)

	assert.NoError(t, err, "Upload should not return an error")
	assert.Equal(t, response.Success, false, "Success should be false")
	assert.Contains(t, response.Message, "user not authenticated", "Error should be user not authenticated")

	assert.Equal(t, result, (builder.Upload{}), "Result should be nil", result)
}

func TestAnonymousCanAccessAllowed(t *testing.T)      {}
func TestAnonymousCanNotAccessForbidden(t *testing.T) {}

func TestAnonymousCanDeleteAllowed(t *testing.T)      {}
func TestAnonymousCanNotDeleteForbidden(t *testing.T) {}

func TestMaxSize(t *testing.T) {}

func TestSupportedMediaType(t *testing.T) {}