package server

import (
	"encoding/json"
	"net/http"

	"suggest-runtime/internal/category/tree"
	"suggest-runtime/internal/config"
	"suggest-runtime/internal/context"
	"suggest-runtime/internal/suggester"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Server struct {
	suggestContext *context.SuggestContext
}

func NewServer(config *config.Config) ServerInterface {
	return &Server{
		suggestContext: context.InitContext(config),
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
	indexItems := h.suggestContext.Blender.Suggest(request)

	suggestItems := make([]SuggestItem, 0, len(indexItems))
	for _, item := range indexItems {
		toAdd := SuggestItem{
			Title: string(item.NormalizedQuery),
		}
		suggestItems = append(suggestItems, toAdd)
	}

	response := SuggestResponse{
		SuggestId: uuid.New().String(),
		Items:     suggestItems,
	}

	return ctx.JSONPretty(http.StatusOK, response, "\t")
}

func (h *Server) GetV1ApiCategoryTree(ctx echo.Context, params GetV1ApiCategoryTreeParams) error {
	root := tree.RootNodeId
	if params.Node != nil {
		root = *params.Node
	}

	children := h.suggestContext.Tree.GetChildren(root)
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

	err = h.suggestContext.QueryLogger.LogRequest(params.UserId, requestBody.Query)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return nil
}
