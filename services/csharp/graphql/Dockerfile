FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /src

COPY services/csharp/graphql/src/graphql.csproj services/csharp/graphql/src/
RUN dotnet restore services/csharp/graphql/src/graphql.csproj

COPY services/csharp/graphql/src/ services/csharp/graphql/src/
RUN dotnet publish services/csharp/graphql/src/graphql.csproj -c Release -o /app/publish /p:UseAppHost=false

FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS runtime
WORKDIR /app

ENV ASPNETCORE_URLS=http://+:5207
COPY --from=build /app/publish .

ENTRYPOINT ["dotnet", "graphql.dll"]
