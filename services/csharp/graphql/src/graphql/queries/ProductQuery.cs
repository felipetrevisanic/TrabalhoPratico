using src.application.interfaces;
using src.domain.entities;

namespace src.graphql.queries;

public class ProductQuery
{
    public Product ProductById([Service] IProductService productService, int id)
    {
        return productService.GetProductById(id);
    }

    public IEnumerable<Product> AllProducts([Service] IProductService productService)
    {
        return productService.GetAllProducts();
    }
}
