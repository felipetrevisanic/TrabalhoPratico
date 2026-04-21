FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /src

COPY services/csharp/rest/src/src.csproj services/csharp/rest/src/
RUN dotnet restore services/csharp/rest/src/src.csproj

COPY services/csharp/rest/src/ services/csharp/rest/src/
RUN dotnet publish services/csharp/rest/src/src.csproj -c Release -o /app/publish /p:UseAppHost=false

FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS runtime
WORKDIR /app

ENV ASPNETCORE_URLS=http://+:5107
COPY --from=build /app/publish .

ENTRYPOINT ["dotnet", "src.dll"]
