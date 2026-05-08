package trading_station

type TradingStation struct {
	id        int32
	name      string
	phone     *string
	longitude float32
	latitude  float32
}

func New(id int32, name string, phone *string, longitude, latitude float32) *TradingStation {
	return &TradingStation{
		id:        id,
		name:      name,
		phone:     phone,
		longitude: longitude,
		latitude:  latitude,
	}
}

func (ts *TradingStation) ID() int32          { return ts.id }
func (ts *TradingStation) Name() string       { return ts.name }
func (ts *TradingStation) Phone() *string     { return ts.phone }
func (ts *TradingStation) Longitude() float32 { return ts.longitude }
func (ts *TradingStation) Latitude() float32  { return ts.latitude }
