package types

import (
	"time"
)

/* GigastakeApp Struct Definition and Methods */
type (
	// GigastakeApp represents a single gigastaked application for a given chain
	GigastakeApp struct {
		AATID      ProtocolAppID `json:"aatID"`
		ChainID    RelayChainID  `json:"chainID"`
		ChainAlias string        `json:"chainAlias"`
		Name       string        `json:"name"`
		AAT        AAT           `json:"aat"`
		CreatedAt  time.Time     `json:"createdAt"`
		UpdatedAt  time.Time     `json:"updatedAt"`
		Deleted    bool          `json:"deleted"`
	}

	// AAT contains the data needed to perform relays using a gigastake app
	AAT struct {
		ID              ProtocolAppID `json:"id"`
		Gigastake       bool          `json:"gigastake"`
		Address         string        `json:"address"`
		PublicKey       string        `json:"publicKey"`
		ClientPublicKey string        `json:"clientPublicKey"`
		PrivateKey      string        `json:"privateKey"`
		Signature       string        `json:"signature"`
		Version         string        `json:"version"`

		// PortalAppID only used for non-gigastake AATs, which
		// are currently not used anywhere but kept for compatibility.
		PortalAppID PortalAppID `json:"appID,omitempty"`
	}
)

func (a *GigastakeApp) Table() Table {
	return TableGigastakeApplications
}

func (a *AAT) Table() Table {
	return TableAATs
}
