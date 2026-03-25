using src.domain.entities;

namespace src.domain.interfaces;

public interface IProductRepository
{
    Product? GetById(int id);
    IEnumerable<Product> GetAll();
    Product Add(Product product);
    Product Update(Product product);
    bool Delete(int id);
    bool DeleteAll();
}
