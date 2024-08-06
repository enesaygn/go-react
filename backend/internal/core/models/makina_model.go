package models

type Makina struct {
	MakinaID   int    `json:"makinaID"`
	MakinaName string `json:"makinaName"`
}

type BloklarGetRequest struct {
	BlokTarafi string `json:"blokTarafi"`
}

type Blok struct {
	BlokID   int    `json:"blokID"`
	BlokName string `json:"blokName"`
}

type PozisyonlarGetRequest struct {
	BlokIDs []int `json:"blokIDs"`
}

type Pozisyon struct {
	PozisyonID   int     `json:"pozisyonID"`
	PozisyonName *string `json:"name"`
}

type BlokPozisyonResponse struct {
	Blok     string     `json:"blok"`
	Pozisyon []Pozisyon `json:"pozisyon"`
}
