using src.domain.entities;
using src.domain.interfaces;

namespace src.infraestructure.Repositories;

public class ProductRepository : IProductRepository
{
    private static readonly List<Product> Products =
    [
        new()
        {
            Id = 1,
            Name = "Notebook",
            Description = "Notebook para desenvolvimento",
            Price = 4500.00m,
            StockQuantity = 10,
            CreatedAt = DateTime.UtcNow
        },
        new()
        {
            Id = 2,
            Name = "Mouse",
            Description = "Mouse sem fio",
            Price = 120.50m,
            StockQuantity = 25,
            CreatedAt = DateTime.UtcNow
        }
    ];

    public Product? GetById(int id)
    {
        return Products.FirstOrDefault(product => product.Id == id);
    }

    public IEnumerable<Product> GetAll()
    {
        return Products;
    }

    public Product Add(Product product)
    {
        product.Id = Products.Count == 0 ? 1 : Products.Max(existingProduct => existingProduct.Id) + 1;
        Products.Add(product);
        return product;
    }

    public Product Update(Product product)
    {
        var existingProduct = GetById(product.Id);

        if (existingProduct is null)
        {
            Products.Add(product);
            return product;
        }

        existingProduct.Name = product.Name;
        existingProduct.Description = product.Description;
        existingProduct.Price = product.Price;
        existingProduct.StockQuantity = product.StockQuantity;
        existingProduct.UpdatedAt = product.UpdatedAt;

        return existingProduct;
    }

    public bool Delete(int id)
    {
        var product = GetById(id);

        if (product is null)
        {
            return false;
        }

        Products.Remove(product);
        return true;
    }

    public bool DeleteAll()
    {
        if (Products.Count == 0)
        {
            return false;
        }

        Products.Clear();
        return true;
    }
}
