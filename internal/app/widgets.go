package app

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/luisnquin/nix-search/internal/nix"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgets/text"
	"github.com/mum4k/termdash/widgets/textinput"
	"github.com/samber/lo"
)

func (app *App) initWidgets() error {
	var err error

	app.elements.resultsBoard, err = app.getResultsBoard()
	if err != nil {
		return fmt.Errorf("results board: %w", err)
	}

	app.elements.searchInput, err = app.getSearchTextInput()
	if err != nil {
		return fmt.Errorf("search input: %w", err)
	}

	app.elements.currentStatus, err = app.getCurrentStatusWidget()
	if err != nil {
		return fmt.Errorf("current status widget: %w", err)
	}

	app.elements.currentLabel, err = app.getCurrentLabelWidget()
	if err != nil {
		return fmt.Errorf("current label widget: %w", err)
	}

	app.elements.searchOptions, err = app.getSearchOptionsWidget()
	if err != nil {
		return fmt.Errorf("search options widget: %w", err)
	}

	app.elements.currentSource, err = app.getCurrentSourceWidget()
	if err != nil {
		return fmt.Errorf("current source widget: %w", err)
	}

	return nil
}

func (app App) updateWidgetTexts() error {
	app.elements.resultsBoard.Reset()

	err := app.elements.currentLabel.Write(app.currentSearchTab.Label, text.WriteReplace())
	if err != nil {
		return err
	}

	err = app.elements.currentSource.Write(app.currentSearchTab.Source, text.WriteReplace())
	if err != nil {
		return err
	}

	return app.updateCurrentStatus(app.currentSearchTab.Status)
}

func (app App) updateCurrentStatus(newStatus string) error {
	return app.currentStatus.Write(newStatus, text.WriteReplace())
}

func (a *App) getSearchTextInput() (*textinput.TextInput, error) {
	ctx := context.Background()

	return textinput.New(
		textinput.Label("Search packages/options: ", cell.FgColor(cell.ColorAqua)),
		textinput.Border(linestyle.None),
		textinput.PlaceHolder("enter any text"),
		textinput.FillColor(cell.ColorDefault),
		textinput.ExclusiveKeyboardOnFocus(),
		textinput.OnSubmit(func(text string) error {
			options, err := a.nixClient.SearchHomeManagerOptions(ctx, text)
			if err != nil {
				uerr, ok := err.(*url.Error)
				if ok && uerr.Timeout() {
					return nil
				}

				return err // TODO: send to terminal screen and do not display context cancelled error
			}

			results := strings.Join(lo.Map(options, func(opt *nix.HomeManagerOption, _ int) string {
				return opt.String()
			}), " ")

			a.resultsBoard.Reset()
			a.resultsBoard.Write(results)

			return nil
		}))
}

func (a App) getResultsBoard() (*text.Text, error) {
	return text.New(text.WrapAtWords())
}

func (a App) getCurrentLabelWidget() (*text.Text, error) {
	return a.newTextWidget(a.currentSearchTab.Label)
}

func (a App) getCurrentStatusWidget() (*text.Text, error) {
	return a.newTextWidget(a.currentSearchTab.Status, text.WriteCellOpts(cell.Bold()))
}

func (a App) getCurrentSourceWidget() (*text.Text, error) {
	return a.newTextWidget(a.currentSearchTab.Source, text.WriteCellOpts(cell.Bold()))
}

func (a App) getSearchOptionsWidget() (*text.Text, error) {
	tabs := lo.Map(a.getSearchTabs(), func(tab searchTabConfig, _ int) string {
		return tab.Label
	})

	return a.newTextWidget(strings.Join(tabs, " | "))
}

func (a App) newTextWidget(content string, tOpts ...text.WriteOption) (*text.Text, error) {
	t, err := text.New(text.DisableScrolling())
	if err != nil {
		return nil, err
	}

	return t, t.Write(content, tOpts...)
}
