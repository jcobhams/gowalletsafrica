package gowalletsafrica

import (
	"net/http"
	"time"
)

type (
	Currency        string
	TransactionType int

	base struct {
		HTTPClient *http.Client
		APIURL     string
		secretKey  string
		publicKey  string
	}

	self struct {
		*base
	}

	wallets struct {
		*base
	}

	payouts struct {
		*base
	}

	airtime struct {
		*base
	}

	identity struct {
		*base
	}

	WalletsAfrica struct {
		Self     *self
		Wallets  *wallets
		Payouts  *payouts
		Airtime  *airtime
		Identity *identity
	}

	Config struct {
		Environment    string
		PublicKey      string
		SecretKey      string
		RequestTimeout time.Duration
	}

	Transaction struct {
		Amount          float64
		Currency        string
		Category        string
		Narration       string
		DateTransacted  string
		PreviousBalance float64
		NewBalance      float64
		Type            string
	}

	Wallet struct {
		Username      string
		AccountNumber string
		BVN           string
		City          string
		Country       string
		DateCreated   string
		DateOfBirth   string
		Email         string
		FirstName     string
		LastName      string
		PhoneNumber   string
	}

	AirtimeProvider struct {
		Code string
		Name string
	}

	Bank struct {
		BankCode     string
		BankName     string
		BankSortCode string
		//PaymentGateway string
	}

	BankDetail struct {
		Bank            string
		AccountNumber   string
		DateTransferred string
		Amount          float64
		RecipientName   string
		SessionId       string
		ResponseCode    string
		Message         string
	}

	payloadBody  map[string]interface{}
	responseBody map[string]interface{}

	//Endpoint Results
	CheckBalanceResult struct {
		WalletBalance  float64
		WalletCurrency string
	}

	Transactions     []Transaction
	Wallets          []Wallet
	AirtimeProviders []AirtimeProvider
	Banks            []Bank

	ResolveBVN struct {
		FirstName        string
		LastName         string
		MiddleName       string
		Email            string
		PhoneNumber      string
		BVN              string
		DateOfBirth      string
		EnrollmentBank   string
		EnrollmentBranch string
		Gender           string
		LevelOfAccount   string
		LgaOfOrigin      string
		LgaOfResidence   string
		MaritalStatus    string
		NameOnCard       string
		Nationality      string
		StateOfOrigin    string
		StateOfResidence string
		Title            string
		WatchListed      string
		Picture          string
		ResponseCode     string
		Message          string
	}
)
