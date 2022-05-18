package model

// Users struct which contains
// an array of users
type Targets struct {
	Targets []Target `json:"targets"`
}

// User struct which contains a name
// a type and a list of social links
type Target struct {
	Name    string `json:"name"`
	Macaddr string `json:"macaddr"`
	Ip      string `json:"ip"`
}
