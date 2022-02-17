package coinpayments

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

type WithdrawalService struct {
	sling        *sling.Sling
	ApiPublicKey string
	Params       WithdrawalBodyParams
	FindParams   WithdrawalFindBodyParams
}

type Withdrawal struct {
	ID     string `json:"id"`
	Status int    `json:"status"`
	Amount string `json:"amount"`
}

type WithdrawalFind struct {
	Amount     string `json:"amountf"`
	Address    string `json:"payment_address"`
	Coin       string `json:"coin"`
	Status     int    `json:"status"`
	StatusText string `json:"status_text"`
	Note       string `json:"note"`
}

type WithdrawalResponse struct {
	Error  string      `json:"error"`
	Result *Withdrawal `json:"result"`
}

type WithdrawalFindResponse struct {
	Error  string          `json:"error"`
	Result *WithdrawalFind `json:"result"`
}

type WithdrawalParams struct {
	Amount      float64 `url:"amount"`
	AddTXFree   int     `url:"add_tx_free"`
	Currency    string  `url:"currency"`
	Currency2   string  `url:"currency2"`
	Address     string  `url:"address"`
	IPNUrl      string  `url:"ipn_url"`
	AutoConfirm int     `url:"auto_confirm"`
	Note        string  `url:"note"`
}

type WithdrawalBodyParams struct {
	APIParams
	WithdrawalParams
}

type WithdrawalFindBodyParams struct {
	APIParams
}

func newWithdrawalService(sling *sling.Sling, apiPublicKey string) *WithdrawalService {
	transactionService := &WithdrawalService{
		sling:        sling.Path("api.php"),
		ApiPublicKey: apiPublicKey,
	}
	transactionService.getParams()
	return transactionService
}

func (s *WithdrawalService) getHMAC(params interface{}) string {
	return getHMAC(getPayload(params))
}

func (s *WithdrawalService) CreateWithdrawal(transactionParams *WithdrawalParams) (WithdrawalResponse, *http.Response, error) {
	transactionResponse := new(WithdrawalResponse)
	s.Params.WithdrawalParams = *transactionParams
	fmt.Println(getPayload(s.Params))
	fmt.Println(getHMAC(getPayload(s.Params)))
	resp, err := s.sling.New().Set("HMAC", s.getHMAC(s.Params)).Post(
		"api.php").BodyForm(s.Params).ReceiveSuccess(transactionResponse)

	return *transactionResponse, resp, err
}

func (s WithdrawalService) FindWithdrawal(txnID string) (WithdrawalFindResponse, *http.Response, error) {
	transactionResponse := new(WithdrawalFindResponse)
	s.FindParams.TxnID = txnID
	fmt.Println(getPayload(s.FindParams))
	fmt.Println(getHMAC(getPayload(s.FindParams)))

	resp, err := s.sling.New().Set("HMAC", s.getHMAC(s.FindParams)).Post(
		"api.php").BodyForm(s.FindParams).ReceiveSuccess(transactionResponse)

	return *transactionResponse, resp, err
}

func (s *WithdrawalService) getParams() {
	s.Params.Command = "create_withdrawal"
	s.Params.Key = s.ApiPublicKey
	s.Params.Version = "1"

	s.FindParams.Command = "get_withdrawal_info"
	s.FindParams.Key = s.Params.Key
	s.FindParams.Version = s.Params.Version
}
