package pool

type Pool struct{
	battery float64	
	capacity float64
	emptied int
}

func NewPool(battery, capacity float64 ) *Pool {
	return &Pool{
		battery: battery,
		capacity: capacity,
		emptied: 0,
	}
}

func (p *Pool) AcceptEnergy(Energy float64){
	p.battery += Energy
}

func (p *Pool) GiveEnergy(Energy float64){
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

func (p *Pool) GetCapacity() float64 {
	return p.capacity
}

func (p *Pool) GiveToOtherPool(amount float64, po *Pool){
	p.battery = p.battery - amount
	po.battery = po.battery + amount
	if po.battery > po.capacity{
		po.battery = po.capacity
	}
}

func (p *Pool) CalculateCostPerKWH (date string, solarIndexNextDay float64) float64 {
	const maxSolarIndex = 10.0
	const totalTokens = 100.0

	batteryRatio := p.GetBattery() / p.GetCapacity()
	batteryScore := (1 - batteryRatio) 

	solarRatio := solarIndexNextDay / maxSolarIndex
	solarScore := (1 - solarRatio) 

	weightedSum := (2.0/3.0)*batteryScore + (1.0/3.0)*solarScore

	priceInTokens := weightedSum * totalTokens

	return priceInTokens
}