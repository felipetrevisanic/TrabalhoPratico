using src.application.interfaces;
using src.domain.entities;
using src.domain.interfaces;
using src.DTO.requests;
using src.DTO.response;
using src.mappings;

namespace src.application.service;

public class ProductService : IProductService
{
    private readonly IProductRepository _productRepository;

    public ProductService(IProductRepository productRepository)
    {
        _productRepository = productRepository;
    }

    public ProductResponseDto GetProductById(int id)
    {
        var product = _productRepository.GetById(id)
            ?? new Product
            {
                Id = id,
                Name = $"Product {id}",
                Description = "Product not found in sample list",
                Price = 0,
                StockQuantity = 0,
                CreatedAt = DateTime.UtcNow
            };

        return product.ToResponseDto();
    }

    public IEnumerable<ProductResponseDto> GetAllProducts()
    {
        return _productRepository.GetAll().Select(product => product.ToResponseDto());
    }

    public ProductResponseDto InsertProduct(CreateProductRequestDto request)
    {
        var product = request.ToEntity();
        var createdProduct = _productRepository.Add(product);
        return createdProduct.ToResponseDto();
    }

    public ProductResponseDto UpdateProduct(int id, UpdateProductRequestDto request)
    {
        var product = _productRepository.GetById(id);

        if (product is null)
        {
            product = request.ToEntityFromUpdate(id);
            var createdProduct = _productRepository.Add(product);
            return createdProduct.ToResponseDto();
        }

        request.UpdateEntity(product);
        var updatedProduct = _productRepository.Update(product);
        return updatedProduct.ToResponseDto();
    }

    public bool DeleteProduct(int id)
    {
        return _productRepository.Delete(id);
    }
}
