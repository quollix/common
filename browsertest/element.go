package browsertest

import (
	"fmt"
	"strconv"

	"github.com/chromedp/chromedp"
	u "github.com/quollix/common/utils"
)

func (e *Element) Element(selector string) (*Element, error) {
	return findElement(e.page, e.expr, selector)
}

func (e *Element) MustElement(selector string) *Element {
	element, err := e.Element(selector)
	if err != nil {
		panic(err)
	}
	return element
}

func (e *Element) Elements(selector string) ([]*Element, error) {
	return findElements(e.page, e.expr, selector)
}

func (e *Element) Has(selector string) (bool, *Element, error) {
	return hasElement(e.page, e.expr, selector)
}

func (e *Element) Text() (string, error) {
	var text string
	err := e.evaluate(fmt.Sprintf(`(() => {
		const element = %s;
		if (!element) throw new Error("element not found");
		return element.innerText ?? element.textContent ?? "";
	})()`, e.expr), &text)
	return text, err
}

func (e *Element) MustText() string {
	text, err := e.Text()
	if err != nil {
		panic(err)
	}
	return text
}

func (e *Element) Attribute(name string) (*string, error) {
	var value *string
	err := e.evaluate(fmt.Sprintf(`(() => {
		const element = %s;
		if (!element) throw new Error("element not found");
		return element.getAttribute(%s);
	})()`, e.expr, strconv.Quote(name)), &value)
	return value, err
}

func (e *Element) MustAttribute(name string) *string {
	value, err := e.Attribute(name)
	if err != nil {
		panic(err)
	}
	return value
}

func (e *Element) Property(name string) (*PropertyValue, error) {
	var value any
	err := e.evaluate(fmt.Sprintf(`(() => {
		const element = %s;
		if (!element) throw new Error("element not found");
		return element[%s];
	})()`, e.expr, strconv.Quote(name)), &value)
	return &PropertyValue{value: value}, err
}

func (e *Element) Click() error {
	return e.evaluate(fmt.Sprintf(`(() => {
		const element = %s;
		if (!element) throw new Error("element not found");
		element.scrollIntoView({block: "center", inline: "center"});
		element.click();
	})()`, e.expr), nil)
}

func (e *Element) MustClick() *Element {
	if err := e.Click(); err != nil {
		panic(err)
	}
	return e
}

func (e *Element) SelectAllText() error {
	return e.evaluate(fmt.Sprintf(`(() => {
		const element = %s;
		if (!element) throw new Error("element not found");
		element.focus();
		if (typeof element.select === "function") {
			element.select();
		}
	})()`, e.expr), nil)
}

func (e *Element) MustSelectAllText() *Element {
	if err := e.SelectAllText(); err != nil {
		panic(err)
	}
	return e
}

func (e *Element) Input(value string) error {
	return e.evaluate(fmt.Sprintf(`(() => {
		const element = %s;
		if (!element) throw new Error("element not found");
		element.focus();
		element.value = %s;
		element.dispatchEvent(new Event("input", {bubbles: true}));
		element.dispatchEvent(new Event("change", {bubbles: true}));
	})()`, e.expr, strconv.Quote(value)), nil)
}

func (e *Element) MustInput(value string) *Element {
	if err := e.Input(value); err != nil {
		panic(err)
	}
	return e
}

func (e *Element) Select(label string) error {
	return e.evaluate(fmt.Sprintf(`(() => {
		const element = %s;
		if (!element) throw new Error("element not found");
		const option = Array.from(element.options).find((candidate) =>
			candidate.textContent.trim() === %s || candidate.value === %s
		);
		if (!option) throw new Error("select option not found");
		element.value = option.value;
		element.dispatchEvent(new Event("input", {bubbles: true}));
		element.dispatchEvent(new Event("change", {bubbles: true}));
	})()`, e.expr, strconv.Quote(label), strconv.Quote(label)), nil)
}

func (e *Element) MustSelect(label string) *Element {
	if err := e.Select(label); err != nil {
		panic(err)
	}
	return e
}

func (e *Element) TypeText(value string) error {
	if err := e.evaluate(fmt.Sprintf(`(() => {
		const element = %s;
		if (!element) throw new Error("element not found");
		element.focus();
	})()`, e.expr), nil); err != nil {
		return err
	}
	if err := chromedp.Run(e.page.ctx, chromedp.KeyEvent(value)); err != nil {
		return u.Logger.NewError(err.Error())
	}
	return nil
}

func (e *Element) ensureExists() error {
	exists, err := e.exists()
	if err != nil {
		return err
	}
	if !exists {
		return u.Logger.NewError("element not found")
	}
	return nil
}

func (e *Element) exists() (bool, error) {
	var exists bool
	err := e.evaluate(fmt.Sprintf(`(() => %s !== null && %s !== undefined)()`, e.expr, e.expr), &exists)
	return exists, err
}

func (e *Element) evaluate(expression string, result any) error {
	if err := chromedp.Run(e.page.ctx, chromedp.EvaluateAsDevTools(expression, result)); err != nil {
		return u.Logger.NewError(err.Error())
	}
	return nil
}

func (v *PropertyValue) String() string {
	if v == nil || v.value == nil {
		return ""
	}
	switch value := v.value.(type) {
	case string:
		return value
	default:
		return fmt.Sprint(value)
	}
}

func (v *PropertyValue) Bool() bool {
	if v == nil || v.value == nil {
		return false
	}
	switch value := v.value.(type) {
	case bool:
		return value
	case string:
		return value == "true"
	default:
		return false
	}
}
