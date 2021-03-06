package anime

import (
	"net/http"

	"github.com/animenotifier/notify.moe/components"
	"github.com/animenotifier/notify.moe/utils"

	"github.com/aerogo/aero"
	"github.com/animenotifier/notify.moe/arn"
)

// Comments ...
func Comments(ctx aero.Context) error {
	user := utils.GetUser(ctx)
	id := ctx.Get("id")
	anime, err := arn.GetAnime(id)

	if err != nil {
		return ctx.Error(http.StatusNotFound, "Anime not found", err)
	}

	return ctx.HTML(components.AnimeComments(anime, user, true))
}
