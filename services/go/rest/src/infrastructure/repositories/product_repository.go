package repositories

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"

	"rest/src/domain/entities"
	domaininterfaces "rest/src/domain/interfaces"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetByID(id int) (*entities.Product, error) {
	const query = `
		SELECT "Id", "Name", "Description", "Category", "Images", "Price", "StockQuantity", "CreatedAt", "UpdatedAt"
		FROM public.products
		WHERE "Id" = $1
	`

	row := r.db.QueryRow(query, id)

	var item entities.Product
	var updatedAt sql.NullTime

	if err := row.Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Category,
		pq.Array(&item.Images),
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

func (r *ProductRepository) GetAll() ([]entities.Product, error) {
	const query = `
		SELECT "Id", "Name", "Description", "Category", "Images", "Price", "StockQuantity", "CreatedAt", "UpdatedAt"
		FROM public.products
		ORDER BY "Id"
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entities.Product, 0)

	for rows.Next() {
		var item entities.Product
		var updatedAt sql.NullTime

		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Category,
			pq.Array(&item.Images),
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

func (r *ProductRepository) Create(request entities.CreateRequest) (*entities.Product, error) {
	const query = `
		INSERT INTO public.products ("Name", "Description", "Category", "Images", "Price", "StockQuantity", "CreatedAt")
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING "Id", "Name", "Description", "Category", "Images", "Price", "StockQuantity", "CreatedAt", "UpdatedAt"
	`

	now := time.Now().UTC()
	var item entities.Product
	var updatedAt sql.NullTime

	err := r.db.QueryRow(
		query,
		request.Name,
		request.Description,
		request.Category,
		pq.Array(request.Images),
		request.Price,
		request.StockQuantity,
		now,
	).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Category,
		pq.Array(&item.Images),
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

func (r *ProductRepository) Update(id int, request entities.UpdateRequest) (*entities.Product, error) {
	existing, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		created, err := r.Create(entities.CreateRequest{
			Name:          request.Name,
			Description:   request.Description,
			Category:      request.Category,
			Images:        request.Images,
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
			"Category" = $3,
			"Images" = $4,
			"Price" = $5,
			"StockQuantity" = $6,
			"UpdatedAt" = $7
		WHERE "Id" = $8
		RETURNING "Id", "Name", "Description", "Category", "Images", "Price", "StockQuantity", "CreatedAt", "UpdatedAt"
	`

	updatedAt := time.Now().UTC()
	var item entities.Product
	var nullableUpdatedAt sql.NullTime

	err = r.db.QueryRow(
		query,
		request.Name,
		request.Description,
		request.Category,
		pq.Array(request.Images),
		request.Price,
		request.StockQuantity,
		updatedAt,
		id,
	).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Category,
		pq.Array(&item.Images),
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

var _ domaininterfaces.ProductRepository = (*ProductRepository)(nil)
