// Copyright (c) 2019-2022 The iexcloud developers. All rights reserved.
// Project site: https://github.com/woodstock-tokyo/iexcloud
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE file for the project.

//go:build integration
// +build integration

package iex_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sort"
	"testing"

	"github.com/BurntSushi/toml"
	iex "github.com/woodstock-tokyo/iexcloud"
)

// Config contains the configuration information neecded to program and test
// the adapaters.
type Config struct {
	Token   string
	Secret  string
	BaseURL string
}

// ReadConfig will read the TOML config file.
func readConfig(configFile string) (Config, error) {

	var cfg Config

	// Read config file
	f, err := os.Open(configFile)
	if err != nil {
		return cfg, err
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return cfg, err
	}
	err = toml.Unmarshal(buf, &cfg)
	return cfg, err
}

func TestIntegrationAnnualBalanceSheets(t *testing.T) {
	cfg, err := readConfig("config_test.toml")
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	client := iex.NewClient(cfg.Token, iex.WithBaseURL(cfg.BaseURL))
	bs, err := client.AnnualBalanceSheets(context.Background(), "aapl", 4)
	if err != nil {
		log.Fatalf("Error getting balance sheets: %s", err)
	}
	assertString(t, "symbol", bs.Symbol, "AAPL")
	assertInt(t, "number of years", len(bs.Statements), 4)
	q1 := bs.Statements[0]
	assertString(t, "filing type", q1.FilingType, "10-K")
	assertInt(t, "fiscal quarter", q1.FiscalQuarter, 0)
	isPositiveInt(t, "fiscal year", q1.FiscalYear)
	assertString(t, "currency", q1.Currency, "USD")
}

func TestIntegrationQuarterlyBalanceSheets(t *testing.T) {
	cfg, err := readConfig("config_test.toml")
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	client := iex.NewClient(cfg.Token, iex.WithBaseURL(cfg.BaseURL))
	bs, err := client.QuarterlyBalanceSheets(context.Background(), "aapl", 2)
	if err != nil {
		log.Fatalf("Error getting balance sheets: %s", err)
	}
	assertString(t, "symbol", bs.Symbol, "AAPL")
	assertInt(t, "number of quarters", len(bs.Statements), 2)
	q1 := bs.Statements[0]
	assertString(t, "filing type", q1.FilingType, "10-K")
	isPositiveInt(t, "fiscal quarter", q1.FiscalQuarter)
	isPositiveInt(t, "fiscal year", q1.FiscalYear)
	assertString(t, "currency", q1.Currency, "USD")
}

func TestIntegrationBook(t *testing.T) {
	cfg, err := readConfig("config_test.toml")
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	client := iex.NewClient(cfg.Token, iex.WithBaseURL(cfg.BaseURL))
	got, err := client.Book(context.Background(), "aapl")
	if err != nil {
		log.Fatalf("Error getting book: %s", err)
	}
	assertString(t, "symbol", got.Quote.Symbol, "AAPL")
	assertString(t, "company name", got.Quote.CompanyName, "Apple Inc")
	assertScrambledString(t, "primary exchange", got.Quote.PrimaryExchange, "NASDAQ")
	assertString(t, "calculation price", got.Quote.CalculationPrice, "close")
	isPositiveFloat64(t, "open", got.Quote.Open)
	assertScrambledString(t, "open source", got.Quote.OpenSource, "official")
	isPositiveFloat64(t, "latest price", got.Quote.LatestPrice)
}

func TestIntegrationHistoricalPrices(t *testing.T) {
	cfg, err := readConfig("config_test.toml")
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	client := iex.NewClient(cfg.Token, iex.WithBaseURL(cfg.BaseURL))
	timeframe := iex.OneMonthHistorical
	histPrices, err := client.HistoricalPrices(context.Background(), "aapl", timeframe, nil)
	if err != nil {
		log.Fatalf("Error getting historical prices: %s", err)
	}
	got := histPrices[0]
	isPositiveFloat64(t, "close", got.Close)
	isPositiveFloat64(t, "high", got.High)
	isPositiveFloat64(t, "low", got.Low)
	isPositiveFloat64(t, "open", got.Open)
	assertString(t, "symbol", got.Symbol, "AAPL")
	isPositiveFloat64(t, "volume", got.Volume)
	assertScrambledString(t, "id", got.ID, "HISTORICAL_PRICES")
	assertScrambledString(t, "key", got.Key, "AAPL")
	assertString(t, "subkey", got.Subkey, "")
}

func TestIntegrationCreatRule(t *testing.T) {
	cfg, err := readConfig("config_test.toml")
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	client := iex.NewClient(cfg.Token, iex.WithBaseURL(cfg.BaseURL))

	condtion1 := iex.Condition{"changePercent", ">", 1}
	rule := iex.Rule{
		Token:      cfg.Token,
		RuleSet:    "AAPL",
		Type:       "any",
		RuleName:   "TestCreateByApiNew1",
		Conditions: []iex.Condition{condtion1},
		Outputs: []iex.Output{
			{
				Frequency: 3600,
				Method:    "webhook",
				Url:       "https://test.ngrok.io",
			},
		},
	}

	result, err := client.CreateRulePriceAlert(context.Background(), rule)

	assertInt(t, "weight", int(result.Weight), 1)

	if err != nil {
		log.Fatalf("Error getting historical prices: %s", err)
	}

}

// it seems eix edit rule doesnt work, contact to support

func TestIntegrationEditRule(t *testing.T) {
	cfg, err := readConfig("config_test.toml")
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	client := iex.NewClient(cfg.Token, iex.WithBaseURL(cfg.BaseURL))

	condtion1 := iex.Condition{"changePercent", ">", 1}
	rule := iex.Rule{
		Token:      cfg.Token,
		ID:         "6f6fc56e-516b-4568-b423-a66febca32a5",
		RuleSet:    "AAPL",
		Type:       "any",
		RuleName:   "TestCreateByApiEditRule",
		Conditions: []iex.Condition{condtion1},
		Outputs: []iex.Output{
			{
				Frequency: 3600,
				Method:    "webhook",
				Url:       "https://test.ngrok.io",
			},
		},
	}

	result, err := client.CreateRulePriceAlert(context.Background(), rule)

	assertInt(t, "weight", int(result.Weight), 0)

	if err != nil {
		log.Fatalf("Error getting historical prices: %s", err)
	}

}

func TestIntegrationDeleteRule(t *testing.T) {
	cfg, err := readConfig("config_test.toml")
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	client := iex.NewClient(cfg.Secret, iex.WithBaseURL(cfg.BaseURL))

	result, err := client.DeleteRule(context.Background(), "7395c2cc-dd12-4872-a176-31748b7b0cbb")
	fmt.Println(result)
}

func TestIntegrationPauseRule(t *testing.T) {
	cfg, err := readConfig("config_test.toml")
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	client := iex.NewClient(cfg.Token, iex.WithBaseURL(cfg.BaseURL))

	result, err := client.PauseRule(context.Background(), iex.RequestRule{
		Token:  cfg.Secret,
		RuleID: "4fb997a3-8d5c-4eae-a433-6e8b341f7696",
	})
	fmt.Println(result)
}

func TestIntegrationResumeRule(t *testing.T) {
	cfg, err := readConfig("config_test.toml")
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	client := iex.NewClient(cfg.Token, iex.WithBaseURL(cfg.BaseURL))

	result, err := client.ResumeRule(context.Background(), iex.RequestRule{
		Token:  cfg.Secret,
		RuleID: "4fb997a3-8d5c-4eae-a433-6e8b341f7696",
	})
	fmt.Println(result)
}

func TestIntegrationGetRule(t *testing.T) {
	cfg, err := readConfig("config_test.toml")
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	client := iex.NewClient(cfg.Token, iex.WithBaseURL(cfg.BaseURL))

	result, err := client.GetRule(context.Background(), iex.RequestRule{
		Token:  cfg.Token,
		RuleID: "ef5d5d6d-175e-41e5-9779-3114ba8a0826",
	})
	fmt.Println(result)
}

func assertInt(t *testing.T, label string, got, want int) {
	if got != want {
		t.Errorf("\t got = %d %s\n\t\twant = %d", got, label, want)
	}
}

func assertFloat64(t *testing.T, label string, got, want, tolerance float64) {
	if diff := math.Abs(want - got); diff >= tolerance {
		t.Errorf("\t got = %f %s\n\t\t\twant = %f", got, label, want)
	}
}

func assertBool(t *testing.T, label string, got, want bool) {
	if got != want {
		t.Errorf("\t got = %t %s\n\t\t\twant = %t", got, label, want)
	}
}

func assertString(t *testing.T, label string, got, want string) {
	if got != want {
		t.Errorf("\t got = %s %s\n\t\t\twant = %s", got, label, want)
	}
}

// IEX scrambles their responses when using the testing sandbox. Therefore, the
// best we can do is assert that all the letters are there even if scrambled.
func assertScrambledString(t *testing.T, label string, got, want string) {
	gotSorted := sortString(got)
	wantSorted := sortString(want)
	if gotSorted != wantSorted {
		t.Errorf("\t got = %s %s\n\t\t\twant = %s", got, label, want)
	}
}

func isPositiveInt(t *testing.T, label string, got int) {
	if got <= 0 {
		t.Errorf("\t got = %d %s\n\t\twant int > 0", got, label)
	}
}

func isPositiveFloat64(t *testing.T, label string, got float64) {
	if got <= 0.0 {
		t.Errorf("\t got = %f %s\n\t\twant float64 > 0.0", got, label)
	}
}

func isString(t *testing.T, label string, got string) {
	if got == "" {
		t.Errorf("\t got = %s %s\n\t\twant non-empty string", got, label)
	}
}

type sortRunes []rune

func (s sortRunes) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s sortRunes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortRunes) Len() int {
	return len(s)
}

func sortString(s string) string {
	r := []rune(s)
	sort.Sort(sortRunes(r))
	return string(r)
}
