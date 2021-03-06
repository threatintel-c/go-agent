// Copyright (c) 2016 - 2019 Sqreen. All Rights Reserved.
// Please refer to our terms for more information:
// https://www.sqreen.io/terms.html

package actor

import (
	"net/url"
	"time"

	"github.com/sqreen/go-agent/internal/sqlib/sqerrors"
)

// Action kinds.
const (
	actionKindBlockIP      = "block_ip"
	actionKindBlockUser    = "block_user"
	actionKindRedirectIP   = "redirect_ip"
	actionKindRedirectUser = "redirect_user"
)

// Action is an interface common to each concrete action type stored in the data
// structures, and allowing to type-switch the stored values.
type Action interface {
	// ActionID returns the unique ID of the request.
	ActionID() string
}

type RedirectAction interface {
	Action
	RedirectionURL() string
}

type UserRedirectAction interface {
	RedirectAction
	UserID() map[string]string
}

// Timed is an interface implemented by actions having an expiration time.
type Timed interface {
	Expired() bool
}

type blockAction struct {
	ID string
}

func newBlockAction(id string) blockAction {
	return blockAction{
		ID: id,
	}
}

func (a blockAction) ActionID() string {
	return a.ID
}

type redirectAction struct {
	ID, URL string
}

func newRedirectAction(id, location string) (*redirectAction, error) {
	if _, err := url.ParseRequestURI(location); err != nil {
		return nil, sqerrors.Wrap(err, "validation of the redirection location url")
	}
	return &redirectAction{
		ID:  id,
		URL: location,
	}, nil
}

func (a *redirectAction) ActionID() string {
	return a.ID
}

// timedAction is an Action with a time deadline after which it is considered
// expired.
type timedAction struct {
	Action
	deadline time.Time
}

// withDuration sets a time duration to an action. The returned value implements
// the Action and Timed interfaces.
func withDuration(action Action, duration time.Duration) *timedAction {
	return &timedAction{
		Action:   action,
		deadline: time.Now().Add(duration),
	}
}

// Expired is true when the deadline has expired, false otherwise.
func (a *timedAction) Expired() bool {
	// Is the current time after the deadline?
	return time.Now().After(a.deadline)
}

type whitelistAction struct {
	// The CIDR is used as action ID
	CIDR string
}

func newWhitelistAction(cidr string) whitelistAction {
	return whitelistAction{
		CIDR: cidr,
	}
}

func (a whitelistAction) ActionID() string {
	return a.CIDR
}
