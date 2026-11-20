// Copyright (c) 2026 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package bridgev2

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type testPinnedEventsHandler struct {
	NetworkAPI
	msg *MatrixRoomPinnedEvents
}

func (t *testPinnedEventsHandler) HandleMatrixRoomPinnedEvents(_ context.Context, msg *MatrixRoomPinnedEvents) (bool, error) {
	t.msg = msg
	return false, nil
}

func TestHandleMatrixRoomPinnedEvents(t *testing.T) {
	stateKey := ""
	evt := &event.Event{Type: event.StatePinnedEvents, StateKey: &stateKey}
	evt.Content.VeryRaw = json.RawMessage(`{"pinned":["$one","$two"]}`)
	require.NoError(t, evt.Content.ParseRaw(event.StatePinnedEvents))
	evt.Unsigned.PrevContent = &event.Content{VeryRaw: json.RawMessage(`{"pinned":["$one"]}`)}

	client := &testPinnedEventsHandler{}
	res := handleMatrixRoomMeta(&Portal{}, context.Background(), &UserLogin{Client: client}, nil, evt, false,
		PinnedEventsHandlingNetworkAPI.HandleMatrixRoomPinnedEvents)
	assert.True(t, res.Success)
	assert.True(t, res.SendMSS)
	require.NotNil(t, client.msg)
	assert.Equal(t, []id.EventID{"$one", "$two"}, client.msg.Content.Pinned)
	require.NotNil(t, client.msg.PrevContent)
	assert.Equal(t, []id.EventID{"$one"}, client.msg.PrevContent.Pinned)
}

func TestHandleMatrixRoomPinnedEventsUnsupported(t *testing.T) {
	stateKey := ""
	evt := &event.Event{Type: event.StatePinnedEvents, StateKey: &stateKey}
	evt.Content.VeryRaw = json.RawMessage(`{"pinned":[]}`)
	require.NoError(t, evt.Content.ParseRaw(event.StatePinnedEvents))

	res := handleMatrixRoomMeta(&Portal{}, context.Background(), &UserLogin{Client: &testNoCustomEventHandler{}}, nil, evt, false,
		PinnedEventsHandlingNetworkAPI.HandleMatrixRoomPinnedEvents)
	assert.True(t, res.Ignored)
	assert.True(t, res.SendMSS)
	assert.ErrorIs(t, res.Error, ErrRoomMetadataNotSupported)
}
