package main

import (
	"alphaflow/models"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/attache/attache"
)

// Guard API
func (c *AlphaFlow) GUARD_API() {
	u, err := c.getUser()
	if err != nil {
		log.Println("authorization:", err)
		attache.ErrorMessageJSON(http.StatusUnauthorized, "invalid auth")
	}
	c.User = u
}

func (c *AlphaFlow) GET_APIValidpairs() {
	// TODO: implement retry logic, caching, etc
	pairs, err := getValidPairs()
	if err != nil {
		attache.ErrorFatal(err)
	}
	attache.RenderJSON(c.ResponseWriter(), pairs)
}

func (c *AlphaFlow) POST_APISubscriptionsCreate() {
	var body struct {
		Pair string `json:"pair"`
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		log.Println(err)
		attache.ErrorMessageJSON(http.StatusBadRequest, "unable to parse request body")
	}
	pairs, err := getValidPairs()
	if err != nil {
		attache.ErrorFatal(err)
	}

	// TODO: cache this list somewhere and use an O(1) map lookup
	// Ensure pair is valid
	isValid := false
	for _, p := range pairs {
		if p == body.Pair {
			isValid = true
			break
		}
	}
	if !isValid {
		attache.ErrorMessageJSON(http.StatusUnprocessableEntity, "pair %q is invalid", body.Pair)
	}

	var existing models.Subscription
	// TODO: make pair unique regardless of order (?? maybe)
	if err := c.DB().GetBy(&existing, "user_id = ? AND pair = ?", c.User.ID, body.Pair); err == nil {
		// TODO: perhaps fail early if we saw database errors
		// A nil error indicates that we found the record, so let's make that a no-op
		attache.RenderJSON(c.ResponseWriter(), existing)
	}

	newSub := models.Subscription{
		UserID: c.User.ID,
		Pair:   body.Pair,
	}

	if err := c.DB().Insert(&newSub); err != nil {
		attache.ErrorFatal(err)
	}

	attache.RenderJSON(c.ResponseWriter(), newSub)
}

func (c *AlphaFlow) GET_APISubscriptionsList() {
	list, err := c.DB().Where(&models.Subscription{}, "user_id = ?", c.User.ID)
	if err != nil {
		if err == attache.ErrRecordNotFound {
			attache.RenderJSON(c.ResponseWriter(), []interface{}{})
		}
		attache.ErrorFatal(err)
	}
	attache.RenderJSON(c.ResponseWriter(), list)
}

func (c *AlphaFlow) GET_APISubscriptionsLimit() {
	list, err := c.DB().Where(&models.Subscription{}, "user_id = ?", c.User.ID)
	if err != nil {
		if err == attache.ErrRecordNotFound {
			attache.RenderJSON(c.ResponseWriter(), []interface{}{})
		}
		attache.ErrorFatal(err)
	}
	offline, err := getOfflineCoins()
	if err != nil {
		attache.ErrorFatal(err)
	}

	limits := make([]LimitInfo, 0, len(list))
	for _, r := range list {
		// filter out pairs containing offline coins
		pair := r.(*models.Subscription).Pair
		if pairIsOffline(offline, pair) {
			continue
		}
		// TODO: avoid fan-out by caching and/or performing requests in parallel
		limit, err := getLimitInfo(pair)
		if err != nil {
			// TODO: do we want to just exclude ones we can't get limit info for?
			attache.ErrorFatal(err)
		}
		limits = append(limits, limit)
	}
	attache.RenderJSON(c.ResponseWriter(), limits)
}

func getValidPairs() ([]string, error) {
	resp, err := http.Get(`https://shapeshift.io/validpairs`)
	if err != nil {
		return nil, err
	}
	var results []string
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// LimitInfo is limit information from the ShapeShift API
type LimitInfo struct {
	Pair  string  `json:"pair"`
	Rate  float64 `json:"rate"`
	Limit float64 `json:"limit"`
	Min   float64 `json:"minimum"`
}

func getLimitInfo(pair string) (LimitInfo, error) {
	var result LimitInfo
	// TODO: validade pair is a valid pair (skipping for brevity, we ensure subscriptions are for valid pairs)
	resp, err := http.Get(`https://shapeshift.io/marketinfo/` + pair)
	if err != nil {
		return result, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, err
	}
	return result, nil
}

func getOfflineCoins() (map[string]bool, error) {
	var list []string
	resp, err := http.Get(`https://shapeshift.io/offlinecoins`)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, err
	}
	result := map[string]bool{}
	for _, name := range list {
		result[name] = true
	}
	return result, nil
}

func pairIsOffline(offline map[string]bool, pair string) bool {
	for _, name := range strings.Split(pair, "_") {
		// if any items in the pair are offline, the pair is offline
		if offline[name] {
			return true
		}
	}
	return false
}
