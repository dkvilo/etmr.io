package core

// Response structure
type Response struct {
	Data    	interface{} 		`json:"data,omitempty"`
	Message 	interface{}      		`json:"message,omitempty"`
	Success 	bool      			`json:"success,omitempty"`
}

