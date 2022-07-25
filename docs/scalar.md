# Scalar

The `scalar` package holds constant data in the style of enums, though most primarily indexed as strings rather than integers.  It inherits the name "scalar" from the [graphql resource for replacing types](https://github.com/graphql-go/graphql/blob/master/scalars.go).  This helps the client/user understand the limitation of certain requests.  For example, the coinbase http request for [creating a new order](https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_postorders) requires a "side" body param.  At the time of this README, "side" is not defined in their docs.  The scalar package removes this ambiguity completely, to create an order request using the coinbase web you must use the `model.CoinbaseNewOrderOptions` struct, which types the `Side` field as `scalar.OderSide` which only has two defined constants: `OrderSideBuy` and `OrderSideSell`.

_TL;DR package removes the ambiguity from questionable body/query parameters in http requests._
