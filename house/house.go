package house

type House struct{
	customer string
	location string
	gen_capacity float64
	cl float64
	gc float64
	gg float64
	battery float64
}

func NewHouse(customer, location string, capacity, cl, gc, gg, battery float64) *House {
	return &House{
		customer: customer,
		location: location,
		gen_capacity: capacity,
		cl: cl,
		gc: gc,
		gg: gg,
		battery: battery,
		
	}
}


func (h *House) GetCustomer() string{
	return h.customer
}

func (h *House) SetCustomer(customer string) {
	h.customer = customer
}

func (h *House) GetLocation() string{
	return h.location
}

func (h *House) SetLocation(location string) {
	h.location = location
}

func (h *House) GetGen_capacity() float64{
	return h.gen_capacity
}

func (h *House) SetGen_capacity(gen_capacity float64) {
	h.gen_capacity = gen_capacity
}

func (h *House) GetCL() float64 { 
	return h.cl
}

func (h *House) SetCL(cl float64) {
	h.cl = cl
}

func (h *House) GetGC() float64 { 
	return h.gc
}

func (h *House) SetGC(gc float64) {
	h.gc = gc
}

func (h *House) GetGG() float64 { 
	return h.gg
}

func (h *House) SetGG(gg float64) {
	h.gc = gg
}

func (h *House) GetBattery() float64 { 
	return h.battery
}

func (h *House) AddBattery(battery float64) {
	h.battery =+ battery
}











