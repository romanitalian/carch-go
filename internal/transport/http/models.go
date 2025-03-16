package http

// Request models
type createUserRQ struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type updateUserRQ struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Response models
type errorRS struct {
	Error string `json:"error"`
}
