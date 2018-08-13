package main

import (
	"log"
	"os"
	"strings"
	"testing"
)

func TestGiroDate(t *testing.T) {
	dkbRecord := dkbGiroCsv{}
	dkbRecord.Wertstellung = "23.12.12"

	homebankRecord := convertFromDkbGiro(&dkbRecord)

	if homebankRecord.Date != dkbRecord.Wertstellung {
		t.Errorf("Expected %v, got %v", dkbRecord.Wertstellung, homebankRecord.Date)
	}
}

func TestGiroPaymentTypeAbschluss(t *testing.T) {
	assertPaymentType(t, "ABSCHLUSS", "0")
}

func TestGiroLohnGehaltRente(t *testing.T) {
	assertPaymentType(t, "LOHN, GEHALT, RENTE", "4")
}

func TestGiroOnlineÜberweisung(t *testing.T) {
	assertPaymentType(t, "ONLINE-UEBERWEISUNG", "4")
}

func TestGiroÜberweisung(t *testing.T) {
	assertPaymentType(t, "ÜBERWEISUNG", "4")
}

func TestGiroRueckÜberweisung(t *testing.T) {
	assertPaymentType(t, "RÜCKÜBERWEISUNG", "4")
}

func TestGiroWertpapiere(t *testing.T) {
	assertPaymentType(t, "WERTPAPIERE", "4")
}

func TestGiroOnlineZinsDividende(t *testing.T) {
	assertPaymentType(t, "ZINS/DIVIDENDE", "4")
}

func TestGiroAuftrag(t *testing.T) {
	assertPaymentType(t, "AUFTRAG", "5")
}
func TestGiroUmbuchung(t *testing.T) {
	assertPaymentType(t, "UMBUCHUNG", "5")
}

func TestGiroKartenzahlungAbrechnung(t *testing.T) {
	assertPaymentType(t, "KARTENZAHLUNG/-ABRECHNUNG", "6")
}

func TestGiroSepaElvLastschrift(t *testing.T) {
	assertPaymentType(t, "SEPA-ELV-LASTSCHRIFT", "6")
}

func TestGiroDauerauftrag(t *testing.T) {
	assertPaymentType(t, "DAUERAUFTRAG", "7")
}

func TestGiroGutschrift(t *testing.T) {
	assertPaymentType(t, "GUTSCHRIFT", "8")
}

func TestGiroLastschrift(t *testing.T) {
	assertPaymentType(t, "LASTSCHRIFT", "11")
}

func TestGiroFolgeLastschrift(t *testing.T) {
	assertPaymentType(t, "FOLGELASTSCHRIFT", "11")
}

func assertPaymentType(t *testing.T, paymentString string, expectedValue string) {
	dkbRecord := dkbGiroCsv{}
	dkbRecord.Buchungstext = paymentString

	homebankRecord := convertFromDkbGiro(&dkbRecord)

	if homebankRecord.Payment != expectedValue {
		t.Errorf("Expected %v, got %v", expectedValue, homebankRecord.Payment)
	}
}

func TestGiroIbanBicVsKontoNrBlz(t *testing.T) {
	dkbRecord := dkbGiroCsv{}
	dkbRecord.Kontonummer = "0000202051"
	dkbRecord.Blz = "12030000"

	homebankRecord := convertFromDkbGiro(&dkbRecord)

	if !strings.Contains(homebankRecord.Info, "Konto-Nr.: 0000202051, BLZ: 12030000") {
		t.Errorf("Expected %v, got %v", "Konto-Nr.: 0000202051, BLZ: 12030000", homebankRecord.Payment)
	}

	dkbRecord.Kontonummer = "DE02120300000000202051"
	dkbRecord.Blz = "BYLADEM1001"

	homebankRecord = convertFromDkbGiro(&dkbRecord)

	if !strings.Contains(homebankRecord.Info, "IBAN: DE02120300000000202051, BIC: BYLADEM1001") {
		t.Errorf("Expected %v, got %v", "IBAN: DE02120300000000202051, BIC: BYLADEM1001", homebankRecord.Payment)
	}
}

func TestGiroGlaeubigerIdMandatsreferenzKundenreferenz(t *testing.T) {
	dkbRecord := dkbGiroCsv{}
	dkbRecord.GlaeubigerID = "DE0012345678"
	dkbRecord.Kundenreferenz = "00012345"

	homebankRecord := convertFromDkbGiro(&dkbRecord)

	if homebankRecord.Info != "Gläubiger-ID: DE0012345678\nKundenreferenz: 00012345" {
		t.Errorf(
			"Expected %v, got %v",
			"Gläubiger-ID: DE0012345678\nKundenreferenz: 00012345\n",
			homebankRecord.Info,
		)
	}

	dkbRecord.GlaeubigerID = ""
	dkbRecord.Mandatsreferenz = "MAN007"
	dkbRecord.Kundenreferenz = "00012345"

	homebankRecord = convertFromDkbGiro(&dkbRecord)

	if homebankRecord.Info != "Mandatsreferenz: MAN007\nKundenreferenz: 00012345" {
		t.Errorf(
			"Expected %v, got %v",
			"Mandatsreferenz: MAN007\nKundenreferenz: 00012345",
			homebankRecord.Info,
		)
	}

	dkbRecord.GlaeubigerID = "DE0012345678"

	homebankRecord = convertFromDkbGiro(&dkbRecord)

	if homebankRecord.Info != "Gläubiger-ID: DE0012345678\nMandatsreferenz: MAN007\nKundenreferenz: 00012345" {
		t.Errorf(
			"Expected %v, got %v",
			"Gläubiger-ID: DE0012345678\nMandatsreferenz: MAN007\nKundenreferenz: 00012345",
			homebankRecord.Info,
		)
	}
}

func TestGiroCategoryEmpty(t *testing.T) {
	dkbRecord := dkbGiroCsv{}
	homebankRecord := convertFromDkbGiro(&dkbRecord)

	if homebankRecord.Payment != "" {
		t.Errorf("Expected empty string, got %v", homebankRecord.Category)
	}
}

func TestGiroTagsEmpty(t *testing.T) {
	dkbRecord := dkbGiroCsv{}
	homebankRecord := convertFromDkbGiro(&dkbRecord)

	if homebankRecord.Tags != "" {
		t.Errorf("Expected empty string, got %v", homebankRecord.Tags)
	}
}

func TestGiroAmount(t *testing.T) {
	dkbRecord := dkbGiroCsv{}
	dkbRecord.BetragEur = "12,34"
	homebankRecord := convertFromDkbGiro(&dkbRecord)

	if homebankRecord.Amount != "12,34" {
		t.Errorf("Expected %v, got %v", dkbRecord.BetragEur, homebankRecord.Amount)
	}
}

func TestGiroMemo(t *testing.T) {
	dkbRecord := dkbGiroCsv{}
	dkbRecord.Verwendungszweck = "This is a test!"
	homebankRecord := convertFromDkbGiro(&dkbRecord)

	if homebankRecord.Memo != dkbRecord.Verwendungszweck {
		t.Errorf("Expected %v, got %v", dkbRecord.Verwendungszweck, homebankRecord.Memo)
	}
}

func TestGiroPayee(t *testing.T) {
	dkbRecord := dkbGiroCsv{}
	dkbRecord.AuftraggeberBeguenstigter = "The Shop"
	homebankRecord := convertFromDkbGiro(&dkbRecord)

	if homebankRecord.Payee != dkbRecord.AuftraggeberBeguenstigter {
		t.Errorf("Expected %v, got %v", dkbRecord.AuftraggeberBeguenstigter, homebankRecord.Payee)
	}
}

func TestDetectFiletype(t *testing.T) {
	inputfile, err := os.Open("test_csvs/credit_test.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer inputfile.Close()
	filetype := detectFiletype(inputfile)
	if filetype != "CREDIT_CSV" {
		t.Errorf("Expected %v, got %v", "CREDIT_CSV", filetype)
	}

	inputfile, err = os.Open("test_csvs/giro_test.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer inputfile.Close()
	filetype = detectFiletype(inputfile)
	if filetype != "GIRO_CSV" {
		t.Errorf("Expected %v, got %v", "GIRO_CSV", filetype)
	}
}

func TestCreditTransactionCleared(t *testing.T) {
	dkbRecord := dkbCreditCsv{}
	dkbRecord.UmsatzAbgerechnet = "Nein"
	homebankRecord := convertFromDkbCredit(&dkbRecord)

	if homebankRecord != nil {
		t.Errorf("Expected empty struct, got %v", homebankRecord)
	}

	dkbRecord.UmsatzAbgerechnet = ""
	dkbRecord.UmsatzAbgerechnetOld = "Nein"
	homebankRecord = convertFromDkbCredit(&dkbRecord)

	if homebankRecord != nil {
		t.Errorf("Expected empty struct, got %v", homebankRecord)
	}
}

func TestCreditPayee(t *testing.T) {
	dkbRecord := dkbCreditCsv{}
	dkbRecord.Beschreibung = "AMAZON MKTPLACE PMTSAMAZON.COM"
	homebankRecord := convertFromDkbCredit(&dkbRecord)

	if homebankRecord.Payee != dkbRecord.Beschreibung {
		t.Errorf("Expected %v, got %v", dkbRecord.Beschreibung, homebankRecord.Payee)
	}

	dkbRecord.Beschreibung = ""
	dkbRecord.Umsatzbeschreibung = "Lidl Vertriebs GmbH & Co"
	homebankRecord = convertFromDkbCredit(&dkbRecord)

	if homebankRecord.Payee != dkbRecord.Umsatzbeschreibung {
		t.Errorf("Expected %v, got %v", dkbRecord.Umsatzbeschreibung, homebankRecord.Payee)
	}
}

func TestCreditAmount(t *testing.T) {
	dkbRecord := dkbCreditCsv{}
	dkbRecord.BetragEur = "-11,04"
	homebankRecord := convertFromDkbCredit(&dkbRecord)

	if homebankRecord.Amount != dkbRecord.BetragEur {
		t.Errorf("Expected %v, got %v", dkbRecord.BetragEur, homebankRecord.Amount)
	}
}

func TestCreditInfo(t *testing.T) {
	dkbRecord := dkbCreditCsv{}
	dkbRecord.UrspruenglicherBetrag = "-123,00 SEK"
	homebankRecord := convertFromDkbCredit(&dkbRecord)

	if homebankRecord.Info != dkbRecord.UrspruenglicherBetrag {
		t.Errorf("Expected %v, got %v", dkbRecord.UrspruenglicherBetrag, homebankRecord.Info)
	}
}

func TestCreditPayment(t *testing.T) {
	dkbRecord := dkbCreditCsv{}
	dkbRecord.UmsatzAbgerechnet = "Ja"
	homebankRecord := convertFromDkbCredit(&dkbRecord)

	if homebankRecord.Payment != "1" {
		t.Errorf("Payment type should be 1, found %v", homebankRecord.Payment)
	}
}

func TestCreditDate(t *testing.T) {
	dkbRecord := dkbCreditCsv{}
	dkbRecord.Wertstellung = "1.2.12"
	dkbRecord.Belegdatum = "1.3.12"
	homebankRecord := convertFromDkbCredit(&dkbRecord)

	if homebankRecord.Date != dkbRecord.Wertstellung {
		t.Errorf("Expected %v, got %v", dkbRecord.Wertstellung, homebankRecord.Date)
	}
}

func TestCreditMemo(t *testing.T) {
	dkbRecord := dkbCreditCsv{}
	dkbRecord.Wertstellung = "1.2.12"
	dkbRecord.Belegdatum = "1.3.12"
	homebankRecord := convertFromDkbCredit(&dkbRecord)

	if homebankRecord.Memo != "Belegedatum: "+dkbRecord.Belegdatum {
		t.Errorf("Expected %v, got %v", dkbRecord.Belegdatum, homebankRecord.Memo)
	}
}
