package repository

import (
	"context"
	"fmt"
	"sasa-elterminali-service/internal/core/models"
	"strconv"
)

func (p *DB) GetMakinalar() (*[]models.Makina, error) {
	var makinalar []models.Makina
	query := `SELECT id, alt_makine as makine_name 
	FROM filament_prm.alt_makine
	WHERE isdelete = false;`
	rows, err := p.postgres.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var makina models.Makina
		err := rows.Scan(&makina.MakinaID, &makina.MakinaName)
		if err != nil {
			return nil, err
		}
		makinalar = append(makinalar, makina)
	}

	return &makinalar, nil
}

func (p *DB) GetBloklar(req *models.BloklarGetRequest) (*[]models.Blok, error) {
	var bloklar []models.Blok
	query := `
		SELECT s.id
			,s.blok_numarasi
		FROM (
			SELECT DISTINCT ON (blok_numarasi) id
				,substring(blok_numarasi FROM 2)::INT AS sira
				,blok_numarasi
			FROM filament_prm.prm_flmnt_dty_makine_blok_numarasi x
			WHERE x.blok_numarasi LIKE '%` + req.BlokTarafi + `%'
			ORDER BY blok_numarasi
				,id ASC
				,substring(blok_numarasi FROM 2)::INT ASC
			) s
		ORDER BY sira ASC
	`
	rows, err := p.postgres.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var blok models.Blok
		err := rows.Scan(&blok.BlokID, &blok.BlokName)
		if err != nil {
			return nil, err
		}
		bloklar = append(bloklar, blok)
	}

	return &bloklar, nil
}

func (p *DB) GetPozisyonlar(req *models.PozisyonlarGetRequest) (*[]models.BlokPozisyonResponse, error) {
	query := `
		SELECT blok_id, pozisyon_id, pozisyon_name
		FROM elterminali.blok_pozisyonlari
		WHERE blok_id = ANY($1)
	`
	rows, err := p.postgres.Query(context.Background(), query, req.BlokIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	blokPozisyonlarMap := make(map[int][]models.Pozisyon)
	for rows.Next() {
		var blokID int
		var pozisyon models.Pozisyon
		err := rows.Scan(&blokID, &pozisyon.PozisyonID, &pozisyon.PozisyonName)
		if err != nil {
			return nil, err
		}
		strPozisyonID := strconv.Itoa(pozisyon.PozisyonID)
		pozisyon.PozisyonName = &strPozisyonID

		blokPozisyonlarMap[blokID] = append(blokPozisyonlarMap[blokID], pozisyon)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var blokPozisyonlar []models.BlokPozisyonResponse
	for _, blokID := range req.BlokIDs {
		if pozisyonlar, ok := blokPozisyonlarMap[blokID]; ok {
			blokPozisyon := models.BlokPozisyonResponse{
				Blok:     fmt.Sprintf("%d. Blok", blokID),
				Pozisyon: pozisyonlar,
			}
			blokPozisyonlar = append(blokPozisyonlar, blokPozisyon)
		}
	}
	return &blokPozisyonlar, nil
}
