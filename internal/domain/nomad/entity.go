package nomad

type Nomad struct {
	id    int32
	phone string
}

func New(id int32, phone string) *Nomad {
	return &Nomad{
		id:    id,
		phone: phone,
	}
}

func (n *Nomad) ID() int32     { return n.id }
func (n *Nomad) Phone() string { return n.phone }
