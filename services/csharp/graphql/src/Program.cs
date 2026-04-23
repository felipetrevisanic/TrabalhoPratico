using Microsoft.EntityFrameworkCore;
using src.application.interfaces;
using src.application.service;
using src.domain.interfaces;
using src.infrastructure.Data;
using src.infrastructure.Repositories;
using src.interfaces.graphql.inputs;
using src.interfaces.graphql.mutations;
using src.interfaces.graphql.queries;

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
