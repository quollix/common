package browsertest

import (
	"strings"
	"time"

	u "github.com/quollix/common/utils"
)

func (p *Page) WaitElement(selector string) (*Element, error) {
	return p.WaitElementWithin(selector, 3*time.Second)
}

func (p *Page) WaitElementWithin(selector string, timeout time.Duration) (*Element, error) {
	var element *Element
	err := u.EventuallyWithTimeout(timeout, 50*time.Millisecond, func() error {
		foundElement, err := p.Element(selector)
		if err != nil {
			return err
		}
		element = foundElement
		return nil
	})
	if err != nil {
		return nil, u.Logger.AddContext(err, "selector", selector)
	}
	return element, nil
}

func (p *Page) WaitElementMatchingText(selector, textPattern string) (*Element, error) {
	return p.WaitElementMatchingTextWithin(selector, textPattern, 3*time.Second)
}

func (p *Page) WaitElementMatchingTextWithin(selector, textPattern string, timeout time.Duration) (*Element, error) {
	var element *Element
	err := u.EventuallyWithTimeout(timeout, 50*time.Millisecond, func() error {
		foundElement, err := p.ElementMatchingText(selector, textPattern)
		if err != nil {
			return err
		}
		element = foundElement
		return nil
	})
	if err != nil {
		return nil, u.Logger.AddContext(err, "selector", selector, "text", textPattern)
	}
	return element, nil
}

func (p *Page) ClickElement(selector string) error {
	return p.ClickElementWithin(selector, 3*time.Second)
}

func (p *Page) ClickElementWithin(selector string, timeout time.Duration) error {
	element, err := p.WaitElementWithin(selector, timeout)
	if err != nil {
		return err
	}
	if err := element.Click(); err != nil {
		return err
	}
	return nil
}

func (p *Page) ClickElementMatchingText(selector, textPattern string) error {
	return p.ClickElementMatchingTextWithin(selector, textPattern, 3*time.Second)
}

func (p *Page) ClickElementMatchingTextWithin(selector, textPattern string, timeout time.Duration) error {
	element, err := p.WaitElementMatchingTextWithin(selector, textPattern, timeout)
	if err != nil {
		return err
	}
	if err := element.Click(); err != nil {
		return err
	}
	return nil
}

func (p *Page) SetInputValue(selector, value string) error {
	input, err := p.WaitElement(selector)
	if err != nil {
		return err
	}
	if err := input.SelectAllText(); err != nil {
		return err
	}
	if err := input.Input(value); err != nil {
		return err
	}
	return nil
}

func (p *Page) TypeInputValue(selector, value string) error {
	input, err := p.WaitElement(selector)
	if err != nil {
		return err
	}
	if err := input.SelectAllText(); err != nil {
		return err
	}
	if err := input.TypeText(value); err != nil {
		return err
	}
	inputValue, err := input.Property("value")
	if err != nil {
		return err
	}
	if inputValue.String() != value {
		return u.Logger.NewError("input value mismatch", "selector", selector, "expected", value, "actual", inputValue.String())
	}
	return nil
}

func (p *Page) SetCheckboxChecked(selector string, checked bool) error {
	checkbox, err := p.WaitElement(selector)
	if err != nil {
		return err
	}
	current, err := checkbox.Property("checked")
	if err != nil {
		return err
	}
	if current.Bool() == checked {
		return nil
	}
	if err := checkbox.Click(); err != nil {
		return err
	}
	return nil
}

func (p *Page) BodyText() (string, error) {
	body, err := p.Element("body")
	if err != nil {
		return "", err
	}
	text, err := body.Text()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}
