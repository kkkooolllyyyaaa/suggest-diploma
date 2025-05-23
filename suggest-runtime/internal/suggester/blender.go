package suggester

type SuggestBlender struct {
	trieSuggester    Suggester
	historySuggester Suggester
	annSuggester     Suggester
}

func NewSuggestBlender(
	trieSuggester Suggester,
	historySuggester Suggester,
	annSuggester Suggester,
) SuggestBlender {
	return SuggestBlender{
		trieSuggester:    trieSuggester,
		historySuggester: historySuggester,
		annSuggester:     annSuggester,
	}
}

func (sb *SuggestBlender) Suggest(request SearchRequest) []*IndexItem {
	finalResult := sb.historySuggester.Suggest(request)
	if len(finalResult) >= SuggestLimit {
		return finalResult
	}
	uniqueQueries := make(map[string]bool, len(finalResult))
	for _, el := range finalResult {
		uniqueQueries[string(el.NormalizedQuery)] = true
	}

	trieResult := sb.trieSuggester.Suggest(request)
	for _, el := range trieResult {
		if _, ok := uniqueQueries[string(el.NormalizedQuery)]; ok {
			continue
		}

		finalResult = append(finalResult, el)
		if len(finalResult) >= SuggestLimit {
			break
		}
	}

	return finalResult
}
