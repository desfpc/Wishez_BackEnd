module github.com/desfpc/Wishez_User

go 1.16

replace github.com/desfpc/Wishez_DB => ../db

replace github.com/desfpc/Wishez_Type => ../types

require (
	github.com/desfpc/Wishez_DB v0.0.0-00010101000000-000000000000
	github.com/desfpc/Wishez_Type v0.0.0-00010101000000-000000000000
	github.com/mitchellh/mapstructure v1.4.1
)