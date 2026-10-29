package pgsql

import "embed"

//go:embed *.sql
var EmbedMigrationsPgsql embed.FS
