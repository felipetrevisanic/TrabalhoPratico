using Microsoft.AspNetCore.Server.Kestrel.Core;
using Microsoft.EntityFrameworkCore;
using src.application.interfaces;
using src.application.service;
using src.domain.interfaces;
using src.infraestructure.Data;
using src.infraestructure.Repositories;
using src.services;

var builder = WebApplication.CreateBuilder(args);
var grpcPort = builder.Configuration.GetValue("GRPC_PORT", 5307);

builder.WebHost.ConfigureKestrel(options =>
{
    options.ListenAnyIP(grpcPort, listenOptions =>
    {
        listenOptions.Protocols = HttpProtocols.Http2;
    });
});

builder.Services.AddDbContext<AppDbContext>(options =>
    options.UseNpgsql(builder.Configuration.GetConnectionString("DefaultConnection")));

builder.Services.AddScoped<IProductService, ProductService>();
builder.Services.AddScoped<IProductRepository, ProductRepository>();
builder.Services.AddGrpc();
builder.Services.AddGrpcReflection();

var app = builder.Build();

app.MapGrpcService<ProductGrpcService>();
app.MapGrpcReflectionService();
app.MapGet("/", () => "Use a gRPC client to communicate with this service.");

app.Run();
