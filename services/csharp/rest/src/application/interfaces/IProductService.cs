using src.interfaces.rest.requests;
using src.interfaces.rest.responses;

namespace src.application.interfaces;

public interface IProductService
{
    ProductResponseDto GetProductById(int id);
    IEnumerable<ProductResponseDto> GetAllProducts();
    ProductResponseDto InsertProduct(CreateProductRequestDto request);
    ProductResponseDto UpdateProduct(int id, UpdateProductRequestDto request);
    bool DeleteProduct(int id);
    bool DeleteProduct();
}
