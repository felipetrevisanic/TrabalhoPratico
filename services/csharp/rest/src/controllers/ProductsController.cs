using Microsoft.AspNetCore.Mvc;

namespace src.controllers;

[ApiController]
[Route("[controller]")]
public class ProductController : ControllerBase
{
    [HttpGet]
    public getProductById([FromParam] int id)
    {
        
    }   

    [HttpGet]
    public getAllProduct()
    {
        //Criar com paginação???        
    }


    [HttpPost]
    public insertProduct()
    {
        
    }

    [HttpPut]
    public putProduct()
    {
        
    }

    
   
}
