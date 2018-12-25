package coin

import "github.com/boomstarternetwork/bestore"

func List() []bestore.Coin {
	return []bestore.Coin{
		bestore.BTC,
		bestore.ETH,
	}
}
