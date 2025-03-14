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
	Purchased:      {PickedUp: true, Cancel: false},
	PickedUp:       {OutForDelivery: true, Returned: true, Cancel: false},
	OutForDelivery: {Delivered: true, Returned: true, Cancel: false},
	Delivered:      {Refunded: true, Returned: true, Cancel: false},
	Refunded:       {Returned: true, Cancel: false},
	Returned:       {},
	Cancel:         {},
}
