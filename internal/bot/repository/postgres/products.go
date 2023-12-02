package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
)

func (pg *Pg) CreateProduct(ctx context.Context, product *entities.Product) error {
	query := "INSERT INTO products (category_id, producer, model, operating_system, additional, description, photo, price, old_price, is_sale) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
	id, err := pg.db.ExecContextWithReturnID(ctx, query, product.CategoryID, product.Producer, product.Model, product.OperatingSystem, product.Additional, product.Description, product.Photo, product.Price, product.OldPrice, product.IsSale)

	if err != nil {
		return err
	}

	product.ID = uuid.MustParse(id.(string))
	return nil
}

func (pg *Pg) EditProduct(ctx context.Context, product *entities.Product) error {
	query := "UPDATE products SET category_id = $1, producer = $2, model = $3, operating_system = $4, additional = $5, description = $6, photo = $7, price = $8, old_price = $9, is_sale = $10 WHERE id = $11"
	_, err := pg.db.ExecContext(ctx, query, product.CategoryID, product.Producer, product.Model, product.OperatingSystem, product.Additional, product.Description, product.Photo, product.Price, product.OldPrice, product.IsSale, product.ID)

	return err
}

func (pg *Pg) GetSaleProducts(ctx context.Context) ([]entities.Product, error) {
	var products []entities.Product

	query := "SELECT * FROM products WHERE products.is_sale = true"
	err := pg.db.SelectContext(ctx, &products, query)

	if err != nil {
		return nil, err
	} else if len(products) < 1 {
		return nil, sql.ErrNoRows
	}

	return products, err
}

func (pg *Pg) GetProducersByCategory(ctx context.Context, category string, cType int) ([]string, error) {
	var producers []string

	query := "SELECT DISTINCT(products.producer) FROM categories JOIN products ON categories.id = products.category_id WHERE categories.name = $1 AND categories.c_type = $2"
	err := pg.db.SelectContext(ctx, &producers, query, category, cType)

	if err != nil {
		return nil, err
	} else if len(producers) < 1 {
		return nil, sql.ErrNoRows
	}

	return producers, err
}

func (pg *Pg) GetProductsByProducer(ctx context.Context, producer string, cType int) ([]entities.Product, error) {
	var products []entities.Product

	query := "SELECT products.* FROM categories JOIN products ON products.producer = $1 AND products.category_id = categories.id WHERE categories.c_type = $2"
	err := pg.db.SelectContext(ctx, &products, query, producer, cType)

	if err != nil {
		return nil, err
	} else if len(products) < 1 {
		return nil, sql.ErrNoRows
	}

	return products, err
}

func (pg *Pg) GetProducts(ctx context.Context) ([]entities.Product, error) {
	var products []entities.Product

	query := "SELECT * FROM products"
	err := pg.db.SelectContext(ctx, &products, query)

	return products, err
}

func (pg *Pg) GetProductWithoutCategoryType(ctx context.Context, model, additional string) (*entities.Product, error) {
	var product struct {
		entities.Product
		CategoryType int `db:"c_type"`
	}

	query := "SELECT products.*, categories.c_type AS c_type FROM products JOIN categories ON products.category_id = categories.id WHERE products.model = $1 AND products.additional = $2"
	err := pg.db.GetContext(ctx, &product, query, model, additional)

	if product.CategoryType == 0 {
		product.Additional = fmt.Sprintf("%s [Новый]", product.Additional)
	} else {
		product.Additional = fmt.Sprintf("%s [Б/У]", product.Additional)
	}

	return &product.Product, err
}

func (pg *Pg) GetProduct(ctx context.Context, model, additional string, cType int) (*entities.Product, error) {
	var product entities.Product

	query := "SELECT products.* FROM categories JOIN products ON products.model = $1 AND products.additional = $2 AND products.category_id = categories.id WHERE categories.c_type = $3"
	err := pg.db.GetContext(ctx, &product, query, model, additional, cType)

	return &product, err
}

func (pg *Pg) GetProductByID(ctx context.Context, productID uuid.UUID) (*entities.Product, error) {
	var product entities.Product

	query := "SELECT * FROM products WHERE id = $1"
	err := pg.db.GetContext(ctx, &product, query, productID)

	return &product, err
}

func (pg *Pg) GetProducersByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]string, error) {
	var producers []string

	query := "SELECT DISTINCT(producer) FROM products WHERE category_id = $1"
	err := pg.db.SelectContext(ctx, &producers, query, categoryID)

	if err != nil {
		return nil, err
	} else if len(producers) < 1 {
		return nil, sql.ErrNoRows
	}

	return producers, err
}

func (pg *Pg) GetModelsByCategoryIDAndProducer(ctx context.Context, categoryID uuid.UUID, producer string) ([]string, error) {
	var models []string

	query := "SELECT DISTINCT(model) FROM products WHERE category_id = $1 AND lower(producer) = lower($2)"
	err := pg.db.SelectContext(ctx, &models, query, categoryID, producer)

	if err != nil {
		return nil, err
	} else if len(models) < 1 {
		return nil, sql.ErrNoRows
	}

	return models, err
}

func (pg *Pg) DeleteProductByID(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM products WHERE id = $1"
	_, err := pg.db.ExecContext(ctx, query, id)

	return err
}

func (pg *Pg) DeleteProductsByCategoryID(ctx context.Context, categoryID uuid.UUID) error {
	query := "DELETE FROM products WHERE category_id = $1"
	_, err := pg.db.ExecContext(ctx, query, categoryID)

	return err
}
