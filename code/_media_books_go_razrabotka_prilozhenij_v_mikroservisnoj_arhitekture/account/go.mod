module github.com/yuliapopova/book_all/account

go 1.24

require (
	github.com/lib/pq v1.10.9
	github.com/rs/zerolog v1.31.0
	github.com/segmentio/kafka-go v0.4.40
	github.com/yuliapopova/book_all/contracts/account v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.67.0
)

require (
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

replace github.com/yuliapopova/book_all/contracts/transaction => ../contracts/transaction

replace github.com/yuliapopova/book_all/contracts/account => ../contracts/account

replace github.com/yuliapopova/book_all/contracts/auth => ../contracts/auth
