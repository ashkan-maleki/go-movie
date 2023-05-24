package memory

import (
	"context"
	"github.com/mamalmaleki/go_movie/metadata/internal/repository"
	"github.com/mamalmaleki/go_movie/metadata/pkg/model"
	"go.opentelemetry.io/otel"
	"sync"
)

const tracerID = "metadata-repository-memory"

// Repository defines a memory movie metadata repository.
type Repository struct {
	sync.RWMutex
	data map[string]*model.Metadata
}

// New creates a new memory repository.
func New() *Repository {
	return &Repository{data: map[string]*model.Metadata{}}
}

// Get retrieves movie metadata for by movie id.
func (r *Repository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	r.RLock()
	defer r.RUnlock()
	_, span := otel.Tracer(tracerID).Start(ctx, "Repository/Get")
	defer span.End()
	m, ok := r.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return m, nil
}

// Put adds movie metadata for a given movie id.
func (r *Repository) Put(ctx context.Context, id string, metadata *model.Metadata) error {
	r.Lock()
	defer r.Unlock()
	_, span := otel.Tracer(tracerID).Start(ctx, "Repository/Put")
	defer span.End()
	r.data[id] = metadata
	return nil
}
