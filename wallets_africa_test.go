package gowalletsafrica

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var client *WalletsAfrica

func TestMain(m *testing.M) {
	t := new(testing.T)
	mockAPIServer := MockAPIServer(t)

	client, _ = New(DefaultConfig)
	client.Self.APIURL = mockAPIServer.URL

	os.Exit(m.Run())
}

//Self Tests
func TestSelf_CheckBalance(t *testing.T) {
	r, _ := client.Self.CheckBalance(CurrencyNigeria)
	assert.Equal(t, 880.16, r.WalletBalance)
	assert.Equal(t, "NGN", r.WalletCurrency)
}

func TestSelf_Transactions(t *testing.T) {
	transactions, _ := client.Self.Transactions(CurrencyNigeria, TransactionTypeAll, 1, 0, "2020-01-23", "")
	assert.Equal(t, 2, len(transactions))
	assert.Equal(t, 1.0, transactions[0].Amount)
	assert.Equal(t, "Credit", transactions[0].Type)

	//Test Take Validation
	transactions, err := client.Self.Transactions(CurrencyNigeria, TransactionTypeAll, 0, 0, "", "")
	assert.Empty(t, transactions)
	assert.NotNil(t, err)

	//Test Date Validations
	transactions, err = client.Self.Transactions(CurrencyNigeria, TransactionTypeAll, 1, 0, "2020-23-10", "")
	assert.Empty(t, transactions)
	assert.NotNil(t, err)
}

func TestSelf_GetWallets(t *testing.T) {
	wallets, _ := client.Self.GetWallets()
	assert.Equal(t, 2, len(wallets))
	assert.Equal(t, "22231485915", wallets[0].BVN)
	assert.Equal(t, "Odekuma", wallets[0].LastName)
}

//Identity Tests
func TestIdentity_ResolveBVN(t *testing.T) {
	bvnData, _ := client.Identity.ResolveBVN("22231485915")
	assert.Equal(t, "JOHN", bvnData.FirstName)
	assert.Equal(t, "DOE", bvnData.LastName)
}

func TestIdentity_ResolveBVNDetails(t *testing.T) {
	bvnData, _ := client.Identity.ResolveBVNDetails("22231485915")
	assert.Equal(t, "JOHN", bvnData.FirstName)
	assert.Equal(t, "DOE", bvnData.LastName)
}

//Airtime Tests
func TestAirtime_GetProviders(t *testing.T) {
	providers, _ := client.Airtime.GetProviders()
	assert.Equal(t, 4, len(providers))
	assert.Equal(t, "airtel", providers[0].Code)
	assert.Equal(t, "Airtel", providers[0].Name)
}

func TestGetResponseCodeAndMessage(t *testing.T) {
	r := responseBody{
		"Response": map[string]interface{}{
			"ResponseCode": "200",
			"Message":      "Balance Retrieved successfully",
		},
	}

	assert.Equal(t, "200", client.Self.getResponseCode(r))
	assert.Equal(t, "Balance Retrieved successfully", client.Self.getResponseMessage(r))
}

//Payouts Tests
func TestPayouts_GetBanks(t *testing.T) {
	banks, _ := client.Payouts.GetBanks()
	assert.Equal(t, 2, len(banks))
	assert.Equal(t, "044", banks[0].BankCode)
	assert.Equal(t, "000017", banks[1].BankSortCode)
}

func TestPayouts_BankDetails(t *testing.T) {
	details, _ := client.Payouts.BankDetails("2578615312")
	assert.Equal(t, "Gtbank Plc", details.Bank)
	assert.Equal(t, "0200556677", details.AccountNumber)
	assert.Equal(t, 10.00, details.Amount)
}

//StartServer initializes a test HTTP server useful for request mocking, Integration tests and Client configuration
func MockAPIServer(t *testing.T) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		switch r.URL.Path {

		case "/self/balance":
			successBody := `{
"Response": {"ResponseCode": "200","Message": "Balance Retrieved successfully"},
"Data": {"WalletBalance": 880.16,"WalletCurrency": "NGN"}
}`
			w.WriteHeader(200)
			fmt.Fprintf(w, successBody)

		case "/self/transactions":
			successBody := `{
    "Response": {
        "ResponseCode": "200",
        "Message": "Transactions Retrieved successfully"
    },
    "Data": {
        "Transactions": [
            {
                "Amount": 1.00,
                "Currency": "NGN",
                "Category": "Wallet Transfer",
                "Narration": "Sent money to Eduvie Agada",
                "DateTransacted": "7/18/2020 6:28:59 PM",
                "PreviousBalance": 7806789.16,
                "NewBalance": 7806790.16,
                "Type": "Credit"
            },
            {
                "Amount": 4250.00,
                "Currency": "NGN",
                "Category": "Dollar Card Withdrawal",
                "Narration": "Dollar Card Withdrawal at exchange rate: 425",
                "DateTransacted": "7/18/2020 11:38:27 AM",
                "PreviousBalance": 7803412.78,
                "NewBalance": 7807662.78,
                "Type": "Credit"
            }
        ]
    }
}`
			w.WriteHeader(200)
			fmt.Fprintf(w, successBody)

		case "/self/users":
			successBody := `{
    "Response": {
        "ResponseCode": "200",
        "Message": "Transactions Retrieved successfully"
    },
    "Data": [
		{
			"Username": null,
			"AccountNumber": "1023236949",
			"BVN": "22231485915",
			"City": null,
			"Country": null,
			"DateCreated": "2020-01-15T11:51:29.207",
			"DateOfBirth": "01-JAN-1990",
			"Email": "okiemuteodekuma@gmail.com",
			"FirstName": "Okiemute",
			"LastName": "Odekuma",
			"PhoneNumber": "2348057998539",
			"AvailableBalance": 3396.00
		},
		{
			"Username": "jCobhams",
			"AccountNumber": null,
			"BVN": null,
			"City": null,
			"Country": null,
			"DateCreated": "2020-01-15T15:00:30.867",
			"DateOfBirth": null,
			"Email": "brucewayne@wayneenterprises.com",
			"FirstName": "Bruce",
			"LastName": "Wayne",
			"PhoneNumber": "10706391833",
			"AvailableBalance": 0.00
		}
	]
}`
			w.WriteHeader(200)
			fmt.Fprintf(w, successBody)

		case "/account/resolvebvn":
			successBody := `{
    "FirstName": "JOHN",
    "LastName": "DOE",
    "MiddleName": null,
    "Email": "test@example.com",
    "PhoneNumber": "0706657415",
    "BVN": "22231485915",
    "DateOfBirth": "11-04-1992",
    "EnrollmentBank": "Access Bank",
    "EnrollmentBranch": "Heaven",
    "Gender": "Male",
    "LevelOfAccount": null,
    "LgaOfOrigin": null,
    "LgaOfResidence": null,
    "MaritalStatus": "Married",
    "NameOnCard": null,
    "Nationality": null,
    "StateOfOrigin": null,
    "StateOfResidence": null,
    "Title": "Chief",
    "WatchListed": null,
    "Picture": null,
    "ResponseCode": "200",
    "Message": "Successful"
}`
			w.WriteHeader(200)
			fmt.Fprintf(w, successBody)

		case "/bills/airtime/providers":
			successBody := `{
    "ResponseCode": "200",
    "Providers": [
        {
            "Code": "airtel",
            "Name": "Airtel"
        },
        {
            "Code": "mtn",
            "Name": "MTN"
        },
        {
            "Code": "glo",
            "Name": "GLO"
        },
        {
            "Code": "etisalat",
            "Name": "Etisalat"
        }
    ]
}`
			w.WriteHeader(200)
			fmt.Fprintf(w, successBody)

		case "/transfer/banks/all":
			successBody := `[
    {
        "BankCode": "044",
        "BankName": "Access Bank Nigeria",
        "BankSortCode": "000014",
        "PaymentGateway": null
    },
    {
        "BankCode": "035A",
        "BankName": "Alat By Wema",
        "BankSortCode": "000017",
        "PaymentGateway": null
    }]`
			w.WriteHeader(200)
			fmt.Fprintf(w, successBody)

		case "/transfer/bank/details":
			successBody := `{
    "Bank": "Gtbank Plc",
    "AccountNumber": "0200556677",
    "DateTransferred": "1/15/2020 1:45:31 PM",
    "Amount": 10.00,
    "RecipientName": "JOHN DOE",
    "SessionId": null,
    "ResponseCode": "200",
    "Message": null
}`
			w.WriteHeader(200)
			fmt.Fprintf(w, successBody)

		}
	}))
	return server
}
