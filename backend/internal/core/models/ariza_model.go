package models

type Ariza struct {
	Ekip                               *string        `json:",omitempty"`
	DetayStatusu                       string         `json:"DetayStatusu"`
	ArizaID                            int            `json:"ArizaID"`
	ArizaDetayID                       *int           `json:",omitempty"`
	ArizaStatusu                       *string        `json:",omitempty"`
	TekrarCount                        *int           `json:",omitempty"`
	Asama                              *string        `json:",omitempty"`
	YetkiGrupID                        *int           `json:",omitempty"`
	MakineID                           *int           `json:",omitempty"`
	Pozisyonlar                        []int          `json:",omitempty"`
	MakineAdi                          *string        `json:",omitempty"`
	BlokTarafi                         *string        `json:",omitempty"`
	ArizaBloklar                       []ArizaBloklar `json:",omitempty"`
	ArizaTipi                          *string        `json:",omitempty"`
	ArizaTipiTanimID                   *int           `json:",omitempty"`
	ArizaTipiTanimAdi                  *string        `json:",omitempty"`
	ArizaNedenID                       *int           `json:",omitempty"`
	ArizaNedenAdi                      *string        `json:",omitempty"`
	KisiSayisi                         *int           `json:",omitempty"`
	YapilanMudahele                    *string        `json:",omitempty"`
	TestTalebiKaynagi                  *string        `json:",omitempty"`
	IsletmeUzmaniEmployeeID            *int           `json:",omitempty"`
	IsletmeUzmaniFullName              *string        `json:",omitempty"`
	BakimOperatoruEmployeeID           *int           `json:",omitempty"`
	BakimOperatoruFullName             *string        `json:",omitempty"`
	BakimUzmaniEmployeeID              *int           `json:",omitempty"`
	BakimUzmaniFullName                *string        `json:",omitempty"`
	BaslamaTarihi                      *string        `json:",omitempty"`
	BitisTarihi                        *string        `json:",omitempty"`
	Suresi                             int            `json:",omitempty"`
	TesteGonderenEmployeeName          *string        `json:",omitempty"`
	ProsesCalisaniEmployeeID           *int           `json:",omitempty"`
	ProsesCalisaniFullName             *string        `json:",omitempty"`
	KaliteKontrolVePaketlemeEmployeeID *int           `json:",omitempty"`
	KaliteKontrolVePaketlemeFullName   *string        `json:",omitempty"`
	Aciklama                           *string        `json:",omitempty"`
	TestAciklamasi                     *string        `json:",omitempty"`
	ArizaLogID                         *int           `json:",omitempty"`
	CreateUserID                       *int           `json:"CreateUserID"`
	CreateUserFullName                 *string        `json:"CreateUserFullName"`
	UpdateUserID                       *int           `json:",omitempty"`
	UpdateUserFullName                 *string        `json:",omitempty"`
	CreateDate                         *string        `json:"CreateDate"`
	UpdateDate                         *string        `json:",omitempty"`
	ToplamEtkilenenPozisyonSayisi      *int           `json:",omitempty"`
	KarsilikliBloklarEtkilenmisMi      *bool          `json:",omitempty"`
}

type ArizaDetay struct {
	Pozisyonlar           []int       `json:"Pozisyonlar"`
	Bloklar               []int       `json:"Bloklar"`
	MakineID              *int        `json:"MakineID"`
	ArizaTipiTanimID      *int        `json:"ArizaTipiTanimID"`
	ArizaTipi             *string     `json:"ArizaTipi"`
	Aciklama              *string     `json:"Aciklama"`
	BlokTarafi            *string     `json:"BlokTarafi"`
	MakineAdi             *string     `json:"MakineAdi"`
	IsEmriOlusturmaTarihi string      `json:"IsEmriOlusturmaTarihi"`
	Detay                 interface{} `json:"Detay"`
	ArizaID               int         `json:"ArizaID"`
	ArizaDurumu           *string     `json:"ArizaDurumu"`
	CreateUserName        *string     `json:"CreateUserName"`
	ArizaTipiTanimAdi     *string     `json:"ArizaTipiTanimAdi"`
}

type ArizaCreateRequestBody struct {
	MakineID         int    `json:"MakineID"`
	BlokTarafi       string `json:"BlokTarafi"`
	ArizaTipi        string `json:"ArizaTipi"`
	ArizaTipiTanimID int    `json:"ArizaTipiTanimID"`
	Aciklama         string `json:"Aciklama"`
	Bloklar          []int  `json:"Bloklar"`
	Pozisyonlar      []int  `json:"Pozisyonlar"`
	UserID           int    `json:"UserID"`
}

type GetArizalarRequestBody struct {
	ArizaStatusu string `json:"ArizaStatusu"`
	UserID       int    `json:"UserID"`
}
type GetArizaRequestBody struct {
	ArizaID      int    `json:"ArizaID"`
	ArizaStatusu string `json:"ArizaStatusu"`
	UserID       int    `json:"UserID"`
}

type ArizaTipTanimlar struct {
	ArizaTipiTanimAdi string `json:"ArizaTipiTanimAdi"`
	ArizaTipiTanimID  int    `json:"ArizaTipiTanimID"`
}

type ArizaNedenleri struct {
	ArizaNedenAdi string `json:"ArizaNedenAdi"`
	ArizaNedenID  int    `json:"ArizaNedenID"`
}

type KullanilanMalzemeler struct {
	MalzemeAdi string `json:"MalzemeAdi"`
	MalzemeID  int    `json:"MalzemeID"`
}
type ArizaBloklar struct {
	BlokNo      int   `json:"BlokNo"`
	Pozisyonlar []int `json:"Pozisyonlar"`
}

type RequestBodyAriza struct {
	ArizaID int `json:"ArizaID"`
}

type IsTransferiRequest struct {
	ArizaID int `json:"ArizaID"`
	UserID  int `json:"UserID"`
}

type IsletmeUzmaniTestTalebiRequest struct {
	ArizaID int `json:"ArizaID"`
	UserID  int `json:"UserID"`
}

type ArizaYapildiRequest struct {
	ArizaID                       int                 `json:"ArizaID"`
	ArizaNedenID                  int                 `json:"ArizaNedenID"`
	YapilanMudahele               string              `json:"YapilanMudahele"`
	KullanilanMalzemeler          []KullanilanMalzeme `json:"KullanilanMalzemeler"`
	UserID                        int                 `json:"UserID"`
	KarsilikliBloklarEtkilenmisMi bool                `json:"KarsilikliBloklarEtkilenmisMi"`
	ToplamEtkilenenPozisyonSayisi int                 `json:"ToplamEtkilenPozisyonSayisi"`
}

type BakimUzmaniTarafindanAcilanArizaYapildiRequest struct {
	ArizaID              int                 `json:"ArizaID"`
	ArizaNedenID         int                 `json:"ArizaNedenID"`
	YapilanMudahele      string              `json:"YapilanMudahele"`
	KullanilanMalzemeler []KullanilanMalzeme `json:"KullanilanMalzemeler"`
	UserID               int                 `json:"UserID"`
}

type KullanilanMalzeme struct {
	MalzemeID int `json:"MalzemeID"`
	Miktar    int `json:"Miktar"`
}
type ArizaOnayiRequest struct {
	ArizaID int `json:"ArizaID"`
	UserID  int `json:"UserID"`
}

type BakimUzmaniBildirimRequest struct {
	ArizaID int `json:"ArizaID"`
	UserID  int `json:"UserID"`
}

type AtamaRequest struct {
	ArizaID int `json:"ArizaID"`
	UserID  int `json:"UserID"`
}

type KaliteKontrolVePaketlemeTestOnayiRequest struct {
	ArizaID  int    `json:"ArizaID"`
	UserID   int    `json:"UserID"`
	Aciklama string `json:"Aciklama"`
}

type ArizaDuzenlemeRequest struct {
	ArizaID              int                 `json:"ArizaID"`
	MakineID             int                 `json:"MakineID"`
	BlokTarafi           string              `json:"BlokTarafi"`
	ArizaTipi            string              `json:"ArizaTipi"`
	ArizaTipiTanimID     int                 `json:"ArizaTipiTanimID"`
	Aciklama             string              `json:"Aciklama"`
	Bloklar              []int               `json:"Bloklar"`
	Pozisyonlar          []int               `json:"Pozisyonlar"`
	ArizaNedenID         int                 `json:"ArizaNedenID"`
	YapilanMudahale      string              `json:"YapilanMudahale"`
	KullanilanMalzemeler []KullanilanMalzeme `json:"KullanilanMalzemeler"`
	UserID               int                 `json:"UserID"`
}

type SAPGonderRequest struct {
	ArizaID int `json:"ArizaID"`
	UserID  int `json:"UserID"`
}

type BildirimRequest struct {
	ArizaID int `json:"ArizaID"`
}

type ArizaTekrarRequest struct {
	ArizaID           int    `json:"ArizaID"`
	ArizaTekrarNedeni string `json:"ArizaTekrarNedeni"`
	UserID            int    `json:"UserID"`
}

type KaliteKontrolVePaketlemeOlumsuzTestSonucuRequest struct {
	ArizaID      int    `json:"ArizaID"`
	TestAciklama string `json:"TestAciklama"`
	UserID       int    `json:"UserID"`
}
