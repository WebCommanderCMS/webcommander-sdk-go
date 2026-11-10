package giftcards

// GiftCardsPath is the path template; the client fills {api_version}.
const GiftCardsPath = "/api/{api_version}/admin/gift_cards"

// GiftCardsKey is the key the listing's rows arrive under.
const GiftCardsKey = "gift_cards"

// GiftCardsPageSizeCap is the largest number limit takes; a larger one returns that many, and 0 is 400.
const GiftCardsPageSizeCap = 20
