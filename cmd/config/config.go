package config

import "github.com/Solar-2020/GoUtils/common"

var (
	Config config
)

type config struct {
	common.SharedConfig
	PostsDataBaseConnectionString  string `envconfig:"POSTS_DB_CONNECTION_STRING" default:"-"`
	UploadDataBaseConnectionString string `envconfig:"UPLOAD_DB_CONNECTION_STRING" default:"-"`
	DomainName                     string `envconfig:"DOMAIN_NAME" default:"solar.ru"` //for static file prefix
	PhotoPath                      string `envconfig:"PHOTO_PATH" default:"/storage/photos"`
	FilePath                       string `envconfig:"FILE_PATH" default:"/storage/files"`
	ServerSecret                   string `envconfig:"SERVER_SECRET" default:"Basic secret"`
	PaymentServiceAddress          string `envconfig:"PAYMENT_SERVICE_ADDRESS" default:"develop.pay-together.ru"`
	GroupServiceHost               string `envconfig:"GROUP_SERVICE_HOST" default:"develop.pay-together.ru"`
	AccountServiceHost             string `envconfig:"ACCOUNT_SERVICE_HOST" default:"develop.pay-together.ru"`
}
