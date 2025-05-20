package subtree

import (
	"cmp"
	"slices"

	"suggest-runtime/internal/category/stats"
	"suggest-runtime/internal/category/tree"
)

type QueryStatsCategoryNode struct {
	Info     stats.CatStats
	Parent   *QueryStatsCategoryNode
	Children []*QueryStatsCategoryNode
}

type QueryCategoryNodesSubtreeBuilder struct {
	categoryTree         tree.CategoryTree
	statsAccessorForSort stats.CatStatsAccessor
}

func NewQueryCategoryNodesSubtreeBuilder(
	categoryTree tree.CategoryTree,
	statsAccessorForSort stats.CatStatsAccessor,
) QueryCategoryNodesSubtreeBuilder {
	return QueryCategoryNodesSubtreeBuilder{
		categoryTree:         categoryTree,
		statsAccessorForSort: statsAccessorForSort,
	}
}

func (b QueryCategoryNodesSubtreeBuilder) Build(nodesStats []stats.CatStats) *QueryStatsCategoryNode {
	if len(nodesStats) == 0 {
		return nil
	}

	root := &QueryStatsCategoryNode{
		Info: stats.CatStats{
			Category:            tree.RootCategoryId,
			Contacts:            nodesStats[0].Contacts,
			Searches:            nodesStats[0].Searches,
			Score:               nodesStats[0].Score,
			CategoryContactRate: 1.0,
			CategorySearchRate:  1.0,
			CategoryScoreRate:   1.0,
			CategoryContacts:    nodesStats[0].Contacts,
			CategorySearches:    nodesStats[0].Searches,
			CategoryScore:       nodesStats[0].Score,
		},
	}

	// Сортировка по глубине. Сначала root, затем вертикали, затем категории 2-го уровня и.т.п.
	// В рамках одной глубины сортировка идет по убыванию контактов и поисков
	slices.SortFunc(nodesStats, func(first, second stats.CatStats) int {
		depthDiff := b.categoryTree.Depth(first.Category) - b.categoryTree.Depth(second.Category)
		if depthDiff == 0 {
			return cmp.Compare(b.statsAccessorForSort.CategoryRate(second), b.statsAccessorForSort.CategoryRate(first))
		}
		return depthDiff
	})

	CategoryIdToNode := map[string]*QueryStatsCategoryNode{root.Info.Category: root}

	for _, nodeStat := range nodesStats {
		parentIdPointer := b.categoryTree.Parent(nodeStat.Category)
		parentId := tree.RootCategoryId
		if parentIdPointer != nil {
			parentId = *parentIdPointer
		}

		parentNode, ok := CategoryIdToNode[parentId]
		if !ok {
			continue
		}

		currentNode := &QueryStatsCategoryNode{
			Parent: parentNode,
			Info: stats.CatStats{
				Category:            nodeStat.Category,
				Contacts:            nodeStat.Contacts,
				Searches:            nodeStat.Searches,
				Score:               nodeStat.Score,
				CategoryContactRate: nodeStat.CategoryContactRate,
				CategorySearchRate:  nodeStat.CategorySearchRate,
				CategoryScoreRate:   nodeStat.CategoryScoreRate,
				CategoryContacts:    nodeStat.CategoryContacts,
				CategorySearches:    nodeStat.CategorySearches,
				CategoryScore:       nodeStat.CategorySearches,
			},
		}

		parentNode.Children = append(parentNode.Children, currentNode)
		CategoryIdToNode[nodeStat.Category] = currentNode
	}

	return root
}

func (b QueryCategoryNodesSubtreeBuilder) SortChildren(root *QueryStatsCategoryNode, statsAccessorForSort stats.CatStatsAccessor) *QueryStatsCategoryNode {
	var dfs func(node *QueryStatsCategoryNode)
	dfs = func(node *QueryStatsCategoryNode) {
		if len(node.Children) == 0 {
			return
		}

		slices.SortFunc(node.Children, func(first, second *QueryStatsCategoryNode) int {
			return cmp.Compare(statsAccessorForSort.CategoryRate(second.Info), statsAccessorForSort.CategoryRate(first.Info))
		})

		for _, nextNode := range node.Children {
			dfs(nextNode)
		}
	}

	return root
}
