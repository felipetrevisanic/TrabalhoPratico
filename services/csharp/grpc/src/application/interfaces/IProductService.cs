using src.domain.entities;

namespace src.application.interfaces;

public interface IProductService
{
    Product? GetProductById(int id);
    IEnumerable<Product> GetAllProducts();
    Product InsertProduct(string name, string description, string category, string[] images, decimal price, int stockQuantity);
    Product UpdateProduct(int id, string name, string description, string category, string[] images, decimal price, int stockQuantity);
    bool DeleteProduct(int id);
    bool DeleteProduct();
}
