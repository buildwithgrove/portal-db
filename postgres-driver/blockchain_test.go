package postgresdriver

import (
	"github.com/pokt-foundation/portal-db/types"
)

func (ts *PGDriverTestSuite) Test_ReadBlockchains() {
	tests := []struct {
		name        string
		blockchains []*types.Blockchain
		err         error
	}{
		{
			name: "Should return all Load Balancers from the database ordered by blockchain_id",
			blockchains: []*types.Blockchain{
				{
					ID:                "0001",
					Altruist:          "https://test:test_93uhfniu23f8@shared-test2.nodes.pokt.network:12345",
					Blockchain:        "pokt-mainnet",
					Description:       "POKT Network Mainnet",
					EnforceResult:     "JSON",
					Network:           "POKT-mainnet",
					Ticker:            "POKT",
					BlockchainAliases: []string{"pokt-mainnet"},
					LogLimitBlocks:    100_000,
					Active:            true,
					Redirects: []types.Redirect{
						{
							Alias:          "test-mainnet",
							Domain:         "test-rpc1.testnet.pokt.network",
							LoadBalancerID: "test_lb_34gg4g43g34g5hh",
						},
						{
							Alias:          "test-mainnet",
							Domain:         "test-rpc2.testnet.pokt.network",
							LoadBalancerID: "test_lb_34gg4g43g34g5hh",
						},
					},
					SyncCheckOptions: types.SyncCheckOptions{
						Body:      `{}`,
						Path:      "/v1/query/height",
						ResultKey: "height",
						Allowance: 1,
					},
				},
				{
					ID:                "0021",
					Altruist:          "https://test:test_u32fh239hf@shared-test2.nodes.eth.network:12345",
					Blockchain:        "eth-mainnet",
					ChainID:           "1",
					ChainIDCheck:      `{\"method\":\"eth_chainId\",\"id\":1,\"jsonrpc\":\"2.0\"}`,
					Description:       "Ethereum Mainnet",
					EnforceResult:     "JSON",
					Network:           "ETH-1",
					Ticker:            "ETH",
					BlockchainAliases: []string{"eth-mainnet"},
					LogLimitBlocks:    100_000,
					Active:            true,
					Redirects: []types.Redirect{
						{
							Alias:          "eth-mainnet",
							Domain:         "test-rpc.testnet.eth.network",
							LoadBalancerID: "test_lb_34gg4g43g34g5hh",
						},
					},
					SyncCheckOptions: types.SyncCheckOptions{
						Body:      `{\"method\":\"eth_blockNumber\",\"id\":1,\"jsonrpc\":\"2.0\"}`,
						ResultKey: "result",
						Allowance: 5,
					},
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		blockchains, err := ts.driver.ReadBlockchains(testCtx)
		ts.Equal(test.err, err)
		for i, blockchain := range blockchains {
			ts.Equal(test.blockchains[i].ID, blockchain.ID)
			ts.Equal(test.blockchains[i].ID, blockchain.ID)
			ts.Equal(test.blockchains[i].Altruist, blockchain.Altruist)
			ts.Equal(test.blockchains[i].Blockchain, blockchain.Blockchain)
			ts.Equal(test.blockchains[i].ChainID, blockchain.ChainID)
			ts.Equal(test.blockchains[i].ChainIDCheck, blockchain.ChainIDCheck)
			ts.Equal(test.blockchains[i].Description, blockchain.Description)
			ts.Equal(test.blockchains[i].EnforceResult, blockchain.EnforceResult)
			ts.Equal(test.blockchains[i].Network, blockchain.Network)
			ts.Equal(test.blockchains[i].Path, blockchain.Path)
			ts.Equal(test.blockchains[i].SyncCheck, blockchain.SyncCheck)
			ts.Equal(test.blockchains[i].Ticker, blockchain.Ticker)
			ts.Equal(test.blockchains[i].BlockchainAliases, blockchain.BlockchainAliases)
			ts.Equal(test.blockchains[i].LogLimitBlocks, blockchain.LogLimitBlocks)
			ts.Equal(test.blockchains[i].RequestTimeout, blockchain.RequestTimeout)
			ts.Equal(test.blockchains[i].SyncAllowance, blockchain.SyncAllowance)
			ts.Equal(test.blockchains[i].Active, blockchain.Active)
			ts.Equal(test.blockchains[i].Redirects, blockchain.Redirects)
			ts.Equal(test.blockchains[i].SyncCheckOptions, blockchain.SyncCheckOptions)
			ts.NotEmpty(blockchain.CreatedAt)
			ts.NotEmpty(blockchain.UpdatedAt)
		}
	}
}

func (ts *PGDriverTestSuite) Test_WriteBlockchain() {
	tests := []struct {
		name                string
		chainInput          *types.Blockchain
		expectedNumOfChains int
		err                 error
	}{
		{
			name: "Should create a single load balancer successfully with correct input",
			chainInput: &types.Blockchain{
				ID:                "003",
				Altruist:          "https://test:test_fg332f@shared-test3.nodes.pol.network:12345",
				Blockchain:        "pol-mainnet",
				Description:       "Polygon Mainnet",
				EnforceResult:     "JSON",
				Network:           "POL-mainnet",
				Ticker:            "POL",
				BlockchainAliases: []string{"pol-mainnet"},
				LogLimitBlocks:    100000,
				Active:            true,
				SyncCheckOptions: types.SyncCheckOptions{
					Body:      "{}",
					ResultKey: "result",
					Allowance: 3,
				},
			},
			expectedNumOfChains: 3,
			err:                 nil,
		},
	}

	for _, test := range tests {
		createdChain, err := ts.driver.WriteBlockchain(testCtx, test.chainInput)
		ts.Equal(test.err, err)
		ts.Equal(test.chainInput.ID, createdChain.ID)
		ts.NotEmpty(createdChain.CreatedAt)
		ts.NotEmpty(createdChain.UpdatedAt)

		chains, err := ts.driver.ReadBlockchains(testCtx)
		ts.Equal(test.err, err)
		ts.Len(chains, test.expectedNumOfChains)
		for _, blockchain := range chains {
			if blockchain.ID == test.chainInput.ID {
				ts.Equal(test.chainInput.ID, blockchain.ID)
				ts.Equal(test.chainInput.Altruist, blockchain.Altruist)
				ts.Equal(test.chainInput.Blockchain, blockchain.Blockchain)
				ts.Equal(test.chainInput.ChainID, blockchain.ChainID)
				ts.Equal(test.chainInput.ChainIDCheck, blockchain.ChainIDCheck)
				ts.Equal(test.chainInput.Description, blockchain.Description)
				ts.Equal(test.chainInput.EnforceResult, blockchain.EnforceResult)
				ts.Equal(test.chainInput.Network, blockchain.Network)
				ts.Equal(test.chainInput.Path, blockchain.Path)
				ts.Equal(test.chainInput.SyncCheck, blockchain.SyncCheck)
				ts.Equal(test.chainInput.Ticker, blockchain.Ticker)
				ts.Equal(test.chainInput.BlockchainAliases, blockchain.BlockchainAliases)
				ts.Equal(test.chainInput.LogLimitBlocks, blockchain.LogLimitBlocks)
				ts.Equal(test.chainInput.RequestTimeout, blockchain.RequestTimeout)
				ts.Equal(test.chainInput.SyncAllowance, blockchain.SyncAllowance)
				ts.Equal(test.chainInput.Active, blockchain.Active)
				ts.Equal(test.chainInput.SyncCheckOptions, blockchain.SyncCheckOptions)
				ts.NotEmpty(blockchain.CreatedAt)
				ts.NotEmpty(blockchain.UpdatedAt)
			}
			break
		}
	}
}

func (ts *PGDriverTestSuite) Test_WriteRedirect() {
	tests := []struct {
		name                   string
		redirectInput          *types.Redirect
		expectedNumOfRedirects int
		err                    error
	}{
		{
			name: "Should add a single redirect to an existing blockchain",
			redirectInput: &types.Redirect{
				BlockchainID:   "0021",
				Alias:          "eth-mainnet",
				Domain:         "test-rpc2.testnet.eth.network",
				LoadBalancerID: "test_lb_34gg4g43g34g5hh",
			},
			expectedNumOfRedirects: 2,
			err:                    nil,
		},
	}

	for _, test := range tests {
		createdRedirect, err := ts.driver.WriteRedirect(testCtx, test.redirectInput)
		ts.Equal(test.err, err)
		ts.Equal(test.redirectInput.BlockchainID, createdRedirect.BlockchainID)

		chains, err := ts.driver.ReadBlockchains(testCtx)
		ts.Equal(test.err, err)
		for _, blockchain := range chains {
			if blockchain.ID == test.redirectInput.BlockchainID {
				ts.Len(blockchain.Redirects, test.expectedNumOfRedirects)
				for i, redirect := range blockchain.Redirects {
					ts.Equal(test.redirectInput.BlockchainID, redirect.BlockchainID)
					ts.Equal(test.redirectInput.Alias, redirect.Alias)
					ts.Equal(test.redirectInput.LoadBalancerID, redirect.LoadBalancerID)
					if i == len(blockchain.Redirects)-1 {
						ts.Equal(test.redirectInput.Domain, redirect.Domain)
					}
				}
			}
			break
		}
	}
}

func (ts *PGDriverTestSuite) Test_ActivateBlockchain() {
	tests := []struct {
		name         string
		blockchainID string
		active       bool
		err          error
	}{
		{
			name:         "Should successfully deactivate a blockchain",
			blockchainID: "0001",
			active:       false,
			err:          nil,
		},
		{
			name:         "Should successfully activate a blockchain",
			blockchainID: "0001",
			active:       true,
			err:          nil,
		},
	}

	for _, test := range tests {
		err := ts.driver.ActivateChain(testCtx, test.blockchainID, test.active)
		ts.Equal(test.err, err)

		chains, err := ts.driver.ReadBlockchains(testCtx)
		ts.Equal(test.err, err)
		for _, blockchain := range chains {
			if blockchain.ID == test.blockchainID {
				ts.Equal(test.active, blockchain.Active)
			}
			break
		}
	}
}
