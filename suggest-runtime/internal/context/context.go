package context

import (
	"fmt"
	"sync"
	"time"

	"suggest-runtime/internal/artifact"
	categoryTree "suggest-runtime/internal/category/tree"
	"suggest-runtime/internal/config"
	"suggest-runtime/internal/history"
	"suggest-runtime/internal/suggester"
	history2 "suggest-runtime/internal/suggester/history"
	"suggest-runtime/internal/suggester/radixtrie"
)

type SuggestContext struct {
	Config      *config.Config
	Tree        categoryTree.CategoryTree
	QueryLogger history.QueryLogger
	Blender     suggester.SuggestBlender

	indexItems       []*suggester.IndexItem
	trieSuggester    suggester.Suggester
	historySuggester suggester.Suggester
}

func InitContext(config *config.Config) *SuggestContext {
	sc := &SuggestContext{
		Config: config,
	}

	jobsPhases := []map[string]func(){
		{
			"index items":   sc.readIndexItems,
			"category tree": sc.categoryTree,
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
	sc.indexItems = nil

	return sc
}

func (c *SuggestContext) readIndexItems() {
	items, _ := artifact.ReadQueriesFromJson()
	c.indexItems = items
}

func (c *SuggestContext) trie() {
	trieSuggester := radixtrie.NewTrieSuggester()
	trieSuggester.Build(c.indexItems)
	c.trieSuggester = trieSuggester
}

func (c *SuggestContext) categoryTree() {
	nodes, _ := artifact.ReadNodesFromJson()
	tree := categoryTree.NewCategoryTree(nodes)
	c.Tree = tree
}

func (c *SuggestContext) history() {
	c.QueryLogger = history.NewQueryLogger(c.Config.Redis.Host)
	c.historySuggester = history2.NewHistorySuggester(c.QueryLogger)
}
