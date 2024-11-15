package resources

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-shiori/dom"
	"golang.org/x/net/html"
)

func GetScalar(url string, querySelector string) (string, error) {
	resp, err := http.Get(url)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)

	if err != nil {
		return "", errors.New("couldn't parse body")
	}

	el := dom.QuerySelector(doc, querySelector)

	if el == nil {
		return "", errors.New("couldn't match selector")
	}

	value := strings.TrimSpace(dom.TextContent(el))

	return value, nil
}

func GetVector(url string, querySelectorAll string) ([]string, error) {
	var values []string

	resp, err := http.Get(url)

	if err != nil {
		return values, err
	}

	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)

	if err != nil {
		return values, errors.New("couldn't parse body")
	}

	els := dom.QuerySelectorAll(doc, querySelectorAll)

	for _, el := range els {
		values = append(values, dom.InnerText(el))
	}

	return values, nil
}
