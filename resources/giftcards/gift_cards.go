package giftcards

import (
	"context"
	"fmt"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// GiftCards is the gift card endpoints, bound to one client.
type GiftCards struct {
	requester http.Requester
}

// ListGiftCards lists gift cards, newest first — GET /admin/gift_cards; code matches a whole code and email a whole address, in any case; is_active takes true, false, 1 or 0, and limit caps at 20 with offset or page to move on. Any other name is 400 unknown_parameter.
func (g *GiftCards) ListGiftCards(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return g.requester.SendGetRequest(ctx, GiftCardsPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: GiftCardPageShape})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "gift_cards"

func oneResponseFormat(responseFormat []string) (string, error) {
	switch len(responseFormat) {
	case 0:
		return "", nil
	case 1:
		return responseFormat[0], nil
	default:
		return "", wcerrors.NewUsageError(fmt.Sprintf(
			"responseFormat: pass one format or none, got %d — the trailing argument is optional, "+
				"and omitting it defers to the client", len(responseFormat)))
	}
}

// New binds the gift card endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *GiftCards {
	return &GiftCards{requester: requester}
}
