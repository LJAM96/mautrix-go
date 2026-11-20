// Copyright (c) 2026 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package event

// Live location beacon events are defined in MSC3489 and MSC3672:
// https://github.com/matrix-org/matrix-spec-proposals/pull/3489
// https://github.com/matrix-org/matrix-spec-proposals/pull/3672
//
// Element sends the unstable event types and the unstable MSC3488 field names,
// so those are the primary fields below. The stable m.* alternatives are kept
// for when the MSCs land; use the Get* helpers to read either form.

// LocationAssetType identifies the asset being tracked in a beacon.
type LocationAssetType string

const (
	LocationAssetTypeSelf LocationAssetType = "m.self"
	LocationAssetTypePin  LocationAssetType = "m.pin"
)

// BeaconAsset is the asset object from MSC3488.
type BeaconAsset struct {
	Type LocationAssetType `json:"type"`
}

// BeaconInfoEventContent is the content of an m.beacon_info state event.
type BeaconInfoEventContent struct {
	// Description is the same as the description of an m.location event.
	Description string `json:"description,omitempty"`
	// Live is true while the beacon is being shared.
	Live bool `json:"live"`
	// Timeout is the number of milliseconds that the beacon stays live after Timestamp.
	Timeout int64 `json:"timeout"`
	// Timestamp is the time when sharing started, in milliseconds since the epoch.
	Timestamp int64        `json:"org.matrix.msc3488.ts,omitempty"`
	Asset     *BeaconAsset `json:"org.matrix.msc3488.asset,omitempty"`

	// StableTimestamp and StableAsset hold the stable m.ts and m.asset fields.
	StableTimestamp int64        `json:"m.ts,omitempty"`
	StableAsset     *BeaconAsset `json:"m.asset,omitempty"`
}

// GetTimestamp returns the timestamp from either the unstable or stable field.
func (content *BeaconInfoEventContent) GetTimestamp() int64 {
	if content.Timestamp != 0 {
		return content.Timestamp
	}
	return content.StableTimestamp
}

// GetAsset returns the asset from either the unstable or stable field.
func (content *BeaconInfoEventContent) GetAsset() *BeaconAsset {
	if content.Asset != nil {
		return content.Asset
	}
	return content.StableAsset
}

// BeaconLocation is the location object from MSC3488. URI is an RFC 5870 geo: URI.
type BeaconLocation struct {
	URI         string `json:"uri"`
	Description string `json:"description,omitempty"`
}

// BeaconEventContent is the content of an m.beacon message event.
type BeaconEventContent struct {
	// RelatesTo points at the beacon_info event that started the share.
	RelatesTo RelatesTo       `json:"m.relates_to"`
	Location  *BeaconLocation `json:"org.matrix.msc3488.location,omitempty"`
	Timestamp int64           `json:"org.matrix.msc3488.ts,omitempty"`

	// StableLocation and StableTimestamp hold the stable m.location and m.ts fields.
	StableLocation  *BeaconLocation `json:"m.location,omitempty"`
	StableTimestamp int64           `json:"m.ts,omitempty"`
}

var _ Relatable = (*BeaconEventContent)(nil)

// GetLocation returns the location from either the unstable or stable field.
func (content *BeaconEventContent) GetLocation() *BeaconLocation {
	if content.Location != nil {
		return content.Location
	}
	return content.StableLocation
}

// GetTimestamp returns the timestamp from either the unstable or stable field.
func (content *BeaconEventContent) GetTimestamp() int64 {
	if content.Timestamp != 0 {
		return content.Timestamp
	}
	return content.StableTimestamp
}

func (content *BeaconEventContent) GetRelatesTo() *RelatesTo {
	return &content.RelatesTo
}

func (content *BeaconEventContent) OptionalGetRelatesTo() *RelatesTo {
	return &content.RelatesTo
}

func (content *BeaconEventContent) SetRelatesTo(rel *RelatesTo) {
	content.RelatesTo = *rel
}
