using Microsoft.AspNetCore.Mvc;
using src.application.interfaces;
using src.DTO.requests;
using src.DTO.response;

namespace src.controllers;

[ApiController]
[Route("[controller]")]
public class ProductController : ControllerBase
{
    private readonly IProductService _productService;

    public ProductController(IProductService productService)
    {
        _productService = productService;
    }

    [HttpGet]
    public ActionResult<ProductResponseDto> GetProductById([FromQuery] int id)
    {
        var product = _productService.GetProductById(id);
        return Ok(product);
    }

    [HttpGet("all")]
    public ActionResult<IEnumerable<ProductResponseDto>> GetAllProduct()
    {
        var products = _productService.GetAllProducts();
        return Ok(products);
    }

    [HttpPost]
    public ActionResult<ProductResponseDto> InsertProduct([FromBody] CreateProductRequestDto request)
    {
        var result = _productService.InsertProduct(request);
        return Ok(result);
    }

    [HttpPut("{id:int}")]
    public ActionResult<ProductResponseDto> UpdateProduct(int id, [FromBody] UpdateProductRequestDto request)
    {
        var result = _productService.UpdateProduct(id, request);
        return Ok(result);
    }

    [HttpDelete("{id:int}")]
    public IActionResult DeleteProduct(int id)
    {
        var deleted = _productService.DeleteProduct(id);

        if (!deleted)
        {
            return NotFound();
        }

        return NoContent();
    }
}
