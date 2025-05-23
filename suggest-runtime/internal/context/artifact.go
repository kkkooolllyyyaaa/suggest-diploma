package context

import (
	"suggest-runtime/internal/artifact"
	"suggest-runtime/internal/artifact/s3"
	categoryTree "suggest-runtime/internal/category/tree"
)

func (c *SuggestContext) initS3() {
	cfg := s3.Config{
		Endpoint:   c.Config.S3.Endpoint,
		AccessKey:  c.Config.S3.AccessKey,
		SecretKey:  c.Config.S3.SecretKey,
		BucketName: c.Config.S3.BucketName,
		UseSSL:     c.Config.S3.UseSSL,
	}
	minio, _ := s3.NewMinio(cfg)
	c.S3 = minio
}

func (c *SuggestContext) readRemote() {
	c.S3.DownloadFileToFile(
		c.Config.ArtifactRemote.Queries,
		c.Config.Artifact.Queries,
	)
	c.S3.DownloadFileToFile(
		c.Config.ArtifactRemote.QueriesCategories,
		c.Config.Artifact.QueriesCategories,
	)
	c.S3.DownloadFileToFile(
		c.Config.ArtifactRemote.Nodes,
		c.Config.Artifact.Nodes,
	)
}

func (c *SuggestContext) readVectorsRemote() {
	c.S3.DownloadFileToFile(
		c.Config.ArtifactRemote.QueriesVectors,
		c.Config.Artifact.QueriesVectors,
	)
	c.S3.DownloadFileToFile(
		c.Config.ArtifactRemote.TokensVectors,
		c.Config.Artifact.TokensVectors,
	)
	c.S3.DownloadFileToFile(
		c.Config.ArtifactRemote.AnnoyIndex,
		c.Config.Artifact.AnnoyIndex,
	)
}

func (c *SuggestContext) readQueries() {
	queries, _ := artifact.ReadQueriesFromJson(c.Config.Artifact.Queries)
	c.queries = queries
}

func (c *SuggestContext) readQueriesCategories() {
	queriesCategories, _ := artifact.ReadQueriesCategories(c.Config.Artifact.QueriesCategories)
	c.QueriesCategories = queriesCategories
}

func (c *SuggestContext) readCategoryTree() {
	nodes, _ := artifact.ReadNodesFromJson(c.Config.Artifact.Nodes)
	tree := categoryTree.NewCategoryTree(nodes)
	c.Tree = tree
}

func (c *SuggestContext) queriesVectors() {
	qv, _ := artifact.ReadQueriesVectors(c.Config.Artifact.QueriesVectors)
	c.QueriesVectors = qv
}

func (c *SuggestContext) tokensVectors() {
	tv, _ := artifact.ReadTokensVectors(c.Config.Artifact.TokensVectors)
	c.TokensVectors = tv
}
