package test

type Details struct {
	Name           string   `json:"name"`
	Age            int      `json:"age"`
	Address        *Address `json:"address"`
	PassportNumber string   `json:"passport_number"`
	PanNumber      string   `json:"pan_number"`
}

type Address struct {
	City       string `json:"city"`
	PostalCode string `json:"postalcode"`
}

type User struct {
	Details  *Details `json:"details"`
	Verified string   `json:"verified"`
	Type     string   `json:"type"`
}

type Data struct {
	Users []User `json:"users"`
}
