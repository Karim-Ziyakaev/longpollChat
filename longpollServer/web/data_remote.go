package web

type MessageRequest struct {
	Text string `json:"text"`
	To   string `json:"to"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserId string `json:"user_id"`
	Token  string `json:"token"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type RegisterResponse struct {
	UserId string `json:"user_id"`
	Token  string `json:"token"`
}
