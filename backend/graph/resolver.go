package graph

import (
	"github.com/quansolashi/golang-boierplate/backend/ent"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

func NewResolver(params *Params) Resolver {
	return Resolver{
		Client: params.Client,
	}
}

type Params struct {
	Client *ent.Client
}

type Resolver struct {
	Client *ent.Client
}
