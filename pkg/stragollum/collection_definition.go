package stragollum

import "fmt"

// CollectionDefinition represents the configuration for creating a collection.
type CollectionDefinition struct {
	DefaultID *CollectionDefaultIDOptions `json:"defaultId,omitempty"`
	Indexing  map[string]any              `json:"indexing,omitempty"`
	Lexical   *CollectionLexicalOptions   `json:"lexical,omitempty"`
	Rerank    *CollectionRerankOptions    `json:"rerank,omitempty"`
	Vector    *CollectionVectorOptions    `json:"vector,omitempty"`
}

type CollectionDefaultIDOptions struct {
	Type string `json:"type,omitempty"`
}

type CollectionLexicalOptions struct {
	Analyzer string `json:"analyzer,omitempty"`
	Enabled  *bool  `json:"enabled,omitempty"`
}

type CollectionRerankOptions struct {
	Enabled *bool                 `json:"enabled,omitempty"`
	Service *RerankServiceOptions `json:"service,omitempty"`
}

type CollectionVectorOptions struct {
	Dimension   *int                  `json:"dimension,omitempty"`
	Metric      string                `json:"metric,omitempty"`
	Service     *VectorServiceOptions `json:"service,omitempty"`
	SourceModel string                `json:"sourceModel,omitempty"`
}

type RerankServiceOptions struct {
	Authentication map[string]any `json:"authentication,omitempty"`
	ModelName      string         `json:"modelName,omitempty"`
	Parameters     map[string]any `json:"parameters,omitempty"`
	Provider       string         `json:"provider,omitempty"`
}

type VectorServiceOptions struct {
	Authentication map[string]any `json:"authentication,omitempty"`
	ModelName      string         `json:"modelName,omitempty"`
	Parameters     map[string]any `json:"parameters,omitempty"`
	Provider       string         `json:"provider,omitempty"`
}

// NewCollectionDefinition creates an empty CollectionDefinition.
func NewCollectionDefinition() *CollectionDefinition {
	return &CollectionDefinition{}
}

// WithDefaultID sets the defaultId configuration.
func (cd *CollectionDefinition) WithDefaultID(idType string) *CollectionDefinition {
	cd.DefaultID = &CollectionDefaultIDOptions{Type: idType}
	return cd
}

// WithIndexing sets the indexing configuration.
func (cd *CollectionDefinition) WithIndexing(indexing map[string]any) *CollectionDefinition {
	cd.Indexing = indexing
	return cd
}

// WithLexical sets the lexical configuration. Defaults enabled to true if not specified.
func (cd *CollectionDefinition) WithLexical(analyzer string, enabled ...bool) *CollectionDefinition {
	val := true
	if len(enabled) > 0 {
		val = enabled[0]
	}
	cd.Lexical = &CollectionLexicalOptions{
		Analyzer: analyzer,
		Enabled:  &val,
	}
	return cd
}

// WithRerank sets the rerank configuration. Defaults enabled to true if not specified.
func (cd *CollectionDefinition) WithRerank(service *RerankServiceOptions, enabled ...bool) *CollectionDefinition {
	val := true
	if len(enabled) > 0 {
		val = enabled[0]
	}
	cd.Rerank = &CollectionRerankOptions{
		Enabled: &val,
		Service: service,
	}
	return cd
}

// WithVector sets the vector configuration.
func (cd *CollectionDefinition) WithVector(vec *CollectionVectorOptions) *CollectionDefinition {
	cd.Vector = vec
	return cd
}

// WithVectorDimension sets the vector dimension (optional).
func (cd *CollectionDefinition) WithVectorDimension(dimension int) *CollectionDefinition {
	if cd.Vector == nil {
		cd.Vector = &CollectionVectorOptions{}
	}
	cd.Vector.Dimension = &dimension
	return cd
}

// WithVectorMetric sets the vector similarity metric.
func (cd *CollectionDefinition) WithVectorMetric(metric string) *CollectionDefinition {
	if cd.Vector == nil {
		cd.Vector = &CollectionVectorOptions{}
	}
	cd.Vector.Metric = metric
	return cd
}

// WithVectorSourceModel sets the vector source model.
func (cd *CollectionDefinition) WithVectorSourceModel(sourceModel string) *CollectionDefinition {
	if cd.Vector == nil {
		cd.Vector = &CollectionVectorOptions{}
	}
	cd.Vector.SourceModel = sourceModel
	return cd
}

// WithVectorService sets the vector service options.
func (cd *CollectionDefinition) WithVectorService(service *VectorServiceOptions) *CollectionDefinition {
	if cd.Vector == nil {
		cd.Vector = &CollectionVectorOptions{}
	}
	cd.Vector.Service = service
	return cd
}

// Validate ensures the CollectionDefinition has valid configuration.
func (cd *CollectionDefinition) Validate() error {
	if cd.Vector != nil && cd.Vector.Dimension != nil && *cd.Vector.Dimension <= 0 {
		return fmt.Errorf("vector dimension must be positive")
	}
	return nil
}
