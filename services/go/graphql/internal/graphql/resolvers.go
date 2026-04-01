package graphql

import (
	"fmt"
	"time"

	gql "github.com/graphql-go/graphql"

	"graphql/internal/product"
)

func productByIDResolver(repository product.Repository) gql.FieldResolveFn {
	return func(params gql.ResolveParams) (any, error) {
		id := params.Args["id"].(int)
		item, err := repository.GetByID(id)
		if err != nil {
			return nil, err
		}

		if item == nil {
			return mapProduct(product.Product{
				ID:            id,
				Name:          fmt.Sprintf("Product %d", id),
				Description:   "Product not found in sample list",
				Price:         0,
				StockQuantity: 0,
				CreatedAt:     time.Now().UTC(),
			}), nil
		}

		return mapProduct(*item), nil
	}
}

func allProductsResolver(repository product.Repository) gql.FieldResolveFn {
	return func(params gql.ResolveParams) (any, error) {
		items, err := repository.GetAll()
		if err != nil {
			return nil, err
		}

		response := make([]map[string]any, 0, len(items))
		for _, item := range items {
			response = append(response, mapProduct(item))
		}

		return response, nil
	}
}

func createProductResolver(repository product.Repository) gql.FieldResolveFn {
	return func(params gql.ResolveParams) (any, error) {
		input := params.Args["input"].(map[string]any)

		item, err := repository.Create(product.CreateInput{
			Name:          input["name"].(string),
			Description:   input["description"].(string),
			Price:         input["price"].(float64),
			StockQuantity: input["stockQuantity"].(int),
		})
		if err != nil {
			return nil, err
		}

		return mapProduct(*item), nil
	}
}

func updateProductResolver(repository product.Repository) gql.FieldResolveFn {
	return func(params gql.ResolveParams) (any, error) {
		id := params.Args["id"].(int)
		input := params.Args["input"].(map[string]any)

		item, err := repository.Update(id, product.UpdateInput{
			Name:          input["name"].(string),
			Description:   input["description"].(string),
			Price:         input["price"].(float64),
			StockQuantity: input["stockQuantity"].(int),
		})
		if err != nil {
			return nil, err
		}

		return mapProduct(*item), nil
	}
}

func deleteProductResolver(repository product.Repository) gql.FieldResolveFn {
	return func(params gql.ResolveParams) (any, error) {
		id := params.Args["id"].(int)
		return repository.Delete(id)
	}
}

func deleteAllProductsResolver(repository product.Repository) gql.FieldResolveFn {
	return func(params gql.ResolveParams) (any, error) {
		return repository.DeleteAll()
	}
}

func mapProduct(item product.Product) map[string]any {
	response := map[string]any{
		"id":            item.ID,
		"name":          item.Name,
		"description":   item.Description,
		"price":         item.Price,
		"stockQuantity": item.StockQuantity,
		"createdAt":     item.CreatedAt.UTC().Format(time.RFC3339Nano),
		"updatedAt":     nil,
	}

	if item.UpdatedAt != nil {
		response["updatedAt"] = item.UpdatedAt.UTC().Format(time.RFC3339Nano)
	}

	return response
}
