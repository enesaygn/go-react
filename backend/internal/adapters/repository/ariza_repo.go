package repository

import (
	"context"
	"sasa-elterminali-service/internal/core/models"
	"strconv"
	"strings"
	"time"
)

const (
	ARIZACREATELOGMESSAGE         = "Arıza oluşturuldu"
	BAKIMUZMANIATAMALOGMESSAGE    = "Bakım uzmanı atandı"
	ISTRANSFERILOGMESSAGE         = "İş transferi yapıldı"
	ARIZAYAPILDILOGMESSAGE        = "Arıza yapıldı"
	PROSESONAYILOGMESSAGE         = "Proses onayı verildi"
	TESTTALEBIOLUSTURMALOGMESSAGE = "Test talebi oluşturuldu"
)

// Ariza Roller
const (
	ISLETME_UZMANI              = 116
	BAKIM_OPERATORU             = 117
	BAKIM_UZMANI                = 118
	PROSES_CALISANI             = 119
	KALITE_KONTROL_VE_PAKETLEME = 120
)

// Aşamalar
const (
	BAKIM_ASAMASI          = "Bakım Aşaması"
	PROSES_KONTROL_ASAMASI = "Proses Kontrol Aşaması"
	TEST_SONUCLARI         = "Test Sonuçları"
)

// Ariza Statuleri ...
const (
	ARIZA_YENI               = "Yeni Arızalar"          //"Arıza Yeni"
	ARIZA_ISLEMDE            = "İşlemdeki Arızalar"     //"Arıza İşlemde"
	ARIZA_YAPILDI            = "Yapılan Arızalar"       //"Arıza Yapıldı"
	ARIZA_PROSES_ONAYI       = "Proses Onaylı Arızalar" //"Proses Onayı"
	ARIZA_YAPILMADI          = "Yapılmayan Arızalar"    //"Yapılmadı"
	ARIZA_SONLANDI           = "Sonlanan Arızalar"      //"Arıza Sonlandı"
	ARIZA_TEST_TALEBI        = "Yeni Testler"           //"Test Talebi"
	ARIZA_TEST_TALEBI_ALINDI = "İşlemdeki Testler"      //"Test Talebi Alındı"
	ARIZA_OLUMSUZ_TEST       = "Olumsuz Testler"        //"Test Onaylanmadı"

	//Alt statüler
	ARIZA_TEKRAR       = "Tekrarlayan Arızalar" //"Tekrara Düştü"
	ARIZA_TEST_ONAYI   = "Onaylanan Testler"    //"Test Onayı"
	ARIZA_ISTRANSFERI  = "İş Transferi"         //"İş Transferi"
	ARIZA_TEST_ASAMASI = "Testte Olan Arızalar" //"Test Aşaması"
	//SON AŞAMA
	ARIZA_SAP_GONDERILDI = "SAP GÖNDERİLDİ" //"Sap Gönderildi"
	ARIZA_TAMAMLANDI     = "TAMAMLANDI"     //"Tamamlandı"
	//Test Kaynağı
	TEST_KAYNAGI_ARIZADAN_GELEN = "Arızadan Gelen"
)

func (p *DB) IsletmeUzmaniCreateAriza(ariza *models.ArizaCreateRequestBody) (*models.EmployeesAndCreatedUser, error) {

	tx, err := p.postgres.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	var arizaID int
	var makineAdi string
	var arizaTipi string

	createDate := time.Now().Format("2006-01-02 15:04:05")

	query := `
		INSERT INTO elterminali.ariza (
			makine_id
			,blok_tarafi
			,ariza_tipi
			,ariza_tipi_tanim_id
			,aciklama
			,ariza_statusu
			,isletme_uzmani_employee_id
			,create_user_id
			,create_date
			)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
			)
		RETURNING ariza_id, (SELECT alt_makine FROM filament_prm.alt_makine WHERE id = $1) AS makine_adi, 
		(SELECT ariza_tipi_tanim_adi FROM elterminali.ariza_tipi_tanım WHERE ariza_tipi_tanim_id = $4) AS ariza_tipi_tanim_adi;
`
	err = p.postgres.QueryRow(context.Background(), query,
		ariza.MakineID, ariza.BlokTarafi, ariza.ArizaTipi, ariza.ArizaTipiTanimID, ariza.Aciklama, ARIZA_YENI,
		ariza.UserID, ariza.UserID, createDate).Scan(&arizaID, &makineAdi, &arizaTipi)

	if err != nil {
		return nil, err
	}

	var arizaDetailsID int

	query3 := `
		INSERT INTO elterminali.ariza_details (
			ariza_id
			,makine_id
			,blok_tarafi
			,ariza_tipi
			,ariza_tipi_tanim_id
			,aciklama
			,ariza_statusu
			,asama
			,ariza_durumu
			,isletme_uzmani_employee_id
			,create_user_id
			,create_date
			,islem_turu
			)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
			)
		RETURNING ariza_details_id;
	`
	err = p.postgres.QueryRow(context.Background(), query3, arizaID, ariza.MakineID, ariza.BlokTarafi, ariza.ArizaTipi,
		ariza.ArizaTipiTanimID, ariza.Aciklama, ARIZA_YENI, "Root", "Yeni Arıza", ariza.UserID, ariza.UserID, createDate,
		"Düzenleme").Scan(&arizaDetailsID)

	if err != nil {
		return nil, err
	}

	query2 := `
		INSERT INTO elterminali.ariza_blok (
			blok_id
			,pozisyon
			,ariza_id
			,ariza_details_id
			,create_date
			,create_user_id
			)
		VALUES (
			
			(SELECT blok_id FROM elterminali.blok_pozisyonlari WHERE pozisyon_id = $1), $1, $2, $3, $4, $5
			
			);

			`

	for _, pozisyon := range ariza.Pozisyonlar {
		_, err = p.postgres.Exec(context.Background(), query2, pozisyon, arizaID, arizaDetailsID, createDate, ariza.UserID)
		if err != nil {
			return nil, err
		}

	}

	//BILDIRIM KISMI

	employees, err := p.GetEmployeeIDsByPermission([]int{BAKIM_UZMANI, BAKIM_OPERATORU}, ariza.UserID, strconv.Itoa(arizaID),
		ARIZA_YENI, makineAdi, arizaTipi, ariza.Pozisyonlar)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return employees, err
}

func (p *DB) ProsesCalisaniCreateAriza(ariza *models.ArizaCreateRequestBody) (*models.EmployeesAndCreatedUser, error) {

	tx, err := p.postgres.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	var arizaID int
	createDate := time.Now().Format("2006-01-02 15:04:05")

	query := `
		INSERT INTO elterminali.ariza (
			makine_id
			,blok_tarafi
			,ariza_tipi
			,ariza_tipi_tanim_id
			,aciklama
			,ariza_statusu
			,proses_calisani_employee_id
			,create_date
			)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
			)
		RETURNING ariza_id;
`
	err = p.postgres.QueryRow(context.Background(), query,
		ariza.MakineID, ariza.BlokTarafi, ariza.ArizaTipi, ariza.ArizaTipiTanimID, ariza.Aciklama, ARIZA_YENI,
		ariza.UserID,
		createDate).Scan(&arizaID)

	if err != nil {
		return nil, err
	}

	query2 := `
		INSERT INTO elterminali.ariza_blok (
			blok_id
			,ariza_id
			,pozisyon
			,create_date
			,create_user_id
			)
		VALUES (
			(
				SELECT blok_id
				FROM elterminali.blok_pozisyonlari
				WHERE pozisyon_id = $2)
				, $1, $2, $3, $4 
			);
			`

	for _, pozisyon := range ariza.Pozisyonlar {
		_, err = p.postgres.Exec(context.Background(), query2, arizaID, pozisyon, createDate,
			ariza.UserID)
		if err != nil {
			return nil, err
		}

	}

	query3 := `
		INSERT INTO elterminali.ariza_details (
			ariza_id
			,makine_id
			,blok_tarafi
			,ariza_tipi
			,ariza_tipi_tanim_id
			,aciklama
			,ariza_statusu
			,proses_calisani_employee_id
			,create_user_id
			,create_date
			,islem_turu
			)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
			);
	`
	_, err = p.postgres.Exec(context.Background(), query3, arizaID, ariza.MakineID, ariza.BlokTarafi, ariza.ArizaTipi,
		ariza.ArizaTipiTanimID, ariza.Aciklama, ARIZA_YENI, ariza.UserID, ariza.UserID, createDate,
		"Düzenleme")

	if err != nil {
		return nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return nil, err

}

func (p *DB) KaliteKontrolVePaketlemeCreateAriza(req *models.ArizaCreateRequestBody) error {
	createDate := time.Now().Format("2006-01-02 15:04:05")
	query := `
		INSERT INTO elterminali.ariza_olustur (
			makine_id
			,ariza_tipi
			,ariza_tipi_tanim_id
			,aciklama
			,ariza_statusu
			,create_user_id
			,create_date
			,is_deleted
			,kalite_kontrol_ve_paketleme_employee_id
			)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		);
`
	_, err := p.postgres.Exec(context.Background(), query, req.MakineID, req.ArizaTipi, req.ArizaTipiTanimID, req.Aciklama,
		ARIZA_YENI, req.UserID, createDate, false, req.UserID)

	if err != nil {
		return err
	}

	//TODO: BAKIM UZMANI VE OPERATORUNE BILDIRIM (bildirim sutunları kalkabilir!!)
	return nil
}

func (p *DB) BakimUzmaniCreateAriza(ariza *models.ArizaCreateRequestBody) (*models.EmployeesAndCreatedUser, error) {

	createDate := time.Now().Format("2006-01-02 15:04:05")
	query := `
		INSERT INTO elterminali.ariza (
			makine_id
			,ariza_tipi
			,ariza_tipi_tanim_id
			,aciklama
			,ariza_statusu
			,bakim_uzmani_employee_id
			,create_date
			)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7
			)
	`

	_, err := p.postgres.Exec(context.Background(), query, ariza.MakineID, ariza.ArizaTipi, ariza.ArizaTipiTanimID,
		ariza.Aciklama, ARIZA_YENI, ariza.UserID, createDate)

	return nil, err
}

func (p *DB) ArizaDuzenleme(ariza *models.ArizaDuzenlemeRequest) error {

	tx, err := p.postgres.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	query1 := `
		UPDATE elterminali.ariza
		SET
			makine_id = $1
			,blok_tarafi = $2
			,ariza_tipi = $3
			,ariza_tipi_tanim_id = $4
			,aciklama = $5
		WHERE ariza_id = $6;
		`
	_, err = tx.Exec(context.Background(), query1, ariza.MakineID, ariza.BlokTarafi, ariza.ArizaTipi,
		ariza.ArizaTipiTanimID, ariza.Aciklama, ariza.ArizaID)
	if err != nil {
		return err
	}

	query2 := `
		INSERT INTO elterminali.ariza_details (
			ariza_id
			,makine_id
			,blok_tarafi
			,ariza_tipi
			,ariza_tipi_tanim_id
			,aciklama
			,ariza_statusu
			,asama
			,ariza_durumu
			,create_user_id
			,create_date
			,islem_turu
			)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
			)
		RETURNING ariza_details_id;
		`
	var arizaDetailsID int
	createDate := time.Now().Format("2006-01-02 15:04:05")

	err = tx.QueryRow(context.Background(), query2, ariza.ArizaID, ariza.MakineID, ariza.BlokTarafi, ariza.ArizaTipi,
		ariza.ArizaTipiTanimID, ariza.Aciklama, ARIZA_YENI, "TODO", "TODO", ariza.UserID, createDate,
		"Düzenleme").Scan(&arizaDetailsID)
	if err != nil {
		return err
	}

	//yeni blokları ekle
	query3 := `
		INSERT INTO elterminali.ariza_blok (
			blok_id
			,pozisyon
			,ariza_id
			,ariza_details_id
			,create_date
			,create_user_id
			)
		VALUES (
			
			(SELECT blok_id FROM elterminali.blok_pozisyonlari WHERE pozisyon_id = $1), $1, $2, $3, $4, $5
			
			);
		`

	for _, pozisyon := range ariza.Pozisyonlar {
		_, err = tx.Exec(context.Background(), query3, pozisyon, ariza.ArizaID, arizaDetailsID, createDate, ariza.UserID)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return err

}

func (p *DB) BakimUzmaniArizaDuzenleme(ariza *models.ArizaDuzenlemeRequest) error {

	query := `
		UPDATE elterminali.ariza
		SET
			makine_id = $1
			,ariza_tipi = $2
			,ariza_tipi_tanim_id = $3
			,aciklama = $4
			,bakim_uzmani_employee_id = $5
		WHERE ariza_id = $6;
		`

	_, err := p.postgres.Exec(context.Background(), query, ariza.MakineID, ariza.ArizaTipi, ariza.ArizaTipiTanimID,
		ariza.Aciklama, ariza.UserID, ariza.ArizaID)
	if err != nil {
		return err
	}

	createDate := time.Now().Format("2006-01-02 15:04:05")
	qInsertDetails := `
		INSERT INTO elterminali.ariza_details (
			ariza_id
			,bakim_uzmani_employee_id
			,create_user_id
			,ariza_statusu
			,asama
			,ariza_durumu
			,create_date
			,islem_turu
			)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
			)
		RETURNING ariza_details_id;
		`

	var arizaDetailsID int

	err = p.postgres.QueryRow(context.Background(), qInsertDetails, ariza.ArizaID, ariza.UserID,
		ariza.UserID, ARIZA_TAMAMLANDI, "Son", "Arıza Sonlandı", createDate, "Düzenleme").Scan(&arizaDetailsID)
	if err != nil {
		return err
	}

	qInsertBlok := `
		INSERT INTO elterminali.ariza_blok (
			blok_id
			,pozisyon
			,ariza_id
			,ariza_details_id
			,create_date
			,create_user_id
			)
		VALUES (
			
			(SELECT blok_id FROM elterminali.blok_pozisyonlari WHERE pozisyon_id = $1), $1, $2, $3, $4, $5
			
			);
		`

	for _, pozisyon := range ariza.Pozisyonlar {
		_, err = p.postgres.Exec(context.Background(), qInsertBlok, pozisyon, ariza.ArizaID, arizaDetailsID,
			createDate, ariza.UserID)
		if err != nil {
			return err
		}
	}

	return err

}

func (p *DB) KaliteKontrolVePaketlemeUzmaniArizaDuzenleme(req *models.ArizaDuzenlemeRequest) error {

	tx, err := p.postgres.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	query1 := `
		UPDATE elterminali.ariza_olustur
		SET
			makine_id = $1
			,ariza_tipi = $2
			,ariza_tipi_tanim_id = $3
			,aciklama = $4
			,ariza_neden_id = $5
			,yapilan_mudahale = $6
		WHERE ariza_id = $7;
		`
	_, err = tx.Exec(context.Background(), query1, req.MakineID, req.ArizaTipi, req.ArizaTipiTanimID, req.Aciklama,
		req.ArizaNedenID, req.YapilanMudahale, req.ArizaID)
	if err != nil {
		return err
	}

	query2 := `
		INSERT INTO elterminali.ariza_kullanilan_malzemeler (
			ariza_id
			,malzeme_id
			,miktar
			)
		VALUES (
			$1, $2, $3
			);
		`
	for _, malzeme := range req.KullanilanMalzemeler {
		_, err = p.postgres.Exec(context.Background(), query2, req.ArizaID, malzeme.MalzemeID,
			malzeme.Miktar)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return err

}

func (p *DB) DeleteAriza(req *models.RequestBodyAriza) error {
	query := `
      UPDATE elterminali.ariza
      SET is_deleted = true
      WHERE ariza_id = $1;
    `
	_, err := p.postgres.Exec(context.Background(), query, req.ArizaID)
	return err
}

// func (p *DB) GetAriza(req *models.RequestBodyAriza) (*models.Ariza, error) {
// 	var ariza models.Ariza
// 	var createDate time.Time
// 	var updateDate *time.Time
// 	query := `
// 		WITH blok_bilgileri
// 		AS (
// 			SELECT ariza_id
// 				,json_build_object('BlokID', blok_id, 'Pozisyonlar', pozisyonlar) AS blok_bilgi
// 			FROM (
// 				SELECT ariza_id
// 					,blok_id
// 					,json_agg(pozisyon) AS pozisyonlar
// 				FROM elterminali.ariza_blok
// 				WHERE is_deleted = false
// 				GROUP BY ariza_id
// 					,blok_id
// 				) AS subquery
// 			)
// 		SELECT ao.ariza_id
// 			,ao.ariza_statusu
// 			,ao.makine_id
// 			,m.alt_makine as makine_adi
// 			,ao.blok_tarafi
// 			,bb.blok_bilgi
// 			,ao.ariza_tipi
// 			,ao.ariza_tipi_tanim_id
// 			,att.ariza_tipi_tanim_adi
// 			,ao.ariza_neden_id
// 			,an.ariza_neden_adi
// 			,ao.kisi_sayisi
// 			,ao.yapilan_mudahele
// 			,ao.test_talebi_kaynagi
// 			,ao.isletme_uzmani_employee_id
// 			,izu.full_name AS isletme_uzmani_full_name
// 			,ao.bakim_operatoru_employee_id
// 			,boe.full_name AS bakim_operatoru_full_name
// 			,ao.bakim_uzmani_employee_id
// 			,bue.full_name AS bakim_uzmani_full_name
// 			,ao.proses_calisani_employee_id
// 			,pce.full_name AS proses_calisani_full_name
// 			,ao.kalite_kontrol_ve_paketleme_employee_id
// 			,kkp.full_name AS kalite_kontrol_ve_paketleme_full_name
// 			,ao.aciklama
// 			,ao.test_aciklamasi
// 			,ao.ariza_log_id
// 			,ao.create_user_id
// 			,cu.full_name AS create_user_full_name
// 			,ao.update_user_id
// 			,uu.full_name AS update_user_full_name
// 			,ao.create_date
// 			,ao.update_date
// 		FROM elterminali.ariza ao
// 		LEFT JOIN filament_prm.alt_makine m ON m.id = ao.makine_id
// 		LEFT JOIN elterminali.ariza_neden an ON an.ariza_neden_id = ao.ariza_neden_id
// 		LEFT JOIN elterminali.ariza_tipi_tanım att ON att.ariza_tipi_tanim_id = ao.ariza_tipi_tanim_id
// 		LEFT JOIN techup.employee izu ON izu.employee_id = ao.isletme_uzmani_employee_id
// 		LEFT JOIN techup.employee boe ON boe.employee_id = ao.bakim_operatoru_employee_id
// 		LEFT JOIN techup.employee bue ON bue.employee_id = ao.bakim_uzmani_employee_id
// 		LEFT JOIN techup.employee pce ON pce.employee_id = ao.proses_calisani_employee_id
// 		LEFT JOIN techup.employee kkp ON kkp.employee_id = ao.kalite_kontrol_ve_paketleme_employee_id
// 		LEFT JOIN techup.employee cu ON cu.employee_id = ao.create_user_id
// 		LEFT JOIN techup.employee uu ON uu.employee_id = ao.update_user_id
// 		LEFT JOIN blok_bilgileri bb ON bb.ariza_id = ao.ariza_id
// 		WHERE ao.ariza_id = $1 AND ao.is_deleted = false

// `
// 	err := p.postgres.QueryRow(context.Background(), query, req.ArizaID).
// 		Scan(
// 			&ariza.ArizaID, &ariza.ArizaStatusu, &ariza.MakineID, &ariza.MakineAdi, &ariza.BlokTarafi, &ariza.ArizaBloklar, &ariza.ArizaTipi,
// 			&ariza.ArizaTipiTanimID, &ariza.ArizaTipiTanimAdi, &ariza.ArizaNedenID, &ariza.ArizaNedenAdi, &ariza.KisiSayisi,
// 			&ariza.YapilanMudahele, &ariza.TestTalebiKaynagi, &ariza.IsletmeUzmaniEmployeeID,
// 			&ariza.IsletmeUzmaniFullName, &ariza.BakimOperatoruEmployeeID, &ariza.BakimOperatoruFullName, &ariza.BakimUzmaniEmployeeID,
// 			&ariza.BakimUzmaniFullName, &ariza.ProsesCalisaniEmployeeID, &ariza.ProsesCalisaniFullName,
// 			&ariza.KaliteKontrolVePaketlemeEmployeeID, &ariza.KaliteKontrolVePaketlemeFullName, &ariza.Aciklama, &ariza.TestAciklamasi,
// 			&ariza.ArizaLogID, &ariza.CreateUserID, &ariza.CreateUserFullName, &ariza.UpdateUserID, &ariza.UpdateUserFullName,
// 			&createDate, &updateDate,
// 		)
// 	if err != nil {
// 		return nil, err

// 	}

// 	ariza.CreateDate = createDate.Format("2006-01-02 15:04:05")

// 	if updateDate != nil {
// 		*ariza.UpdateDate = updateDate.Format("2006-01-02 15:04:05")
// 	}

//		return &ariza, nil
//	}
func (p *DB) GetArizalar(req *models.GetArizalarRequestBody) (map[string][]models.Ariza, error) {

	var ariza models.Ariza
	var createDate *time.Time
	var updateDate *time.Time

	query := `
	WITH pozisyonlar
		AS (
			SELECT row_number() OVER (
					PARTITION BY ariza_id ORDER BY ariza_details_id DESC
					) AS islem_sira_no
				,ariza_id
				,array_agg(pozisyon) AS pozisyonlar
			FROM elterminali.ariza_blok
			GROUP BY ariza_id
				,ariza_details_id
			)
		SELECT 
			e.product_name 
			,ao.ariza_id
			,ao.ariza_statusu
			,ao.makine_id
			,m.alt_makine AS makine_adi
			,ao.blok_tarafi
			,p.pozisyonlar
			,ao.ariza_tipi
			,ao.ariza_tipi_tanim_id
			,att.ariza_tipi_tanim_adi
			,ao.ariza_neden_id
			,an.ariza_neden_adi
			,ao.kisi_sayisi
			,ao.yapilan_mudahele
			,ao.test_talebi_kaynagi
			,ao.isletme_uzmani_employee_id
			,izu.full_name AS isletme_uzmani_full_name
			,ao.bakim_operatoru_employee_id
			,boe.full_name AS bakim_operatoru_full_name
			,ao.bakim_uzmani_employee_id
			,bue.full_name AS bakim_uzmani_full_name
			,ao.proses_calisani_employee_id
			,pce.full_name AS proses_calisani_full_name
			,ao.kalite_kontrol_ve_paketleme_employee_id
			,kkp.full_name AS kalite_kontrol_ve_paketleme_full_name
			,ao.aciklama
			,ao.test_aciklamasi
			,ao.ariza_log_id
			,ao.create_user_id
			,cu.full_name AS create_user_full_name
			,ao.update_user_id
			,uu.full_name AS update_user_full_name
			,ao.create_date
			,ao.update_date
			,0 AS deneme
		FROM elterminali.ariza ao
		LEFT JOIN filament_prm.alt_makine m ON m.id = ao.makine_id
		LEFT JOIN elterminali.ariza_neden an ON an.ariza_neden_id = ao.ariza_neden_id
		LEFT JOIN elterminali.ariza_tipi_tanım att ON att.ariza_tipi_tanim_id = ao.ariza_tipi_tanim_id
		LEFT JOIN techup.employee izu ON izu.employee_id = ao.isletme_uzmani_employee_id
		LEFT JOIN techup.employee boe ON boe.employee_id = ao.bakim_operatoru_employee_id
		LEFT JOIN techup.employee bue ON bue.employee_id = ao.bakim_uzmani_employee_id
		LEFT JOIN techup.employee pce ON pce.employee_id = ao.proses_calisani_employee_id
		LEFT JOIN techup.employee kkp ON kkp.employee_id = ao.kalite_kontrol_ve_paketleme_employee_id
		LEFT JOIN techup.employee cu ON cu.employee_id = ao.create_user_id
		LEFT JOIN techup.employee uu ON uu.employee_id = ao.update_user_id
		LEFT JOIN elterminali.employee_birim eb ON eb.employee_id = ao.create_user_id 
		LEFT JOIN techup.equipment_ng e ON e.equipment_id = eb.birim_id  
		LEFT JOIN pozisyonlar p ON p.ariza_id = ao.ariza_id
			AND p.islem_sira_no = 1
		WHERE CASE 
				WHEN $1 = 'Sonlanan Arızalar'
					THEN ao.ariza_statusu IN (
							'Onaylanan Testler'
							,'Sonlanan Arızalar'
							)
				WHEN $1 = 'Yeni Arızalar'
					THEN ao.ariza_statusu IN (
							'Olumsuz Testler'
							,'Yeni Arızalar'
							,'Tekrarlayan Arızalar'
							)
				WHEN $1 = 'İşlemdeki Arızalar'
					THEN ao.ariza_statusu = $1
				WHEN $1 = 'Yapılan Arızalar'
					THEN ao.ariza_statusu = $1
				WHEN $1 = 'Proses Onaylı Arızalar'
					THEN ao.ariza_statusu = $1
				WHEN $1 = 'Yeni Testler'
					THEN ao.ariza_statusu = $1
				WHEN $1 = 'İşlemdeki Testler'
					THEN ao.ariza_statusu = $1
				WHEN $1 = 'Onaylanan Testler'
					THEN ao.ariza_statusu = $1
				WHEN $1 = 'Olumsuz Testler'
					THEN ao.ariza_statusu = $1
				WHEN $1 = 'Testte Olan Arızalar'
					THEN ao.ariza_statusu IN (
							'İşlemdeki Testler'
							,'Olumsuz Testler'
							,'Yeni Testler'
							)
				ELSE ariza_statusu = ''
				END
			AND ao.is_deleted = false

		`

	rows, err := p.postgres.Query(context.Background(), query, req.ArizaStatusu)
	if err != nil {
		return nil, err
	}

	resultMap := make(map[string][]models.Ariza)

	for rows.Next() {
		var ekip string
		err = rows.Scan(
			&ekip, &ariza.ArizaID, &ariza.ArizaStatusu, &ariza.MakineID, &ariza.MakineAdi, &ariza.BlokTarafi,
			&ariza.Pozisyonlar, &ariza.ArizaTipi, &ariza.ArizaTipiTanimID, &ariza.ArizaTipiTanimAdi, &ariza.ArizaNedenID,
			&ariza.ArizaNedenAdi, &ariza.KisiSayisi, &ariza.YapilanMudahele, &ariza.TestTalebiKaynagi,
			&ariza.IsletmeUzmaniEmployeeID, &ariza.IsletmeUzmaniFullName, &ariza.BakimOperatoruEmployeeID,
			&ariza.BakimOperatoruFullName, &ariza.BakimUzmaniEmployeeID, &ariza.BakimUzmaniFullName,
			&ariza.ProsesCalisaniEmployeeID, &ariza.ProsesCalisaniFullName, &ariza.KaliteKontrolVePaketlemeEmployeeID,
			&ariza.KaliteKontrolVePaketlemeFullName,
			&ariza.Aciklama, &ariza.TestAciklamasi, &ariza.ArizaLogID, &ariza.CreateUserID, &ariza.CreateUserFullName,
			&ariza.UpdateUserID, &ariza.UpdateUserFullName, &createDate, &updateDate, &ariza.TekrarCount,
		)
		if err != nil {
			return nil, err
		}

		if createDate != nil {
			date := createDate.Format("2006-01-02 15:04:05")
			ariza.CreateDate = &date
		}
		ariza.Ekip = &ekip

		if containts := strings.Contains(req.ArizaStatusu, "Test"); containts {
			resultMap["Test"] = append(resultMap["Test"], ariza)
		} else {
			resultMap[ekip] = append(resultMap[ekip], ariza)
			if *ariza.ArizaTipi == "Elektrik" {
				resultMap["ELEKTRİK"] = append(resultMap["ELEKTRİK"], ariza)
			} else if *ariza.ArizaTipi == "Mekanik" {
				resultMap["MEKANİK"] = append(resultMap["MEKANİK"], ariza)
			}
		}

	}

	return resultMap, nil

}

func (p *DB) GetAriza(req *models.GetArizaRequestBody) (*models.ArizaDetay, error) {

	var createDate time.Time

	query := `
		WITH blok_bilgileri
		AS (
			SELECT s.ariza_details_id
				,jsonb_agg(s.blok_bilgi) AS blok_bilgi
			FROM (
				SELECT ariza_details_id
					,json_build_object('BlokNo', blok_id, 'Pozisyonlar', pozisyonlar) AS blok_bilgi
				FROM (
					SELECT ariza_details_id
						,blok_id
						,json_agg(pozisyon) AS pozisyonlar
					FROM elterminali.ariza_blok
					WHERE is_deleted = false
					GROUP BY ariza_details_id
						,blok_id
					) AS subquery
				) s
			GROUP BY ariza_details_id
			)
			,pozisyonlar
		AS (
			SELECT row_number() OVER (
					PARTITION BY ariza_id ORDER BY ariza_details_id DESC
					) AS islem_sira_no
				,ariza_id
				,array_agg(pozisyon) AS pozisyonlar
				,array_agg(distinct blok_id) as bloklar
			FROM elterminali.ariza_blok
			GROUP BY ariza_id
				,ariza_details_id
			)
			,asamalar
		AS (
			SELECT CONCAT (
					row_number() OVER (
						PARTITION BY islem_turu
						,asama ORDER BY ariza_details_id
						)
					,'. '
					,asama
					) AS asama_son
				,islem_turu
				,ariza_details_id
			FROM elterminali.ariza_details
			WHERE ariza_id = $1
				AND islem_turu != 'Atama 1'
			ORDER BY ariza_details_id ASC
			)
		SELECT coalesce(CASE 
					WHEN ad.islem_turu = 'Atama 1'
						THEN LEAD(asm.asama_son) OVER (
								ORDER BY ad.ariza_details_id
								)
					ELSE asm.asama_son
					END, CONCAT (
					left(lag(asm.asama_son) OVER (
							ORDER BY ad.ariza_details_id
							), 1)
					,'. '
					,asama
					))
			,ad.ariza_id
			,ad.ariza_details_id
			,ao.ariza_statusu
			,ad.ariza_statusu AS detail_status
			,ad.asama
			,ad.makine_id
			,ao.makine_id
			,m1.alt_makine AS makine_adi
			,p.pozisyonlar 	   			-- güncel pozisyonlar
			,p.bloklar
			,m2.alt_makine AS guncel_makine_adi
			,ad.blok_tarafi    			-- detay blok tarafı
			,ao.blok_tarafi              
			,bb.blok_bilgi
			,ad.ariza_tipi     			-- detay ariza_tipi
			,ao.ariza_tipi 
			,ad.ariza_tipi_tanim_id  	-- detay tanim tipi
			,ao.ariza_tipi_tanim_id 
			,att.ariza_tipi_tanim_adi 	-- detay
			,att2.ariza_tipi_tanim_adi 	-- 
			,ad.ariza_neden_id
			,an.ariza_neden_adi
			,ad.kisi_sayisi
			,ad.yapilan_mudahele
			,ad.test_talebi_kaynagi
			,ad.aciklama     --
			,ao.aciklama
			,ad.test_aciklamasi
			,ad.ariza_log_id
			,ad.create_user_id
			,cu.full_name AS create_user_full_name
			,ad.create_date
			,yetki.yetki_grup_id
			,ad.toplam_etkilenen_pozisyon_sayisi
			,ad.karsilikli_bloklar_etkilenmis
			,ad.ariza_durumu 
		FROM elterminali.ariza_details ad
		LEFT JOIN elterminali.ariza ao ON ao.ariza_id = ad.ariza_id
		LEFT JOIN blok_bilgileri bb ON bb.ariza_details_id = ad.ariza_details_id
		LEFT JOIN elterminali.ariza_neden an ON an.ariza_neden_id = ad.ariza_neden_id
		LEFT JOIN elterminali.ariza_tipi_tanım att ON att.ariza_tipi_tanim_id = ad.ariza_tipi_tanim_id
		LEFT JOIN elterminali.ariza_tipi_tanım att2 ON att2.ariza_tipi_tanim_id = ao.ariza_tipi_tanim_id
		LEFT JOIN filament_prm.alt_makine m1 ON m1.id = ao.makine_id
		LEFT JOIN filament_prm.alt_makine m2 ON m2.id = ad.makine_id
		LEFT JOIN asamalar asm ON ad.ariza_details_id = asm.ariza_details_id
		LEFT JOIN techup.employee cu ON cu.employee_id = ad.create_user_id
		LEFT JOIN pozisyonlar p ON p.ariza_id = ao.ariza_id
			AND islem_sira_no = 1
		LEFT JOIN (
			SELECT er.employee_id
				,er.yetki_grup_id
			FROM techup.employee_role er
			LEFT JOIN techup.permission_group pg ON pg.id = er.yetki_grup_id
			WHERE er.yetki_grup_id IN (
					116
					,117
					,118
					,119
					,120
					)
			) yetki ON yetki.employee_id = ad.create_user_id
		WHERE ad.ariza_id = $1
			AND ad.is_deleted = false
		ORDER BY ad.ariza_details_id ASC



		`

	rows, err := p.postgres.Query(context.Background(), query, req.ArizaID)
	if err != nil {
		return nil, err
	}

	var ariza models.ArizaDetay
	detayMap := make(map[string]models.Ariza)

	for rows.Next() {

		var row models.Ariza
		var asamaSon string
		var guncelMakineAdi *string
		var asama string
		var aciklama, blokTarafi, arizaTipi, arizaTipiTanimAdi, arizaDurumu string
		var arizaTipiTanimId, makineID int
		var bloklar []int
		err = rows.Scan(
			&asamaSon, &row.ArizaID, &row.ArizaDetayID, &row.ArizaStatusu, &row.ArizaStatusu, &asama, &row.MakineID,
			&makineID,
			&row.MakineAdi, &row.Pozisyonlar, &bloklar, &guncelMakineAdi, &row.BlokTarafi, &blokTarafi, &row.ArizaBloklar,
			&row.ArizaTipi, &arizaTipi,
			&row.ArizaTipiTanimID, &arizaTipiTanimId, &row.ArizaTipiTanimAdi, &arizaTipiTanimAdi, &row.ArizaNedenID,
			&row.ArizaNedenAdi, &row.KisiSayisi,
			&row.YapilanMudahele, &row.TestTalebiKaynagi, &row.Aciklama, &aciklama, &row.TestAciklamasi, &row.ArizaLogID,
			&row.CreateUserID, &row.CreateUserFullName, &createDate, &row.YetkiGrupID,
			&row.ToplamEtkilenenPozisyonSayisi, &row.KarsilikliBloklarEtkilenmisMi, &arizaDurumu,
		)
		if err != nil {
			return nil, err
		}

		if asama == "Root" {

			ariza.ArizaDurumu = &arizaDurumu
			ariza.Pozisyonlar = row.Pozisyonlar
			ariza.Bloklar = bloklar
			ariza.MakineAdi = guncelMakineAdi
			ariza.IsEmriOlusturmaTarihi = createDate.Format("2006-01-02 15:04:05")
			ariza.ArizaID = row.ArizaID
			ariza.MakineID = &makineID
			ariza.ArizaTipi = &arizaTipi
			ariza.ArizaTipiTanimID = &arizaTipiTanimId
			ariza.BlokTarafi = &blokTarafi
			ariza.Aciklama = &aciklama
			ariza.CreateUserName = row.CreateUserFullName
			ariza.ArizaTipiTanimAdi = &arizaTipiTanimAdi

		} else {

			if arizaDurumu == "Onaylanan Test" {

				if req.ArizaStatusu == "Sonlanan Arızalar" {
					durum := "Sonlanan Arıza"
					ariza.ArizaDurumu = &durum
				} else if req.ArizaStatusu == "Onaylanan Testler" {
					durum := "Onaylanan Test"
					ariza.ArizaDurumu = &durum
				} else {
					ariza.ArizaDurumu = &arizaDurumu
				}

			} else {
				ariza.ArizaDurumu = &arizaDurumu
			}

			a := time.Now().Format("2006-01-02 15:04:05")
			row.BaslamaTarihi = &a
			row.BitisTarihi = &a
			row.Suresi = 1
			row.Asama = &asamaSon
			row.Pozisyonlar = nil
			detayMap[asamaSon] = row

		}

	}

	objects := detayMap
	objectList := make([]models.Ariza, 0)
	for _, value := range objects {
		objectList = append(objectList, value)
	}

	ariza.Detay = objectList

	return &ariza, nil

}

func (p *DB) IsTransferi(req *models.IsTransferiRequest) error {

	qUpdateArizaStatusu := `
		UPDATE elterminali.ariza
		SET ariza_statusu = $1
		WHERE ariza_id = $2;
	`

	_, err := p.postgres.Exec(context.Background(), qUpdateArizaStatusu, ARIZA_ISTRANSFERI, req.ArizaID)
	if err != nil {
		return err
	}

	queryDetails := `
		INSERT INTO elterminali.ariza_details (
			ariza_id
			,ariza_statusu
			,asama
			,ariza_durumu
			,bakim_operatoru_employee_id
			,create_user_id
			,create_date
			,islem_turu
			)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
			);
		`

	createDate := time.Now().Format("2006-01-02 15:04:05")

	_, err = p.postgres.Exec(context.Background(), queryDetails, req.ArizaID, ARIZA_ISTRANSFERI, "Bakım Aşaması",
		"Transferden Gelen Arıza", req.UserID, req.UserID, createDate, "Düzenleme")
	if err != nil {
		return err
	}

	return err

}

// 1. senaryo olumlu
func (p *DB) ArizaYapildi(req *models.ArizaYapildiRequest) (*models.EmployeesAndCreatedUser, error) {

	var makineAdi string
	var pozisyonlar []int

	qUpdateArizaStatusu := `
		WITH updated
		AS (
			UPDATE elterminali.ariza
			SET ariza_statusu = $1
				,ariza_neden_id = $2
				,yapilan_mudahele = $3
				,karsilikli_bloklar_etkilenmis = $4
				,toplam_etkilenen_pozisyon_sayisi = $5
			WHERE ariza_id = $6 RETURNING *
			)
		SELECT (
				SELECT alt_makine
				FROM filament_prm.alt_makine
				WHERE id = updated.makine_id LIMIT 1
				) AS makine_adi
			,(
				SELECT array_agg(ab.pozisyon) AS pozisyonlar
				FROM elterminali.ariza_blok ab
				WHERE ab.ariza_id = updated.ariza_id
				)
		FROM updated;

	`
	err := p.postgres.QueryRow(context.Background(), qUpdateArizaStatusu, ARIZA_YAPILDI, req.ArizaNedenID,
		req.YapilanMudahele, req.KarsilikliBloklarEtkilenmisMi, req.ToplamEtkilenenPozisyonSayisi, req.ArizaID).
		Scan(&makineAdi, &pozisyonlar)
	if err != nil {
		return nil, err

	}

	createDate := time.Now().Format("2006-01-02 15:04:05")
	var arizaDetailsID int
	qInsertArizaDetails := `
		INSERT INTO elterminali.ariza_details (
			ariza_id
			,ariza_neden_id
			,yapilan_mudahele
			,karsilikli_bloklar_etkilenmis
			,toplam_etkilenen_pozisyon_sayisi
			,ariza_statusu
			,asama
			,ariza_durumu
			,bakim_operatoru_employee_id
			,create_user_id
			,create_date
			,islem_turu
			)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
			)
		RETURNING ariza_details_id;
		`
	err = p.postgres.QueryRow(context.Background(), qInsertArizaDetails, req.ArizaID, req.ArizaNedenID, req.YapilanMudahele,
		req.KarsilikliBloklarEtkilenmisMi, req.ToplamEtkilenenPozisyonSayisi, ARIZA_YAPILDI, "Bakım Aşaması",
		"Tamamlanan Arıza", req.UserID, req.UserID, createDate, "Düzenleme").Scan(&arizaDetailsID)
	if err != nil {
		return nil, err
	}

	qInsertKullanilanMalzemeler := `
		INSERT INTO elterminali.ariza_kullanilan_malzemeler (
			ariza_id
			,ariza_details_id
			,kullanilan_malzeme_id
			,miktar
			)
		VALUES (
			$1, $2, $3, $4
			);
		`
	for _, malzeme := range req.KullanilanMalzemeler {
		_, err = p.postgres.Exec(context.Background(), qInsertKullanilanMalzemeler, req.ArizaID, malzeme.MalzemeID,
			malzeme.Miktar, arizaDetailsID)
		if err != nil {
			return nil, err
		}
	}

	//Bildirim

	employees, err := p.GetEmployeeIDsByPermission([]int{PROSES_CALISANI}, req.UserID, strconv.Itoa(req.ArizaID), ARIZA_YAPILDI,
		makineAdi, "", pozisyonlar)
	if err != nil {
		return nil, err
	}

	return employees, nil

}

// senaryo3.6 ariza yapildi işareti için
func (p *DB) BakimUzmaniTarafindanAcilanArizaYapildi(req *models.BakimUzmaniTarafindanAcilanArizaYapildiRequest) error {
	qUpdateArizaStatusu := `
		UPDATE elterminali.ariza
		SET ariza_statusu = $1, ariza_neden_id = $2, yapilan_mudahele = $3
		WHERE ariza_id = $4;
		`
	_, err := p.postgres.Exec(context.Background(), qUpdateArizaStatusu, ARIZA_SONLANDI, req.ArizaNedenID,
		req.YapilanMudahele, req.ArizaID)
	if err != nil {
		return err
	}

	qInsertKullanilanMalzemeler := `
		INSERT INTO elterminali.ariza_kullanilan_malzemeler (
			ariza_id
			,kullanilan_malzeme_id
			,miktar
			)
		VALUES (
			$1, $2, $3
			);
		`
	for _, malzeme := range req.KullanilanMalzemeler {
		_, err = p.postgres.Exec(context.Background(), qInsertKullanilanMalzemeler, req.ArizaID, malzeme.MalzemeID,
			malzeme.Miktar)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *DB) ProsesPersoneliArizaOnayi(req *models.ArizaOnayiRequest) (*models.EmployeesAndCreatedUser, error) {

	var makineAdi string
	var pozisyonlar []int
	var arizaTipiTanimAdi string
	query := `
		WITH updated
		AS (
		UPDATE elterminali.ariza
		SET ariza_statusu = $1, proses_calisani_employee_id = $2
		WHERE ariza_id = $3
		RETURNING *
			)
		SELECT (
				SELECT alt_makine
				FROM filament_prm.alt_makine
				WHERE id = updated.makine_id LIMIT 1
				) AS makine_adi
			,(
				SELECT array_agg(ab.pozisyon) AS pozisyonlar
				FROM elterminali.ariza_blok ab
				WHERE ab.ariza_id = updated.ariza_id
				)
			,(
				SELECT att.ariza_tipi_tanim_adi
				FROM elterminali.ariza_tipi_tanım att
				WHERE att.ariza_tipi_tanim_id = updated.ariza_tipi_tanim_id
				)
		FROM updated;
		`

	err := p.postgres.QueryRow(context.Background(), query, ARIZA_PROSES_ONAYI, req.UserID, req.ArizaID).Scan(
		&makineAdi, &pozisyonlar, &arizaTipiTanimAdi)
	if err != nil {
		return nil, err
	}

	qInsertDetails := `
		INSERT INTO elterminali.ariza_details (
			ariza_id
			,ariza_statusu
			,asama
			,ariza_durumu
			,proses_calisani_employee_id
			,create_user_id
			,create_date
			,islem_turu
			)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
			);
		`
	createDate := time.Now().Format("2006-01-02 15:04:05")

	_, err = p.postgres.Exec(context.Background(), qInsertDetails, req.ArizaID, ARIZA_PROSES_ONAYI, "Proses Kontrol Aşaması",
		"Proses Onaylı Arıza", req.UserID, req.UserID, createDate, "Atama 2")

	if err != nil {
		return nil, err
	}

	//Bildirim
	employees, err := p.GetEmployeeIDsByPermission([]int{ISLETME_UZMANI}, req.UserID,
		strconv.Itoa(req.ArizaID), ARIZA_PROSES_ONAYI, makineAdi, arizaTipiTanimAdi, pozisyonlar)
	if err != nil {
		return nil, err
	}

	return employees, err
}

func (p *DB) IsletmeUzmaniTestTalebi(req *models.IsletmeUzmaniTestTalebiRequest) (*models.EmployeesAndCreatedUser, error) {

	var makineAdi string
	var pozisyonlar []int
	var arizaTipiTanimAdi string

	query := `
		WITH updated
		AS (
			UPDATE elterminali.ariza
			SET ariza_statusu = $1, test_talebi_kaynagi = $2, isletme_uzmani_employee_id = $3
			WHERE ariza_id = $4
			RETURNING *
			)
		SELECT (
				SELECT alt_makine
				FROM filament_prm.alt_makine
				WHERE id = updated.makine_id LIMIT 1
				) AS makine_adi
			,(
				SELECT array_agg(ab.pozisyon) AS pozisyonlar
				FROM elterminali.ariza_blok ab
				WHERE ab.ariza_id = updated.ariza_id
				)
			,(
				SELECT att.ariza_tipi_tanim_adi
				FROM elterminali.ariza_tipi_tanım att
				WHERE att.ariza_tipi_tanim_id = updated.ariza_tipi_tanim_id
				)
		FROM updated;
	
		`

	err := p.postgres.QueryRow(context.Background(), query, ARIZA_TEST_TALEBI, TEST_KAYNAGI_ARIZADAN_GELEN, req.UserID,
		req.ArizaID).Scan(&makineAdi, &pozisyonlar, &arizaTipiTanimAdi)
	if err != nil {
		return nil, err
	}

	qInsertDetails := `
		INSERT INTO elterminali.ariza_details (
			ariza_id
			,ariza_statusu
			,asama
			,ariza_durumu
			,isletme_uzmani_employee_id
			,create_user_id
			,create_date
			,islem_turu
			)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
			);
		`
	createDate := time.Now().Format("2006-01-02 15:04:05")

	_, err = p.postgres.Exec(context.Background(), qInsertDetails, req.ArizaID, ARIZA_TEST_TALEBI, "Test Talebi",
		"Yeni Test", req.UserID, req.UserID, createDate, "Düzenleme")
	if err != nil {
		return nil, err
	}

	//Bildirim

	employees, err := p.GetEmployeeIDsByPermission([]int{KALITE_KONTROL_VE_PAKETLEME}, req.UserID, strconv.Itoa(req.ArizaID),
		ARIZA_TEST_TALEBI, makineAdi, arizaTipiTanimAdi, pozisyonlar)
	if err != nil {
		return nil, err
	}

	return employees, err

}

func (p *DB) KaliteKontrolVePaketlemeTestOnayi(req *models.KaliteKontrolVePaketlemeTestOnayiRequest) (*models.EmployeesAndCreatedUser, error) {

	var makineAdi string
	var pozisyonlar []int
	var arizaTipiTanimAdi string
	query := `
			WITH updated
				AS (
					UPDATE elterminali.ariza
					SET ariza_statusu = $1, kalite_kontrol_ve_paketleme_employee_id = $2, test_aciklamasi = $3
					WHERE ariza_id = $4
					RETURNING *
					)
				SELECT (
						SELECT alt_makine
						FROM filament_prm.alt_makine
						WHERE id = updated.makine_id LIMIT 1
						) AS makine_adi
					,(
						SELECT array_agg(ab.pozisyon) AS pozisyonlar
						FROM elterminali.ariza_blok ab
						WHERE ab.ariza_id = updated.ariza_id
						)
					,(
						SELECT att.ariza_tipi_tanim_adi
						FROM elterminali.ariza_tipi_tanım att
						WHERE att.ariza_tipi_tanim_id = updated.ariza_tipi_tanim_id
						)
				FROM updated;
		`

	err := p.postgres.QueryRow(context.Background(), query, ARIZA_TEST_ONAYI, req.UserID, req.Aciklama, req.ArizaID).Scan(
		&makineAdi, &pozisyonlar, &arizaTipiTanimAdi)
	if err != nil {
		return nil, err
	}

	qInsertDetails := `
		INSERT INTO elterminali.ariza_details (
			ariza_id
			,ariza_statusu
			,asama
			,ariza_durumu
			,kalite_kontrol_ve_paketleme_employee_id
			,create_user_id
			,create_date
			,islem_turu
			)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
			);
		`
	createDate := time.Now().Format("2006-01-02 15:04:05")

	_, err = p.postgres.Exec(context.Background(), qInsertDetails, req.ArizaID, ARIZA_TEST_ONAYI, "Test Sonuçları",
		"Onaylanan Test", req.UserID, req.UserID, createDate, "Düzenleme")
	if err != nil {
		return nil, err
	}

	//Bildirim
	employees, err := p.GetEmployeeIDsByPermission([]int{ISLETME_UZMANI, BAKIM_UZMANI}, req.UserID, strconv.Itoa(req.ArizaID),
		ARIZA_TEST_ONAYI, makineAdi, arizaTipiTanimAdi, pozisyonlar)
	if err != nil {
		return nil, err
	}

	return employees, err
}

func (p *DB) SAPGonder(req *models.SAPGonderRequest) error {
	query := `
		UPDATE elterminali.ariza
		SET ariza_statusu = $1
		WHERE ariza_id = $2;
		`

	_, err := p.postgres.Exec(context.Background(), query, ARIZA_SAP_GONDERILDI, req.ArizaID)
	return err
}

func (p *DB) BakimOperatoruAtama(req *models.AtamaRequest) (*models.EmployeesAndCreatedUser, error) {

	var makineAdi string
	var pozisyonlar []int
	var arizaTipi string
	qUpdateArizaStatusu := `
		WITH updated
		AS (
			UPDATE elterminali.ariza
			SET bakim_operatoru_employee_id = $1
				,ariza_statusu = $2
			WHERE ariza_id = $3 
			RETURNING *
			)
		SELECT (
				SELECT alt_makine
				FROM filament_prm.alt_makine
				WHERE id = updated.makine_id LIMIT 1
				) AS makine_adi
			,(
				SELECT array_agg(ab.pozisyon) AS pozisyonlar
				FROM elterminali.ariza_blok ab
				WHERE ab.ariza_id = updated.ariza_id
				)
			,(
				SELECT att.ariza_tipi_tanim_adi
				FROM elterminali.ariza_tipi_tanım att
				WHERE att.ariza_tipi_tanim_id = updated.ariza_tipi_tanim_id
				)
		FROM updated;
	`

	err := p.postgres.QueryRow(context.Background(), qUpdateArizaStatusu, req.UserID,
		ARIZA_ISLEMDE, req.ArizaID).Scan(&makineAdi, &pozisyonlar, &arizaTipi)
	if err != nil {
		return nil, err
	}

	createDate := time.Now().Format("2006-01-02 15:04:05")

	qInserArizaDetails := `
		INSERT INTO elterminali.ariza_details (
			ariza_id
			,bakim_operatoru_employee_id
			,create_user_id
			,ariza_statusu
			,asama
			,ariza_durumu
			,create_date
			,islem_turu
			)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
			)
		RETURNING ariza_details_id;
		`
	_, err = p.postgres.Exec(context.Background(), qInserArizaDetails, req.ArizaID, req.UserID,
		req.UserID, ARIZA_ISLEMDE, "Bakım Aşaması", "İşlemdeki Arıza", createDate, "Atama 1")
	if err != nil {
		return nil, err
	}

	//BILDIRIM
	employees, err := p.GetEmployeeIDsByPermission([]int{BAKIM_UZMANI}, req.UserID, strconv.Itoa(req.ArizaID),
		ARIZA_ISLEMDE, makineAdi, arizaTipi, pozisyonlar)
	if err != nil {
		return nil, err
	}

	return employees, err
}

func (p *DB) BakimUzmaniAtama(req *models.AtamaRequest) (*models.EmployeesAndCreatedUser, error) {
	query := `
		UPDATE elterminali.ariza
		SET bakim_uzmani_employee_id = $1
		WHERE ariza_id = $2;
	`
	_, err := p.postgres.Exec(context.Background(), query, req.UserID, req.ArizaID)
	if err != nil {
		return nil, err
	}

	createDate := time.Now().Format("2006-01-02 15:04:05")

	qInserArizaDetails := `
		INSERT INTO elterminali.ariza_details (
			ariza_id
			,bakim_operatoru_employee_id
			,create_user_id
			,ariza_statusu
			,asama
			,ariza_durumu
			,create_date
			,islem_turu
			)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
			)
		RETURNING ariza_details_id;
		`
	_, err = p.postgres.Exec(context.Background(), qInserArizaDetails, req.ArizaID, req.UserID,
		req.UserID, ARIZA_ISLEMDE, "Bakım Aşaması", "İşlemdeki Arıza", createDate, "Atama 1")
	if err != nil {
		return nil, err
	}

	return nil, err

}

func (p *DB) KaliteKontrolVePaketlemeUzmaniAtama(req *models.AtamaRequest) (*models.EmployeesAndCreatedUser, error) {
	query := `
		UPDATE elterminali.ariza
		SET ariza_statusu = $1 , kalite_kontrol_ve_paketleme_employee_id = $2
		WHERE ariza_id = $3;
	`

	_, err := p.postgres.Exec(context.Background(), query, ARIZA_TEST_TALEBI_ALINDI, req.UserID,
		req.ArizaID)
	if err != nil {
		return nil, err
	}

	qInsertDetails := `
		INSERT INTO elterminali.ariza_details (
			ariza_id
			,ariza_statusu
			,asama
			,ariza_durumu
			,kalite_kontrol_ve_paketleme_employee_id
			,create_user_id
			,create_date
			,islem_turu
			)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
			);
		`
	createDate := time.Now().Format("2006-01-02 15:04:05")

	_, err = p.postgres.Exec(context.Background(), qInsertDetails, req.ArizaID, ARIZA_TEST_TALEBI_ALINDI, "Test Sonuçları",
		"Açık Test", req.UserID, req.UserID, createDate, "Atama 1")
	if err != nil {
		return nil, err
	}

	return nil, err

}

func (p *DB) GetArizaTipTanimlar() ([]*models.ArizaTipTanimlar, error) {
	var arizaTipTanimlar []*models.ArizaTipTanimlar
	var arizaTipTanim models.ArizaTipTanimlar

	query := `
		SELECT ariza_tipi_tanim_id, ariza_tipi_tanim_adi
		FROM elterminali.ariza_tipi_tanım;
	`

	rows, err := p.postgres.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&arizaTipTanim.ArizaTipiTanimID, &arizaTipTanim.ArizaTipiTanimAdi)
		if err != nil {
			return nil, err
		}
		arizaTipTanimlar = append(arizaTipTanimlar, &arizaTipTanim)
	}

	return arizaTipTanimlar, nil
}

func (p *DB) GetArizaNedenleri() ([]*models.ArizaNedenleri, error) {
	var arizaNedenleri []*models.ArizaNedenleri
	var arizaNeden models.ArizaNedenleri

	query := `
		SELECT ariza_neden_id, ariza_neden_adi
		FROM elterminali.ariza_neden;
	`

	rows, err := p.postgres.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&arizaNeden.ArizaNedenID, &arizaNeden.ArizaNedenAdi)
		if err != nil {
			return nil, err
		}
		arizaNedenleri = append(arizaNedenleri, &arizaNeden)
	}

	return arizaNedenleri, nil
}

func (p *DB) GetKullanilanMalzemeler() ([]*models.KullanilanMalzemeler, error) {
	var kullanilanMalzemeler []*models.KullanilanMalzemeler
	var kullanilanMalzeme models.KullanilanMalzemeler

	query := `
		SELECT malzeme_id, malzeme_adi
		FROM elterminali.kullanilan_malzemeler;
	`

	rows, err := p.postgres.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&kullanilanMalzeme.MalzemeID, &kullanilanMalzeme.MalzemeAdi)
		if err != nil {
			return nil, err
		}
		kullanilanMalzemeler = append(kullanilanMalzemeler, &kullanilanMalzeme)
	}

	return kullanilanMalzemeler, nil
}

func (p *DB) ArizaTekrari(req *models.ArizaTekrarRequest) (*models.EmployeesAndCreatedUser, error) {

	var makineAdi string
	var pozisyonlar []int
	var arizaTipiTanimAdi string

	query := `
		WITH updated
		AS (
			UPDATE elterminali.ariza
			SET ariza_statusu = $1, ariza_tekrar_nedeni = $2
			WHERE ariza_id = $3
			RETURNING *
			)
		SELECT (
				SELECT alt_makine
				FROM filament_prm.alt_makine
				WHERE id = updated.makine_id LIMIT 1
				) AS makine_adi
			,(
				SELECT array_agg(ab.pozisyon) AS pozisyonlar
				FROM elterminali.ariza_blok ab
				WHERE ab.ariza_id = updated.ariza_id
				)
			,(
				SELECT att.ariza_tipi_tanim_adi
				FROM elterminali.ariza_tipi_tanım att
				WHERE att.ariza_tipi_tanim_id = updated.ariza_tipi_tanim_id
				)
		FROM updated;
		;
	`
	err := p.postgres.QueryRow(context.Background(), query, req.ArizaTekrarNedeni, ARIZA_TEKRAR, req.ArizaID).Scan(
		&makineAdi, &pozisyonlar, &arizaTipiTanimAdi)
	if err != nil {
		return nil, err
	}

	qInsertDetails := `
		INSERT INTO elterminali.ariza_details (
			ariza_id
			,ariza_statusu
			,asama
			,ariza_durumu
			,isletme_uzmani_employee_id
			,create_user_id
			,create_date
			,islem_turu
			)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
			);
		`
	createDate := time.Now().Format("2006-01-02 15:04:05")

	_, err = p.postgres.Exec(context.Background(), qInsertDetails, req.ArizaID, ARIZA_TEKRAR, "Proses Kontrol Aşaması",
		"Tekrar Açılmış Arıza (Proses Redli)", req.UserID, req.UserID, createDate, "Düzenleme")
	if err != nil {
		return nil, err
	}

	//bildirim

	employees, err := p.GetEmployeeIDsByPermission([]int{BAKIM_UZMANI, BAKIM_OPERATORU}, req.UserID, strconv.Itoa(req.ArizaID),
		ARIZA_TEKRAR, makineAdi, arizaTipiTanimAdi, pozisyonlar)
	if err != nil {
		return nil, err
	}

	return employees, nil

}

func (p *DB) KaliteKontrolVePaketlemeOlumsuzTestSonucu(req *models.KaliteKontrolVePaketlemeOlumsuzTestSonucuRequest) (
	*models.EmployeesAndCreatedUser, error) {

	var makineAdi string
	var pozisyonlar []int
	var arizaTipiTanimAdi string

	query := `
		WITH updated
		AS (
			UPDATE elterminali.ariza
			SET ariza_statusu = $1, test_aciklamasi = $2
			WHERE ariza_id = $3
			RETURNING *
			)
		SELECT (
				SELECT alt_makine
				FROM filament_prm.alt_makine
				WHERE id = updated.makine_id LIMIT 1
				) AS makine_adi
			,(
				SELECT array_agg(ab.pozisyon) AS pozisyonlar
				FROM elterminali.ariza_blok ab
				WHERE ab.ariza_id = updated.ariza_id
				)
			,(
				SELECT att.ariza_tipi_tanim_adi
				FROM elterminali.ariza_tipi_tanım att
				WHERE att.ariza_tipi_tanim_id = updated.ariza_tipi_tanim_id
				)
		FROM updated;
		`

	err := p.postgres.QueryRow(context.Background(), query, ARIZA_OLUMSUZ_TEST, req.TestAciklama, req.ArizaID).Scan(
		&makineAdi, &pozisyonlar, &arizaTipiTanimAdi)
	if err != nil {
		return nil, err
	}

	qInsertDetails := `
		INSERT INTO elterminali.ariza_details (
			ariza_id
			,ariza_statusu
			,asama
			,ariza_durumu
			,isletme_uzmani_employee_id
			,create_user_id
			,create_date
			,islem_turu
			)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
			);
		`
	createDate := time.Now().Format("2006-01-02 15:04:05")

	_, err = p.postgres.Exec(context.Background(), qInsertDetails, req.ArizaID, ARIZA_OLUMSUZ_TEST, "Test Sonuçları",
		"Tekrar Açılmış Arıza (LAB Redli)", req.UserID, req.UserID, createDate, "Düzenleme")
	if err != nil {
		return nil, err
	}

	//bildirim

	employees, err := p.GetEmployeeIDsByPermission([]int{BAKIM_UZMANI, BAKIM_OPERATORU}, req.UserID, strconv.Itoa(req.ArizaID),
		ARIZA_OLUMSUZ_TEST, makineAdi, arizaTipiTanimAdi, pozisyonlar)
	if err != nil {
		return nil, err
	}

	return employees, nil

}

// TODO: BAKILACAK MANTIKSAL BIR HATA VAR
func (p *DB) KaliteKontrolVePaketlemeArizaOnayladi(req *models.ArizaOnayiRequest) (*models.EmployeesAndCreatedUser, error) {
	var makineAdi string
	var pozisyonlar []int
	var arizaTipiTanimAdi string

	query := `
		WITH updated
		AS (
			UPDATE elterminali.ariza
			SET ariza_statusu = $1, kalite_kontrol_ve_paketleme_employee_id = $2
			WHERE ariza_id = $3
			RETURNING *
			)
		SELECT (
				SELECT alt_makine
				FROM filament_prm.alt_makine
				WHERE id = updated.makine_id LIMIT 1
				) AS makine_adi
			,(
				SELECT array_agg(ab.pozisyon) AS pozisyonlar
				FROM elterminali.ariza_blok ab
				WHERE ab.ariza_id = updated.ariza_id
				)
			,(
				SELECT att.ariza_tipi_tanim_adi
				FROM elterminali.ariza_tipi_tanım att
				WHERE att.ariza_tipi_tanim_id = updated.ariza_tipi_tanim_id
				)
		FROM updated;

`
	err := p.postgres.QueryRow(context.Background(), query, ARIZA_SONLANDI, req.UserID, req.ArizaID).Scan(
		&makineAdi, &pozisyonlar, &arizaTipiTanimAdi)
	if err != nil {
		return nil, err
	}

	queryDetails := `
		INSERT INTO elterminali.ariza_details (
			ariza_id
			,ariza_statusu	
			,asama	
			,ariza_durumu
			,isletme_uzmani_employee_id
			,create_user_id
			,create_date
			,islem_turu
			)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
			);
		`

	createDate := time.Now().Format("2006-01-02 15:04:05")

	_, err = p.postgres.Exec(context.Background(), queryDetails, req.ArizaID, ARIZA_SONLANDI, "Test Sonuçları",
		"Sonlanan Arıza", req.UserID, req.UserID, createDate, "Düzenleme")
	if err != nil {
		return nil, err
	}

	//bildirim

	employees, err := p.GetEmployeeIDsByPermission([]int{BAKIM_UZMANI, BAKIM_OPERATORU}, req.UserID, strconv.Itoa(req.ArizaID),
		ARIZA_SONLANDI, makineAdi, arizaTipiTanimAdi, pozisyonlar)
	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (p *DB) GetUserYetkiID(userID int) (int, error) {

	var yetkiID int
	query1 := `
		SELECT er.yetki_grup_id
		FROM techup.employee_role er
		WHERE er.employee_id = $1
			AND er.is_deleted = false
			AND er.yetki_grup_id IN (
				116
				,117
				,118
				,119
				,120
				)

	`

	err := p.postgres.QueryRow(context.Background(), query1, userID).Scan(&yetkiID)
	if err != nil {
		return 0, err
	}

	return yetkiID, nil
}
