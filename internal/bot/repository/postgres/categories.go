package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
)

func (pg *Pg) CreateCategory(ctx context.Context, category *entities.Category) error {
	query := "INSERT INTO categories (name, c_type) VALUES ($1, $2)"
	id, err := pg.db.ExecContextWithReturnID(ctx, query, category.Name, category.Type)

	if err != nil {
		return err
	}

	category.ID = uuid.MustParse(id.(string))
	return nil
}

func (pg *Pg) EditCategory(ctx context.Context, category *entities.Category) error {
	query := "UPDATE categories SET name = $1, c_type = $2 WHERE id = $3"
	_, err := pg.db.ExecContext(ctx, query, category.Name, category.Type, category.ID)

	return err
}

func (pg *Pg) FindCategory(ctx context.Context, name string, cType int) (*entities.Category, error) {
	var category entities.Category

	query := "SELECT * FROM categories WHERE lower(name) = lower($1) AND c_type = $2"
	err := pg.db.GetContext(ctx, &category, query, name, cType)

	return &category, err
}

func (pg *Pg) GetCategoryByID(ctx context.Context, id uuid.UUID) (*entities.Category, error) {
	var category entities.Category

	query := "SELECT * FROM categories WHERE id = $1"
	err := pg.db.GetContext(ctx, &category, query, id)

	return &category, err
}

func (pg *Pg) GetCategoriesByType(ctx context.Context, cType int) ([]entities.Category, error) {
	var categories []entities.Category

	query := "SELECT * FROM categories WHERE c_type = $1"
	err := pg.db.SelectContext(ctx, &categories, query, cType)

	if err != nil {
		return nil, err
	} else if len(categories) < 1 {
		return nil, sql.ErrNoRows
	}

	return categories, err
}

func (pg *Pg) GetCategories(ctx context.Context) ([]entities.Category, error) {
	var categories []entities.Category

	query := "SELECT * FROM categories"
	err := pg.db.SelectContext(ctx, &categories, query)

	if err != nil {
		return nil, err
	} else if len(categories) < 1 {
		return nil, sql.ErrNoRows
	}

	return categories, err
}

func (pg *Pg) DeleteCategory(ctx context.Context, categoryID uuid.UUID) error {
	query := "DELETE FROM categories WHERE id = $1"
	_, err := pg.db.ExecContext(ctx, query, categoryID)

	return err
}
