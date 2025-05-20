package context

import (
	"fmt"
	"sync"
	"time"

	"suggest-runtime/internal/artifact"
	"suggest-runtime/internal/category/stats"
	categoryTree "suggest-runtime/internal/category/tree"
	"suggest-runtime/internal/config"
	historyLogger "suggest-runtime/internal/history"
	"suggest-runtime/internal/suggester"
	"suggest-runtime/internal/suggester/history"
	"suggest-runtime/internal/suggester/radixtrie"
)

type SuggestContext struct {
	Config            *config.Config
	Tree              categoryTree.CategoryTree
	QueryLogger       historyLogger.QueryLogger
	Blender           suggester.SuggestBlender
	QueriesCategories stats.QueriesCategoriesDict
	queries           []*suggester.IndexItem
	trieSuggester     suggester.Suggester
	historySuggester  suggester.Suggester
}

func InitContext(config *config.Config) *SuggestContext {
	sc := &SuggestContext{
		Config: config,
	}

	jobsPhases := []map[string]func(){
		{
			"queries":            sc.readQueries,
			"queries categories": sc.readQueriesCategories,
			"category tree":      sc.categoryTree,
		},
		{
			"trie": sc.trie,
		},
		{
			"history": sc.history,
		},
	}

	for _, jobs := range jobsPhases {
		wg := sync.WaitGroup{}
		wg.Add(len(jobs))

		for name, job := range jobs {
			job := job
			name := name

			go func() {
				fmt.Printf("Starting job `%s`\n", name)
				start := time.Now()
				job()
				fmt.Printf("Running job `%s` took %.4f seconds\n", name, time.Now().Sub(start).Seconds())
				wg.Done()
			}()
		}

		wg.Wait()
	}

	sc.Blender = suggester.NewSuggestBlender(
		sc.trieSuggester,
		sc.historySuggester,
	)
	sc.queries = nil

	return sc
}

func (c *SuggestContext) readQueries() {
	queries, _ := artifact.ReadQueriesFromJson(c.Config.Artifact.Queries)
	c.queries = queries
}

func (c *SuggestContext) readQueriesCategories() {
	queriesCategories, _ := artifact.ReadQueriesCategories(c.Config.Artifact.QueriesCategories)
	c.QueriesCategories = queriesCategories
}

func (c *SuggestContext) trie() {
	trieSuggester := radixtrie.NewTrieSuggester()
	trieSuggester.Build(c.queries)
	c.trieSuggester = trieSuggester
}

func (c *SuggestContext) categoryTree() {
	nodes, _ := artifact.ReadNodesFromJson(c.Config.Artifact.Nodes)
	tree := categoryTree.NewCategoryTree(nodes)
	c.Tree = tree
}

func (c *SuggestContext) history() {
	c.QueryLogger = historyLogger.NewQueryLogger(c.Config.Redis.Host)
	c.historySuggester = history.NewHistorySuggester(c.QueryLogger)
}
