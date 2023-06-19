package entity

type Asset struct {
	ID           string
	Name         string
	MarketVolume int
}

func NewAsset(ID string, name string, marketVolume int) *Asset {
	return &Asset{
		ID:           ID,
		Name:         name,
		MarketVolume: marketVolume,
	}
}
