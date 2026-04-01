package postgres

import (
	"database/sql"
	"errors"
	"time"

	"grpc/internal/product"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetByID(id int) (*product.Product, error) {
	const query = `
		SELECT "Id", "Name", "Description", "Price", "StockQuantity", "CreatedAt", "UpdatedAt"
		FROM public.products
		WHERE "Id" = $1
	`

	row := r.db.QueryRow(query, id)
	var item product.Product
	var updatedAt sql.NullTime

	if err := row.Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Price,
		&item.StockQuantity,
		&item.CreatedAt,
		&updatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	item.UpdatedAt = toTimePointer(updatedAt)
	return &item, nil
}

func (r *ProductRepository) GetAll() ([]product.Product, error) {
	const query = `
		SELECT "Id", "Name", "Description", "Price", "StockQuantity", "CreatedAt", "UpdatedAt"
		FROM public.products
		ORDER BY "Id"
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]product.Product, 0)
	for rows.Next() {
		var item product.Product
		var updatedAt sql.NullTime

		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Price,
			&item.StockQuantity,
			&item.CreatedAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		item.UpdatedAt = toTimePointer(updatedAt)
		products = append(products, item)
	}

	return products, rows.Err()
}

func (r *ProductRepository) Create(input product.CreateInput) (*product.Product, error) {
	const query = `
		INSERT INTO public.products ("Name", "Description", "Price", "StockQuantity", "CreatedAt")
		VALUES ($1, $2, $3, $4, $5)
		RETURNING "Id", "Name", "Description", "Price", "StockQuantity", "CreatedAt", "UpdatedAt"
	`

	now := time.Now().UTC()
	var item product.Product
	var updatedAt sql.NullTime

	err := r.db.QueryRow(
		query,
		input.Name,
		input.Description,
		input.Price,
		input.StockQuantity,
		now,
	).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Price,
		&item.StockQuantity,
		&item.CreatedAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	item.UpdatedAt = toTimePointer(updatedAt)
	return &item, nil
}

func (r *ProductRepository) Update(id int, input product.UpdateInput) (*product.Product, error) {
	existing, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return r.Create(product.CreateInput{
			Name:          input.Name,
			Description:   input.Description,
			Price:         input.Price,
			StockQuantity: input.StockQuantity,
		})
	}

	const query = `
		UPDATE public.products
		SET "Name" = $1,
			"Description" = $2,
			"Price" = $3,
			"StockQuantity" = $4,
			"UpdatedAt" = $5
		WHERE "Id" = $6
		RETURNING "Id", "Name", "Description", "Price", "StockQuantity", "CreatedAt", "UpdatedAt"
	`

	var item product.Product
	var updatedAt sql.NullTime

	err = r.db.QueryRow(
		query,
		input.Name,
		input.Description,
		input.Price,
		input.StockQuantity,
		time.Now().UTC(),
		id,
	).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Price,
		&item.StockQuantity,
		&item.CreatedAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	item.UpdatedAt = toTimePointer(updatedAt)
	return &item, nil
}

func (r *ProductRepository) Delete(id int) (bool, error) {
	result, err := r.db.Exec(`DELETE FROM public.products WHERE "Id" = $1`, id)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func (r *ProductRepository) DeleteAll() (bool, error) {
	result, err := r.db.Exec(`DELETE FROM public.products`)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func toTimePointer(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}

	timestamp := value.Time
	return &timestamp
}
