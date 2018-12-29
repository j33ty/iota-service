package controllers

import (
	"errors"

	api "github.com/iotaledger/iota.go/api"
	bundle "github.com/iotaledger/iota.go/bundle"
	validators "github.com/iotaledger/iota.go/guards/validators"
	transaction "github.com/iotaledger/iota.go/transaction"
	trinary "github.com/iotaledger/iota.go/trinary"
	app "github.com/j33ty/iota-service/includes/app"
	iota "github.com/j33ty/iota-service/includes/iota"
)

// GetNodeInfo - GetNodeInfo
func GetNodeInfo(c *app.Context) (*api.GetNodeInfoResponse, error) {
	return c.IOTA.GetNodeInfo()
}

// GetNodeVersion - GetNodeVersion
func GetNodeVersion(c *app.Context) (string, error) {
	resp, err := c.IOTA.GetNodeInfo()
	if err != nil {
		return "", err
	}
	return resp.AppVersion, nil
}

// GetLastBlock - GetLastBlock
func GetLastBlock(c *app.Context) (int64, error) {
	resp, err := c.IOTA.GetNodeInfo()
	if err != nil {
		return 0, err
	}
	return resp.LatestMilestoneIndex, nil
}

// GetBalances - fetches confirmed balances of the given addresses at the latest solid milestone.
func GetBalances(c *app.Context, addresses []string) (*api.Balances, error) {
	return c.IOTA.GetBalances(addresses, 100)
}

// IsAddrValid - IsAddrValid
func IsAddrValid(c *app.Context, addr string) error {
	if err := validators.Validate(validators.ValidateHashes(trinary.Hash(addr))); err != nil {
		return errors.New("Invalid IOTA Address")
	}
	if b, err := c.IOTA.IsAddressUsed(addr); err != nil || b {
		return errors.New("IOTA Address already in use")
	}
	return nil
}

// IsSeedValid - IsSeedValid
func IsSeedValid(seed string) error {
	if err := validators.Validate(validators.ValidateSeed(trinary.Trytes(seed))); err != nil {
		return errors.New("Invalid IOTA Seed")
	}
	return nil
}

// AccountHistory - AccountHistory
// TODO: Complete this method
func AccountHistory(c *app.Context) (*api.AccountData, error) {
	return c.IOTA.GetAccountData(c.Config.Seed, api.GetAccountDataOptions{})
}

// FindTransByAddr - FindTransByAddr
func FindTransByAddr(c *app.Context, addr string) ([]string, error) {
	return c.IOTA.FindTransactions(api.FindTransactionsQuery{Addresses: trinary.Hashes([]string{addr})})
}

// FindTransObjByAddr - FindTransObjByAddr
func FindTransObjByAddr(c *app.Context, addr string) (transaction.Transactions, error) {
	return c.IOTA.FindTransactionObjects(api.FindTransactionsQuery{Addresses: []trinary.Hash{addr}})
}

// FindTransByTag - FindTransByTag
func FindTransByTag(c *app.Context, tag string) ([]string, error) {
	return c.IOTA.FindTransactions(api.FindTransactionsQuery{Tags: []trinary.Trytes{tag}})
}

// FindTransObjByTag - FindTransObjByTag
func FindTransObjByTag(c *app.Context, tag string) (transaction.Transactions, error) {
	return c.IOTA.FindTransactionObjects(api.FindTransactionsQuery{Tags: []trinary.Trytes{tag}})
}

// AreTransComplete - AreTransComplete
func AreTransComplete(c *app.Context, txHashes []string) ([]bool, error) {
	return c.IOTA.GetLatestInclusion(trinary.Hashes(txHashes))
}

// FetchTransByHash - FetchTransByHash
func FetchTransByHash(c *app.Context, txHash string) (transaction.Transactions, error) {
	return c.IOTA.GetTransactionObjects(trinary.Hash(txHash))
}

// TransferBal - TransferBal
func TransferBal(c *app.Context, senderAddr, recipientAddr string, initialBal, amount uint64) (bundle.Bundle, error) {
	var b bundle.Bundle
	transfers := iota.CreateTransfer(recipientAddr, "", "", amount)
	if len(transfers) < 1 {
		return b, errors.New("Failed to create the transfer")
	}

	inputs := iota.CreateTransferInput(senderAddr, 0, initialBal)

	// Create an address for the remainder.
	remainderAddr, err := iota.GenAddr(c)
	if err != nil {
		return b, errors.New("Failed to generate the address")
	}

	prepTransferOpts := iota.PrepTransferOptions(inputs, remainderAddr)
	trytes, err := iota.PrepareTransfers(c, transfers, prepTransferOpts)
	if err != nil {
		return b, errors.New("Failed to prepare the tranfer")
	}

	spent, err := iota.IsAddrSentFrom(c, transfers[0].Address)
	if err != nil {
		return b, errors.New("Failed to check address for SpentFrom")
	}

	if len(spent) > 0 && spent[0] {
		return b, errors.New("recipient address is spent from, aborting transfer")
	}

	b, err = iota.BroadcastTransaction(c, trytes)
	if err != nil {
		return b, errors.New("Failed to Send Trytes on the network")
	}

	return b, nil
}
