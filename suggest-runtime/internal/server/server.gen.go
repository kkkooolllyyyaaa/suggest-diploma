// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by unknown module path version unknown version DO NOT EDIT.
package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

// CategoryTreeNode defines model for CategoryTreeNode.
type CategoryTreeNode struct {
	HasChildren bool   `json:"has_children"`
	Id          string `json:"id"`
	Name        string `json:"name"`
}

// CategoryTreeNodesResponse defines model for CategoryTreeNodesResponse.
type CategoryTreeNodesResponse = []CategoryTreeNode

// SearchRequest defines model for SearchRequest.
type SearchRequest struct {
	Query string `json:"query"`
}

// SuggestItem defines model for SuggestItem.
type SuggestItem struct {
	// CategoryId Category id for search engine
	CategoryId *string `json:"categoryId,omitempty"`

	// CategoryName Category name for search engine
	CategoryName *string `json:"categoryName,omitempty"`

	// LocationId Location id for search engine
	LocationId *string `json:"locationId,omitempty"`

	// Query Text query for search engine
	Query *string `json:"query,omitempty"`

	// Title Title visible for user
	Title string `json:"title"`
}

// SuggestResponse defines model for SuggestResponse.
type SuggestResponse struct {
	Items     []SuggestItem `json:"items"`
	SuggestId string        `json:"suggestId"`
}

// GetV1ApiCategoryTreeParams defines parameters for GetV1ApiCategoryTree.
type GetV1ApiCategoryTreeParams struct {
	Node   *string `form:"node,omitempty" json:"node,omitempty"`
	UserId string  `json:"userId"`
}

// PostV1ApiSearchParams defines parameters for PostV1ApiSearch.
type PostV1ApiSearchParams struct {
	UserId string `json:"userId"`
}

// PostV1ApiSuggestParams defines parameters for PostV1ApiSuggest.
type PostV1ApiSuggestParams struct {
	UserId string `json:"userId"`
}

// PostV1ApiSearchJSONRequestBody defines body for PostV1ApiSearch for application/json ContentType.
type PostV1ApiSearchJSONRequestBody = SearchRequest

// PostV1ApiSuggestJSONRequestBody defines body for PostV1ApiSuggest for application/json ContentType.
type PostV1ApiSuggestJSONRequestBody = SearchRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /v1/api/category/tree)
	GetV1ApiCategoryTree(ctx echo.Context, params GetV1ApiCategoryTreeParams) error

	// (POST /v1/api/search)
	PostV1ApiSearch(ctx echo.Context, params PostV1ApiSearchParams) error

	// (POST /v1/api/suggest)
	PostV1ApiSuggest(ctx echo.Context, params PostV1ApiSuggestParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetV1ApiCategoryTree converts echo context to params.
func (w *ServerInterfaceWrapper) GetV1ApiCategoryTree(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetV1ApiCategoryTreeParams
	// ------------- Optional query parameter "node" -------------

	err = runtime.BindQueryParameter("form", true, false, "node", ctx.QueryParams(), &params.Node)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter node: %s", err))
	}

	headers := ctx.Request().Header
	// ------------- Required header parameter "userId" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("userId")]; found {
		var UserId string
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for userId, got %d", n))
		}

		err = runtime.BindStyledParameterWithOptions("simple", "userId", valueList[0], &UserId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: true})
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userId: %s", err))
		}

		params.UserId = UserId
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Header parameter userId is required, but not found"))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetV1ApiCategoryTree(ctx, params)
	return err
}

// PostV1ApiSearch converts echo context to params.
func (w *ServerInterfaceWrapper) PostV1ApiSearch(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params PostV1ApiSearchParams

	headers := ctx.Request().Header
	// ------------- Required header parameter "userId" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("userId")]; found {
		var UserId string
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for userId, got %d", n))
		}

		err = runtime.BindStyledParameterWithOptions("simple", "userId", valueList[0], &UserId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: true})
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userId: %s", err))
		}

		params.UserId = UserId
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Header parameter userId is required, but not found"))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostV1ApiSearch(ctx, params)
	return err
}

// PostV1ApiSuggest converts echo context to params.
func (w *ServerInterfaceWrapper) PostV1ApiSuggest(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params PostV1ApiSuggestParams

	headers := ctx.Request().Header
	// ------------- Required header parameter "userId" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("userId")]; found {
		var UserId string
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for userId, got %d", n))
		}

		err = runtime.BindStyledParameterWithOptions("simple", "userId", valueList[0], &UserId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: true})
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userId: %s", err))
		}

		params.UserId = UserId
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Header parameter userId is required, but not found"))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostV1ApiSuggest(ctx, params)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/v1/api/category/tree", wrapper.GetV1ApiCategoryTree)
	router.POST(baseURL+"/v1/api/search", wrapper.PostV1ApiSearch)
	router.POST(baseURL+"/v1/api/suggest", wrapper.PostV1ApiSuggest)

}
