package mysql

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mamalmaleki/go-movie/metadata/internal/repository"
	"github.com/mamalmaleki/go-movie/metadata/pkg/model"
	"go.opentelemetry.io/otel"
)

const tracerID = "metadata-repository-mysql"

// Repository defines a MySQL-based movie metadata repository.
type Repository struct {
	db *sql.DB
}

// New creates a new MySQL-based repository.
func New() (*Repository, error) {
	dataSourceName := "root:mauFJcuf5dhRMQrjj@/movie"
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &Repository{db}, nil
}

func (r *Repository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	_, span := otel.Tracer(tracerID).Start(ctx, "Repository/Get")
	defer span.End()
	var title, description, director string
	row := r.db.QueryRowContext(ctx,
		"SELECT title, description, director FROM movies WHERE id = ?", id)

	if err := row.Scan(&title, &description, &director); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &model.Metadata{
		ID:          id,
		Title:       title,
		Description: description,
		Director:    director,
	}, nil
}

// Put adds movie metadata for a given movie id.
func (r *Repository) Put(ctx context.Context, id string, metadata *model.Metadata) error {
	_, span := otel.Tracer(tracerID).Start(ctx, "Repository/Put")
	defer span.End()
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO movies (id, title, description, director) VALUES (?, ?, ?, ?)",
		id, metadata.Title, metadata.Description, metadata.Director)
	return err
}
