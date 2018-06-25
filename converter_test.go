package main

import "testing"

func TestAbschluss(t *testing.T) {
	dkbRecord := DkbCsv{
		Buchungstag:               "30.12.13",
		Wertstellung:              "01.01.14",
		Buchungstext:              "ABSCHLUSS",
		AuftraggeberBeguenstigter: "",
		Verwendungszweck:          "Abrechnung 30.12.2013      siehe Anlage<br />Abrechnung 30.12.2013<br />Information zur Abrechnung<br />Kontostand am 30.12.2013                                        4.321,12 +<br />Abrechnungszeitraum vom 01.10.2013 bis 31.12.2013<br />Zinsen für Guthaben                                                  1,23+<br /> 0,2000 v.H. Haben-Zins bis 30.12.2013<br />Abrechnung 31.12.2013                                                1,23+<br />Sollzinssätze am 30.12.2013<br /> 7,9000 v.H. für Dispositionskredit<br />(aktuelle Kreditlinie       1.000,00)<br />12,0000 v.H. für Kontoüberziehungen<br />über die Kreditlinie hinaus<br />Kontostand/Rechnungsabschluss am 30.12.2013                     1.234,56 +<br />Rechnungsnummer: 20131230-BY111-00001234567",
		Kontonummer:               "0000202051",
		Blz:                       "12030000",
		BetragEur:                 "1,23",
		GlaeubigerId:              "",
		Mandatsreferenz:           "",
		Kundenreferenz:            "",
	}

	HomebankRecord := HomebankCsv{
		Date:     "01.01.14",
		Payment:  "0",
		Info:     "",
		Payee:    "",
		Memo:     "Abrechnung 30.12.2013      siehe Anlage<br />Abrechnung 30.12.2013<br />Information zur Abrechnung<br />Kontostand am 30.12.2013                                        4.321,12 +<br />Abrechnungszeitraum vom 01.10.2013 bis 31.12.2013<br />Zinsen für Guthaben                                                  1,23+<br /> 0,2000 v.H. Haben-Zins bis 30.12.2013<br />Abrechnung 31.12.2013                                                1,23+<br />Sollzinssätze am 30.12.2013<br /> 7,9000 v.H. für Dispositionskredit<br />(aktuelle Kreditlinie       1.000,00)<br />12,0000 v.H. für Kontoüberziehungen<br />über die Kreditlinie hinaus<br />Kontostand/Rechnungsabschluss am 30.12.2013                     1.234,56 +<br />Rechnungsnummer: 20131230-BY111-00001234567",
		Amount:   "1,23",
		Category: "",
		Tags:     "",
	}

	NewRecord := ConvertFromDkb(&dkbRecord)

	if NewRecord != HomebankRecord {
		t.Errorf("Expected %v, got %v", HomebankRecord, NewRecord)
	}
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

	NewRecord := ConvertFromDkb(&dkbRecord)

	if NewRecord != HomebankRecord {
		t.Errorf("Expected %v, got %v", HomebankRecord, NewRecord)
	}
}
