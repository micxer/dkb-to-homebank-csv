package main

import (
  "flag"
  "fmt"
  "log"
  "os"
  "io"
  "encoding/csv"
  "github.com/gocarina/gocsv"
)

type Config struct {
  InputFilename string
  OutputFilename string
}

var config *Config

type DkbCsv struct {
  Buchungstag               string `csv:"Buchungstag"`
  Wertstellung              string `csv:"Wertstellung"`
  Buchungstext              string `csv:"Buchungstext"`
  AuftraggeberBeguenstigter string `csv:"Auftraggeber / Begünstigter"`
  Verwendungszweck          string `csv:"Verwendungszweck"`
  Kontonummer               string `csv:"Kontonummer"`
  Blz                       string `csv:"BLZ"`
  BetragEur                 string `csv:"Betrag (EUR)"`
  GlaeubigerId              string `csv:"Gläubiger-ID"`
  Mandatsreferenz           string `csv:"Mandatsreferenz"`
  Kundenreferenz            string `csv:"Kundenreferenz"`
}

type HomebankCsv struct {
  Date      string `csv:"date"`
  Payment   string `csv:"payment"`
  Info      string `csv:"info"`
  Payee     string `csv:"payee"`
  Memo      string `csv:"memo"`
  Amount    string `csv:"amount"`
  Category  string `csv:"category"`
  Tags      string `csv:"tags"`
}

func init() {
  config = &Config{}

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

  convert_file(inputfile, outputfile)
}

func convert_file(inputfile *os.File, outputfile *os.File) {
  DkbRecords := []*DkbCsv{}
  HomebankRecords := []*HomebankCsv{}

  gocsv.SetCSVReader(func(in io.Reader) *csv.Reader {
    read := gocsv.DefaultCSVReader(in)
    read.TrimLeadingSpace = true
    read.Comma = ';'
    return read
  })

  if err := gocsv.UnmarshalFile(inputfile, &DkbRecords); err != nil {
    log.Fatal(err)
  }

  for _, record := range DkbRecords {
    NewRecord := ConvertFromDkb(record)
    HomebankRecords = append(HomebankRecords, &NewRecord)
  }

  if err := gocsv.MarshalFile(&HomebankRecords, outputfile); err != nil {
    log.Fatal(err)
  }
}

func ConvertFromDkb(DkbRecord *DkbCsv) HomebankCsv {
  result := HomebankCsv{}
  result.Amount = DkbRecord.BetragEur;

  return result
}
