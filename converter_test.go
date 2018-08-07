package main

import (
	"encoding/csv"

	"os"
	"testing"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func TestAbschluss(t *testing.T) {
	dkbRecord, homebankRecord := loadCsv(t, "abschluss")
	convertRecordAndVerify(t, dkbRecord, homebankRecord)
}

func TestLohnGehaltRente(t *testing.T) {
	dkbRecord, homebankRecord := loadCsv(t, "lohn_gehalt_rente")
	convertRecordAndVerify(t, dkbRecord, homebankRecord)
}

func TestDauerauftrag(t *testing.T) {
	dkbRecord, homebankRecord := loadCsv(t, "dauerauftrag")
	convertRecordAndVerify(t, dkbRecord, homebankRecord)
}

func TestLastschrift(t *testing.T) {
	dkbRecord, homebankRecord := loadCsv(t, "lastschrift")
	convertRecordAndVerify(t, dkbRecord, homebankRecord)
}

func TestAuftrag(t *testing.T) {
	dkbRecord, homebankRecord := loadCsv(t, "auftrag")
	convertRecordAndVerify(t, dkbRecord, homebankRecord)
}

func convertRecordAndVerify(t *testing.T, dkbRecord dkbCsv, homebankRecord homebankCsv) {
	NewRecord := convertFromDkb(&dkbRecord)

	if NewRecord != homebankRecord {
		t.Errorf("Expected %v, got %v", homebankRecord, NewRecord)
	}
}

func loadCsv(t *testing.T, filename string) (dkbCsv, homebankCsv) {
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
	dkbRecord := dkbCsv{}

	dkbRecord.Buchungstag = rows[0][0]
	dkbRecord.Wertstellung = rows[0][1]
	dkbRecord.Buchungstext = rows[0][2]
	dkbRecord.AuftraggeberBeguenstigter = rows[0][3]
	dkbRecord.Verwendungszweck = rows[0][4]
	dkbRecord.Kontonummer = rows[0][5]
	dkbRecord.Blz = rows[0][6]
	dkbRecord.BetragEur = rows[0][7]
	dkbRecord.GlaeubigerID = rows[0][8]
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
	homebankRecord := homebankCsv{}

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
