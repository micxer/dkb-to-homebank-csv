/**
 * Copyright (c) 2017-2018, Michael Weinrich
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/gocarina/gocsv"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type configuration struct {
	InputFilename  string
	OutputFilename string
}

var config *configuration

type dkbGiroCsv struct {
	Buchungstag               string `csv:"Buchungstag"`
	Wertstellung              string `csv:"Wertstellung"`
	Buchungstext              string `csv:"Buchungstext"`
	AuftraggeberBeguenstigter string `csv:"Auftraggeber / Begünstigter"`
	Verwendungszweck          string `csv:"Verwendungszweck"`
	Kontonummer               string `csv:"Kontonummer"`
	Blz                       string `csv:"BLZ"`
	BetragEur                 string `csv:"Betrag (EUR)"`
	GlaeubigerID              string `csv:"Gläubiger-ID"`
	Mandatsreferenz           string `csv:"Mandatsreferenz"`
	Kundenreferenz            string `csv:"Kundenreferenz"`
}

type dkbCreditCsv struct {
	UmsatzAbgerechnet     string `csv:"Umsatz abgerechnet und nicht im Saldo enthalten"`
	UmsatzAbgerechnetOld  string `csv:"Umsatz abgerechnet"`
	Wertstellung          string `csv:"Wertstellung"`
	Belegdatum            string `csv:"Belegdatum"`
	Beschreibung          string `csv:"Beschreibung"`
	Umsatzbeschreibung    string `csv:"Umsatzbeschreibung"`
	BetragEur             string `csv:"Betrag (EUR)"`
	UrspruenglicherBetrag string `csv:"Ursprünglicher Betrag"`
}

type homebankCsv struct {
	Date     string `csv:"date"`
	Payment  string `csv:"payment"`
	Info     string `csv:"info"`
	Payee    string `csv:"payee"`
	Memo     string `csv:"memo"`
	Amount   string `csv:"amount"`
	Category string `csv:"category"`
	Tags     string `csv:"tags"`
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	config = &configuration{}

	flag.StringVar(&config.InputFilename, "input", "", "Input CSV file in DKB format (Giro or Creditcard")
	flag.StringVar(&config.OutputFilename, "output", "", "Output CSV file in Homebank format")
}

func main() {
	flag.Parse()

	if config.InputFilename == "" || config.OutputFilename == "" {
		log.Errorln("Input and output file must be given")
		flag.Usage()
		os.Exit(1)
	}

	inputfile, err := os.Open(config.InputFilename)
	if err != nil {
		log.Fatalln(err)
	}
	defer inputfile.Close()

	outputfile, err := os.OpenFile(config.OutputFilename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalln(err)
	}
	defer outputfile.Close()

	convertFile(inputfile, outputfile)
}

func detectFiletype(inputfile *os.File) string {
	reader := csv.NewReader(transform.NewReader(inputfile, charmap.ISO8859_15.NewDecoder()))
	reader.Comma = ';'
	record, err := reader.Read()
	if err == io.EOF || err != nil {
		log.Fatal(err)
	}
	if record[0] == "Kreditkarte:" {
		return "CREDIT_CSV"
	}
	if record[0] == "Kontonummer:" {
		return "GIRO_CSV"
	}

	return "FILETYPE_UNKNOWN"
}

func convertFile(inputfile *os.File, outputfile *os.File) {
	DkbRecords := []*dkbGiroCsv{}
	HomebankRecords := []*homebankCsv{}

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		reader := csv.NewReader(transform.NewReader(in, charmap.ISO8859_15.NewDecoder()))
		reader.TrimLeadingSpace = true
		reader.Comma = ';'
		reader.FieldsPerRecord = -1
		seekToStart(reader)
		return reader
	})

	if err := gocsv.UnmarshalFile(inputfile, &DkbRecords); err != nil {
		log.Fatal(err)
	}

	for _, record := range DkbRecords {
		NewRecord := convertFromDkbGiro(record)
		HomebankRecords = append(HomebankRecords, &NewRecord)
	}

	// TODO: Homebank wants to have values enclosed in quotes
	gocsv.SetCSVWriter(func(out io.Writer) *gocsv.SafeCSVWriter {
		writer := csv.NewWriter(out)
		writer.Comma = ';'
		return gocsv.NewSafeCSVWriter(writer)
	})

	if err := gocsv.MarshalFile(&HomebankRecords, outputfile); err != nil {
		log.Fatal(err)
	}
}

func seekToStart(r *csv.Reader) {
	for i := 0; i < 4; i++ {
		_, _ = r.Read()
	}
}

func convertFromDkbGiro(DkbRecord *dkbGiroCsv) homebankCsv {

	paymentTypes := map[string]string{
		"abschluss":                 "0",
		"lohn, gehalt, rente":       "4",
		"online-ueberweisung":       "4",
		"überweisung":               "4",
		"rücküberweisung":           "4",
		"wertpapiere":               "4",
		"zins/dividende":            "4",
		"auftrag":                   "5",
		"umbuchung":                 "5",
		"kartenzahlung/-abrechnung": "6",
		"sepa-elv-lastschrift":      "6",
		"dauerauftrag":              "7",
		"gutschrift":                "8",
		"lastschrift":               "11",
		"folgelastschrift":          "11",
	}

	result := homebankCsv{}
	result.Date = DkbRecord.Wertstellung
	result.Payment = paymentTypes[strings.ToLower(DkbRecord.Buchungstext)]
	info := ""
	if DkbRecord.Kontonummer != "" {
		_, err := strconv.ParseFloat(DkbRecord.Kontonummer, 64)
		if err != nil {
			info = info + fmt.Sprintf("IBAN: %v, BIC: %v\n", DkbRecord.Kontonummer, DkbRecord.Blz)
		} else {
			info = info + fmt.Sprintf("Konto-Nr.: %v, BLZ: %v\n", DkbRecord.Kontonummer, DkbRecord.Blz)
		}
	}
	if DkbRecord.GlaeubigerID != "" {
		info = info + fmt.Sprintf("Gläubiger-ID: %v\n", DkbRecord.GlaeubigerID)
	}
	if DkbRecord.Mandatsreferenz != "" {
		info = info + fmt.Sprintf("Mandatsreferenz: %v\n", DkbRecord.Mandatsreferenz)
	}
	if DkbRecord.Kundenreferenz != "" {
		info = info + fmt.Sprintf("Kundenreferenz: %v\n", DkbRecord.Kundenreferenz)
	}
	result.Info = strings.TrimSpace(info)

	result.Payee = DkbRecord.AuftraggeberBeguenstigter
	result.Memo = DkbRecord.Verwendungszweck
	result.Amount = DkbRecord.BetragEur
	result.Category = ""
	result.Tags = ""

	return result
}

func convertFromDkbCredit(DkbRecord *dkbCreditCsv) *homebankCsv {
	if DkbRecord.UmsatzAbgerechnet == "Nein" || DkbRecord.UmsatzAbgerechnetOld == "Nein" {
		return nil
	}

	result := homebankCsv{}
	result.Date = DkbRecord.Wertstellung
	result.Payment = "1"
	result.Info = "Belegedatum: " + DkbRecord.Belegdatum
	if DkbRecord.Beschreibung != "" {
		result.Payee = DkbRecord.Beschreibung
	} else {
		result.Payee = DkbRecord.Umsatzbeschreibung
	}
	result.Memo = DkbRecord.UrspruenglicherBetrag
	result.Amount = DkbRecord.BetragEur
	result.Category = ""
	result.Tags = ""

	return &result
}
