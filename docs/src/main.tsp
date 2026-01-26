import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";
import "./models.tsp";
import "./operations.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

@service(#{
  title: "Go Boilerplate API"
})
@server("http://localhost:8080", "Servidor de desenvolvimento")
namespace GrowthAPI;
