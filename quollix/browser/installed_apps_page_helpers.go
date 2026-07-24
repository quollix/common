package browser

import (
	"strings"

	"github.com/quollix/common/browsertest"
	u "github.com/quollix/common/utils"
)

type InstalledAppsPageHelpers struct {
	Page *browsertest.Page
}

func (p *InstalledAppsPageHelpers) ClickOpenButton(appName string) error {
	row, err := p.FindRowByAppName(appName)
	if err != nil {
		return err
	}

	accessCell, err := row.Element(".app-access")
	if err != nil {
		return u.Logger.NewError(err.Error())
	}

	openButton, err := accessCell.Element("button.open-btn")
	if err != nil {
		return u.Logger.NewError(err.Error())
	}

	if err := openButton.Click(); err != nil {
		return u.Logger.NewError(err.Error())
	}
	return nil
}

func (p *InstalledAppsPageHelpers) OpenAppInNewTab(appName string) (*browsertest.Page, error) {
	waitForNewTab := p.Page.WaitOpen()
	if err := p.ClickOpenButton(appName); err != nil {
		return nil, err
	}
	page, err := waitForNewTab()
	if err != nil {
		return nil, u.Logger.NewError(err.Error())
	}
	return page, nil
}

func (p *InstalledAppsPageHelpers) FindRowByAppName(appName string) (*browsertest.Element, error) {
	rows, err := p.Page.Elements("#installed-apps-tbody tr")
	if err != nil {
		return nil, u.Logger.NewError(err.Error())
	}
	for _, row := range rows {
		appNameAttr, err := row.Attribute("data-app-name")
		if err != nil {
			return nil, u.Logger.NewError(err.Error())
		}
		if appNameAttr != nil && strings.TrimSpace(*appNameAttr) == appName {
			return row, nil
		}
	}
	return nil, u.Logger.NewError("app row not found", "app_name", appName)
}
