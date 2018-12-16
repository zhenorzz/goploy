package model

import (
	"log"
)

type Ceshi struct {
	id   int32
	name string
}

type Ceshis []Ceshi

func (ceshi *Ceshi) Query() {
	db := NewDB()
	err := db.QueryRow("SELECT * FROM ceshi").Scan(&ceshi.id, &ceshi.name)
	if err != nil {
		log.Fatal(err)
	}
}

func (ceshis *Ceshis) QueryMany() {
	db := NewDB()
	rows, err := db.Query("SELECT * FROM ceshi")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var ceshi Ceshi

		if err := rows.Scan(&ceshi.id, &ceshi.name); err != nil {
			log.Fatal(err)
		}

		*ceshis = append(*ceshis, ceshi)
	}
}
