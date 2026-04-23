using src.application.interfaces;
using src.domain.entities;
using src.interfaces.graphql.inputs;

namespace src.interfaces.graphql.mutations;

public class ProductMutation
{
    public Product CreateProduct([Service] IProductService productService, CreateProductInput input)
    {
        return productService.InsertProduct(input);
    }

    public Product UpdateProduct([Service] IProductService productService, int id, UpdateProductInput input)
    {
        return productService.UpdateProduct(id, input);
    }

    public bool DeleteProduct([Service] IProductService productService, int id)
    {
        return productService.DeleteProduct(id);
    }

    public bool DeleteAllProducts([Service] IProductService productService)
    {
        return productService.DeleteProduct();
    }
}
