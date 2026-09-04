package browsertest

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	u "github.com/quollix/common/utils"
)

type Browser struct {
	ctx    context.Context
	cancel context.CancelFunc
}

type Page struct {
	browser *Browser
	ctx     context.Context
	cancel  context.CancelFunc
}

type PageInfo struct {
	URL string
}

type Element struct {
	page *Page
	expr string
}

type PropertyValue struct {
	value any
}

func queryExpression(parentExpression, selector string) string {
	return fmt.Sprintf(`(%s).querySelector(%s)`, parentExpression, strconv.Quote(selector))
}

func queryAllExpression(parentExpression, selector string) string {
	return fmt.Sprintf(`(%s).querySelectorAll(%s)`, parentExpression, strconv.Quote(selector))
}

func indexedQueryExpression(parentExpression, selector string, index int) string {
	return fmt.Sprintf(`(%s)[%d]`, queryAllExpression(parentExpression, selector), index)
}

func findElement(page *Page, parentExpression, selector string) (*Element, error) {
	element := &Element{page: page, expr: queryExpression(parentExpression, selector)}
	if err := element.ensureExists(); err != nil {
		return nil, err
	}
	return element, nil
}

func findElements(page *Page, parentExpression, selector string) ([]*Element, error) {
	count, err := page.count(queryAllExpression(parentExpression, selector))
	if err != nil {
		return nil, err
	}

	elements := make([]*Element, 0, count)
	for index := 0; index < count; index++ {
		elements = append(elements, &Element{
			page: page,
			expr: indexedQueryExpression(parentExpression, selector, index),
		})
	}
	return elements, nil
}

func findElementMatchingText(page *Page, selector string, textPattern string) (*Element, error) {
	pattern, err := regexp.Compile(textPattern)
	if err != nil {
		return nil, u.Logger.NewError(err.Error(), "text", textPattern)
	}

	elements, err := findElements(page, "document", selector)
	if err != nil {
		return nil, err
	}
	for _, element := range elements {
		text, err := element.Text()
		if err != nil {
			return nil, err
		}
		if pattern.MatchString(text) {
			return element, nil
		}
	}
	return nil, u.Logger.NewError("element not found")
}

func hasElement(page *Page, parentExpression, selector string) (bool, *Element, error) {
	element := &Element{page: page, expr: queryExpression(parentExpression, selector)}
	exists, err := element.exists()
	if err != nil {
		return false, nil, err
	}
	if !exists {
		return false, nil, nil
	}
	return true, element, nil
}

func cookiePath(cookie *http.Cookie) string {
	if strings.TrimSpace(cookie.Path) == "" {
		return "/"
	}
	return cookie.Path
}
