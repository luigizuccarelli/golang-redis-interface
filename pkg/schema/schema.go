package schema

// SchemaInterface - acts as an interface wrapper for our profile schema
// All the go microservices will using this schema
type SchemaInterface struct {
}

// Response schema
type Response struct {
	Name       string `json:"name"`
	StatusCode string `json:"statuscode"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Payload    string `json:"payload"`
}
