package futures

import (
	"errors"
	"okx-bot/exchange/model"
)

type IsolatedPrvApi struct {
	*PrvApi
}

func (f *IsolatedPrvApi) CreateOrder(pair model.CurrencyPair, qty, price float64, side model.OrderSide, orderTy model.OrderType, opts ...model.OptionParameter) (*model.Order, []byte, error) {
	if side != model.Futures_OpenBuy &&
		side != model.Futures_OpenSell &&
		side != model.Futures_CloseBuy &&
		side != model.Futures_CloseSell {
		return nil, nil, errors.New("futures side only is Futures_OpenBuy or Futures_OpenSell or Futures_CloseBuy or Futures_CloseSell")
	}

	opts = append(opts,
		model.OptionParameter{
			Key:   "tdMode",
			Value: "isolated",
		})

	return f.Prv.CreateOrder(pair, qty, price, side, orderTy, opts...)
}

func (f *IsolatedPrvApi) PlaceOrder(order model.PlaceOrderRequest, opts ...model.OptionParameter) (*model.Order, []byte, error) {

	opts = append(opts,
		model.OptionParameter{
			Key:   "tdMode",
			Value: "isolated",
		})

	return f.Prv.PlaceOrder(order, opts...)
}

func (f *IsolatedPrvApi) SetLeverage(order model.SetLeverageRequest, opts ...model.OptionParameter) (*model.SetLeverageResponse, []byte, error) {

	opts = append(opts,
		model.OptionParameter{
			Key:   "tdMode",
			Value: "isolated",
		})

	return f.Prv.SetLeverage(order, opts...)
}

func (f *IsolatedPrvApi) CancelOrder(req *model.BaseOrderRequest, opts ...model.OptionParameter) (model.CancelOrderResponse, []byte, error) {

	opts = append(opts,
		model.OptionParameter{
			Key:   "tdMode",
			Value: "isolated",
		})

	return f.Prv.CancelOrder(req, opts...)
}

func (f *IsolatedPrvApi) ClosePositions(req *model.ClosePositionsRequest, opts ...model.OptionParameter) (model.ClosePositionResponse, []byte, error) {

	opts = append(opts,
		model.OptionParameter{
			Key:   "mgnMode",
			Value: "isolated",
		})

	return f.Prv.ClosePositions(req, opts...)
}
