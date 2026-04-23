using src.domain.entities;
using src.interfaces.graphql.inputs;

namespace src.application.interfaces;

public interface IProductService
{
    Product GetProductById(int id);
    IEnumerable<Product> GetAllProducts();
    Product InsertProduct(CreateProductInput request);
    Product UpdateProduct(int id, UpdateProductInput request);
    bool DeleteProduct(int id);
    bool DeleteProduct();
}
