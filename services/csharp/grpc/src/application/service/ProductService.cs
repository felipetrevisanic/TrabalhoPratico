using src.application.interfaces;
using src.domain.entities;
using src.domain.interfaces;

namespace src.application.service;

public class ProductService : IProductService
{
    private readonly IProductRepository _productRepository;

    public ProductService(IProductRepository productRepository)
    {
        _productRepository = productRepository;
    }

    public Product GetProductById(int id)
    {
        return _productRepository.GetById(id)
            ?? new Product
            {
                Id = id,
                Name = $"Product {id}",
                Description = "Product not found in sample list",
                Price = 0,
                StockQuantity = 0,
                CreatedAt = DateTime.UtcNow
            };
    }

    public IEnumerable<Product> GetAllProducts()
    {
        return _productRepository.GetAll();
    }

    public Product InsertProduct(string name, string description, decimal price, int stockQuantity)
    {
        var product = new Product
        {
            Name = name,
            Description = description,
            Price = price,
            StockQuantity = stockQuantity,
            CreatedAt = DateTime.UtcNow
        };

        return _productRepository.Add(product);
    }

    public Product UpdateProduct(int id, string name, string description, decimal price, int stockQuantity)
    {
        var product = _productRepository.GetById(id);

        if (product is null)
        {
            product = new Product
            {
                Id = id,
                Name = name,
                Description = description,
                Price = price,
                StockQuantity = stockQuantity,
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            };

            return _productRepository.Add(product);
        }

        product.Name = name;
        product.Description = description;
        product.Price = price;
        product.StockQuantity = stockQuantity;
        product.UpdatedAt = DateTime.UtcNow;

        return _productRepository.Update(product);
    }

    public bool DeleteProduct(int id)
    {
        return _productRepository.Delete(id);
    }

    public bool DeleteProduct()
    {
        return _productRepository.DeleteAll();
    }
}
