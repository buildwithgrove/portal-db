package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateBlockchain(t *testing.T) {
	tests := []struct {
		name     string
		chain    Chain
		update   UpdateChain
		expected Chain
	}{
		{
			name: "Should update blockchain description",
			chain: Chain{
				Description: "Old Description",
			},
			update: UpdateChain{
				Description: "New Description",
			},
			expected: Chain{
				Description: "New Description",
			},
		},
		{
			name:     "Should update multiple fields",
			chain:    Chain{Blockchain: "Old Blockchain", Description: "Old Description"},
			update:   UpdateChain{Blockchain: "New Blockchain", Description: "New Description"},
			expected: Chain{Blockchain: "New Blockchain", Description: "New Description"},
		},
		{
			name: "Should update chain checks",
			chain: Chain{
				Checks: map[ChainCheckType]Check{
					ChainCheckTypeChain: {Payload: "Old Payload"},
				},
			},
			update: UpdateChain{
				Checks: map[ChainCheckType]UpdateCheck{
					ChainCheckTypeChain: {Payload: "New Payload"},
				},
			},
			expected: Chain{
				Checks: map[ChainCheckType]Check{
					ChainCheckTypeChain: {Payload: "New Payload"},
				},
			},
		},
		{
			name:     "Should not update with empty update",
			chain:    Chain{Blockchain: "Old Blockchain", Description: "Old Description"},
			update:   UpdateChain{},
			expected: Chain{Blockchain: "Old Blockchain", Description: "Old Description"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			updatedChain := test.chain.UpdateBlockchain(&test.update)
			assert.Equal(t, test.expected, *updatedChain)
		})
	}
}
