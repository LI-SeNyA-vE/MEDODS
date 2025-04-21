package authDB

const (
	querySearchRefresh = `SELECT refreshtoken FROM token WHERE guid = $1`
	addRefreshToken    = `INSERT INTO token (guid, refreshtoken) VALUES ($1, $2)`
	queryDeleteRefresh = `DELETE FROM token WHERE guid = $1`
)
