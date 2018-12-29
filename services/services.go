package services

import (
	"fmt"
	"os"

	bundle "github.com/iotaledger/iota.go/bundle"
	controllers "github.com/j33ty/iota-service/controllers"
	app "github.com/j33ty/iota-service/includes/app"
	iota "github.com/j33ty/iota-service/includes/iota"
)

const sampleAddr = "LXPACSDCMQSPLOHVXYRWWZJXKAVJXLODYCJVMZR9TCBIZWLUMXWVNDURMJSTZLVQHRMRVOPEVOYIZKOPXCPAHWL9BC"

// Run - Run
func Run(c *app.Context) {

	version, err := controllers.GetNodeVersion(c)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("App Version: %s\n", version)
	}

	index, err := controllers.GetLastBlock(c)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("LastBlock: %d\n", index)
	}

	bal, err := controllers.GetBalances(c, []string{sampleAddr})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Balance: %d Address: %s\n", bal.Balances, sampleAddr)
	}

	addr, err := iota.GenAddr(c)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("New Address: %s\n", addr)
	}

	if controllers.IsAddrValid(c, sampleAddr) != nil {
		fmt.Printf("Invalid Address: %s\n", sampleAddr)
	} else {
		fmt.Printf("Valid Address: %s\n", sampleAddr)
	}

	seed := c.Config.Seed
	if controllers.IsSeedValid(seed) != nil {
		fmt.Printf("Invalid Seed: %s\n", seed)
	} else {
		fmt.Printf("Valid Seed: %s\n", seed)
	}

	data, err := controllers.FindTransObjByAddr(c, "LXPACSDCMQSPLOHVXYRWWZJXKAVJXLODYCJVMZR9TCBIZWLUMXWVNDURMJSTZLVQHRMRVOPEVOYIZKOPXCPAHWL9BC")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Find Trans Objects by Address")
		for i, h := range data {
			fmt.Printf("Index: %d Hash: %s\n", i, h.Hash)
		}
	}

	os.Exit(0)
	bndl, err := controllers.TransferBal(c, sampleAddr, "MWSK9HM9VSBAWKJNTUIWDAPEUJQRRQVXWBUMJRWKJLXSJUGNIEFLEFTFRHKGTLHONODSGQMMAPGOIIJRW", 9500000, 1)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("broadcasted bundle with tail tx hash: ", bundle.TailTransactionHash(bndl))
	}
}
