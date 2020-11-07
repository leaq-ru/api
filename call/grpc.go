package call

import (
	"github.com/nnqq/scr-api/config"
	"github.com/nnqq/scr-api/logger"
	"github.com/nnqq/scr-proto/codegen/go/user"
	"google.golang.org/grpc"
)

var User user.UserClient

func init() {
	connUser, err := grpc.Dial(config.Env.Service.User, grpc.WithInsecure())
	logger.Must(err)
	User = user.NewUserClient(connUser)
}
