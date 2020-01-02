/**
 * Copyright (c) 2017-2020, Michael Weinrich
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
	AuftraggeberBeguenstigter string
	BetragEur                 string
	Blz                       string
	Buchungstag               string
	Buchungstext              string
	GlaeubigerID              string
	Kontonummer               string
	Kundenreferenz            string
	Mandatsreferenz           string
	Verwendungszweck          string
	Wertstellung              string
}

type dkbCreditCsv struct {
	Belegdatum            string
	Beschreibung          string
	BetragEur             string
	UmsatzAbgerechnet     string
	UrspruenglicherBetrag string
	Wertstellung          string
}

type homebankCsv struct {
	Date     string
	Payment  string
	Info     string
	Payee    string
	Memo     string
	Amount   string
	Category string
	Tags     string
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		QuoteEmptyFields:       true,
	})

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
	homebankRecords := convertFile(inputfile)
	writeHomebankFile(homebankRecords, outputfile)
}

func convertFile(inputfile *os.File) []*homebankCsv {

	fileType := detectFiletype(inputfile)
	homebankRecords := []*homebankCsv{}

	if fileType == "GIRO_CSV" {
		homebankRecords = processGiroFile(inputfile)
	}
	if fileType == "CREDIT_CSV" {
		homebankRecords = processCreditFile(inputfile)
	}

	return homebankRecords
}

func (h homebankCsv) getRecord() []string {
	record := make([]string, 8)

	record[0] = h.Date
	record[1] = h.Payment
	record[2] = h.Info
	record[3] = h.Payee
	record[4] = h.Memo
	record[5] = h.Amount
	record[6] = h.Category
	record[7] = h.Tags

	return record
}

func writeHomebankFile(homebankRecords []*homebankCsv, outputfile *os.File) {
	outputfile.WriteString("\"" + strings.Join([]string{"date", "paymode", "info", "payee", "memo", "amount", "category", "tags"}, "\";\"") + "\"\n")

	for _, homebankRecord := range homebankRecords {
		csvRecord := homebankRecord.getRecord()
		if _, err := outputfile.WriteString("\"" + strings.Join(csvRecord, "\";\"") + "\"\n"); err != nil {
			log.Fatal("error writing record to csv:", err)
		}
	}
}

func detectFiletype(inputfile *os.File) string {
	reader := csv.NewReader(transform.NewReader(inputfile, charmap.ISO8859_15.NewDecoder()))
	reader.Comma = ';'
	record, err := reader.Read()
	if err == io.EOF || err != nil {
		log.Fatal(err)
	}
	inputfile.Seek(0, 0)
	if record[0] == "Kreditkarte:" {
		log.Info("Detected file format: CREDIT_CSV")
		return "CREDIT_CSV"
	}
	if record[0] == "Kontonummer:" {
		log.Info("Detected file format: GIRO_CSV")
		return "GIRO_CSV"
	}

	log.Warning("Unknown filetype")
	return "FILETYPE_UNKNOWN"
}

func processGiroFile(inputfile *os.File) []*homebankCsv {
	homebankRecords := []*homebankCsv{}

	dkbRecords := readGiroFile(inputfile)

	for _, record := range dkbRecords {
		newRecord := convertFromDkbGiro(record)
		log.Debugf("processGiroFile(): %v", newRecord)
		homebankRecords = append(homebankRecords, &newRecord)
	}

	return homebankRecords
}

func processCreditFile(inputfile *os.File) []*homebankCsv {
	homebankRecords := []*homebankCsv{}

	dkbRecords := readCreditFile(inputfile)

	for _, record := range dkbRecords {
		newRecord := convertFromDkbCredit(record)
		log.Debugf("processCreditFile(): %v", newRecord)
		if newRecord.Date != "" {
			homebankRecords = append(homebankRecords, &newRecord)
		}
	}

	return homebankRecords
}

func readGiroFile(input io.Reader) []*dkbGiroCsv {
	dkbRecords := []*dkbGiroCsv{}

	firstLineFound := false

	reader := getCsvReader(input)
	for {
		csvRecord, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if !firstLineFound {
			if csvRecord[0] == "Buchungstag" {
				firstLineFound = true
			}
			continue
		}

		log.Debugf("readGiroFile(): %v", csvRecord)

		dkbRecords = append(dkbRecords, mapGiroCsvToStruct(csvRecord))
	}

	return dkbRecords
}

func readCreditFile(input io.Reader) []*dkbCreditCsv {
	dkbRecords := []*dkbCreditCsv{}

	firstLineFound := false

	reader := getCsvReader(input)
	for {
		csvRecord, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if !firstLineFound {
			if strings.HasPrefix(csvRecord[0], "Umsatz abgerechnet") {
				firstLineFound = true
			}
			continue
		}

		log.Debugf("readCreditFile(): %v", csvRecord)

		dkbRecords = append(dkbRecords, mapCreditCsvToStruct(csvRecord))
	}

	return dkbRecords
}

func mapGiroCsvToStruct(csvRecord []string) *dkbGiroCsv {
	dkbRecord := dkbGiroCsv{}

	dkbRecord.Buchungstag = csvRecord[0]
	dkbRecord.Wertstellung = csvRecord[1]
	dkbRecord.Buchungstext = csvRecord[2]
	dkbRecord.AuftraggeberBeguenstigter = csvRecord[3]
	dkbRecord.Verwendungszweck = csvRecord[4]
	dkbRecord.Kontonummer = csvRecord[5]
	dkbRecord.Blz = csvRecord[6]
	dkbRecord.BetragEur = csvRecord[7]
	dkbRecord.GlaeubigerID = csvRecord[8]
	dkbRecord.Mandatsreferenz = csvRecord[9]
	dkbRecord.Kundenreferenz = csvRecord[10]

	return &dkbRecord
}

func mapCreditCsvToStruct(csvRecord []string) *dkbCreditCsv {
	dkbRecord := dkbCreditCsv{}

	dkbRecord.UmsatzAbgerechnet = csvRecord[0]
	dkbRecord.Wertstellung = csvRecord[1]
	dkbRecord.Belegdatum = csvRecord[2]
	dkbRecord.Beschreibung = csvRecord[3]
	dkbRecord.BetragEur = csvRecord[4]
	dkbRecord.UrspruenglicherBetrag = csvRecord[5]

	return &dkbRecord
}

func getCsvReader(input io.Reader) *csv.Reader {
	reader := csv.NewReader(transform.NewReader(input, charmap.ISO8859_15.NewDecoder()))
	reader.TrimLeadingSpace = true
	reader.Comma = ';'
	reader.FieldsPerRecord = -1

	return reader
}

func convertFromDkbGiro(dkbRecord *dkbGiroCsv) homebankCsv {

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
	result.Date = dkbRecord.Wertstellung
	result.Payment = paymentTypes[strings.ToLower(dkbRecord.Buchungstext)]
	info := ""
	if dkbRecord.Kontonummer != "" {
		_, err := strconv.ParseFloat(dkbRecord.Kontonummer, 64)
		if err != nil {
			info = info + fmt.Sprintf("IBAN: %v, BIC: %v,", dkbRecord.Kontonummer, dkbRecord.Blz)
		} else {
			info = info + fmt.Sprintf("Konto-Nr.: %v, BLZ: %v,", dkbRecord.Kontonummer, dkbRecord.Blz)
		}
	}
	if dkbRecord.GlaeubigerID != "" {
		info = info + fmt.Sprintf("Gläubiger-ID: %v, ", dkbRecord.GlaeubigerID)
	}
	if dkbRecord.Mandatsreferenz != "" {
		info = info + fmt.Sprintf("Mandatsreferenz: %v, ", dkbRecord.Mandatsreferenz)
	}
	if dkbRecord.Kundenreferenz != "" {
		info = info + fmt.Sprintf("Kundenreferenz: %v, ", dkbRecord.Kundenreferenz)
	}
	result.Info = strings.TrimSuffix(strings.TrimSpace(info), ",")

	result.Payee = dkbRecord.AuftraggeberBeguenstigter
	result.Memo = dkbRecord.Verwendungszweck
	result.Amount = dkbRecord.BetragEur
	result.Category = ""
	result.Tags = ""

	return result
}

func convertFromDkbCredit(DkbRecord *dkbCreditCsv) homebankCsv {
	result := homebankCsv{}

	result.Date = DkbRecord.Wertstellung
	result.Payment = "1"
	result.Info = "Belegedatum: " + DkbRecord.Belegdatum
	result.Payee = DkbRecord.Beschreibung
	result.Memo = DkbRecord.UrspruenglicherBetrag
	result.Amount = DkbRecord.BetragEur
	result.Category = ""
	result.Tags = ""

	return result
}
