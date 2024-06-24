package constant

type DbType string
type OssType string

const (
	MySql    DbType = "mysql"
	Postgres DbType = "postgres"
)

const (
	aliyun OssType = "aliyun"
)
