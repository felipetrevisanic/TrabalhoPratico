using src.domain.entities;
using src.DTO.requests;
using src.DTO.response;

namespace src.mappings;

public static class ProductMappings
{
    public static ProductResponseDto ToResponseDto(this Product product)
    {
        return new ProductResponseDto
        {
            Id = product.Id,
            Name = product.Name,
            Description = product.Description,
            Price = product.Price,
            StockQuantity = product.StockQuantity,
            CreatedAt = product.CreatedAt,
            UpdatedAt = product.UpdatedAt
        };
    }

    public static Product ToEntity(this CreateProductRequestDto request)
    {
        return new Product
        {
            Name = request.Name,
            Description = request.Description,
            Price = request.Price,
            StockQuantity = request.StockQuantity,
            CreatedAt = DateTime.UtcNow
        };
    }

    public static Product ToEntityFromUpdate(this UpdateProductRequestDto request, int id)
    {
        return new Product
        {
            Id = id,
            Name = request.Name,
            Description = request.Description,
            Price = request.Price,
            StockQuantity = request.StockQuantity,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };
    }

    public static void UpdateEntity(this UpdateProductRequestDto request, Product product)
    {
        product.Name = request.Name;
        product.Description = request.Description;
        product.Price = request.Price;
        product.StockQuantity = request.StockQuantity;
        product.UpdatedAt = DateTime.UtcNow;
    }
}
