package imux

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

type InteractionRequest struct {
	session     *discordgo.Session
	interaction *discordgo.InteractionCreate
	Context     context.Context
}

// InteractionHandler is standard signature for interaction handler functions.
// should feel similar to http.Handler.
type InteractionHandler = func(*discordgo.InteractionResponse, *InteractionRequest)
type InteractionMiddleware = func(InteractionHandler) InteractionHandler

func chain(mw []InteractionMiddleware, endpoint InteractionHandler) InteractionHandler {

	if len(mw) == 0 {
		return endpoint
	}

	h := endpoint
	// applys middleware in order
	for i := len(mw) - 1; i >= 0; i-- {
		h = mw[i](h)
	}

	return h
}
