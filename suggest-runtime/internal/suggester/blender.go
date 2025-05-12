package suggester

type SuggestBlender struct {
	TrieSuggester    Suggester
	HistorySuggester Suggester
}

func NewSuggestBlender(
	trieSuggester Suggester,
	historySuggester Suggester,
) SuggestBlender {
	return SuggestBlender{
		TrieSuggester:    trieSuggester,
		HistorySuggester: historySuggester,
	}
}

func (sb *SuggestBlender) Suggest(request SearchRequest) []*IndexItem {
	finalResult := sb.HistorySuggester.Suggest(request)
	if len(finalResult) >= SuggestLimit {
		return finalResult
	}

	trieResult := sb.TrieSuggester.Suggest(request)
	for _, el := range trieResult {
		finalResult = append(finalResult, el)
		if len(finalResult) >= SuggestLimit {
			break
		}
	}

	return finalResult
}
