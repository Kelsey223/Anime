package listimportkitsu

import (
	"net/http"

	"github.com/aerogo/aero"
	"github.com/animenotifier/arn"
	"github.com/animenotifier/kitsu"
	"github.com/animenotifier/notify.moe/components"
	"github.com/animenotifier/notify.moe/utils"
)

// Preview shows an import preview.
func Preview(ctx *aero.Context) string {
	user := utils.GetUser(ctx)

	if user == nil {
		return ctx.Error(http.StatusBadRequest, "Not logged in", nil)
	}

	matches, response := getMatches(ctx)

	if response != "" {
		return response
	}

	return ctx.HTML(components.ImportKitsu(user, matches))
}

// Finish ...
func Finish(ctx *aero.Context) string {
	user := utils.GetUser(ctx)

	if user == nil {
		return ctx.Error(http.StatusBadRequest, "Not logged in", nil)
	}

	matches, response := getMatches(ctx)

	if response != "" {
		return response
	}

	animeList := user.AnimeList()

	for _, match := range matches {
		if match.ARNAnime == nil || match.KitsuItem == nil {
			continue
		}

		rating := match.KitsuItem.Attributes.RatingTwenty

		if rating < 2 {
			rating = 2
		}

		if rating > 20 {
			rating = 20
		}

		// Convert rating
		convertedRating := (float64(rating-2) / 18.0) * 10.0

		item := &arn.AnimeListItem{
			AnimeID:  match.ARNAnime.ID,
			Status:   arn.KitsuStatusToARNStatus(match.KitsuItem.Attributes.Status),
			Episodes: match.KitsuItem.Attributes.Progress,
			Notes:    match.KitsuItem.Attributes.Notes,
			Rating: &arn.AnimeRating{
				Overall: convertedRating,
			},
			RewatchCount: match.KitsuItem.Attributes.ReconsumeCount,
			Created:      arn.DateTimeUTC(),
			Edited:       arn.DateTimeUTC(),
		}

		animeList.Import(item)
	}

	err := animeList.Save()

	if err != nil {
		return ctx.Error(http.StatusInternalServerError, "Error saving your anime list", err)
	}

	return ctx.Redirect("/+" + user.Nick + "/animelist")
}

// getMatches finds and returns all matches for the logged in user.
func getMatches(ctx *aero.Context) ([]*arn.KitsuMatch, string) {
	user := utils.GetUser(ctx)

	if user == nil {
		return nil, ctx.Error(http.StatusBadRequest, "Not logged in", nil)
	}

	kitsuUser, err := kitsu.GetUser(user.Accounts.Kitsu.Nick)

	if err != nil {
		return nil, ctx.Error(http.StatusBadRequest, "Couldn't load your user info from Kitsu", err)
	}

	library := kitsuUser.StreamLibraryEntries()
	matches := findAllMatches(library)

	return matches, ""
}

// findAllMatches returns all matches for the anime inside an anilist anime list.
func findAllMatches(library chan *kitsu.LibraryEntry) []*arn.KitsuMatch {
	matches := []*arn.KitsuMatch{}

	for item := range library {
		// Ignore non-anime entries
		if item.Anime == nil {
			continue
		}

		var anime *arn.Anime
		connection, err := arn.GetKitsuToAnime(item.Anime.ID)

		if err == nil {
			anime, _ = arn.GetAnime(connection.AnimeID)
		}

		matches = append(matches, &arn.KitsuMatch{
			KitsuItem: item,
			ARNAnime:  anime,
		})
	}

	return matches
}