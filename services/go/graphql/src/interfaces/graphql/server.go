package graphql

import (
	"net/http"

	gql "github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	appinterfaces "graphql/src/application/interfaces"
)

func NewServer(service appinterfaces.ProductService) (http.Handler, error) {
	schema, err := newSchema(service)
	if err != nil {
		return nil, err
	}

	return handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	}), nil
}

func newSchema(service appinterfaces.ProductService) (gql.Schema, error) {
	productType := gql.NewObject(gql.ObjectConfig{
		Name: "Product",
		Fields: gql.Fields{
			"id":            &gql.Field{Type: gql.NewNonNull(gql.Int)},
			"name":          &gql.Field{Type: gql.NewNonNull(gql.String)},
			"description":   &gql.Field{Type: gql.NewNonNull(gql.String)},
			"category":      &gql.Field{Type: gql.NewNonNull(gql.String)},
			"images":        &gql.Field{Type: gql.NewNonNull(gql.NewList(gql.NewNonNull(gql.String)))},
			"price":         &gql.Field{Type: gql.NewNonNull(gql.Float)},
			"stockQuantity": &gql.Field{Type: gql.NewNonNull(gql.Int)},
			"createdAt":     &gql.Field{Type: gql.NewNonNull(gql.String)},
			"updatedAt":     &gql.Field{Type: gql.String},
		},
	})

	createInputType := gql.NewInputObject(gql.InputObjectConfig{
		Name: "CreateProductInput",
		Fields: gql.InputObjectConfigFieldMap{
			"name":          &gql.InputObjectFieldConfig{Type: gql.NewNonNull(gql.String)},
			"description":   &gql.InputObjectFieldConfig{Type: gql.NewNonNull(gql.String)},
			"category":      &gql.InputObjectFieldConfig{Type: gql.NewNonNull(gql.String)},
			"images":        &gql.InputObjectFieldConfig{Type: gql.NewNonNull(gql.NewList(gql.NewNonNull(gql.String)))},
			"price":         &gql.InputObjectFieldConfig{Type: gql.NewNonNull(gql.Float)},
			"stockQuantity": &gql.InputObjectFieldConfig{Type: gql.NewNonNull(gql.Int)},
		},
	})

	updateInputType := gql.NewInputObject(gql.InputObjectConfig{
		Name: "UpdateProductInput",
		Fields: gql.InputObjectConfigFieldMap{
			"name":          &gql.InputObjectFieldConfig{Type: gql.NewNonNull(gql.String)},
			"description":   &gql.InputObjectFieldConfig{Type: gql.NewNonNull(gql.String)},
			"category":      &gql.InputObjectFieldConfig{Type: gql.NewNonNull(gql.String)},
			"images":        &gql.InputObjectFieldConfig{Type: gql.NewNonNull(gql.NewList(gql.NewNonNull(gql.String)))},
			"price":         &gql.InputObjectFieldConfig{Type: gql.NewNonNull(gql.Float)},
			"stockQuantity": &gql.InputObjectFieldConfig{Type: gql.NewNonNull(gql.Int)},
		},
	})

	queryType := gql.NewObject(gql.ObjectConfig{
		Name: "Query",
		Fields: gql.Fields{
			"productById": &gql.Field{
				Type: productType,
				Args: gql.FieldConfigArgument{
					"id": &gql.ArgumentConfig{Type: gql.NewNonNull(gql.Int)},
				},
				Resolve: productByIDResolver(service),
			},
			"allProducts": &gql.Field{
				Type:    gql.NewNonNull(gql.NewList(gql.NewNonNull(productType))),
				Resolve: allProductsResolver(service),
			},
		},
	})

	mutationType := gql.NewObject(gql.ObjectConfig{
		Name: "Mutation",
		Fields: gql.Fields{
			"createProduct": &gql.Field{
				Type: productType,
				Args: gql.FieldConfigArgument{
					"input": &gql.ArgumentConfig{Type: gql.NewNonNull(createInputType)},
				},
				Resolve: createProductResolver(service),
			},
			"updateProduct": &gql.Field{
				Type: productType,
				Args: gql.FieldConfigArgument{
					"id":    &gql.ArgumentConfig{Type: gql.NewNonNull(gql.Int)},
					"input": &gql.ArgumentConfig{Type: gql.NewNonNull(updateInputType)},
				},
				Resolve: updateProductResolver(service),
			},
			"deleteProduct": &gql.Field{
				Type: gql.NewNonNull(gql.Boolean),
				Args: gql.FieldConfigArgument{
					"id": &gql.ArgumentConfig{Type: gql.NewNonNull(gql.Int)},
				},
				Resolve: deleteProductResolver(service),
			},
			"deleteAllProducts": &gql.Field{
				Type:    gql.NewNonNull(gql.Boolean),
				Resolve: deleteAllProductsResolver(service),
			},
		},
	})

	return gql.NewSchema(gql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
}
