package graphql

import (
	"time"

	gql "github.com/graphql-go/graphql"

	appinterfaces "graphql/src/application/interfaces"
	"graphql/src/domain/entities"
)

func productByIDResolver(service appinterfaces.ProductService) gql.FieldResolveFn {
	return func(params gql.ResolveParams) (any, error) {
		id := params.Args["id"].(int)
		item, err := service.GetByID(id)
		if err != nil {
			return nil, err
		}

		if item == nil {
			return nil, nil
		}

		return mapProduct(*item), nil
	}
}

func allProductsResolver(service appinterfaces.ProductService) gql.FieldResolveFn {
	return func(params gql.ResolveParams) (any, error) {
		items, err := service.GetAll()
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

func createProductResolver(service appinterfaces.ProductService) gql.FieldResolveFn {
	return func(params gql.ResolveParams) (any, error) {
		input := params.Args["input"].(map[string]any)

		item, err := service.Create(entities.CreateInput{
			Name:          input["name"].(string),
			Description:   input["description"].(string),
			Category:      input["category"].(string),
			Images:        toStringSlice(input["images"]),
			Price:         input["price"].(float64),
			StockQuantity: input["stockQuantity"].(int),
		})
		if err != nil {
			return nil, err
		}

		return mapProduct(*item), nil
	}
}

func updateProductResolver(service appinterfaces.ProductService) gql.FieldResolveFn {
	return func(params gql.ResolveParams) (any, error) {
		id := params.Args["id"].(int)
		input := params.Args["input"].(map[string]any)

		item, err := service.Update(id, entities.UpdateInput{
			Name:          input["name"].(string),
			Description:   input["description"].(string),
			Category:      input["category"].(string),
			Images:        toStringSlice(input["images"]),
			Price:         input["price"].(float64),
			StockQuantity: input["stockQuantity"].(int),
		})
		if err != nil {
			return nil, err
		}

		return mapProduct(*item), nil
	}
}

func deleteProductResolver(service appinterfaces.ProductService) gql.FieldResolveFn {
	return func(params gql.ResolveParams) (any, error) {
		id := params.Args["id"].(int)
		return service.Delete(id)
	}
}

func deleteAllProductsResolver(service appinterfaces.ProductService) gql.FieldResolveFn {
	return func(params gql.ResolveParams) (any, error) {
		return service.DeleteAll()
	}
}

func mapProduct(item entities.Product) map[string]any {
	response := map[string]any{
		"id":            item.ID,
		"name":          item.Name,
		"description":   item.Description,
		"category":      item.Category,
		"images":        item.Images,
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

func toStringSlice(value any) []string {
	items, ok := value.([]any)
	if !ok {
		return []string{}
	}

	result := make([]string, 0, len(items))
	for _, item := range items {
		text, ok := item.(string)
		if ok {
			result = append(result, text)
		}
	}

	return result
}
