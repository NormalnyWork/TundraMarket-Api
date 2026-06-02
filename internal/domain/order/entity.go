package order

import "time"

type Status string

const (
	StatusCreated    Status = "CREATED"
	StatusProcessing Status = "PROCESSING"
	StatusSent       Status = "SENT"
	StatusCompleted  Status = "COMPLETED"
	StatusCancelled  Status = "CANCELLED"
	StatusDenied     Status = "DENIED"
)

type OrderCategory string

const (
	OrderCategoryNew        OrderCategory = "NEW"
	OrderCategoryProcessing OrderCategory = "PROCESSING"
	OrderCategoryHistory    OrderCategory = "HISTORY"
)

type ProductCount struct {
	ProductID int32
	Quantity  int32
}

type StatusHistory struct {
	status    Status
	createdAt time.Time
}

type Order struct {
	id               int32
	nomadID          int32
	tradingStationID int32
	status           Status
	history          []StatusHistory
	comment          string
	longitude        float32
	latitude         float32
	products         []ProductCount
	createdAt        time.Time
}

func New(nomadID, tradingStationID int32, comment string, longitude, latitude float32, products []ProductCount) (*Order, error) {
	if len(products) == 0 {
		return nil, ErrEmptyCart
	}
	return &Order{
		nomadID:          nomadID,
		tradingStationID: tradingStationID,
		status:           StatusCreated,
		comment:          comment,
		longitude:        longitude,
		latitude:         latitude,
		products:         products,
	}, nil
}

func NewStatusHistory(status Status, createdAt time.Time) StatusHistory {
	return StatusHistory{
		status:    status,
		createdAt: createdAt,
	}
}

func Restore(id, nomadID, tradingStationID int32, status Status, comment string, longitude, latitude float32, products []ProductCount, history []StatusHistory, createdAt time.Time) *Order {
	return &Order{
		id:               id,
		nomadID:          nomadID,
		tradingStationID: tradingStationID,
		status:           status,
		history:          history,
		comment:          comment,
		longitude:        longitude,
		latitude:         latitude,
		products:         products,
		createdAt:        createdAt,
	}
}

func (o *Order) ID() int32                { return o.id }
func (o *Order) NomadID() int32           { return o.nomadID }
func (o *Order) TradingStationID() int32  { return o.tradingStationID }
func (o *Order) Status() Status           { return o.status }
func (o *Order) History() []StatusHistory { return o.history }
func (o *Order) Comment() string          { return o.comment }
func (o *Order) Longitude() float32       { return o.longitude }
func (o *Order) Latitude() float32        { return o.latitude }
func (o *Order) Products() []ProductCount { return o.products }
func (o *Order) CreatedAt() time.Time     { return o.createdAt }

func (h StatusHistory) Status() Status       { return h.status }
func (h StatusHistory) CreatedAt() time.Time { return h.createdAt }
