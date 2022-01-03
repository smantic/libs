package imux

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InteractionMux handles interaction commands.
type InteractionMux struct {

	// starting ctx
	parentCtx context.Context

	//middleware stack
	middleware []InteractionMiddleware

	mux map[string]InteractionHandler
}

func NewInteractionMux(ctx context.Context) InteractionMux {
	if ctx == nil {
		ctx = context.Background()
	}

	return InteractionMux{parentCtx: ctx}
}

// Use appends a middleware to the middleware stack
func (i *InteractionMux) Use(middlewares ...InteractionMiddleware) {

	for _, m := range middlewares {
		if m == nil {
			panic("nil interaction middleware")
		}
	}

	i.middleware = append(i.middleware, middlewares...)
}

func (i *InteractionMux) Add(command string, handler InteractionHandler) {

	if command == "" {
		return
	}

	if handler == nil {
		panic("nil interaction handler")
	}

	h := chain(i.middleware, handler)

	i.mux[command] = h
}

// Serve will serve interaction responses
// to be used by discordgo session.AddHandler()
func (m *InteractionMux) Serve() func(s *discordgo.Session, i *discordgo.InteractionCreate) {

	ctx := m.parentCtx
	f := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		cmd := i.ApplicationCommandData()

		// 3 second timeout for an inital response.
		ctx, cancel := context.WithTimeout(ctx, time.Second*3)
		request := &InteractionRequest{
			Session:     s,
			Interaction: i,
			Context:     ctx,
		}

		response := &discordgo.InteractionResponse{}

		handle, ok := m.mux[cmd.Name]
		if !ok {
			m.notFoundHandler(response, request)
			return
		}

		handle(response, request)

		cancel()
	}

	return f
}

// Respond will take the response and send it to discord.
func Respond(response *discordgo.InteractionResponse, request *InteractionRequest) {

	// send the response
	err := request.Session.InteractionRespond(request.Interaction.Interaction, response)
	if err != nil {
		log.Ctx(request.Context).UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.AnErr("discord", err)
		})
	}
}

func (i *InteractionMux) notFoundHandler(response *discordgo.InteractionResponse, request *InteractionRequest) {

	endpoint := func(response *discordgo.InteractionResponse, request *InteractionRequest) {

		c := fmt.Sprintf(
			"command (%s) not found.",
			request.Interaction.ApplicationCommandData().Name,
		)

		response = &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: c,
			},
		}

		Respond(response, request)
	}

	h := chain(i.middleware, endpoint)

	h(response, request)
}
