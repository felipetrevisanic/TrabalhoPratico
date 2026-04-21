FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /src

COPY services/csharp/grpc/src/grpc.csproj services/csharp/grpc/src/
COPY shared/proto/ shared/proto/
RUN dotnet restore services/csharp/grpc/src/grpc.csproj

COPY services/csharp/grpc/src/ services/csharp/grpc/src/
RUN dotnet publish services/csharp/grpc/src/grpc.csproj -c Release -o /app/publish /p:UseAppHost=false

FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS runtime
WORKDIR /app

ENV GRPC_PORT=5307
COPY --from=build /app/publish .

ENTRYPOINT ["dotnet", "grpc.dll"]
