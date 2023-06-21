package entity

type Order struct {
	ID            string
	Investor      *Investor
	Asset         *Asset
	Shares        int
	PendingShares int
	Price         float64
	OrderType     string
	Status        string
	Transactions  []*Transaction
}

func NewOrder(ID string, investor *Investor, asset *Asset, shares int, price float64, orderType string) *Order {
	return &Order{
		ID:           ID,
		Investor:     investor,
		Asset:        asset,
		Shares:       shares,
		Price:        price,
		OrderType:    orderType,
		Status:       "OPEN",
		Transactions: []*Transaction{},
	}
}
