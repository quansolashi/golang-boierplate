package cmd

import (
	"context"

	web "github.com/quansolashi/message-extractor/backend/internal/web/controller"
)

func (a *app) inject(ctx context.Context) error {
	a.web = web.NewController(&web.Params{})

	return nil
}
