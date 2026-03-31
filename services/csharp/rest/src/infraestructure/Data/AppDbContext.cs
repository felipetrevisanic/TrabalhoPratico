using Microsoft.EntityFrameworkCore;
using src.domain.entities;

namespace src.infraestructure.Data;

public class AppDbContext : DbContext
{
    public AppDbContext(DbContextOptions<AppDbContext> options) : base(options)
    {
    }

    public DbSet<Product> Products => Set<Product>();

}
