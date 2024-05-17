package sql

const (
	InsertUser = "INSERT INTO auth_service.users (ID, name, email, password, role) VALUES (?, ?, ?, ?, ?)"
	GetUser    = "SELECT id, email, password FROM auth_service.users WHERE email = ? ALLOW FILTERING"
)
