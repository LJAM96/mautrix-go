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
)

type testJoinRuleHandler struct {
	NetworkAPI
	msg *MatrixJoinRuleChange
}

func (t *testJoinRuleHandler) HandleMatrixJoinRule(_ context.Context, msg *MatrixJoinRuleChange) (bool, error) {
	t.msg = msg
	return false, nil
}

func TestHandleMatrixJoinRule(t *testing.T) {
	stateKey := ""
	evt := &event.Event{Type: event.StateJoinRules, StateKey: &stateKey}
	evt.Content.VeryRaw = json.RawMessage(`{"join_rule":"invite"}`)
	require.NoError(t, evt.Content.ParseRaw(event.StateJoinRules))
	evt.Unsigned.PrevContent = &event.Content{VeryRaw: json.RawMessage(`{"join_rule":"public"}`)}

	client := &testJoinRuleHandler{}
	res := handleMatrixRoomMeta(&Portal{}, context.Background(), &UserLogin{Client: client}, nil, evt, false,
		JoinRuleHandlingNetworkAPI.HandleMatrixJoinRule)
	assert.True(t, res.Success)
	assert.True(t, res.SendMSS)
	require.NotNil(t, client.msg)
	assert.Equal(t, event.JoinRuleInvite, client.msg.Content.JoinRule)
	require.NotNil(t, client.msg.PrevContent)
	assert.Equal(t, event.JoinRulePublic, client.msg.PrevContent.JoinRule)
}

func TestHandleMatrixJoinRuleUnsupported(t *testing.T) {
	stateKey := ""
	evt := &event.Event{Type: event.StateJoinRules, StateKey: &stateKey}
	evt.Content.VeryRaw = json.RawMessage(`{"join_rule":"invite"}`)
	require.NoError(t, evt.Content.ParseRaw(event.StateJoinRules))

	res := handleMatrixRoomMeta(&Portal{}, context.Background(), &UserLogin{Client: &testNoCustomEventHandler{}}, nil, evt, false,
		JoinRuleHandlingNetworkAPI.HandleMatrixJoinRule)
	assert.True(t, res.Ignored)
	assert.True(t, res.SendMSS)
	assert.ErrorIs(t, res.Error, ErrRoomMetadataNotSupported)
}
