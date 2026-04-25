namespace src.interfaces.rest.requests;

public class CreateProductRequestDto
{
    public string Name { get; set; } = string.Empty;
    public string Description { get; set; } = string.Empty;
    public string Category { get; set; } = string.Empty;
    public string[] Images { get; set; } = Array.Empty<string>();
    public decimal Price { get; set; }
    public int StockQuantity { get; set; }
}
