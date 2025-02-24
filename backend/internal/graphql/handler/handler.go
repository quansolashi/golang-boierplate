package handler

import (
	"entgo.io/contrib/entgql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/quansolashi/golang-boierplate/backend/ent"
	bgraph "github.com/quansolashi/golang-boierplate/backend/graph"
	egraph "github.com/quansolashi/golang-boierplate/backend/internal/entity/graph"
	"github.com/quansolashi/golang-boierplate/backend/pkg/auth"
	"github.com/vektah/gqlparser/v2/ast"
)

type Graph interface {
	Handler(rg *gin.RouterGroup)
}

type graph struct {
	resolver bgraph.Resolver
	auth     auth.LocalClient
}

type Params struct {
	LocalTokenSecret string
	Ent              *ent.Client
}

func NewGraph(params *Params) Graph {
	resolver := bgraph.NewResolver(&bgraph.Params{
		Client: params.Ent,
	})

	return &graph{
		resolver: resolver,
		auth:     auth.NewLocalClient(params.LocalTokenSecret),
	}
}

func (g *graph) Handler(rg *gin.RouterGroup) {
	rg.Use(egraph.TokenAuthMiddleware(g.auth))

	rg.POST("", graphqlHandler(g.resolver))
	rg.GET("/playground", playgroundHandler())
}

func graphqlHandler(resolver bgraph.Resolver) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.New(bgraph.NewExecutableSchema(bgraph.Config{Resolvers: &resolver}))

	// Server setup:
	h.AddTransport(transport.Options{})
	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})

	h.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	h.Use(extension.Introspection{})
	h.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})
	h.Use(entgql.Transactioner{
		TxOpener: resolver.Client,
	})

	return func(c *gin.Context) {
		c.Request.Header.Set("Content-Type", "application/json")
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/graphql")

	return func(c *gin.Context) {
		c.Request.Header.Set("Content-Type", "application/json")
		h.ServeHTTP(c.Writer, c.Request)
	}
}
