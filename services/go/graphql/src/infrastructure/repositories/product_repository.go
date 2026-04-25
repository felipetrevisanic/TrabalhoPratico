package repositories

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"

	"graphql/src/domain/entities"
	domaininterfaces "graphql/src/domain/interfaces"
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

	item.UpdatedAt = toTimePointer(updatedAt)
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

	products := make([]entities.Product, 0)
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

		item.UpdatedAt = toTimePointer(updatedAt)
		products = append(products, item)
	}

	return products, rows.Err()
}

func (r *ProductRepository) Create(input entities.CreateInput) (*entities.Product, error) {
	const query = `
		INSERT INTO public.products ("Name", "Description", "Category", "Images", "Price", "StockQuantity", "CreatedAt")
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING "Id", "Name", "Description", "Category", "Images", "Price", "StockQuantity", "CreatedAt", "UpdatedAt"
	`

	var item entities.Product
	var updatedAt sql.NullTime
	now := time.Now().UTC()

	err := r.db.QueryRow(
		query,
		input.Name,
		input.Description,
		input.Category,
		pq.Array(input.Images),
		input.Price,
		input.StockQuantity,
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

	item.UpdatedAt = toTimePointer(updatedAt)
	return &item, nil
}

func (r *ProductRepository) Update(id int, input entities.UpdateInput) (*entities.Product, error) {
	existing, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return r.Create(entities.CreateInput{
			Name:          input.Name,
			Description:   input.Description,
			Category:      input.Category,
			Images:        input.Images,
			Price:         input.Price,
			StockQuantity: input.StockQuantity,
		})
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

	var item entities.Product
	var updatedAt sql.NullTime

	err = r.db.QueryRow(
		query,
		input.Name,
		input.Description,
		input.Category,
		pq.Array(input.Images),
		input.Price,
		input.StockQuantity,
		time.Now().UTC(),
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

var _ domaininterfaces.ProductRepository = (*ProductRepository)(nil)
