using Grpc.Core;
using ProductGrpc.V1;
using src.application.interfaces;
using src.domain.entities;

namespace src.services;

public class ProductGrpcService : ProductService.ProductServiceBase
{
    private readonly IProductService _productService;

    public ProductGrpcService(IProductService productService)
    {
        _productService = productService;
    }

    public override Task<ProductResponse> GetProductById(GetProductByIdRequest request, ServerCallContext context)
    {
        var product = _productService.GetProductById(request.Id);
        return Task.FromResult(ToResponse(product));
    }

    public override Task<GetAllProductsResponse> GetAllProducts(GetAllProductsRequest request, ServerCallContext context)
    {
        var response = new GetAllProductsResponse();
        response.Products.AddRange(_productService.GetAllProducts().Select(ToResponse));
        return Task.FromResult(response);
    }

    public override Task<ProductResponse> CreateProduct(CreateProductRequest request, ServerCallContext context)
    {
        var product = _productService.InsertProduct(request.Name, request.Description, Convert.ToDecimal(request.Price), request.StockQuantity);
        return Task.FromResult(ToResponse(product));
    }

    public override Task<ProductResponse> UpdateProduct(UpdateProductRequest request, ServerCallContext context)
    {
        var product = _productService.UpdateProduct(request.Id, request.Name, request.Description, Convert.ToDecimal(request.Price), request.StockQuantity);
        return Task.FromResult(ToResponse(product));
    }

    public override Task<DeleteProductResponse> DeleteProduct(DeleteProductRequest request, ServerCallContext context)
    {
        return Task.FromResult(new DeleteProductResponse
        {
            Deleted = _productService.DeleteProduct(request.Id)
        });
    }

    public override Task<DeleteAllProductsResponse> DeleteAllProducts(DeleteAllProductsRequest request, ServerCallContext context)
    {
        return Task.FromResult(new DeleteAllProductsResponse
        {
            Deleted = _productService.DeleteProduct()
        });
    }

    private static ProductResponse ToResponse(Product product)
    {
        return new ProductResponse
        {
            Id = product.Id,
            Name = product.Name,
            Description = product.Description,
            Price = Convert.ToDouble(product.Price),
            StockQuantity = product.StockQuantity,
            CreatedAt = product.CreatedAt.ToUniversalTime().ToString("O"),
            UpdatedAt = product.UpdatedAt?.ToUniversalTime().ToString("O") ?? string.Empty
        };
    }
}
