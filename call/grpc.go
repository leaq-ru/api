package call

import (
	"github.com/nnqq/scr-api/config"
	"github.com/nnqq/scr-api/logger"
	"github.com/nnqq/scr-proto/codegen/go/billing"
	"github.com/nnqq/scr-proto/codegen/go/user"
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
