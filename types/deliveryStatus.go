package types

type DeliveryStatus string

const (
	Pending        DeliveryStatus = "Pending"
	Purchased      DeliveryStatus = "Purchased"
	PickedUp       DeliveryStatus = "PickedUp"
	OutForDelivery DeliveryStatus = "OutForDelivery"
	Delivered      DeliveryStatus = "Delivered"
	Cancel         DeliveryStatus = "Cancel"
	Returned       DeliveryStatus = "Returned"
	Refunded       DeliveryStatus = "Refunded"
	Opening        DeliveryStatus = "Opening"
)
