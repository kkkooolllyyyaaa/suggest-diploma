package stats

import (
	"cmp"
	"slices"
	"strings"

	"suggest-runtime/internal/category/stats/subtree"
	"suggest-runtime/internal/category/tree"
)

type DrilldownEngineParams struct {
	DrilldownThreshold float64
	MinQueryFreq       int64
	MinCategoryRate    float64
	MaxCategories      int
	StatsAccessor      CatStatsAccessor
}

type categoryTreeTraversalEngine struct {
	categoryTree          tree.CategoryTree
	queriesCategoriesDict QueriesCategoriesDict
	params                DrilldownEngineParams
	builder               subtree.QueryCategoryNodesSubtreeBuilder
}

func NewCategoryTreeTraversalEngine(
	categoryTree tree.CategoryTree,
	queryInfoProvider QueriesCategoriesDict,
	params DrilldownEngineParams,
	builder subtree.QueryCategoryNodesSubtreeBuilder,
) CatEngine {
	return &categoryTreeTraversalEngine{
		categoryTree:          categoryTree,
		queriesCategoriesDict: queryInfoProvider,
		params:                params,
		builder:               builder,
	}
}

func (c *categoryTreeTraversalEngine) Suggest(query string) []string {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil
	}

	// Для запроса нет статистики
	categoriesStats := c.queriesCategoriesDict[query]
	if len(categoriesStats) == 0 {
		return nil
	}

	// Запрос недостаточно частотный
	if c.params.StatsAccessor.QueryFreq(categoriesStats[0]) < c.params.MinQueryFreq {
		return nil
	}

	tmp := make([]CatStats, len(categoriesStats))
	copy(tmp, categoriesStats)
	queryInfmNodesSubtree := c.builder.Build(tmp)
	return c.drilldown(queryInfmNodesSubtree)
}

// drilldownDetails Сколько статистики мы потеряем (Loss) заменив Node на Children
type drilldownDetails struct {
	Loss                 float64
	Node                 *subtree.QueryStatsCategoryNode
	Children             []*subtree.QueryStatsCategoryNode
	ParentCategoryRate   float64
	ChildrenCategoryRate float64
}

func (c *categoryTreeTraversalEngine) drilldown(root *subtree.QueryStatsCategoryNode) []string {
	var result = []*subtree.QueryStatsCategoryNode{root}
	drilldownDetailsMap := map[string]drilldownDetails{}
	completeness := 1.0

	for {
		// В этом цикле пробуем задрилдаунить и считаем потерю статистики для каждой из категорий
		// Это нужно, чтобы выбрать лучшую категорию, задрилдаунив которую мы потеряем меньше всего статистики
		for _, drilldownCandidate := range result {
			if _, ok := drilldownDetailsMap[drilldownCandidate.Info.Category]; ok {
				continue
			}

			remainingSlots := min(len(drilldownCandidate.Children), c.params.MaxCategories-len(result)+1)
			if remainingSlots == 0 {
				continue
			}

			childrenCategoryRateSum := 0.0
			nodesCnt := 0

			for i := 0; i < remainingSlots; i++ {
				rate := c.params.StatsAccessor.CategoryRate(drilldownCandidate.Children[i].Info)
				if rate < c.params.MinCategoryRate {
					break
				}

				childrenCategoryRateSum += rate
				nodesCnt += 1
			}

			loss := c.params.StatsAccessor.CategoryRate(drilldownCandidate.Info) - childrenCategoryRateSum
			drilldownDetailsMap[drilldownCandidate.Info.Category] = drilldownDetails{
				Loss:                 loss,
				Node:                 drilldownCandidate,
				Children:             drilldownCandidate.Children[:nodesCnt],
				ParentCategoryRate:   c.params.StatsAccessor.CategoryRate(drilldownCandidate.Info),
				ChildrenCategoryRate: c.params.StatsAccessor.CategoryRate(drilldownCandidate.Info) - loss,
			}
		}

		// Проверяем, какие ноды можно задрилдаунить, выбираем вариант с наименьшей потерей
		bestCandidate, bestCandidateDetails := c.tryDrilldownCategory(result, drilldownDetailsMap)

		// Если есть подходящий кандидат на дрилдаун, заменяем ноду его детьми и пробуем еще в новой итерации
		if bestCandidate != nil && bestCandidateDetails != nil {
			result = c.replaceNodeToHisChildren(result, bestCandidate.Info.Category, bestCandidateDetails.Children)
			completeness -= bestCandidateDetails.Loss
		} else {
			break
		}
	}

	if len(result) == 0 || (len(result) == 1 && result[0].Info.Category == tree.RootCategoryId) {
		return nil
	}

	slices.SortFunc(result, func(a, b *subtree.QueryStatsCategoryNode) int {
		return cmp.Compare(c.params.StatsAccessor.CategoryRate(b.Info), c.params.StatsAccessor.CategoryRate(a.Info))
	})

	categoryIdsResult := make([]string, 0, len(result))
	for _, categoryInfo := range result {
		categoryIdsResult = append(categoryIdsResult, categoryInfo.Info.Category)
	}
	return categoryIdsResult
}

func (c *categoryTreeTraversalEngine) tryDrilldownCategory(
	currentNodes []*subtree.QueryStatsCategoryNode,
	drilldownDetailsMap map[string]drilldownDetails,
) (*subtree.QueryStatsCategoryNode, *drilldownDetails) {
	childrenRate := func(details drilldownDetails) float64 {
		return details.ChildrenCategoryRate / details.ParentCategoryRate
	}

	var bestCandidate *subtree.QueryStatsCategoryNode
	var bestCandidateDetails drilldownDetails

	for _, resultNode := range currentNodes {
		details, ok := drilldownDetailsMap[resultNode.Info.Category]
		if !ok || childrenRate(details) < c.params.DrilldownThreshold {
			continue
		}

		if (len(currentNodes) - 1 + len(details.Children)) > c.params.MaxCategories {
			continue
		}

		if bestCandidate == nil || childrenRate(details) >= childrenRate(bestCandidateDetails) {
			bestCandidate = resultNode
			bestCandidateDetails = drilldownDetailsMap[bestCandidate.Info.Category]
		}
	}

	return bestCandidate, &bestCandidateDetails
}

func (c *categoryTreeTraversalEngine) replaceNodeToHisChildren(
	source []*subtree.QueryStatsCategoryNode,
	categoryToReplace string,
	categoriesReplaceBy []*subtree.QueryStatsCategoryNode,
) []*subtree.QueryStatsCategoryNode {
	newResult := make([]*subtree.QueryStatsCategoryNode, 0, len(source)-1+len(categoriesReplaceBy))
	for _, node := range source {
		if node.Info.Category != categoryToReplace {
			newResult = append(newResult, node)
		}
	}
	newResult = append(newResult, categoriesReplaceBy...)
	return newResult
}
