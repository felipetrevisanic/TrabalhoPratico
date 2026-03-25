using Microsoft.AspNetCore.Mvc;
using src.application.interfaces;

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
    public ActionResult<string> GetProductById([FromQuery] int id)
    {
        var product = _productService.GetProductById(id);
        return Ok(product);
    }

    [HttpGet("all")]
    public ActionResult<IEnumerable<string>> GetAllProduct()
    {
        var products = _productService.GetAllProducts();
        return Ok(products);
    }

    [HttpPost]
    public ActionResult<string> InsertProduct()
    {
        var result = _productService.InsertProduct();
        return Ok(result);
    }

    [HttpPut]
    public ActionResult<string> UpdateProduct()
    {
        var result = _productService.UpdateProduct();
        return Ok(result);
    }
}
