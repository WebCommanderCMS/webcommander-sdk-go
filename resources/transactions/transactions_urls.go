package transactions

// The path templates; the client fills {api_version}.
const (
	TransactionsPath = "/api/{api_version}/admin/transactions"

	TransactionPath = "/api/{api_version}/admin/transactions/{transaction_no}"

	TransactionItemsPath = "/api/{api_version}/admin/transactions/{transaction_no}/items"
)

// TransactionsKey is the key the listing's rows arrive under.
const TransactionsKey = "transactions"

// TransactionItemsKey is the key a number's records of one type arrive under.
const TransactionItemsKey = "items"

// TransactionsPageSize is the page size applied when no limit is sent on either listing.
const TransactionsPageSize = 20

// TransactionsMaxLimit is the largest limit either listing applies: the number listing holds a larger one to it, the items listing refuses one with 400.
const TransactionsMaxLimit = 200
