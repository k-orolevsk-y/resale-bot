package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
)

func (pg *Pg) CreateProduct(ctx context.Context, product *entities.Product) error {
	query := "INSERT INTO products (category_id, producer, model, operating_system, additional, description, photo, price, old_price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
	id, err := pg.db.ExecContextWithReturnID(ctx, query, product.CategoryID, product.Producer, product.Model, product.OperatingSystem, product.Additional, product.Description, product.Photo, product.Price, product.OldPrice)

	if err != nil {
		return err
	}

	product.ID = uuid.MustParse(id.(string))
	return nil
}

func (pg *Pg) GetSaleProducts(ctx context.Context) ([]entities.Product, error) {
	var products []entities.Product

	query := "SELECT products.* FROM products JOIN categories ON products.category_id = categories.id WHERE products.is_sale = true"
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
