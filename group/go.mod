module github.com/desfpc/Wishez_Group

go 1.16

replace github.com/desfpc/Wishez_Type => ../types

replace github.com/desfpc/Wishez_Helpers => ../helpers

replace github.com/desfpc/Wishez_DB => ../db

replace github.com/desfpc/Wishez_User => ../user

require (
	github.com/desfpc/Wishez_DB v0.0.0-00010101000000-000000000000
	github.com/desfpc/Wishez_Helpers v0.0.0-00010101000000-000000000000
	github.com/desfpc/Wishez_Type v0.0.0-00010101000000-000000000000
	github.com/desfpc/Wishez_User v0.0.0-00010101000000-000000000000
)
