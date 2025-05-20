package pool

type Pool struct{
	battery float64	
	emptied int
}

func NewPool(battery float64) *Pool {
	return &Pool{
		battery: battery,
		emptied: 0,
	}
}

func (p *Pool) ContributeEnergy(Energy float64){
	p.battery += Energy
}

func (p *Pool) WithdrawEnergy(Energy float64){
	p.battery -= Energy
}

func (p *Pool) EmptyPool() {
	p.emptied ++
}

func (p *Pool) GetEmptiedCount() int {
	return p.emptied
}

func (p *Pool) GetBattery() float64 {
	return p.battery
}


