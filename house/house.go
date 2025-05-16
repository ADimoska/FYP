package house

type House struct{
	customer string
	location string
	city string
	gen_capacity float64
	cl float64
	gc float64
	gg float64
	battery float64
	energyData   map[string]map[string]map[string]float64
	blackouts int
}

func NewHouse(customer, location string, capacity, cl, gc, gg, battery float64) *House {
	return &House{
		customer: customer,
		location: location,
		city: "not defined",
		gen_capacity: capacity,
		cl: cl,
		gc: gc,
		gg: gg,
		battery: battery,
		energyData:   make(map[string]map[string]map[string]float64),
		blackouts: 0,		
		
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

func (h *House) GetCity() string{
	return h.city
}

func (h *House) SetCity(city string) {
	h.city = city
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

func (h *House) ResetBattery() { 
	h.battery = 0
}

func (h *House) AddBattery(battery float64) {
	h.battery += battery
}

func (h *House) AddBlackout() {
	h.blackouts ++
}
func (h *House) GetBlackouts() int {
	return h.blackouts
}
func (h *House) StoreEnergyData(date, time, consumptionType string, value float64) {
	if _, ok := h.energyData[date]; !ok {
		h.energyData[date] = make(map[string]map[string]float64)
	}
	if _, ok := h.energyData[date][time]; !ok {
		h.energyData[date][time] = make(map[string]float64)
	}
	h.energyData[date][time][consumptionType] = value
}

func (h *House) GetEnergyData() map[string]map[string]map[string]float64 {
	return h.energyData
}

func (h *House) GetCurrentEnergy(date, time string) {
	h.cl = h.energyData[date][time]["CL"]
	h.gc = h.energyData[date][time]["GC"]
	h.gg = h.energyData[date][time]["GG"]
}

func (h *House) SetBatteryCapacity() {

}








