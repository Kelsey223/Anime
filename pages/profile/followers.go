package profile

import (
	"net/http"
	"sort"

	"github.com/aerogo/aero"
	"github.com/animenotifier/arn"
	"github.com/animenotifier/notify.moe/components"
	"github.com/animenotifier/notify.moe/utils"
)

// GetFollowers shows the followers of a particular user.
func GetFollowers(ctx *aero.Context) string {
	nick := ctx.Get("nick")
	viewUser, err := arn.GetUserByNick(nick)

	if err != nil {
		return ctx.Error(http.StatusNotFound, "User not found", err)
	}

	followers := viewUser.Followers()

	sort.Slice(followers, func(i, j int) bool {
		return followers[i].LastSeen > followers[j].LastSeen
	})

	return ctx.HTML(components.ProfileFollowers(followers, viewUser, utils.GetUser(ctx), ctx.URI()))

}