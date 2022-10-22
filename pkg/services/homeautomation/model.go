package homeautomation

type Device struct {
	Id        string `json:"id" gorm:"primary_key"`
	Name      string `json:"name"`
	IpAddress string `json:"ip_address"`
	Port      int    `json:"port"`
	Type      string `json:"type"`
	State     string `json:"state"`
	User      string `json:"user"`
	Password  string `json:"password"`
	Unit      string `json:"unit"`

	Value string `json:"value"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

type Scenario struct {
	Id        string `json:"id" gorm:"primary_key"`
	Name      string `json:"name"`
	Condition string `json:"condition"`
}
