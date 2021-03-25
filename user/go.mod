module github.com/desfpc/Wishez_User

go 1.16

replace github.com/desfpc/Wishez_DB => ../db

replace github.com/desfpc/Wishez_Type => ../types

replace github.com/desfpc/Wishez_Authorize => ../authorize

require (
	github.com/desfpc/Wishez_Authorize v0.0.0-00010101000000-000000000000
	github.com/desfpc/Wishez_DB v0.0.0-00010101000000-000000000000
	github.com/desfpc/Wishez_Type v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
)
