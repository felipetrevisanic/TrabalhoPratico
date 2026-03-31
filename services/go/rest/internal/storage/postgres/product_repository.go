package postgres

import (
	"database/sql"
	"errors"
	"time"

	"rest/internal/product"
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

	item.UpdatedAt = nullTimeToPointer(updatedAt)
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

	items := make([]product.Product, 0)

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

		item.UpdatedAt = nullTimeToPointer(updatedAt)
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *ProductRepository) Create(request product.CreateRequest) (*product.Product, error) {
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
		request.Name,
		request.Description,
		request.Price,
		request.StockQuantity,
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

	item.UpdatedAt = nullTimeToPointer(updatedAt)
	return &item, nil
}

func (r *ProductRepository) Update(id int, request product.UpdateRequest) (*product.Product, error) {
	existing, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		created, err := r.Create(product.CreateRequest{
			Name:          request.Name,
			Description:   request.Description,
			Price:         request.Price,
			StockQuantity: request.StockQuantity,
		})
		if err != nil {
			return nil, err
		}

		return created, nil
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

	updatedAt := time.Now().UTC()
	var item product.Product
	var nullableUpdatedAt sql.NullTime

	err = r.db.QueryRow(
		query,
		request.Name,
		request.Description,
		request.Price,
		request.StockQuantity,
		updatedAt,
		id,
	).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Price,
		&item.StockQuantity,
		&item.CreatedAt,
		&nullableUpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	item.UpdatedAt = nullTimeToPointer(nullableUpdatedAt)
	return &item, nil
}

func (r *ProductRepository) Delete(id int) (bool, error) {
	const query = `DELETE FROM public.products WHERE "Id" = $1`

	result, err := r.db.Exec(query, id)
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
	const query = `DELETE FROM public.products`

	result, err := r.db.Exec(query)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func nullTimeToPointer(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}

	timestamp := value.Time
	return &timestamp
}
