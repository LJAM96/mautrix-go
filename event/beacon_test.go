// Copyright (c) 2026 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package event_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

func TestBeaconTypes(t *testing.T) {
	for _, tc := range []struct {
		name  string
		typ   event.Type
		class event.TypeClass
	}{
		{"stable info", event.StateBeaconInfo, event.StateEventType},
		{"unstable info", event.StateUnstableBeaconInfo, event.StateEventType},
		{"stable beacon", event.EventBeacon, event.MessageEventType},
		{"unstable beacon", event.EventUnstableBeacon, event.MessageEventType},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.class, tc.typ.Class)
			assert.Equal(t, tc.class, (&event.Type{Type: tc.typ.Type}).GuessClass())
			var content event.Content
			content.VeryRaw = json.RawMessage(`{}`)
			require.NoError(t, content.ParseRaw(tc.typ))
			assert.NotNil(t, content.Parsed)
		})
	}
}

func TestBeaconInfoEventContent(t *testing.T) {
	t.Run("unstable fields", func(t *testing.T) {
		var content event.BeaconInfoEventContent
		require.NoError(t, json.Unmarshal([]byte(`{
			"description": "Kylie's live location",
			"live": true,
			"timeout": 3600000,
			"org.matrix.msc3488.asset": {"type": "m.self"},
			"org.matrix.msc3488.ts": 651066205621
		}`), &content))
		assert.Equal(t, "Kylie's live location", content.Description)
		assert.True(t, content.Live)
		assert.EqualValues(t, 3600000, content.Timeout)
		assert.EqualValues(t, 651066205621, content.GetTimestamp())
		require.NotNil(t, content.GetAsset())
		assert.Equal(t, event.LocationAssetTypeSelf, content.GetAsset().Type)
	})

	t.Run("stable fields", func(t *testing.T) {
		var content event.BeaconInfoEventContent
		require.NoError(t, json.Unmarshal([]byte(`{
			"live": false,
			"timeout": 86400000,
			"m.asset": {"type": "m.pin"},
			"m.ts": 1436829458432
		}`), &content))
		assert.False(t, content.Live)
		assert.EqualValues(t, 1436829458432, content.GetTimestamp())
		require.NotNil(t, content.GetAsset())
		assert.Equal(t, event.LocationAssetTypePin, content.GetAsset().Type)
	})

	t.Run("type map", func(t *testing.T) {
		var content event.Content
		content.VeryRaw = json.RawMessage(`{"live":true,"timeout":1000}`)
		require.NoError(t, content.ParseRaw(event.StateUnstableBeaconInfo))
		assert.True(t, content.AsBeaconInfo().Live)
	})
}

func TestBeaconEventContent(t *testing.T) {
	t.Run("unstable fields", func(t *testing.T) {
		var content event.BeaconEventContent
		require.NoError(t, json.Unmarshal([]byte(`{
			"m.relates_to": {"event_id": "$info", "rel_type": "m.reference"},
			"org.matrix.msc3488.location": {"uri": "geo:8.95,12.49;u=10"},
			"org.matrix.msc3488.ts": 651066206460
		}`), &content))
		assert.Equal(t, event.RelReference, content.RelatesTo.Type)
		assert.Equal(t, id.EventID("$info"), content.GetRelatesTo().GetReferenceID())
		assert.EqualValues(t, 651066206460, content.GetTimestamp())
		require.NotNil(t, content.GetLocation())
		assert.Equal(t, "geo:8.95,12.49;u=10", content.GetLocation().URI)
	})

	t.Run("stable fields", func(t *testing.T) {
		var content event.BeaconEventContent
		require.NoError(t, json.Unmarshal([]byte(`{
			"m.relates_to": {"event_id": "$info", "rel_type": "m.reference"},
			"m.location": {"uri": "geo:51.5,0.12", "description": "here"},
			"m.ts": 1636829458432
		}`), &content))
		assert.EqualValues(t, 1636829458432, content.GetTimestamp())
		require.NotNil(t, content.GetLocation())
		assert.Equal(t, "here", content.GetLocation().Description)
	})

	t.Run("marshal uses unstable names", func(t *testing.T) {
		content := event.BeaconEventContent{
			RelatesTo: event.RelatesTo{Type: event.RelReference, EventID: "$info"},
			Location:  &event.BeaconLocation{URI: "geo:1,2"},
			Timestamp: 123,
		}
		data, err := json.Marshal(&content)
		require.NoError(t, err)
		assert.JSONEq(t, `{
			"m.relates_to": {"event_id": "$info", "rel_type": "m.reference"},
			"org.matrix.msc3488.location": {"uri": "geo:1,2"},
			"org.matrix.msc3488.ts": 123
		}`, string(data))
	})

	t.Run("type map", func(t *testing.T) {
		var content event.Content
		content.VeryRaw = json.RawMessage(`{"m.relates_to":{"rel_type":"m.reference"}}`)
		require.NoError(t, content.ParseRaw(event.EventBeacon))
		assert.NotNil(t, content.AsBeacon().GetRelatesTo())
	})
}
