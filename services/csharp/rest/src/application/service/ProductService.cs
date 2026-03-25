using src.application.interfaces;

namespace src.application.service;

public class ProductService : IProductService
{
    public string GetProductById(int id)
    {
        return $"Product {id}";
    }

    public IEnumerable<string> GetAllProducts()
    {
        return
        [
            "Product 1",
            "Product 2",
            "Product 3"
        ];
    }

    public string InsertProduct()
    {
        return "Product inserted successfully";
    }

    public string UpdateProduct()
    {
        return "Product updated successfully";
    }
}
