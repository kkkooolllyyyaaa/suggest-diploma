package context

import (
	"fmt"
	"sync"
	"time"

	"suggest-runtime/internal/artifact/s3"
	"suggest-runtime/internal/category/stats"
	categoryTree "suggest-runtime/internal/category/tree"
	"suggest-runtime/internal/config"
	historyLogger "suggest-runtime/internal/history"
	"suggest-runtime/internal/suggester"
	"suggest-runtime/internal/suggester/ann"
	"suggest-runtime/internal/suggester/history"
	"suggest-runtime/internal/suggester/radixtrie"
	"suggest-runtime/internal/vector"
)

type SuggestContext struct {
	Config            *config.Config
	Tree              categoryTree.CategoryTree
	QueryLogger       historyLogger.QueryLogger
	Blender           suggester.SuggestBlender
	QueriesCategories stats.QueriesCategoriesDict
	CatEngine         stats.CatEngine
	S3                *s3.Minio
	QueriesVectors    vector.QueriesVectors
	TokensVectors     vector.TokensVectors
	AnnIndex          vector.AnnIndex
	queries           []*suggester.IndexItem
	trieSuggester     suggester.Suggester
	historySuggester  suggester.Suggester
	annSuggester      suggester.Suggester
}

func InitContext(config *config.Config) *SuggestContext {
	sc := &SuggestContext{
		Config: config,
	}

	jobsPhases := []map[string]func(){
		//{
		//	"init s3": sc.initS3,
		//},
		//{
		//	"remote artifacts": sc.readRemote,
		//},
		{
			"queries":            sc.readQueries,
			"queries categories": sc.readQueriesCategories,
			"category tree":      sc.readCategoryTree,
		},
		{
			"suggesters": sc.suggesters,
		},
		{
			"category": sc.category,
		},
	}

	startAll := time.Now()
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
	fmt.Printf("Running all jobs took %.4f seconds\n", time.Now().Sub(startAll).Seconds())

	sc.Blender = suggester.NewSuggestBlender(
		sc.trieSuggester,
		sc.historySuggester,
		nil,
	)
	sc.queries = nil

	return sc
}

func (c *SuggestContext) suggesters() {
	trieSuggester := radixtrie.NewTrieSuggester()
	trieSuggester.Build(c.queries)
	c.trieSuggester = trieSuggester

	c.QueryLogger = historyLogger.NewQueryLogger(c.Config.Redis.Host)
	c.historySuggester = history.NewHistorySuggester(c.QueryLogger)
}

func (c *SuggestContext) category() {
	c.CatEngine = stats.NewCategoryEngine(
		c.QueriesCategories,
		stats.NewCategoryContactsAccessor(),
		c.Config.CategoryEngine.Threshold,
	)
}

func (c *SuggestContext) ann() {
	c.AnnIndex = vector.NewIndex(c.Config, c.QueriesVectors, c.TokensVectors)
	annSuggester := ann.NewAnnSuggester(c.AnnIndex)
	annSuggester.Build(c.queries)
	c.annSuggester = annSuggester
}
