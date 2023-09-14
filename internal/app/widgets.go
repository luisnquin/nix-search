package app

import (
	"context"
	"fmt"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgets/text"
	"github.com/mum4k/termdash/widgets/textinput"
)

func (app *App) initWidgets() error {
	var err error

	app.widgets.resultsBoard, err = app.getResultsBoard()
	if err != nil {
		return fmt.Errorf("results board: %w", err)
	}

	app.widgets.searchInput, err = app.getSearchTextInput()
	if err != nil {
		return fmt.Errorf("search input: %w", err)
	}

	app.widgets.currentStatus, err = app.getCurrentStatusWidget()
	if err != nil {
		return fmt.Errorf("current status widget: %w", err)
	}

	app.widgets.currentLabel, err = app.getCurrentLabelWidget()
	if err != nil {
		return fmt.Errorf("current label widget: %w", err)
	}

	app.widgets.currentSource, err = app.getCurrentSourceWidget()
	if err != nil {
		return fmt.Errorf("current source widget: %w", err)
	}

	app.widgets.currentChannelId, err = app.getCurrentChannelWidget()
	if err != nil {
		return fmt.Errorf("current channel id widget: %w", err)
	}

	return nil
}

func (app *App) updateWidgetTexts() error {
	app.widgets.resultsBoard.Reset()

	err := app.widgets.currentLabel.Write(app.tabs.search.Label, text.WriteReplace())
	if err != nil {
		return err
	}

	err = app.widgets.currentSource.Write(app.tabs.search.Source, text.WriteReplace())
	if err != nil {
		return err
	}

	if err := app.updateCurrentChannelID(); err != nil {
		return err
	}

	return app.updateCurrentStatus(app.tabs.search.Status)
}

func (app *App) updateCurrentChannelID() error {
	channelId := app.tabs.search.CurrentChannelID
	if channelId == "" {
		channelId = "No channel"
	}

	app.widgets.currentChannelId.Reset()

	return app.widgets.currentChannelId.Write(channelId)
}

func (app *App) updateCurrentStatus(newStatus string) error {
	return app.widgets.currentStatus.Write(newStatus, text.WriteReplace())
}

func (app *App) getSearchTextInput() (*textinput.TextInput, error) {
	return textinput.New(
		textinput.Label("Search packages/options: ", cell.FgColor(cell.ColorAqua)),
		textinput.Border(linestyle.None),
		textinput.PlaceHolder("enter any text"),
		textinput.FillColor(cell.ColorDefault),
		textinput.ExclusiveKeyboardOnFocus(),
		textinput.OnChange(app.handleSearchInputChange),
		textinput.OnSubmit(app.handleSearchInputSubmit))
}

func (app *App) handleSearchInputChange(input string) {
	if app.tabs.search.WaitForEnter {
		return
	}

	app.performSearch(context.Background(), input)
}

func (app *App) handleSearchInputSubmit(input string) error {
	app.performSearch(context.Background(), input)

	return nil
}

func (a App) getResultsBoard() (*text.Text, error) {
	return text.New(text.WrapAtWords())
}

func (app App) getCurrentLabelWidget() (*text.Text, error) {
	return app.newTextWidget(app.tabs.search.Label)
}

func (app App) getCurrentStatusWidget() (*text.Text, error) {
	return app.newTextWidget(app.tabs.search.Status, text.WriteCellOpts(cell.Bold()))
}

func (app App) getCurrentSourceWidget() (*text.Text, error) {
	return app.newTextWidget(app.tabs.search.Source, text.WriteCellOpts(cell.Bold()))
}

func (app App) getCurrentChannelWidget() (*text.Text, error) {
	return app.newTextWidget(app.tabs.search.CurrentChannelID, text.WriteCellOpts(cell.Bold()))
}

func (a *App) newTextWidget(content string, tOpts ...text.WriteOption) (*text.Text, error) {
	t, err := text.New(text.DisableScrolling())
	if err != nil {
		return nil, err
	}

	return t, t.Write(content, tOpts...)
}
