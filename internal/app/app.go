package app

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/luisnquin/nix-search/internal/config"
	"github.com/luisnquin/nix-search/internal/nix"
	nix_search "github.com/luisnquin/nix-search/internal/nix/search"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/align"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/text"
	"github.com/mum4k/termdash/widgets/textinput"
	"github.com/samber/lo"
)

type (
	App struct {
		nixClient *nix_search.Client

		currentSearchTab *searchTabConfig
		elements
	}

	elements struct {
		searchInput   *textinput.TextInput
		resultsBoard  *text.Text
		currentLabel  *text.Text
		currentStatus *text.Text
		currentSource *text.Text
		searchOptions *text.Text
	}
)

func New(config *config.Config) (App, error) {
	app := App{
		nixClient: nix_search.NewClient(config),
	}

	app.currentSearchTab = app.getDefaultSearchTab()

	var err error

	app.elements.resultsBoard, err = app.getResultsBoard()
	if err != nil {
		return App{}, fmt.Errorf("results board: %w", err)
	}

	app.elements.searchInput, err = app.getSearchTextInput()
	if err != nil {
		return App{}, fmt.Errorf("search input: %w", err)
	}

	app.elements.currentStatus, err = app.getCurrentStatusWidget()
	if err != nil {
		return App{}, fmt.Errorf("current status widget: %w", err)
	}

	app.elements.currentLabel, err = app.getCurrentLabelWidget()
	if err != nil {
		return App{}, fmt.Errorf("current label widget: %w", err)
	}

	app.elements.searchOptions, err = app.getSearchOptionsWidget()
	if err != nil {
		return App{}, fmt.Errorf("search options widget: %w", err)
	}

	app.elements.currentSource, err = app.getCurrentSourceWidget()
	if err != nil {
		return App{}, fmt.Errorf("current source widget: %w", err)
	}

	return app, nil
}

func (a App) Run(ctx context.Context) error {
	return a.run(ctx)
}

func (a App) run(ctx context.Context) error {
	t, err := tcell.New()
	if err != nil {
		return err
	}

	defer t.Close()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	builder := grid.New()

	builder.Add(
		grid.RowHeightPerc(99,
			grid.RowHeightFixed(5,
				grid.Widget(
					a.elements.searchInput,
					container.Border(linestyle.Light),
					container.AlignHorizontal(align.HorizontalCenter),
					container.Focused(),
					container.PaddingLeft(1),
				),
			),
			grid.RowHeightPerc(6,
				grid.ColWidthPerc(30,
					grid.Widget(
						a.elements.currentLabel,
						container.BorderTitle("Current tab"),
						container.Border(linestyle.Light),
						container.BorderColor(cell.ColorMagenta),
					)),
				grid.ColWidthPerc(15,
					grid.Widget(
						a.elements.currentStatus,
						container.BorderTitle("Status"),
						container.Border(linestyle.Light),
						container.BorderColor(cell.ColorGreen),
					)),
				grid.ColWidthPerc(55,
					grid.Widget(
						a.elements.currentSource,
						container.BorderTitle("Source"),
						container.Border(linestyle.Light),
						container.BorderColor(cell.ColorNavy),
					),
				)),
			grid.RowHeightPerc(87,
				grid.Widget(
					a.elements.resultsBoard,
					container.Border(linestyle.Light),
					container.BorderTitle("Search results"),
					container.BorderColor(cell.ColorAqua),
				),
			),
			grid.RowHeightPerc(3,
				grid.Widget(
					a.elements.searchOptions,
					container.Border(linestyle.Light),
					container.BorderTitle("Nix options"),
					container.BorderColor(cell.ColorAqua),
				),
			),
		),
	)

	gridOpts, err := builder.Build()
	if err != nil {
		return err
	}

	c, err := container.New(t, gridOpts...)
	if err != nil {
		return err
	}

	termOptions := []termdash.Option{
		termdash.KeyboardSubscriber(func(k *terminalapi.Keyboard) {
			if lo.Contains([]keyboard.Key{keyboard.KeyEsc, keyboard.KeyCtrlC}, k.Key) {
				cancel()
			}
		}),
		termdash.RedrawInterval(500 * time.Millisecond),
	}

	return termdash.Run(ctx, t, c, termOptions...)
}

func (a *App) getSearchTextInput() (*textinput.TextInput, error) {
	ctx := context.Background()

	

	return textinput.New(
		textinput.Label(a.currentSearchTab.Prompt, cell.FgColor(cell.ColorAqua)),
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
	return a.newTextWidget(a.currentSearchTab.State, text.WriteCellOpts(cell.Bold()))
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
