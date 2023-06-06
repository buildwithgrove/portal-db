package types

import (
	"time"
)

/* GigastakeApp Struct Definition and Methods */
type (
	// GigastakeApp represents a single gigastaked application for a given chain
	GigastakeApp struct {
		ID              GigastakeAppID            `json:"id"`
		ChainIDs        map[RelayChainID]struct{} `json:"chainIDs"`
		Name            string                    `json:"name"`
		Address         string                    `json:"address"`
		PublicKey       string                    `json:"publicKey"`
		ClientPublicKey string                    `json:"clientPublicKey"`
		Signature       string                    `json:"signature"`
		Version         string                    `json:"version"`
		CreatedAt       time.Time                 `json:"createdAt"`
		UpdatedAt       time.Time                 `json:"updatedAt"`
		Deleted         bool                      `json:"deleted"`

		// PrivateKey used when read from the DB, will always be ""
		// Only used for saving to DB
		// TODO remove when decided to not support saving private key to DB
		PrivateKey string `json:"privateKey,omitempty"`
		// TODO remove legacy field when migration to V2 schema complete
		LegacyLBID string `json:"legacyLBID"`
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
