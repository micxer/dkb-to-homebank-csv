package main

import (
	"strings"
	"testing"
)

func TestPaymentTypeAbschluss(t *testing.T) {
	assertPaymentType(t, "ABSCHLUSS", "0")
}

func TestLohnGehaltRente(t *testing.T) {
	assertPaymentType(t, "LOHN, GEHALT, RENTE", "4")
}

func TestOnlineÜberweisung(t *testing.T) {
	assertPaymentType(t, "ONLINE-UEBERWEISUNG", "4")
}

func TestÜberweisung(t *testing.T) {
	assertPaymentType(t, "ÜBERWEISUNG", "4")
}

func TestRueckÜberweisung(t *testing.T) {
	assertPaymentType(t, "RÜCKÜBERWEISUNG", "4")
}

func TestWertpapiere(t *testing.T) {
	assertPaymentType(t, "WERTPAPIERE", "4")
}

func TestOnlineZinsDividende(t *testing.T) {
	assertPaymentType(t, "ZINS/DIVIDENDE", "4")
}

func TestAuftrag(t *testing.T) {
	assertPaymentType(t, "AUFTRAG", "5")
}
func TestUmbuchung(t *testing.T) {
	assertPaymentType(t, "UMBUCHUNG", "5")
}

func TestKartenzahlungAbrechnung(t *testing.T) {
	assertPaymentType(t, "KARTENZAHLUNG/-ABRECHNUNG", "6")
}

func TestSepaElvLastschrift(t *testing.T) {
	assertPaymentType(t, "SEPA-ELV-LASTSCHRIFT", "6")
}

func TestDauerauftrag(t *testing.T) {
	assertPaymentType(t, "DAUERAUFTRAG", "7")
}

func TestGutschrift(t *testing.T) {
	assertPaymentType(t, "GUTSCHRIFT", "8")
}

func TestLastschrift(t *testing.T) {
	assertPaymentType(t, "LASTSCHRIFT", "11")
}

func TestFolgeLastschrift(t *testing.T) {
	assertPaymentType(t, "FOLGELASTSCHRIFT", "11")
}

func assertPaymentType(t *testing.T, paymentString string, expectedValue string) {
	dkbRecord := dkbCsv{}
	dkbRecord.Buchungstext = paymentString

	homebankRecord := convertFromDkb(&dkbRecord)

	if homebankRecord.Payment != expectedValue {
		t.Errorf("Expected %v, got %v", expectedValue, homebankRecord.Payment)
	}
}

func TestIbanBicVsKontoNrBlz(t *testing.T) {
	dkbRecord := dkbCsv{}
	dkbRecord.Kontonummer = "0000202051"
	dkbRecord.Blz = "12030000"

	homebankRecord := convertFromDkb(&dkbRecord)

	if !strings.Contains(homebankRecord.Info, "Konto-Nr.: 0000202051, BLZ: 12030000") {
		t.Errorf("Expected %v, got %v", "Konto-Nr.: 0000202051, BLZ: 12030000", homebankRecord.Payment)
	}

	dkbRecord.Kontonummer = "DE02120300000000202051"
	dkbRecord.Blz = "BYLADEM1001"

	homebankRecord = convertFromDkb(&dkbRecord)

	if !strings.Contains(homebankRecord.Info, "IBAN: DE02120300000000202051, BIC: BYLADEM1001") {
		t.Errorf("Expected %v, got %v", "IBAN: DE02120300000000202051, BIC: BYLADEM1001", homebankRecord.Payment)
	}
}

func TestGlaeubigerIdMandatsreferenzKundenreferenz(t *testing.T) {
	dkbRecord := dkbCsv{}
	dkbRecord.GlaeubigerID = "DE0012345678"
	dkbRecord.Kundenreferenz = "00012345"

	homebankRecord := convertFromDkb(&dkbRecord)

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

	homebankRecord = convertFromDkb(&dkbRecord)

	if homebankRecord.Info != "Mandatsreferenz: MAN007\nKundenreferenz: 00012345" {
		t.Errorf(
			"Expected %v, got %v",
			"Mandatsreferenz: MAN007\nKundenreferenz: 00012345",
			homebankRecord.Info,
		)
	}

	dkbRecord.GlaeubigerID = "DE0012345678"

	homebankRecord = convertFromDkb(&dkbRecord)

	if homebankRecord.Info != "Gläubiger-ID: DE0012345678\nMandatsreferenz: MAN007\nKundenreferenz: 00012345" {
		t.Errorf(
			"Expected %v, got %v",
			"Gläubiger-ID: DE0012345678\nMandatsreferenz: MAN007\nKundenreferenz: 00012345",
			homebankRecord.Info,
		)
	}
}

func TestCategoryEmpty(t *testing.T) {
	dkbRecord := dkbCsv{}
	homebankRecord := convertFromDkb(&dkbRecord)

	if homebankRecord.Payment != "" {
		t.Errorf("Expected empty string, got %v", homebankRecord.Category)
	}
}

func TestTagsEmpty(t *testing.T) {
	dkbRecord := dkbCsv{}
	homebankRecord := convertFromDkb(&dkbRecord)

	if homebankRecord.Tags != "" {
		t.Errorf("Expected empty string, got %v", homebankRecord.Tags)
	}
}

func TestAmount(t *testing.T) {
	dkbRecord := dkbCsv{}
	dkbRecord.BetragEur = "12,34"
	homebankRecord := convertFromDkb(&dkbRecord)

	if homebankRecord.Amount != "12,34" {
		t.Errorf("Expected %v, got %v", dkbRecord.BetragEur, homebankRecord.Amount)
	}
}

func TestMemo(t *testing.T) {
	dkbRecord := dkbCsv{}
	dkbRecord.Verwendungszweck = "This is a test!"
	homebankRecord := convertFromDkb(&dkbRecord)

	if homebankRecord.Memo != dkbRecord.Verwendungszweck {
		t.Errorf("Expected %v, got %v", dkbRecord.Verwendungszweck, homebankRecord.Memo)
	}
}

func TestPayee(t *testing.T) {
	dkbRecord := dkbCsv{}
	dkbRecord.AuftraggeberBeguenstigter = "The Shop"
	homebankRecord := convertFromDkb(&dkbRecord)

	if homebankRecord.Payee != dkbRecord.AuftraggeberBeguenstigter {
		t.Errorf("Expected %v, got %v", dkbRecord.AuftraggeberBeguenstigter, homebankRecord.Payee)
	}
}
