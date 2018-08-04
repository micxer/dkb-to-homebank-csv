package main

import (
	"encoding/csv"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"os"
	"testing"
)

func TestAbschluss(t *testing.T) {
	dkbRecord, homebankRecord := LoadCsv(t, "abschluss")
	convertRecordAndVerify(t, dkbRecord, homebankRecord)
}

func TestLohnGehaltRente(t *testing.T) {
	dkbRecord := DkbCsv{
		Buchungstag:               "24.06.13",
		Wertstellung:              "24.06.13",
		Buchungstext:              "LOHN, GEHALT, RENTE",
		AuftraggeberBeguenstigter: "ACME GMBH",
		Verwendungszweck:          "LOHN / GEHALT         06/13",
		Kontonummer:               "0000202051",
		Blz:                       "12030000",
		BetragEur:                 "1234,56",
		GlaeubigerId:              "",
		Mandatsreferenz:           "",
		Kundenreferenz:            "",
	}

	HomebankRecord := HomebankCsv{
		Date:     "24.06.13",
		Payment:  "4",
		Info:     "",
		Payee:    "ACME GMBH",
		Memo:     "LOHN / GEHALT         06/13",
		Amount:   "1234,56",
		Category: "",
		Tags:     "",
	}

	convertRecordAndVerify(t, dkbRecord, HomebankRecord)
}

func TestDauerauftrag(t *testing.T) {
	dkbRecord := DkbCsv{
		Buchungstag:               "01.10.13",
		Wertstellung:              "01.10.13",
		Buchungstext:              "DAUERAUFTRAG",
		AuftraggeberBeguenstigter: "RUNDFUNKGEBUEHRENABWICKLUNG",
		Verwendungszweck:          "",
		Kontonummer:               "2024100",
		Blz:                       "70050000",
		BetragEur:                 "-51,84",
		GlaeubigerId:              "",
		Mandatsreferenz:           "",
		Kundenreferenz:            "",
	}

	HomebankRecord := HomebankCsv{
		Date:     "01.10.13",
		Payment:  "7",
		Info:     "",
		Payee:    "RUNDFUNKGEBUEHRENABWICKLUNG",
		Memo:     "",
		Amount:   "-51,84",
		Category: "",
		Tags:     "",
	}

	convertRecordAndVerify(t, dkbRecord, HomebankRecord)
}

func convertRecordAndVerify(t *testing.T, dkbRecord DkbCsv, homebankRecord HomebankCsv) {
	NewRecord := ConvertFromDkb(&dkbRecord)

	if NewRecord != homebankRecord {
		t.Errorf("Expected %v, got %v", homebankRecord, NewRecord)
	}
}

func LoadCsv(t *testing.T, filename string) (DkbCsv, HomebankCsv) {
	f, err := os.Open("testdata/" + filename + "_dkb.csv")
	if err != nil {
		t.Fatal(err)
	}
	reader := csv.NewReader(transform.NewReader(f, charmap.ISO8859_15.NewDecoder()))
	reader.Comma = ';'
	rows, err := reader.ReadAll()
	f.Close()
	if err != nil {
		t.Fatal(err)
	}
	dkbRecord := DkbCsv{}

	dkbRecord.Buchungstag = rows[0][0]
	dkbRecord.Wertstellung = rows[0][1]
	dkbRecord.Buchungstext = rows[0][2]
	dkbRecord.AuftraggeberBeguenstigter = rows[0][3]
	dkbRecord.Verwendungszweck = rows[0][4]
	dkbRecord.Kontonummer = rows[0][5]
	dkbRecord.Blz = rows[0][6]
	dkbRecord.BetragEur = rows[0][7]
	dkbRecord.GlaeubigerId = rows[0][8]
	dkbRecord.Mandatsreferenz = rows[0][9]
	dkbRecord.Kundenreferenz = rows[0][10]

	f, err = os.Open("testdata/" + filename + "_homebank.csv")
	if err != nil {
		t.Fatal(err)
	}
	reader = csv.NewReader(transform.NewReader(f, charmap.ISO8859_15.NewDecoder()))
	reader.Comma = ';'
	rows, err = reader.ReadAll()
	f.Close()
	if err != nil {
		t.Fatal(err)
	}
	homebankRecord := HomebankCsv{}

	homebankRecord.Date = rows[0][0]
	homebankRecord.Payment = rows[0][1]
	homebankRecord.Info = rows[0][2]
	homebankRecord.Payee = rows[0][3]
	homebankRecord.Memo = rows[0][4]
	homebankRecord.Amount = rows[0][5]
	homebankRecord.Category = rows[0][6]
	homebankRecord.Tags = rows[0][7]


	return dkbRecord, homebankRecord
}
