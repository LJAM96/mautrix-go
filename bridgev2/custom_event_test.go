// Copyright (c) 2026 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package bridgev2

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"maunium.net/go/mautrix/event"
)

type testCustomEventHandler struct {
	NetworkAPI
	called   int
	received *MatrixCustomEvent
	err      error
}

func (t *testCustomEventHandler) HandleMatrixCustomEvent(_ context.Context, msg *MatrixCustomEvent) error {
	t.called++
	t.received = msg
	return t.err
}

type testNoCustomEventHandler struct {
	NetworkAPI
}

type testCustomTypesProvider struct {
	NetworkConnector
	types []event.Type
}

func (t *testCustomTypesProvider) GetCustomMatrixEventTypes() []event.Type {
	return t.types
}

func TestCustomEventTypesProvider(t *testing.T) {
	types := []event.Type{{Type: "com.example.custom", Class: event.MessageEventType}}
	var provider CustomEventTypesProvider = &testCustomTypesProvider{types: types}
	assert.Equal(t, types, provider.GetCustomMatrixEventTypes())
}

func TestHandleMatrixCustomEvent(t *testing.T) {
	evtType := event.Type{Type: "com.example.custom", Class: event.MessageEventType}
	newEvent := func() *event.Event {
		evt := &event.Event{Type: evtType, ID: "$custom"}
		evt.Content.Raw = map[string]any{"hello": "world"}
		return evt
	}

	t.Run("supported", func(t *testing.T) {
		client := &testCustomEventHandler{}
		portal := &Portal{}
		origSender := &OrigSender{}
		res := portal.handleMatrixCustomEvent(context.Background(), &UserLogin{Client: client}, origSender, newEvent())
		assert.True(t, res.Success)
		assert.True(t, res.SendMSS)
		assert.Equal(t, 1, client.called)
		require.NotNil(t, client.received)
		assert.Same(t, portal, client.received.Portal)
		assert.Same(t, origSender, client.received.OrigSender)
		assert.Equal(t, evtType, client.received.Event.Type)
		assert.Equal(t, "world", client.received.Content.Raw["hello"])
	})

	t.Run("unsupported", func(t *testing.T) {
		res := (&Portal{}).handleMatrixCustomEvent(context.Background(), &UserLogin{Client: &testNoCustomEventHandler{}}, nil, newEvent())
		assert.True(t, res.Ignored)
		assert.False(t, res.SendMSS)
	})

	t.Run("error", func(t *testing.T) {
		expectedErr := errors.New("custom handler failed")
		client := &testCustomEventHandler{err: expectedErr}
		res := (&Portal{}).handleMatrixCustomEvent(context.Background(), &UserLogin{Client: client}, nil, newEvent())
		assert.False(t, res.Success)
		assert.True(t, res.SendMSS)
		assert.ErrorIs(t, res.Error, expectedErr)
		assert.Equal(t, 1, client.called)
	})
}
