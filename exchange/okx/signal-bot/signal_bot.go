package signal_bot

import (
	"okx-bot/exchange/okx/common"
	"okx-bot/exchange/options"
)

type SignalBot struct {
	*common.OKxV5
}

func (g *SignalBot) NewPrvApi(apiOps ...options.ApiOption) *PrvApi {
	prv := new(PrvApi)
	prv.Prv = g.OKxV5.NewPrvApi(apiOps...)
	prv.Prv.OKxV5 = g.OKxV5
	return prv
}
