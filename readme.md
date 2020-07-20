#Wallets Africa Go API Wrapper
[![Build Status](https://travis-ci.org/jcobhams/gowalletsafrica.svg?branch=master)](https://travis-ci.org/jcobhams/gowalletsafrica)
[![codecov](https://codecov.io/gh/jcobhams/gowalletsafrica/branch/master/graph/badge.svg)](https://codecov.io/gh/jcobhams/gowalletsafrica)


### Installation
`$ go get github.com/jcobhams/gowalletsafrica`

### Usage

### Concerns
* `Payouts.GetBanks()` ignores the `PaymentGateway` field of the result since we don't know what the data structure could possibly be.
To avoid a runtime panic if wallets.africa ever returns something else apart from `null`.

### Not Covered
* `Self - Verify BVN`: This implementation is a bit confusing. The verify BVN endpoint performs an update operation. 
`https://documenter.getpostman.com/view/10058163/SWLk4RPL?version=latest#feb190a5-53e2-45b7-84a5-77e11ea341a0`

* `Self - Wallet To Wallet Transactions`: Could not test this endpoint on POSTMAN. Can't implement.
`https://documenter.getpostman.com/view/10058163/SWLk4RPL?version=latest#44de9ef6-c97b-498c-8074-4b2c3c76c706`

* `Airtime - Purchase`: The API documentation is not very helpful and makes it a bit hard to design/test the function.
`https://documenter.getpostman.com/view/10058163/SWLk4RPL?version=latest#f698015a-71a5-4fe6-8c24-6677d530baa0` 

### Run Tests
`$ go test -v ./... -coverprofile cover.out`

### View Coverage
`$ go tool cover -html=cover.out`