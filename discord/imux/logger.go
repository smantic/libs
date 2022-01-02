package imux

import (
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var base zerolog.Logger

func init() {
	base = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

// LogInteraction is a middleware to print structured log lines for interactions
//
// use "github.com/rs/zerolog/log".Ctx(request.Context) to get the logger.
//
func LogInteraction(next InteractionHandler) InteractionHandler {

	m := func(response *discordgo.InteractionResponse, request *InteractionRequest) {

		request.Context = base.WithContext(request.Context)

		start := time.Now()

		next(response, request)

		dur := time.Since(start)
		e := log.Ctx(request.Context).
			Info().
			Dur("interaction duration", dur).
			Str("interaction type", iTypeString(response.Type)).
			Int("content-length", len(response.Data.Content))

		if d := response.Data; d != nil {
			e.Int("num choices", len(d.Choices)).
				Int("num components", len(d.Components)).
				Int("num embeds", len(d.Embeds))
		}
		e.Msg("interaction finished")
	}

	return m
}

// iTypeString  gives the strng representation of the interaction type.
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object-interaction-callback-type
func iTypeString(i discordgo.InteractionResponseType) string {

	switch i {
	case 1:
		return "PONG"
	case 4:
		return "CHANNEL_MESSAGE_WITH_SOURCE"
	case 5:
		return "DEFERRED_CHANNEL_MESSAGE_WITH_SOURCE"
	case 6:
		return "DEFERRED_UPDATE_MESSAGE"
	case 7:
		return "UPDATE_MESSAGE"
	case 8:
		return "APPLICATION_COMMAND_AUTOCOMPLETE_RESULT"
	default:
		return ""
	}
}
