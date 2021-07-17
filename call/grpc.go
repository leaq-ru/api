package call

import (
	"github.com/leaq-ru/api/config"
	"github.com/leaq-ru/api/logger"
	"github.com/leaq-ru/proto/codegen/go/billing"
	"github.com/leaq-ru/proto/codegen/go/user"
	"google.golang.org/grpc"
)

var (
	User    user.UserClient
	Billing billing.BillingClient
)

func init() {
	connUser, err := grpc.Dial(config.Env.Service.User, grpc.WithInsecure())
	logger.Must(err)
	User = user.NewUserClient(connUser)

	connBilling, err := grpc.Dial(config.Env.Service.Billing, grpc.WithInsecure())
	logger.Must(err)
	Billing = billing.NewBillingClient(connBilling)
}
