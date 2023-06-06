package types

import (
	"time"
)

/* GigastakeApp Struct Definition and Methods */
type (
	// GigastakeApp represents a single gigastaked application for a given chain
	GigastakeApp struct {
		ID        GigastakeAppID            `json:"id"`
		AATID     ProtocolAppID             `json:"aatID"`
		ChainIDs  map[RelayChainID]struct{} `json:"chainIDs"`
		Name      string                    `json:"name"`
		AAT       AAT                       `json:"aat"`
		CreatedAt time.Time                 `json:"createdAt"`
		UpdatedAt time.Time                 `json:"updatedAt"`
		Deleted   bool                      `json:"deleted"`

		// TODO remove legacy field when migration to V2 schema complete
		LegacyLBID string `json:"legacyLBID"`
	}

	// AAT contains the data needed to perform relays
	AAT struct {
		ID              ProtocolAppID `json:"id"`
		Gigastake       bool          `json:"gigastake"`
		Address         string        `json:"address"`
		PublicKey       string        `json:"publicKey"`
		ClientPublicKey string        `json:"clientPublicKey"`
		Signature       string        `json:"signature"`
		Version         string        `json:"version"`

		// PrivateKey used when read from the DB, will always be ""
		// Only used for saving to DB
		// TODO remove when decided to not support saving private key to DB
		PrivateKey string `json:"privateKey,omitempty"`

		// PortalAppID only used for non-gigastake AATs, which
		// are currently not used anywhere but kept for compatibility.
		PortalAppID PortalAppID `json:"appID,omitempty"`
	}

	// GigastakeApp represents the relationship between a Chain and GigastakeApp
	// Only used by the listener
	ChainGigastakeApp struct {
		ChainID        RelayChainID   `json:"chainID"`
		GigastakeAppID GigastakeAppID `json:"gigastakeAppID"`
	}
)

func (a *GigastakeApp) Table() Table {
	return TableGigastakeApps
}

func (a *AAT) Table() Table {
	return TableAATs
}
