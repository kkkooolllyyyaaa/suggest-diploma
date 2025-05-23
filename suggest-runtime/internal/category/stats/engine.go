package stats

import (
	"cmp"
	"slices"
	"strings"
)

type categoryEngine struct {
	queriesCategoriesDict QueriesCategoriesDict
	accessor              CatStatsAccessor
	threshold             float64
}

func NewCategoryEngine(
	queriesCategoriesDict QueriesCategoriesDict,
	accessor CatStatsAccessor,
	threshold float64,
) CatEngine {
	return &categoryEngine{
		queriesCategoriesDict: queriesCategoriesDict,
		accessor:              accessor,
		threshold:             threshold,
	}
}

func (c *categoryEngine) Suggest(query string) *string {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil
	}

	categoriesStats := c.queriesCategoriesDict[query]
	if len(categoriesStats) == 0 {
		return nil
	}

	slices.SortStableFunc(categoriesStats, func(a, b CatStats) int {
		return cmp.Compare(c.accessor.CategoryRate(b), c.accessor.CategoryRate(a))
	})

	var category *string = nil
	for _, stat := range categoriesStats {
		if c.accessor.CategoryRate(stat) >= c.threshold {
			category = &stat.Category
		} else {
			break
		}
	}

	return category
}
