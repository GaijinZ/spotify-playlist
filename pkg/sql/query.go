package sql

const (
	InsertUser = "INSERT INTO auth_service.users (ID, name, email, password, role) VALUES (?, ?, ?, ?, ?)"
)
