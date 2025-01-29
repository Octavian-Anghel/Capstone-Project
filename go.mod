module github.com/your-org/your-repo

go 1.21

require (
	github.com/hyperledger/fabric v2.5.0
	github.com/hyperledger/fabric-sdk-go v1.0.0
	github.com/hyperledger/fabric-protos-go v0.4.0
	github.com/golangci/golangci-lint v1.55.2
	github.com/gofiber/fiber/v2 v2.48.0 // Example for REST API
	github.com/stretchr/testify v1.8.0   // For unit testing
	github.com/securego/gosec/v2 v2.16.0 // Security scanner
	honnef.co/go/tools v0.12.0           // Staticcheck
)

replace (
	github.com/hyperledger/fabric => github.com/hyperledger/fabric v2.5.0
	github.com/hyperledger/fabric-sdk-go => github.com/hyperledger/fabric-sdk-go v1.0.0
)
