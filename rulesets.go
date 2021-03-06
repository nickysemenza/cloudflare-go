package cloudflare

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const (
	RulesetKindCustom  RulesetKind = "custom"
	RulesetKindManaged RulesetKind = "managed"
	RulesetKindRoot    RulesetKind = "root"
	RulesetKindSchema  RulesetKind = "schema"
	RulesetKindZone    RulesetKind = "zone"

	RulesetPhaseDDoSL7                     RulesetPhase = "ddos_l7"
	RulesetPhaseHTTPRequestFirewallCustom  RulesetPhase = "http_request_firewall_custom"
	RulesetPhaseHTTPRequestFirewallManaged RulesetPhase = "http_request_firewall_managed"
	RulesetPhaseHTTPRequestMain            RulesetPhase = "http_request_main"
	RulesetPhaseHTTPRequestSanitize        RulesetPhase = "http_request_sanitize"
	RulesetPhaseHTTPRequestTransform       RulesetPhase = "http_request_transform"
	RulesetPhaseMagicTransit               RulesetPhase = "magic_transit"

	RulesetRuleActionBlock                RulesetRuleAction = "block"
	RulesetRuleActionChallenge            RulesetRuleAction = "challenge"
	RulesetRuleActionDDoSDynamic          RulesetRuleAction = "ddos_dynamic"
	RulesetRuleActionExecute              RulesetRuleAction = "execute"
	RulesetRuleActionForceConnectionClose RulesetRuleAction = "force_connection_close"
	RulesetRuleActionJSChallenge          RulesetRuleAction = "js_challenge"
	RulesetRuleActionLog                  RulesetRuleAction = "log"
	RulesetRuleActionRewrite              RulesetRuleAction = "rewrite"
	RulesetRuleActionScore                RulesetRuleAction = "score"
	RulesetRuleActionSkip                 RulesetRuleAction = "skip"

	RulesetActionParameterProductBIC           RulesetActionParameterProduct = "bic"
	RulesetActionParameterProductHOT           RulesetActionParameterProduct = "hot"
	RulesetActionParameterProductRateLimit     RulesetActionParameterProduct = "ratelimit"
	RulesetActionParameterProductSecurityLevel RulesetActionParameterProduct = "securityLevel"
	RulesetActionParameterProductUABlock       RulesetActionParameterProduct = "uablock"
	RulesetActionParameterProductWAF           RulesetActionParameterProduct = "waf"
	RulesetActionParameterProductZoneLockdown  RulesetActionParameterProduct = "zonelockdown"

	RulesetRuleActionParametersHTTPHeaderOperationRemove RulesetRuleActionParametersHTTPHeaderOperation = "remove"
	RulesetRuleActionParametersHTTPHeaderOperationSet    RulesetRuleActionParametersHTTPHeaderOperation = "set"
)

// RulesetRuleAction defines a custom type that is used to express allowed
// values for the rule action.
type RulesetRuleAction string

// RulesetKind is the custom type for allowed variances of rulesets.
type RulesetKind string

// RulesetPhase is the custom type for defining at what point the ruleset will
// be applied in the request pipeline.
type RulesetPhase string

type RulesetActionParameterProduct string

// RulesetRuleActionParametersHTTPHeaderOperation defines available options for
// HTTP header operations in actions.
type RulesetRuleActionParametersHTTPHeaderOperation string

// Ruleset contains the structure of a Ruleset.
type Ruleset struct {
	ID          string        `json:"id,omitempty"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Kind        RulesetKind   `json:"kind"`
	Version     string        `json:"version,omitempty"`
	LastUpdated *time.Time    `json:"last_updated,omitempty"`
	Phase       RulesetPhase  `json:"phase"`
	Rules       []RulesetRule `json:"rules"`
}

// RulesetRuleActionParameters specifies the action parameters for a Ruleset
// rule.
type RulesetRuleActionParameters struct {
	ID        string                                           `json:"id,omitempty"`
	Ruleset   string                                           `json:"ruleset,omitempty"`
	Increment int                                              `json:"increment,omitempty"`
	URI       RulesetRuleActionParametersURI                   `json:"uri,omitempty"`
	Headers   map[string]RulesetRuleActionParametersHTTPHeader `json:"headers,omitempty"`
	Products  []RulesetActionParameterProduct                  `json:"products,omitempty"`
}

// RulesetRuleActionParametersURI holds the URI struct for an action parameter.
type RulesetRuleActionParametersURI struct {
	Path   RulesetRuleActionParametersURIPath  `json:"path,omitempty"`
	Query  RulesetRuleActionParametersURIQuery `json:"query,omitempty"`
	Origin bool                                `json:"origin,omitempty"`
}

// RulesetRuleActionParametersURIPath holds the path specific portion of a URI
// action parameter.
type RulesetRuleActionParametersURIPath struct {
	Expression string `json:"expression,omitempty"`
}

// RulesetRuleActionParametersURIQuery holds the query specific portion of a URI
// action parameter.
type RulesetRuleActionParametersURIQuery struct {
	Value      string `json:"value,omitempty"`
	Expression string `json:"expression,omitempty"`
}

// RulesetRuleActionParametersHTTPHeader is the definition for define action
// parameters that involve HTTP headers.
type RulesetRuleActionParametersHTTPHeader struct {
	Operation  string `json:"operation,omitempty"`
	Value      string `json:"value,omitempty"`
	Expression string `json:"expression,omitempty"`
}

// RulesetRule contains information about a single Ruleset Rule.
type RulesetRule struct {
	ID               string                       `json:"id,omitempty"`
	Version          string                       `json:"version,omitempty"`
	Action           RulesetRuleAction            `json:"action"`
	ActionParameters *RulesetRuleActionParameters `json:"action_parameters,omitempty"`
	Expression       string                       `json:"expression"`
	Description      string                       `json:"description"`
	LastUpdated      *time.Time                   `json:"last_updated,omitempty"`
	Ref              string                       `json:"ref,omitempty"`
	Enabled          bool                         `json:"enabled"`
	Categories       []string                     `json:"categories,omitempty"`
	ScoreThreshold   int                          `json:"score_threshold,omitempty"`
}

// UpdateRulesetRequest is the representation of a Ruleset update.
type UpdateRulesetRequest struct {
	Description string        `json:"description"`
	Rules       []RulesetRule `json:"rules"`
}

// ListRulesetResponse contains all Rulesets.
type ListRulesetResponse struct {
	Response
	Result []Ruleset `json:"result"`
}

// GetRulesetResponse contains a single Ruleset.
type GetRulesetResponse struct {
	Response
	Result Ruleset `json:"result"`
}

// CreateRulesetResponse contains response data when creating a new Ruleset.
type CreateRulesetResponse struct {
	Response
	Result Ruleset `json:"result"`
}

// UpdateRulesetResponse contains response data when updating an existing
// Ruleset.
type UpdateRulesetResponse struct {
	Response
	Result Ruleset `json:"result"`
}

// ListZoneRulesets fetches all rulesets for a zone.
//
// API reference: https://api.cloudflare.com/#zone-rulesets-list-zone-rulesets
func (api *API) ListZoneRulesets(ctx context.Context, zoneID string) ([]Ruleset, error) {
	return api.listRulesets(ctx, ZoneRouteRoot, zoneID)
}

// ListAccountRulesets fetches all rulesets for an account.
//
// API reference: https://api.cloudflare.com/#account-rulesets-list-account-rulesets
func (api *API) ListAccountRulesets(ctx context.Context, accountID string) ([]Ruleset, error) {
	return api.listRulesets(ctx, AccountRouteRoot, accountID)
}

// listRulesets lists all Rulesets for a given zone or account depending on the
// identifier type provided.
func (api *API) listRulesets(ctx context.Context, identifierType RouteRoot, identifier string) ([]Ruleset, error) {
	uri := fmt.Sprintf("/%s/%s/rulesets", identifierType, identifier)

	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return []Ruleset{}, err
	}

	result := ListRulesetResponse{}
	if err := json.Unmarshal(res, &result); err != nil {
		return []Ruleset{}, errors.Wrap(err, errUnmarshalError)
	}

	return result.Result, nil
}

// GetZoneRuleset fetches a single ruleset for a zone.
//
// API reference: https://api.cloudflare.com/#zone-rulesets-get-a-zone-ruleset
func (api *API) GetZoneRuleset(ctx context.Context, zoneID, rulesetID string) (Ruleset, error) {
	return api.getRuleset(ctx, ZoneRouteRoot, zoneID, rulesetID)
}

// GetAccountRuleset fetches a single ruleset for an account.
//
// API reference: https://api.cloudflare.com/#account-rulesets-get-an-account-ruleset
func (api *API) GetAccountRuleset(ctx context.Context, accountID, rulesetID string) (Ruleset, error) {
	return api.getRuleset(ctx, AccountRouteRoot, accountID, rulesetID)
}

// getRuleset fetches a single ruleset based on the zone or account, the
// identifer and the ruleset ID.
func (api *API) getRuleset(ctx context.Context, identifierType RouteRoot, identifier, rulesetID string) (Ruleset, error) {
	uri := fmt.Sprintf("/%s/%s/rulesets/%s", identifierType, identifier, rulesetID)
	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return Ruleset{}, err
	}

	result := GetRulesetResponse{}
	if err := json.Unmarshal(res, &result); err != nil {
		return Ruleset{}, errors.Wrap(err, errUnmarshalError)
	}

	return result.Result, nil
}

// CreateZoneRuleset creates a new ruleset for a zone.
//
// API reference: https://api.cloudflare.com/#zone-rulesets-create-zone-ruleset
func (api *API) CreateZoneRuleset(ctx context.Context, zoneID string, ruleset Ruleset) (Ruleset, error) {
	return api.createRuleset(ctx, ZoneRouteRoot, zoneID, ruleset)
}

// CreateAccountRuleset creates a new ruleset for an account.
//
// API reference: https://api.cloudflare.com/#account-rulesets-create-account-ruleset
func (api *API) CreateAccountRuleset(ctx context.Context, accountID string, ruleset Ruleset) (Ruleset, error) {
	return api.createRuleset(ctx, AccountRouteRoot, accountID, ruleset)
}

func (api *API) createRuleset(ctx context.Context, identifierType RouteRoot, identifier string, ruleset Ruleset) (Ruleset, error) {
	uri := fmt.Sprintf("/%s/%s/rulesets", identifierType, identifier)
	res, err := api.makeRequestContext(ctx, http.MethodPost, uri, ruleset)

	if err != nil {
		return Ruleset{}, err
	}

	result := CreateRulesetResponse{}
	if err := json.Unmarshal(res, &result); err != nil {
		return Ruleset{}, errors.Wrap(err, errUnmarshalError)
	}

	return result.Result, nil
}

// DeleteZoneRuleset deletes a single ruleset for a zone.
//
// API reference: https://api.cloudflare.com/#zone-rulesets-delete-zone-ruleset
func (api *API) DeleteZoneRuleset(ctx context.Context, zoneID, rulesetID string) error {
	return api.deleteRuleset(ctx, ZoneRouteRoot, zoneID, rulesetID)
}

// DeleteAccountRuleset deletes a single ruleset for an account.
//
// API reference: https://api.cloudflare.com/#account-rulesets-delete-account-ruleset
func (api *API) DeleteAccountRuleset(ctx context.Context, accountID, rulesetID string) error {
	return api.deleteRuleset(ctx, AccountRouteRoot, accountID, rulesetID)
}

// deleteRuleset removes a ruleset based on the ruleset ID.
func (api *API) deleteRuleset(ctx context.Context, identifierType RouteRoot, identifier, rulesetID string) error {
	uri := fmt.Sprintf("/%s/%s/rulesets/%s", identifierType, identifier, rulesetID)
	res, err := api.makeRequestContext(ctx, http.MethodDelete, uri, nil)

	if err != nil {
		return err
	}

	// The API is not implementing the standard response blob but returns an
	// empty response (204) in case of a success. So we are checking for the
	// response body size here.
	if len(res) > 0 {
		return errors.Wrap(errors.New(string(res)), errMakeRequestError)
	}

	return nil
}

// UpdateZoneRuleset updates a single ruleset for a zone.
//
// API reference: https://api.cloudflare.com/#zone-rulesets-update-a-zone-ruleset
func (api *API) UpdateZoneRuleset(ctx context.Context, zoneID, rulesetID, description string, rules []RulesetRule) (Ruleset, error) {
	return api.updateRuleset(ctx, ZoneRouteRoot, zoneID, rulesetID, description, rules)
}

// UpdateAccountRuleset updates a single ruleset for an account.
//
// API reference: https://api.cloudflare.com/#account-rulesets-update-account-ruleset
func (api *API) UpdateAccountRuleset(ctx context.Context, accountID, rulesetID, description string, rules []RulesetRule) (Ruleset, error) {
	return api.updateRuleset(ctx, AccountRouteRoot, accountID, rulesetID, description, rules)
}

// updateRuleset updates a ruleset based on the ruleset ID.
func (api *API) updateRuleset(ctx context.Context, identifierType RouteRoot, identifier, rulesetID, description string, rules []RulesetRule) (Ruleset, error) {
	uri := fmt.Sprintf("/%s/%s/rulesets/%s", identifierType, identifier, rulesetID)
	payload := UpdateRulesetRequest{Description: description, Rules: rules}
	res, err := api.makeRequestContext(ctx, http.MethodPut, uri, payload)
	if err != nil {
		return Ruleset{}, err
	}

	result := UpdateRulesetResponse{}
	if err := json.Unmarshal(res, &result); err != nil {
		return Ruleset{}, errors.Wrap(err, errUnmarshalError)
	}

	return result.Result, nil
}
