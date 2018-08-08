package main

import (
	"encoding/csv"

	"os"
	"testing"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func TestAbschluss(t *testing.T) {
	dkbRecord, homebankRecord := loadCsvSingleRow(t, "abschluss")
	convertRecordAndVerify(t, dkbRecord, homebankRecord)
}

func TestLohnGehaltRente(t *testing.T) {
	dkbRecord, homebankRecord := loadCsvSingleRow(t, "lohn_gehalt_rente")
	convertRecordAndVerify(t, dkbRecord, homebankRecord)
}

func TestDauerauftrag(t *testing.T) {
	dkbRecord, homebankRecord := loadCsvSingleRow(t, "dauerauftrag")
	convertRecordAndVerify(t, dkbRecord, homebankRecord)
}

func TestLastschrift(t *testing.T) {
	dkbRecord, homebankRecord := loadCsvSingleRow(t, "lastschrift")
	convertRecordAndVerify(t, dkbRecord, homebankRecord)
}

func TestAuftrag(t *testing.T) {
	dkbRecord, homebankRecord := loadCsvSingleRow(t, "auftrag")
	convertRecordAndVerify(t, dkbRecord, homebankRecord)
}

func TestGutschrift(t *testing.T) {
	dkbRecord, homebankRecord := loadCsvSingleRow(t, "gutschrift")
	convertRecordAndVerify(t, dkbRecord, homebankRecord)
}

func TestIbanBicVsKontoNrBlz(t *testing.T) {
	dkbRecords, homebankRecords := loadCsv(t, "ibankontonr")
	for index, dkbRecord := range dkbRecords {
		convertRecordAndVerify(t, dkbRecord, homebankRecords[index])
	}
}

func convertRecordAndVerify(t *testing.T, dkbRecord dkbCsv, homebankRecord homebankCsv) {
	NewRecord := convertFromDkb(&dkbRecord)

	if NewRecord != homebankRecord {
		t.Errorf("Expected %v, got %v", homebankRecord, NewRecord)
	}
}

func loadCsvSingleRow(t *testing.T, filename string) (dkbCsv, homebankCsv) {
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

func loadCsv(t *testing.T, filename string) ([]dkbCsv, []homebankCsv) {
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

	dkbRecords := make([]dkbCsv, len(rows))
	for _, row := range rows {
		dkbRecord := dkbCsv{}

		dkbRecord.Buchungstag = row[0]
		dkbRecord.Wertstellung = row[1]
		dkbRecord.Buchungstext = row[2]
		dkbRecord.AuftraggeberBeguenstigter = row[3]
		dkbRecord.Verwendungszweck = row[4]
		dkbRecord.Kontonummer = row[5]
		dkbRecord.Blz = row[6]
		dkbRecord.BetragEur = row[7]
		dkbRecord.GlaeubigerID = row[8]
		dkbRecord.Mandatsreferenz = row[9]
		dkbRecord.Kundenreferenz = row[10]
		dkbRecords = append(dkbRecords, dkbRecord)
	}

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
	homebankRecords := make([]homebankCsv, len(rows))
	for _, row := range rows {
		homebankRecord := homebankCsv{}

		homebankRecord.Date = row[0]
		homebankRecord.Payment = row[1]
		homebankRecord.Info = row[2]
		homebankRecord.Payee = row[3]
		homebankRecord.Memo = row[4]
		homebankRecord.Amount = row[5]
		homebankRecord.Category = row[6]
		homebankRecord.Tags = row[7]

		homebankRecords = append(homebankRecords, homebankRecord)
	}

	return dkbRecords, homebankRecords
}
