package boc

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-shiori/dom"
	"golang.org/x/net/html"
)

func GetPrime() (float32, error) {
	resp, err := http.Get("https://www.bochk.com/whk/rates/hkDollarPrimeRate/hkDollarPrimeRate-enquiry.action?lang=en")

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)

	if err != nil {
		return 0, errors.New("couldn't parse body")
	}

	td := dom.QuerySelector(doc, ".best-rate td:nth-child(2)")

	prime, err := strconv.ParseFloat(dom.InnerText(td.FirstChild), 32)

	if err != nil {
		return 0, errors.New("couldn't convert prime to float")
	}

	return float32(prime), nil
}