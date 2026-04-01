using src.application.interfaces;
using src.domain.entities;
using src.domain.interfaces;
using src.graphql.inputs;
using src.mappings;

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

    public Product InsertProduct(CreateProductInput request)
    {
        var product = request.ToEntity();
        return _productRepository.Add(product);
    }

    public Product UpdateProduct(int id, UpdateProductInput request)
    {
        var product = _productRepository.GetById(id);

        if (product is null)
        {
            product = request.ToEntityFromUpdate(id);
            return _productRepository.Add(product);
        }

        request.UpdateEntity(product);
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
