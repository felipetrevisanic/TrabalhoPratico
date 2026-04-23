using src.domain.entities;
using src.interfaces.graphql.inputs;

namespace src.interfaces.graphql.mappings;

public static class ProductMappings
{
    public static Product ToEntity(this CreateProductInput request)
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

    public static Product ToEntityFromUpdate(this UpdateProductInput request, int id)
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

    public static void UpdateEntity(this UpdateProductInput request, Product product)
    {
        product.Name = request.Name;
        product.Description = request.Description;
        product.Price = request.Price;
        product.StockQuantity = request.StockQuantity;
        product.UpdatedAt = DateTime.UtcNow;
    }
}
