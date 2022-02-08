package coinpayments

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

type TransactionService struct {
	sling        *sling.Sling
	ApiPublicKey string
	Params       TransactionBodyParams
	FindParams   TransactionFindBodyParams
}

type Transaction struct {
	Amount         string `json:"amount"`
	Address        string `url:"address"`
	TXNId          string `json:"txn_id"`
	ConfirmsNeeded string `json:"confirms_needed"`
	Timeout        uint32 `json:"timeout"`
	StatusUrl      string `json:"status_url"`
	QRCodeUrl      string `json:"qrcode_url"`
}

type TransactionFind struct {
	Amount     string `json:"amountf"`
	Address    string `json:"payment_address"`
	Coin       string `json:"coin"`
	Status     int    `json:"status"`
	StatusText string `json:"status_text"`
}

type TransactionResponse struct {
	Error  string       `json:"error"`
	Result *Transaction `json:"result"`
}

type TransactionFindResponse struct {
	Error  string           `json:"error"`
	Result *TransactionFind `json:"result"`
}

type TransactionParams struct {
	Amount     float64 `url:"amount"`
	Currency1  string  `url:"currency1"`
	Currency2  string  `url:"currency2"`
	Address    string  `url:"address"`
	BuyerEmail string  `url:"buyer_email"`
	BuyerName  string  `url:"buyer_name"`
	ItemName   string  `url:"item_name"`
	ItemNumber string  `url:"item_number"`
	Invoice    string  `url:"invoice"`
	Custom     string  `url:"custom"`
	IPNUrl     string  `url:"ipn_url"`
}

type TransactionBodyParams struct {
	APIParams
	TransactionParams
}

type TransactionFindBodyParams struct {
	APIParams
}

func newTransactionService(sling *sling.Sling, apiPublicKey string) *TransactionService {
	transactionService := &TransactionService{
		sling:        sling.Path("api.php"),
		ApiPublicKey: apiPublicKey,
	}
	transactionService.getParams()
	return transactionService
}

func (s *TransactionService) getHMAC(params interface{}) string {
	return getHMAC(getPayload(params))
}

func (s *TransactionService) CreateTransaction(transactionParams *TransactionParams) (TransactionResponse, *http.Response, error) {
	transactionResponse := new(TransactionResponse)
	s.Params.TransactionParams = *transactionParams
	fmt.Println(getPayload(s.Params))
	fmt.Println(getHMAC(getPayload(s.Params)))
	resp, err := s.sling.New().Set("HMAC", s.getHMAC(s.Params)).Post(
		"api.php").BodyForm(s.Params).ReceiveSuccess(transactionResponse)
	return *transactionResponse, resp, err
}

func (s TransactionService) FindTransaction(txnID string) (TransactionFindResponse, *http.Response, error) {
	transactionResponse := new(TransactionFindResponse)
	s.FindParams.TXNID = txnID

	fmt.Println(getPayload(s.FindParams))
	fmt.Println(getHMAC(getPayload(s.FindParams)))
	resp, err := s.sling.New().Set("HMAC", s.getHMAC(s.FindParams)).Post(
		"api.php").BodyForm(s.FindParams).ReceiveSuccess(transactionResponse)

	return *transactionResponse, resp, err
}

func (s *TransactionService) getParams() {
	s.Params.Command = "create_transaction"
	s.Params.Key = s.ApiPublicKey
	s.Params.Version = "1"

	s.FindParams.Command = "get_tx_info"
	s.FindParams.Key = s.Params.Key
	s.FindParams.Version = s.Params.Version
}
