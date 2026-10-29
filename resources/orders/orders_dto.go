package orders

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// Reference is a country or a state as the store answers it: {"id", "name", "code"}.
type Reference struct {
	ID   *int64
	Name *string
	Code *string
}

// OrderCustomer is the order's customer section, also sent as customer_summary; a listing row carries five of these keys.
type OrderCustomer struct {
	FirstName     *string
	LastName      *string
	CustomerGroup *string
	Gender        *string
	// CustomerID is the public id; InternalID is the one the listing's customerId filter takes.
	CustomerID   *int64
	InternalID   *int64
	Email        *string
	AddressLine1 *string
	AddressLine2 *string
	City         *string
	Country      *Reference
	State        *Reference
	PostCode     *string
	Phone        *string
	Mobile       *string
	Fax          *string
	CompanyName  *string
}

// OrderAddress is a billing or shipping address; CompanyName reads "" whatever was written.
type OrderAddress struct {
	ID           *int64
	FirstName    *string
	LastName     *string
	Email        *string
	CompanyName  *string
	AddressLine1 *string
	AddressLine2 *string
	Country      *Reference
	State        *Reference
	City         *string
	PostCode     *string
	Mobile       *string
	Phone        *string
	Fax          *string
}

// OrderItem is one order line, under items or order_line_details.
type OrderItem struct {
	ItemID      *int64
	ProductName *string
	ProductID   *int64
	SKU         *string
	Quantity    *int64
	Price       *float64
	TotalAmount *float64
	Tax         *float64
	Discount    *float64
	TaxDiscount *float64
	IsTaxable   *bool
	IsShippable *bool
	ImageURL    *string
	// Variations is [] everywhere except a create's answer, which sends null.
	Variations        []any
	SerialBatchNumber *string
}

// OrderShipping is the shipping section; ShippingType is the first shipment's method.
type OrderShipping struct {
	DeliveryType   *string
	ShippingType   *string
	ShippingStatus *string
	ShippingCost   *float64
	ShippingTax    *float64
}

// OrderDiscounts is the discounts section.
type OrderDiscounts struct {
	TotalDiscount *float64
	CouponCode    *string
	GiftCardCode  *string
}

// OrderChargesTotals is the charges_totals section; a payment refund lowers Paid and an order refund does not.
type OrderChargesTotals struct {
	SubTotal       *float64
	TotalTax       *float64
	HandlingCost   *float64
	TotalSurcharge *float64
	OtherCharges   *float64
	GrandTotal     *float64
	Paid           *float64
	Due            *float64
	Refunded       *float64
	RefundStatus   *string
	ItemsTotal     *int64
}

// OrderPayment is one entry of an order's payments; GatewayResponse comes only in the flat shape.
type OrderPayment struct {
	ID                   *int64
	Status               *string
	Amount               *float64
	PaymentMethod        *string
	TrackInfo            *string
	GatewayName          *string
	GatewayCode          *string
	GatewayResponse      *string
	PayerInfo            *string
	PayingDate           *string
	TransactionDate      *string
	TransactionReference *string
}

// OrderFulfillment is the fulfillment section.
type OrderFulfillment struct {
	FulfillmentMode          *string
	FulfillmentTerm          *string
	InvoiceMode              *string
	PaymentTerm              *string
	CustomerPONumber         *string
	ExternalReferenceOrderID *string
	ExternalReferenceStoreID *string
	ExternalReferenceOrderNo *string
	ExternalReferenceStatus  *string
	ExternalReferenceURL     *string
	CustomFields             []any
}

// SectionMeta is one group_meta entry.
type SectionMeta struct {
	Kind     *string
	Writable *bool
	Href     *string
}

// GroupMeta is the group_meta block; the listing's has no entry for the two addresses.
type GroupMeta struct {
	Customer        *SectionMeta
	BillingAddress  *SectionMeta
	ShippingAddress *SectionMeta
	Items           *SectionMeta
	Shipping        *SectionMeta
	Discounts       *SectionMeta
	ChargesTotals   *SectionMeta
	Payments        *SectionMeta
	Fulfillment     *SectionMeta
}

// Order is one order in sections: the row path's body, a PUT or PATCH answer, and a listing row, which carries no addresses.
type Order struct {
	OrderID         *int64
	OrderNo         *int64
	InternalID      *int64
	UUID            *string
	OrderStatus     *string
	PaymentStatus   *string
	ShippingStatus  *string
	CurrencyCode    *string
	TransactionNo   *string
	IPAddress       *string
	OrderChannel    *string
	CreatedOn       *string
	LastUpdatedOn   *string
	Customer        *OrderCustomer
	BillingAddress  *OrderAddress
	ShippingAddress *OrderAddress
	Items           []*OrderItem
	Shipping        *OrderShipping
	Discounts       *OrderDiscounts
	ChargesTotals   *OrderChargesTotals
	Payments        []*OrderPayment
	Fulfillment     *OrderFulfillment
	// GroupMeta and FieldMetadata are read from beside the envelope.
	GroupMeta     *GroupMeta
	FieldMetadata any
}

// EnvelopeKey is the key the row path, PUT and PATCH wrap one order in, and the wrapper those writes read beside a bare body.
func (Order) EnvelopeKey() string { return OrderKey }

// Pagination is the listing's ten-key pagination block; Limit is the size applied, at most 20.
type Pagination struct {
	Total        *int64
	Limit        *int64
	Offset       *int64
	Count        *int64
	CurrentPage  *int64
	TotalPages   *int64
	HasNext      *bool
	HasPrevious  *bool
	PreviousPage *string
	NextPage     *string
}

// OrderPage is one page of the listing, with the field_metadata block when it was asked for.
type OrderPage struct {
	Orders        []*Order
	Pagination    *Pagination
	GroupMeta     *GroupMeta
	FieldMetadata any
}

// OrderCount is the count sub-path's {"orders": {"count": N}}.
type OrderCount struct {
	Count *int64
}

// EnvelopeKey is the plural key the count record arrives under.
func (OrderCount) EnvelopeKey() string { return OrdersKey }

// FlatOrder is one order with its totals and fulfilment at the top: what a create, a cancel and a status move answer, and a per-customer listing row.
type FlatOrder struct {
	OrderID                  *int64
	OrderNo                  *int64
	InternalID               *int64
	UUID                     *string
	TransactionNo            *string
	CurrencyCode             *string
	OrderStatus              *string
	PaymentStatus            *string
	ShippingStatus           *string
	SubTotal                 *float64
	ShippingCost             *float64
	ShippingTax              *float64
	HandlingCost             *float64
	TotalSurcharge           *float64
	TotalDiscount            *float64
	TotalTax                 *float64
	GrandTotal               *float64
	Paid                     *float64
	Due                      *float64
	ItemsTotal               *int64
	IPAddress                *string
	OrderChannel             *string
	CreatedAt                *string
	UpdatedAt                *string
	CustomerSummary          *OrderCustomer
	OrderLineDetails         []*OrderItem
	Payments                 []*OrderPayment
	ExternalReferenceOrderID *string
	ExternalReferenceStoreID *string
	ExternalReferenceOrderNo *string
	ExternalReferenceStatus  *string
	ExternalReferenceURL     *string
	CustomerPONumber         *string
	PaymentTerm              *string
	FulfillmentMode          *string
	FulfillmentTerm          *string
	InvoiceMode              *string
	BillingAddress           *OrderAddress
	ShippingAddress          *OrderAddress
	CustomFields             []any
	// Status sits beside the envelope, on a create's answer only.
	Status *string
}

// EnvelopeKey is the key a create, a cancel and a status move wrap the order in; a per-customer row has none.
func (FlatOrder) EnvelopeKey() string { return OrderKey }

// CustomerOrders is one customer's orders, unpaged, cancelled ones left out.
type CustomerOrders struct {
	CustomerID *int64
	InternalID *int64
	Orders     []*FlatOrder
}

// OrderAddresses is what the addresses sub-path answers: both addresses as they now stand.
type OrderAddresses struct {
	OrderID         *int64
	InternalID      *int64
	BillingAddress  *OrderAddress
	ShippingAddress *OrderAddress
}

// OrderComment is one comment; Created is M/D/YYYY and a time two spaces apart, not ISO 8601.
type OrderComment struct {
	ID                *int64
	Content           *string
	AdminName         *string
	VisibleToCustomer *bool
	Admin             *bool
	Created           *string
}

// EnvelopeKey is the key one comment arrives under; the comment write takes a bare record.
func (OrderComment) EnvelopeKey() string { return OrderCommentKey }

// OrderCommentList is an order's comments, oldest first and unpaged.
type OrderCommentList struct {
	Comments []*OrderComment
}

// OrderCustomFields is an order's custom field values; OrderID echoes the internal_id asked for, and a delivery order adds DeliveryCost.
type OrderCustomFields struct {
	OrderID      *int64
	CustomFields []any
	DeliveryCost *float64
}

// PaymentReceiptCustomer is the customer a payment receipt names; ID is the public customer id.
type PaymentReceiptCustomer struct {
	ID         *int64
	InternalID *int64
	Name       *string
}

// PaymentReceipt is what recording a payment answers; ID is the payment_id a payment refund takes.
type PaymentReceipt struct {
	Status               *string
	ID                   *int64
	Date                 *string
	TransactionDate      *string
	TransactionReference *string
	PaymentApplied       *float64
	PaymentMethod        *string
	PayerInfo            *string
	OrderID              *int64
	OrderInternalID      *int64
	OrderTotal           *float64
	DueAmount            *float64
	PaidAmount           *float64
	Customer             *PaymentReceiptCustomer
	CreatedBy            *string
	CreatedAt            *string
}

// EnvelopeKey is the key the receipt arrives under.
func (PaymentReceipt) EnvelopeKey() string { return OrderPaymentKey }

// OrderRefund is what an order refund answers; Amount is the text that was sent.
type OrderRefund struct {
	Status     *string
	RefundID   *int64
	OrderID    *int64
	InternalID *int64
	Amount     *string
}

// OrderRefundRecordLine is one order line a refund was spread over; the money fields are decimal text such as "1.00". ExtItemReferer was null on every line read, so it is untyped.
type OrderRefundRecordLine struct {
	OrderItemID       *int64
	ProductSKU        *string
	ProductName       *string
	ExtItemReferer    any
	PluginLine        *bool
	Quantity          *int64
	LinePaid          *string
	Amount            *string
	LineRefundedTotal *string
	LineFullyRefunded *bool
}

// OrderRefundRecord is one refund an order holds. OrderID is the order's internal id. A payment refund fills RecordOnly to Lines; an order refund leaves them nil and Lines empty. GatewayChargeID, GatewayRefundID and FailureCause were null on every refund read, so they are untyped.
type OrderRefundRecord struct {
	Gateway            *string
	GatewayChargeID    any
	GatewayRefundID    any
	RecordOnly         *bool
	InitiatedBy        *string
	CancelledWithOrder *bool
	NotifyCustomer     *bool
	Allocation         *string
	OtherAmount        *string
	Lines              []*OrderRefundRecordLine
	RefundID           *int64
	PaymentID          *int64
	OrderID            *int64
	Amount             *float64
	Status             *string
	TransactionStatus  *string
	TransactionID      *string
	Reason             *string
	FailureCause       any
	CreatedAt          *string
	UpdatedAt          *string
}

// OrderRefundPage is one page of an order's refunds, newest first. OrderID is the order's internal id and OrderNo the {order_id} the path took.
type OrderRefundPage struct {
	OrderID    *int64
	OrderNo    *int64
	Refunds    []*OrderRefundRecord
	Pagination *Pagination
}

// PaymentRefund is what a payment refund answers; a replay drops four flags and adds IdempotentReplay.
type PaymentRefund struct {
	Status    *string
	RefundID  *int64
	PaymentID *int64
	// OrderID holds the order's internal id and OrderNo its order_id.
	OrderID              *int64
	OrderNo              *int64
	Currency             *string
	Amount               *float64
	Reason               *string
	GatewayRefund        *bool
	GatewayTransactionID *string
	TransactionStatus    *string
	PaymentStatus        *string
	PaymentRefundStatus  *string
	OrderRefundStatus    *string
	OrderRefundedTotal   *float64
	OrderAdjustment      *float64
	OrderTotal           *float64
	AvailableForRefund   *float64
	AdjustWithOrder      *bool
	RefundEmailQueued    *bool
	InventoryReadjusted  *bool
	RefundedAt           *string
	IdempotentReplay     *bool
}

// InvoiceDispatch is what sending an invoice answers: who it went to and when.
type InvoiceDispatch struct {
	Status          *string
	OrderID         *int64
	InternalID      *int64
	SentTo          *string
	Template        *string
	InvoiceDocument any
	SentAt          *string
}

// ShipmentTracking is one shipment as the shipment status lists it.
type ShipmentTracking struct {
	ShippingProvider *string
	TrackingNumber   *string
	LastUpdated      *string
}

// ShipmentStatus is unshipped, partial or shipped, where the order itself reads not_shipped for the first.
type ShipmentStatus struct {
	OrderID        *int64
	InternalID     *int64
	ShipmentStatus *string
	ShipmentsInfo  []*ShipmentTracking
}

// ShipmentItem is one line a shipment carries.
type ShipmentItem struct {
	OrderItemID     *int64
	ProductID       *int64
	ProductName     *string
	ProductSKU      *string
	ShippingClass   *string
	OrderQuantity   *int64
	ShippedQuantity *int64
	TotalAmount     *float64
}

// OrderShipment is one shipment; a write's method, tracking_info and shipping_date read back as ShippingMethod, TrackingNumber and DateCreated.
type OrderShipment struct {
	ShippingID         *int64
	OrderID            *int64
	OrderInternalID    *int64
	CustomerID         *int64
	CustomerInternalID *int64
	OrderAddressID     *int64
	DateCreated        *string
	TrackingNumber     *string
	ShippingMethod     *string
	ShippingAddress    *OrderAddress
	BillingAddress     *OrderAddress
	Items              []*ShipmentItem
	// Status sits beside the envelope.
	Status *string
}

// EnvelopeKey is the key a shipment write's answer wraps the shipment in; a listed shipment has none.
func (OrderShipment) EnvelopeKey() string { return OrderShipmentKey }

// OrderShipmentList is an order's shipments, unpaged.
type OrderShipmentList struct {
	Shipments []*OrderShipment
}

// OrderAddressInput is an address to write, replacing the whole one.
type OrderAddressInput struct {
	FirstName    string
	LastName     string
	Email        string
	CompanyName  string
	AddressLine1 string
	AddressLine2 string
	City         string
	PostCode     string
	Phone        string
	Mobile       string
	Fax          string
	CountryCode  string
	StateCode    string
}

// OrderProductInput is one line to order, by the product's id.
type OrderProductInput struct {
	ProductID *int64
	Quantity  *int64
}

// OrderInput is the create body, with no envelope. PaymentGateway takes a code, PIS for pay in store, and DeliveryType is lower case; PAY_IN_STORE and SHIPPING are 400.
type OrderInput struct {
	CustomerID               *int64
	PaymentGateway           string
	DeliveryType             string
	IsSendEmail              *bool
	TransactionNo            string
	ExternalReferenceOrderID string
	BillingAddress           *OrderAddressInput
	ShippingAddress          *OrderAddressInput
	Products                 []*OrderProductInput
	OrderChannel             string
	IPAddress                string
}

// ToMap renders the request body.
func (o OrderInput) ToMap() map[string]any { return dto.ToMap(o) }

// InitialPaymentInput is the payment a create-with-payment records, filed under method API; PayerInfo, TransactionReference and TransactionDate are not kept.
type InitialPaymentInput struct {
	Amount               *float64
	Surcharge            *float64
	Status               string
	PaymentMethod        string
	TrackInfo            string
	PayerInfo            string
	TransactionReference string
	TransactionDate      string
}

// OrderWithPaymentInput is the create-with-payment body.
type OrderWithPaymentInput struct {
	CustomerID               *int64
	PaymentGateway           string
	DeliveryType             string
	IsSendEmail              *bool
	TransactionNo            string
	ExternalReferenceOrderID string
	BillingAddress           *OrderAddressInput
	ShippingAddress          *OrderAddressInput
	Products                 []*OrderProductInput
	OrderChannel             string
	IPAddress                string
	Payment                  *InitialPaymentInput
}

func (o OrderWithPaymentInput) ToMap() map[string]any { return dto.ToMap(o) }

// OrderUpdateInput is the PUT body, wrapped in order; it writes the customer and the two addresses, and nothing else.
type OrderUpdateInput struct {
	CustomerID      *int64
	BillingAddress  *OrderAddressInput
	ShippingAddress *OrderAddressInput
}

func (OrderUpdateInput) EnvelopeKey() string { return OrderKey }

func (o OrderUpdateInput) ToMap() map[string]any { return dto.ToMap(o) }

// OrderPatchInput is the PATCH body, wrapped in order: the seven fields it writes; an address sent replaces that address.
type OrderPatchInput struct {
	OrderStatus      string
	CustomerPONumber string
	PaymentTerm      string
	BillingAddress   *OrderAddressInput
	ShippingAddress  *OrderAddressInput
	OrderChannel     string
	IPAddress        string
}

func (OrderPatchInput) EnvelopeKey() string { return OrderKey }

func (o OrderPatchInput) ToMap() map[string]any { return dto.ToMap(o) }

// OrderAddressesInput is the addresses sub-path's body; at least one address, and one left out is untouched.
type OrderAddressesInput struct {
	TransactionNo   string
	BillingAddress  *OrderAddressInput
	ShippingAddress *OrderAddressInput
}

func (o OrderAddressesInput) ToMap() map[string]any { return dto.ToMap(o) }

// OrderCancelInput is the cancel body; every field is optional and an empty request cancels too.
type OrderCancelInput struct {
	TransactionNo       string
	Reason              string
	RefundAutomatically *bool
	IsGatewayRefund     *bool
}

func (o OrderCancelInput) ToMap() map[string]any { return dto.ToMap(o) }

// OrderCommentInput is the comment body, with no envelope; Content is required.
type OrderCommentInput struct {
	Content       string
	SaveAndSend   *bool
	TransactionNo string
}

func (o OrderCommentInput) ToMap() map[string]any { return dto.ToMap(o) }

// OrderPaymentInput is the payment body: {"payment": {...}} with transaction_no beside the envelope.
type OrderPaymentInput struct {
	Amount               *float64
	PaymentMethod        string
	PaymentStatus        string
	TrackInfo            string
	PayerInfo            string
	TransactionDate      string
	TransactionReference string
	Surcharge            *float64
	TransactionNo        string
}

func (OrderPaymentInput) EnvelopeKey() string { return OrderPaymentKey }

// TopLevel names the field that sits beside that envelope rather than inside it.
func (OrderPaymentInput) TopLevel() []string { return []string{"transaction_no"} }

func (o OrderPaymentInput) ToMap() map[string]any { return dto.ToMap(o) }

// RefundLineInput is one line an order refund names, by the order line's item_id.
type RefundLineInput struct {
	OrderItemID *int64
	Quantity    *int64
	Amount      *float64
}

// OrderRefundInput is the order refund body, {"refund": {...}}; RefundDate is required, and Amount may not pass what was paid less earlier refunds.
type OrderRefundInput struct {
	Amount     string
	RefundDate string
	Reason     string
	LineItems  []*RefundLineInput
}

func (OrderRefundInput) EnvelopeKey() string { return "refund" }

func (o OrderRefundInput) ToMap() map[string]any { return dto.ToMap(o) }

// PaymentRefundInput is the payment refund body: {"refund": {...}} with transaction_no beside it.
type PaymentRefundInput struct {
	Amount          string
	Reason          string
	IsGatewayRefund *bool
	SendRefundEmail *bool
	AdjustWithOrder *bool
	TransactionNo   string
}

func (PaymentRefundInput) EnvelopeKey() string { return "refund" }

// TopLevel names the field that sits beside that envelope rather than inside it.
func (PaymentRefundInput) TopLevel() []string { return []string{"transaction_no"} }

func (p PaymentRefundInput) ToMap() map[string]any { return dto.ToMap(p) }

// ShipmentItemInput is one line to ship, by the order line's item_id.
type ShipmentItemInput struct {
	OrderItemID *int64
	Quantity    *int64
}

// OrderShipmentInput is the shipment body; a create needs Method, ShippingDate and ShipmentItems. An update changes only what it names, and new ShipmentItems need ChangeNote.
type OrderShipmentInput struct {
	TransactionNo string
	Method        string
	TrackingInfo  string
	ShippingDate  string
	ShipmentItems []*ShipmentItemInput
	ChangeNote    string
}

func (o OrderShipmentInput) ToMap() map[string]any { return dto.ToMap(o) }

// OrderStatusInput is the status body; Status is required: cancelled, completed, draft, in_progress, paused or reopened.
type OrderStatusInput struct {
	Status              string
	TransactionNo       string
	RefundAutomatically *bool
}

func (o OrderStatusInput) ToMap() map[string]any { return dto.ToMap(o) }

// CustomOrderStatusInput is the custom status body; OrderID is required, a missing one is 500.
type CustomOrderStatusInput struct {
	OrderID          *int64
	OrderStatusValue string
	OrderStatusLabel string
	OrderType        string
	StatusType       string
}

func (c CustomOrderStatusInput) ToMap() map[string]any { return dto.ToMap(c) }

// OrderShape is the typed view of a body shaped {"order": {…}} in sections.
func OrderShape(decoded any) any { return dto.Read[Order](decoded) }

// OrderPageShape is the typed view of the listing body.
func OrderPageShape(decoded any) any { return dto.Read[OrderPage](decoded) }

// OrderCountShape is the typed view of the count sub-path's body.
func OrderCountShape(decoded any) any { return dto.Read[OrderCount](decoded) }

// FlatOrderShape is the typed view of a body shaped {"order": {…}} in the flat shape.
func FlatOrderShape(decoded any) any { return dto.Read[FlatOrder](decoded) }

// CustomerOrdersShape is the typed view of one customer's orders.
func CustomerOrdersShape(decoded any) any { return dto.Read[CustomerOrders](decoded) }

// OrderAddressesShape is the typed view of the addresses sub-path's body.
func OrderAddressesShape(decoded any) any { return dto.Read[OrderAddresses](decoded) }

// OrderCommentShape is the typed view of a body shaped {"comment": {…}}.
func OrderCommentShape(decoded any) any { return dto.Read[OrderComment](decoded) }

// OrderCommentListShape is the typed view of the comments sub-path's body.
func OrderCommentListShape(decoded any) any { return dto.Read[OrderCommentList](decoded) }

// OrderCustomFieldsShape is the typed view of the custom_fields sub-path's body.
func OrderCustomFieldsShape(decoded any) any { return dto.Read[OrderCustomFields](decoded) }

// PaymentReceiptShape is the typed view of a body shaped {"payment": {…}}.
func PaymentReceiptShape(decoded any) any { return dto.Read[PaymentReceipt](decoded) }

// OrderRefundShape is the typed view of an order refund's answer.
func OrderRefundShape(decoded any) any { return dto.Read[OrderRefund](decoded) }

// OrderRefundPageShape is the typed view of an order's refund listing.
func OrderRefundPageShape(decoded any) any { return dto.Read[OrderRefundPage](decoded) }

// PaymentRefundShape is the typed view of a payment refund's answer.
func PaymentRefundShape(decoded any) any { return dto.Read[PaymentRefund](decoded) }

// InvoiceDispatchShape is the typed view of the send-invoice answer.
func InvoiceDispatchShape(decoded any) any { return dto.Read[InvoiceDispatch](decoded) }

// ShipmentStatusShape is the typed view of the shipment-status body.
func ShipmentStatusShape(decoded any) any { return dto.Read[ShipmentStatus](decoded) }

// OrderShipmentShape is the typed view of a body shaped {"shipment": {…}}.
func OrderShipmentShape(decoded any) any { return dto.Read[OrderShipment](decoded) }

// OrderShipmentListShape is the typed view of the shipments sub-path's body.
func OrderShipmentListShape(decoded any) any { return dto.Read[OrderShipmentList](decoded) }
