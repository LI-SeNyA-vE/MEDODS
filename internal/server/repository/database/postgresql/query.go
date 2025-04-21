package postgresql

const (
	queryTokenTableIsNot = `
	CREATE TABLE IF NOT EXISTS "token" (
		id BIGSERIAL PRIMARY KEY,    
		guid uuid UNIQUE NOT NULL,  
		refreshToken BYTEA NOT NULL
	);`
)
