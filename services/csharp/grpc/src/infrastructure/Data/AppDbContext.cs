using Microsoft.EntityFrameworkCore;
using src.domain.entities;

namespace src.infrastructure.Data;

public class AppDbContext : DbContext
{
    public AppDbContext(DbContextOptions<AppDbContext> options) : base(options)
    {
    }

    public DbSet<Product> Products => Set<Product>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.Entity<Product>(entity =>
        {
            entity.ToTable("products");
            entity.HasKey(product => product.Id);

            entity.Property(product => product.Id)
                .ValueGeneratedOnAdd();

            entity.Property(product => product.Name)
                .HasMaxLength(150)
                .IsRequired();

            entity.Property(product => product.Description)
                .HasMaxLength(500)
                .IsRequired();

            entity.Property(product => product.Price)
                .HasPrecision(18, 2)
                .IsRequired();

            entity.Property(product => product.StockQuantity)
                .IsRequired();

            entity.Property(product => product.CreatedAt)
                .IsRequired();
        });
    }
}
