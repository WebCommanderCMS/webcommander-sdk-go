package transactions

import (
	"context"
	"fmt"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
)

// Transactions is the transaction-number endpoints, bound to one client; every path is read-only, a write is 405.
type Transactions struct {
	requester http.Requester
}

// ListTransactions lists the numbers in ascending order — GET /admin/transactions; limit or max from 1, held to 200, and transaction_no matches exactly.
func (c *Transactions) ListTransactions(ctx context.Context, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendGetRequest(ctx, TransactionsPath, &http.RequestOptions{
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: TransactionPageShape})
}

// HeadTransactions asks whether the collection is reachable — HEAD /admin/transactions; no body.
func (c *Transactions) HeadTransactions(ctx context.Context, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendApiRequest(ctx, "HEAD", TransactionsPath, &http.RequestOptions{
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// GetTransaction counts the records one number tagged, by type — GET /admin/transactions/{transaction_no}; matched exactly and case-sensitively.
func (c *Transactions) GetTransaction(ctx context.Context, transactionNo string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendGetRequest(ctx, TransactionPath, &http.RequestOptions{
		PathParams:   map[string]any{"transaction_no": transactionNo},
		RaiseOnError: raiseOnError, ResponseFormat: format, DTOShape: TransactionShape})
}

// HeadTransaction asks whether one number has records — HEAD /admin/transactions/{transaction_no}; no body.
func (c *Transactions) HeadTransaction(ctx context.Context, transactionNo string, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendApiRequest(ctx, "HEAD", TransactionPath, &http.RequestOptions{
		PathParams:   map[string]any{"transaction_no": transactionNo},
		RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ListTransactionItems lists the records under a transaction number, newest first — GET /admin/transactions/{transaction_no}/items; queryParams may be nil, as type is optional (left out, every type; sent, one or more names comma-separated), limit 1 to 200, detail slim or full, and an unknown number or type returns an empty page.
func (c *Transactions) ListTransactionItems(ctx context.Context, transactionNo string, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendGetRequest(ctx, TransactionItemsPath, &http.RequestOptions{
		PathParams:  map[string]any{"transaction_no": transactionNo},
		QueryParams: queryParams, RaiseOnError: raiseOnError,
		ResponseFormat: format, DTOShape: TransactionItemPageShape})
}

// HeadTransactionItems asks whether the records under a number are reachable — HEAD /admin/transactions/{transaction_no}/items; no body, and queryParams may be nil, as type is optional here too.
func (c *Transactions) HeadTransactionItems(ctx context.Context, transactionNo string, queryParams map[string]any, raiseOnError *bool, responseFormat ...string) (*http.ApiResponse, error) {
	format, err := oneResponseFormat(responseFormat)
	if err != nil {
		return nil, err
	}
	return c.requester.SendApiRequest(ctx, "HEAD", TransactionItemsPath, &http.RequestOptions{
		PathParams:  map[string]any{"transaction_no": transactionNo},
		QueryParams: queryParams, RaiseOnError: raiseOnError, ResponseFormat: format})
}

// ResourceName is this resource's registry key — what ListResourceNames returns, snake_case.
const ResourceName = "transactions"

func oneResponseFormat(responseFormat []string) (string, error) {
	switch len(responseFormat) {
	case 0:
		return "", nil
	case 1:
		return responseFormat[0], nil
	default:
		return "", wcerrors.NewUsageError(fmt.Sprintf(
			"responseFormat: pass one format or none, got %d — the trailing argument is optional, "+
				"and omitting it defers to the client", len(responseFormat)))
	}
}

// New binds the transaction-number endpoints to any requester (the WebCommanderClient).
func New(requester http.Requester) *Transactions {
	return &Transactions{requester: requester}
}
