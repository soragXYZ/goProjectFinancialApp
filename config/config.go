package config

import (
	"database/sql"
)

// Glopbal env var spreads across every package
// Not ideal, should be modified later
// https://www.alexedwards.net/blog/organising-database-access
var DB *sql.DB
