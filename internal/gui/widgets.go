package gui

import (
	"context"
	"fmt"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgets/text"
	"github.com/mum4k/termdash/widgets/textinput"
)

func (g *GUI) initWidgets() error {
	var err error

	g.widgets.resultsBoard, err = g.getResultsBoard()
	if err != nil {
		return fmt.Errorf("results board: %w", err)
	}

	g.widgets.searchInput, err = g.getSearchTextInput()
	if err != nil {
		return fmt.Errorf("search input: %w", err)
	}

	g.widgets.currentStatus, err = g.getCurrentStatusWidget()
	if err != nil {
		return fmt.Errorf("current status widget: %w", err)
	}

	g.widgets.currentLabel, err = g.getCurrentLabelWidget()
	if err != nil {
		return fmt.Errorf("current label widget: %w", err)
	}

	g.widgets.currentSource, err = g.getCurrentSourceWidget()
	if err != nil {
		return fmt.Errorf("current source widget: %w", err)
	}

	g.widgets.currentChannelId, err = g.getCurrentChannelWidget()
	if err != nil {
		return fmt.Errorf("current channel id widget: %w", err)
	}

	return nil
}

func (g *GUI) updateWidgetTexts() error {
	g.widgets.resultsBoard.Reset()

	err := g.widgets.currentLabel.Write(g.tabs.search.Label, text.WriteReplace())
	if err != nil {
		return err
	}

	err = g.widgets.currentSource.Write(g.tabs.search.Source, text.WriteReplace())
	if err != nil {
		return err
	}

	return g.updateCurrentChannelID()
}

func (g *GUI) updateCurrentChannelID() error {
	channelId := g.tabs.search.CurrentChannelID
	if channelId == "" {
		channelId = "No channel"
	}

	g.widgets.currentChannelId.Reset()

	return g.widgets.currentChannelId.Write(channelId)
}

func (g *GUI) updateCurrentStatus(newStatus string) error {
	return g.widgets.currentStatus.Write(newStatus, text.WriteReplace())
}

func (g GUI) clearSearchInput() {
	g.widgets.searchInput.ReadAndClear()
}

func (g *GUI) getSearchTextInput() (*textinput.TextInput, error) {
	return textinput.New(
		textinput.Label("Search packages/options: ", cell.FgColor(cell.ColorAqua)),
		textinput.Border(linestyle.None),
		textinput.PlaceHolder("enter any text"),
		textinput.FillColor(cell.ColorDefault),
		textinput.ExclusiveKeyboardOnFocus(),
		textinput.OnChange(g.handleSearchInputChange),
		textinput.OnSubmit(g.handleSearchInputSubmit))
}

func (g *GUI) handleSearchInputChange(input string) {
	if g.tabs.search.WaitForEnter {
		return
	}

	if g.tabs.search.Name == HOME_MANAGER_OPTIONS && !g.nixClient.HomeManagerOptionsAlreadyFetched() {
		return
	}

	g.performSearch(context.Background(), input)
}

func (g *GUI) handleSearchInputSubmit(input string) error {
	g.performSearch(context.Background(), input)

	return nil
}

func (GUI) getResultsBoard() (*text.Text, error) {
	return text.New(text.WrapAtWords())
}

func (g GUI) getCurrentLabelWidget() (*text.Text, error) {
	return g.newTextWidget(g.tabs.search.Label)
}

func (g GUI) getCurrentStatusWidget() (*text.Text, error) {
	return g.newTextWidget(WAITING, text.WriteCellOpts(cell.Bold()))
}

func (g GUI) getCurrentSourceWidget() (*text.Text, error) {
	return g.newTextWidget(g.tabs.search.Source, text.WriteCellOpts(cell.Bold()))
}

func (g GUI) getCurrentChannelWidget() (*text.Text, error) {
	return g.newTextWidget(g.tabs.search.CurrentChannelID, text.WriteCellOpts(cell.Bold()))
}

func (GUI) newTextWidget(content string, tOpts ...text.WriteOption) (*text.Text, error) {
	t, err := text.New(text.DisableScrolling())
	if err != nil {
		return nil, err
	}

	return t, t.Write(content, tOpts...)
}
