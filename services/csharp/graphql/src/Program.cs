using Microsoft.EntityFrameworkCore;
using src.application.interfaces;
using src.application.service;
using src.domain.interfaces;
using src.graphql.inputs;
using src.graphql.mutations;
using src.graphql.queries;
using src.infraestructure.Data;
using src.infraestructure.Repositories;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddDbContext<AppDbContext>(options =>
    options.UseNpgsql(builder.Configuration.GetConnectionString("DefaultConnection")));

builder.Services.AddScoped<IProductService, ProductService>();
builder.Services.AddScoped<IProductRepository, ProductRepository>();

builder.Services
    .AddGraphQLServer()
    .AddQueryType<ProductQuery>()
    .AddMutationType<ProductMutation>()
    .AddType<CreateProductInputType>()
    .AddType<UpdateProductInputType>();

var app = builder.Build();

app.MapGraphQL("/graphql");

app.Run();
