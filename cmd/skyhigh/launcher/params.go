package launcher

import (
	"github.com/ethereum/go-ethereum/params"
)

var (
	Bootnodes = []string{
		"enode://eda0c594a92721706f134fe7bb5855324c8a50785509abf48f68d48f80965efb326858dd227d8b29612a3e4a366a41a183c0cbc875c6d12ffd013807c9073eb2@54.170.213.54:5050",
		"enode://a2f99d7a6fc4090ed067951a98e6121a49cb1789364dbd8be653d08b741df75235cab7524eb7f17614f6c51b9755eb4131ef943d5c8cbfb626d0022964cc0f27@34.243.192.59:5050",
		"enode://802704db6d524a2c06c97def5d2865b66a629578347b69bd7f7d8ba6bf44bc7bf6983352b393bd85423868ee16ae8467af218d1de50cd01a42ddd2dec6775d22@54.170.73.3:5050",
		"enode://29709396dfb06374196854e1d4871f947f56deb7b4b8fbe1ce8117240c12da7236abf7f2831d8ce5665069802ec678f1beff4094678a5b5de6a6203bed6f0e94@34.245.45.68:5050",
	}
)

func overrideParams() {
	params.MainnetBootnodes = []string{}
	params.RopstenBootnodes = []string{}
	params.RinkebyBootnodes = []string{}
	params.GoerliBootnodes = []string{}
}
