package product

type Product struct {
	id      int32
	name    string
	details *string
	weight  int32
	volume  int32
}

func New(id int32, name string, details *string, weight, volume int32) *Product {
	return &Product{
		id:      id,
		name:    name,
		details: details,
		weight:  weight,
		volume:  volume,
	}
}

func (p *Product) ID() int32        { return p.id }
func (p *Product) Name() string     { return p.name }
func (p *Product) Details() *string { return p.details }
func (p *Product) Weight() int32    { return p.weight }
func (p *Product) Volume() int32    { return p.volume }
