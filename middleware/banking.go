package middleware

import (
	"github.com/aeolus3000/lendo-sdk/banking"
	"github.com/aeolus3000/lendo-sdk/banking/dnb"
	"github.com/gobuffalo/buffalo"
)


var (
	Banking banking.BankingApi
	CtxBanking = "banking"
	configBankingPrefix = "BANKINGCONF"
)

func init() {
	Banking = dnb.NewDnbBanking(dnb.NewDnbDefaultConfiguration())
}

func BankingMiddleware(api banking.BankingApi) buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			c.Set(CtxBanking, api)
			err := next(c)
			return err
		}
	}
}