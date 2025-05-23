package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"suggest-runtime/internal/category/tree"
	"suggest-runtime/internal/config"
	"suggest-runtime/internal/context"
	"suggest-runtime/internal/suggester"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Server struct {
	context *context.SuggestContext
}

func NewServer(config *config.Config) ServerInterface {
	return &Server{
		context: context.InitContext(config),
	}
}

func (h *Server) PostV1ApiSuggest(ctx echo.Context, params PostV1ApiSuggestParams) error {
	var requestBody SearchRequest
	err := json.NewDecoder(ctx.Request().Body).Decode(&requestBody)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	request := suggester.SearchRequest{
		Query:  requestBody.Query,
		UserId: params.UserId,
	}
	indexItems := h.context.Blender.Suggest(request)

	suggestItems := make([]SuggestItem, 0, len(indexItems))
	for _, item := range indexItems {
		normalizedQ := string(item.NormalizedQuery)
		category := h.context.CatEngine.Suggest(normalizedQ)
		var categoryName *string
		if category != nil {
			categoryNameString := h.context.Tree.Title(*category)
			categoryName = &categoryNameString
		}
		toAdd := SuggestItem{
			Title:        normalizedQ,
			CategoryId:   category,
			CategoryName: categoryName,
		}
		suggestItems = append(suggestItems, toAdd)
	}

	response := SuggestResponse{
		SuggestId: uuid.New().String(),
		Items:     suggestItems,
	}
	fmt.Println(requestBody.Query)

	return ctx.JSONPretty(http.StatusOK, response, "\t")
}

func (h *Server) GetV1ApiCategoryTree(ctx echo.Context, params GetV1ApiCategoryTreeParams) error {
	root := tree.RootCategoryId
	if params.Node != nil {
		root = *params.Node
	}

	children := h.context.Tree.Children(root)
	treeNodes := make([]CategoryTreeNode, 0, len(children))
	for _, node := range children {
		treeNodes = append(treeNodes, CategoryTreeNode{
			Id:          node.Id,
			Name:        node.Title,
			HasChildren: node.HasChildren,
		})
	}
	return ctx.JSONPretty(http.StatusOK, treeNodes, "\t")
}

func (h *Server) PostV1ApiSearch(ctx echo.Context, params PostV1ApiSearchParams) error {
	var requestBody SearchRequest
	err := json.NewDecoder(ctx.Request().Body).Decode(&requestBody)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	if len(params.UserId) == 0 {
		return ctx.String(http.StatusBadRequest, "Empty userId")
	}

	err = h.context.QueryLogger.LogRequest(params.UserId, requestBody.Query)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return nil
}
