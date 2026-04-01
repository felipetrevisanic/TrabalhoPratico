package graphql

import (
	"net/http"

	gql "github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"graphql/internal/product"
)

func NewServer(repository product.Repository) (http.Handler, error) {
	schema, err := newSchema(repository)
	if err != nil {
		return nil, err
	}

	return handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	}), nil
}

func newSchema(repository product.Repository) (gql.Schema, error) {
	productType := gql.NewObject(gql.ObjectConfig{
		Name: "Product",
		Fields: gql.Fields{
			"id":            &gql.Field{Type: gql.NewNonNull(gql.Int)},
			"name":          &gql.Field{Type: gql.NewNonNull(gql.String)},
			"description":   &gql.Field{Type: gql.NewNonNull(gql.String)},
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
			"price":         &gql.InputObjectFieldConfig{Type: gql.NewNonNull(gql.Float)},
			"stockQuantity": &gql.InputObjectFieldConfig{Type: gql.NewNonNull(gql.Int)},
		},
	})

	updateInputType := gql.NewInputObject(gql.InputObjectConfig{
		Name: "UpdateProductInput",
		Fields: gql.InputObjectConfigFieldMap{
			"name":          &gql.InputObjectFieldConfig{Type: gql.NewNonNull(gql.String)},
			"description":   &gql.InputObjectFieldConfig{Type: gql.NewNonNull(gql.String)},
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
				Resolve: productByIDResolver(repository),
			},
			"allProducts": &gql.Field{
				Type:    gql.NewNonNull(gql.NewList(gql.NewNonNull(productType))),
				Resolve: allProductsResolver(repository),
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
				Resolve: createProductResolver(repository),
			},
			"updateProduct": &gql.Field{
				Type: productType,
				Args: gql.FieldConfigArgument{
					"id":    &gql.ArgumentConfig{Type: gql.NewNonNull(gql.Int)},
					"input": &gql.ArgumentConfig{Type: gql.NewNonNull(updateInputType)},
				},
				Resolve: updateProductResolver(repository),
			},
			"deleteProduct": &gql.Field{
				Type: gql.NewNonNull(gql.Boolean),
				Args: gql.FieldConfigArgument{
					"id": &gql.ArgumentConfig{Type: gql.NewNonNull(gql.Int)},
				},
				Resolve: deleteProductResolver(repository),
			},
			"deleteAllProducts": &gql.Field{
				Type:    gql.NewNonNull(gql.Boolean),
				Resolve: deleteAllProductsResolver(repository),
			},
		},
	})

	return gql.NewSchema(gql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
}
