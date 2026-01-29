package domain

type ShippingItem struct {
	Quantity int32
}

type Shipping struct {
	OrderID int64
	Items []ShippingItem
}

func (s Shipping) TotalUnits() int32 {
	var total int32 = 0
	for _, item := range s.Items {
		total += item.Quantity
	}
	return total
}
