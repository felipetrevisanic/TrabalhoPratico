namespace src.application.interfaces;

public interface IProductService
{
    string GetProductById(int id);
    IEnumerable<string> GetAllProducts();
    string InsertProduct();
    string UpdateProduct();
}
