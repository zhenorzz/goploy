package database

import "embed"

//go:embed *
var File embed.FS

const FileExt = ".sql"
const GoploySQL = "goploy.sql"
