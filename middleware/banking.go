package middleware

import (
	"github.com/aeolus3000/lendo-sdk/banking"
	"github.com/aeolus3000/lendo-sdk/banking/dnb"
	"github.com/aeolus3000/lendo-sdk/configuration"
	"github.com/gobuffalo/buffalo"
	log "github.com/sirupsen/logrus"
)


var (
	Banking banking.BankingApi
	CtxBanking = "banking"
	configBankingPrefix = "BANKINGCONF"
)

func init() {
	config := configuration.NewDefaultConfiguration()
	bankingConfiguration := banking.Configuration{}
	confErr := config.Process(configPrefix + configBankingPrefix, &bankingConfiguration)
	if confErr != nil {
		log.Fatal(confErr)
	}
	Banking = dnb.NewDnbBanking(banking.Configuration{})
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