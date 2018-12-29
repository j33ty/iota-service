package iota

import (
	address "github.com/iotaledger/iota.go/address"
	api "github.com/iotaledger/iota.go/api"
	bundle "github.com/iotaledger/iota.go/bundle"
	consts "github.com/iotaledger/iota.go/consts"
	trinary "github.com/iotaledger/iota.go/trinary"
	app "github.com/j33ty/iota-service/includes/app"
)

// GenAddr - GenAddr
func GenAddr(c *app.Context) (string, error) {
	return address.GenerateAddress(c.Config.Seed, 2, consts.SecurityLevelMedium)
}

// CreateTransfer - Create a transfer to the given recipient address
func CreateTransfer(recipientAddr, msg, tag string, amount uint64) bundle.Transfers {
	return bundle.Transfers{
		{
			Address: trinary.Hash(recipientAddr),
			Value:   amount,
			Message: trinary.Trytes(msg),
			Tag:     trinary.Trytes(tag),
		},
	}
}

// CreateTransferInput - Create inputs for the transfer
func CreateTransferInput(inputAddr string, keyIndex uint64, bal uint64) []api.Input {
	return []api.Input{
		{
			Address:  trinary.Hash(inputAddr),
			Security: consts.SecurityLevelMedium,
			KeyIndex: keyIndex,
			Balance:  bal,
		},
	}
}

// PrepTransferOptions - Create inputs for the transfer
func PrepTransferOptions(inputs []api.Input, remainderAddr string) api.PrepareTransfersOptions {
	r := trinary.Trytes(remainderAddr)
	return api.PrepareTransfersOptions{
		Inputs:           inputs,
		RemainderAddress: &r,
	}
}

// PrepareTransfers - Prepare the transfer by creating a bundle with the given transfers and inputs. The result are trytes ready for PoW.
func PrepareTransfers(c *app.Context, transfers bundle.Transfers, prepTransferOpts api.PrepareTransfersOptions) ([]trinary.Trytes, error) {
	return c.IOTA.PrepareTransfers(c.Config.Seed, transfers, prepTransferOpts)
}

// IsAddrSentFrom -	Decrease chances of sending to a spent address by checking the address before broadcasting the bundle.
func IsAddrSentFrom(c *app.Context, addr trinary.Hash) ([]bool, error) {
	return c.IOTA.WereAddressesSpentFrom(addr)
}

// BroadcastTransaction - bundle trytes are signed. Now select two tips, do PoW, broadcast the bundle and store the bundle. SendTrytes() does it all..
func BroadcastTransaction(c *app.Context, trytes []trinary.Trytes) (bundle.Bundle, error) {
	return c.IOTA.SendTrytes(trytes, c.Config.Depth, c.Config.MWM)
}
