package bitflyer

type ProductCode int

const (
	Btcjpy ProductCode = iota
	Ethjpy
	FxbtcJpy
	Ethbtc
	Bchbtc
)

func (code ProductCode) String() string {
	switch code {
	case Btcjpy:
		return "BTC_JPY"
	case Ethjpy:
		return "ETH_JPY"
	case FxbtcJpy:
		return "FX_BTC_JPY"
	case Ethbtc:
		return "ETH_BTC"
	case Bchbtc:
		return "BCH_BTC"
	default:
		return "BTC_JPY"
	}
}
