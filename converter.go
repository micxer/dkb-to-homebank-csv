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
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocarina/gocsv"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type configuration struct {
	InputFilename  string
	OutputFilename string
}

var config *configuration

type dkbCsv struct {
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
	config = &configuration{}

	flag.StringVar(&config.InputFilename, "input", "", "Input CSV file in DKB format")
	flag.StringVar(&config.OutputFilename, "output", "", "Output CSV file in Homebank format")
}

func main() {
	flag.Parse()

	if config.InputFilename == "" || config.OutputFilename == "" {
		log.Fatalln("Input and output file must be given")
	}

	inputfile, err := os.Open(config.InputFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer inputfile.Close()

	outputfile, err := os.OpenFile(config.OutputFilename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer outputfile.Close()

	convertFile(inputfile, outputfile)
}

func convertFile(inputfile *os.File, outputfile *os.File) {
	DkbRecords := []*dkbCsv{}
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
		NewRecord := convertFromDkb(record)
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

func convertFromDkb(DkbRecord *dkbCsv) homebankCsv {

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
