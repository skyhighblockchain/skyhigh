package launcher

import (
	"github.com/ethereum/go-ethereum/params"
)

var (
	Bootnodes = []string{
		"enode://8e47600391415c9bf04e854eb13691d04e100f48c3476fdc3661ebb8b3047099efcc8beeba54f7588763d53a188ea909a7a5323afd81b852f5f888d303a95f8e@38.242.128.146:5050",
		"enode://27fdc5544a0d0a8ab8173826f9dce645ae5e4d5cf5205dc0275b19303b6b01d376c23971093faef9d9e5b852361eddab0da0d78587c1892424da0227f6092ef2@62.171.181.125:5050",
		"enode://1ae643494d71559200d9e4d8d7dfabfe124a9f68fd292f89377b886c9a184318349946f299661a112f6baf5ce504249f8a07381fc5ee4792ddcc2d05bc425e59@38.242.128.202:5050",
		"enode://b5c7e40a1f6b794ab2aa3c248d5a6a6309259fbed12725891b6832f7d42831f1d2ea1bc9033ba3434e35acea727379e12ec40a278380be6af4e45877e79bdec3@62.171.185.88:5050",
	}
)

func overrideParams() {
	params.MainnetBootnodes = []string{}
	params.RopstenBootnodes = []string{}
	params.RinkebyBootnodes = []string{}
	params.GoerliBootnodes = []string{}
}
