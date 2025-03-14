package types

type DeliveryStatus string

const (
	Opening        DeliveryStatus = "Opening"
	Pending        DeliveryStatus = "Pending"
	Purchased      DeliveryStatus = "Purchased"
	PickedUp       DeliveryStatus = "PickedUp"
	OutForDelivery DeliveryStatus = "OutForDelivery"
	Delivered      DeliveryStatus = "Delivered"
	Cancel         DeliveryStatus = "Cancel"
	Returned       DeliveryStatus = "Returned"
	Refunded       DeliveryStatus = "Refunded"
)

var AllowedStatusTransitions = map[DeliveryStatus]map[DeliveryStatus]bool{
	Opening:        {Pending: true, Cancel: true},
	Pending:        {Purchased: true, Cancel: true},
	Purchased:      {PickedUp: true},
	PickedUp:       {OutForDelivery: true, Returned: true},
	OutForDelivery: {Delivered: true, Returned: true},
	Delivered:      {Refunded: true, Returned: true},
	Refunded:       {Returned: true},
	Returned:       {},
	Cancel:         {},
}
