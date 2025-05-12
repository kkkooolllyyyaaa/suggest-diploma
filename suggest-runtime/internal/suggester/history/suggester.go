package history

import (
	"fmt"
	"strings"

	"suggest-runtime/internal/history"
	"suggest-runtime/internal/suggester"
)

type userHistory struct {
	logger history.QueryLogger
}

func NewHistorySuggester(logger history.QueryLogger) suggester.Suggester {
	return &userHistory{
		logger: logger,
	}
}

func (u userHistory) Build(collection []*suggester.IndexItem) {
	return
}

func (u userHistory) Suggest(request suggester.SearchRequest) []*suggester.IndexItem {
	userQueries, err := u.logger.GetUserRequests(request.UserId)
	if err != nil || len(userQueries) == 0 {
		fmt.Println(fmt.Sprintf("Got error fetching user history: %v", err))
		return nil
	}

	return u.filterByUserQuery(request, userQueries)
}

func (u userHistory) filterByUserQuery(request suggester.SearchRequest, userQueries []history.QueryTimestamp) []*suggester.IndexItem {
	result := make([]*suggester.IndexItem, 0, suggester.SuggestLimit)
	requestTokens := strings.Fields(request.Query)
	for _, uq := range userQueries {
		uqTokens := strings.Fields(uq.Query)
		if u.containsAll(requestTokens, uqTokens) {
			queryRunes := []rune(uq.Query)
			toAppend := &suggester.IndexItem{
				Query:           queryRunes,
				NormalizedQuery: queryRunes,
				Score:           0.0,
			}
			result = append(result, toAppend)

			if len(result) == suggester.SuggestLimit {
				break
			}
		}
	}
	return result
}

func (u userHistory) containsAll(userQueryTokens, historyQueryTokens []string) bool {
	if len(userQueryTokens) == 0 {
		return true
	}

	mp := make(map[string]struct{}, len(historyQueryTokens))
	for _, hqt := range historyQueryTokens {
		mp[hqt] = struct{}{}
	}

	containsAll := true
	for _, uqt := range userQueryTokens {
		if _, ok := mp[uqt]; !ok {
			containsAll = false
			break
		}
	}
	return containsAll
}
