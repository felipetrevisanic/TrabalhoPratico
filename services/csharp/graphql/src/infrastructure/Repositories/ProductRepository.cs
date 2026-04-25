using Microsoft.EntityFrameworkCore;
using src.domain.entities;
using src.domain.interfaces;
using src.infrastructure.Data;

namespace src.infrastructure.Repositories;

public class ProductRepository : IProductRepository
{
    private readonly AppDbContext _context;

    public ProductRepository(AppDbContext context)
    {
        _context = context;
    }

    public Product? GetById(int id)
    {
        return _context.Products.AsNoTracking().FirstOrDefault(product => product.Id == id);
    }

    public IEnumerable<Product> GetAll()
    {
        return _context.Products
            .AsNoTracking()
            .OrderBy(product => product.Id)
            .ToList();
    }

    public Product Add(Product product)
    {
        _context.Products.Add(product);
        _context.SaveChanges();
        return product;
    }

    public Product Update(Product product)
    {
        var existingProduct = _context.Products.FirstOrDefault(existingProduct => existingProduct.Id == product.Id);

        if (existingProduct is null)
        {
            _context.Products.Add(product);
            _context.SaveChanges();
            return product;
        }

        existingProduct.Name = product.Name;
        existingProduct.Description = product.Description;
        existingProduct.Category = product.Category;
        existingProduct.Images = product.Images;
        existingProduct.Price = product.Price;
        existingProduct.StockQuantity = product.StockQuantity;
        existingProduct.CreatedAt = product.CreatedAt;
        existingProduct.UpdatedAt = product.UpdatedAt;
        _context.SaveChanges();

        return existingProduct;
    }

    public bool Delete(int id)
    {
        var product = _context.Products.FirstOrDefault(existingProduct => existingProduct.Id == id);

        if (product is null)
        {
            return false;
        }

        _context.Products.Remove(product);
        _context.SaveChanges();
        return true;
    }

    public bool DeleteAll()
    {
        var products = _context.Products.ToList();

        if (products.Count == 0)
        {
            return false;
        }

        _context.Products.RemoveRange(products);
        _context.SaveChanges();
        return true;
    }
}
