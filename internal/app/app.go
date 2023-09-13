package app

import (
	"context"
	"time"

	"github.com/luisnquin/nix-search/internal/config"
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

	if err := app.initWidgets(); err != nil {
		return App{}, err
	}

	return app, nil
}

func (a App) Run(ctx context.Context) error {
	return a.run(ctx)
}

func (app App) run(ctx context.Context) error {
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
					app.elements.searchInput,
					container.Border(linestyle.Light),
					container.AlignHorizontal(align.HorizontalCenter),
					container.Focused(),
					container.PaddingLeft(1),
				),
			),
			grid.RowHeightPerc(6,
				grid.ColWidthPerc(30,
					grid.Widget(
						app.elements.currentLabel,
						container.BorderTitle("Current tab"),
						container.Border(linestyle.Light),
						container.BorderColor(cell.ColorMagenta),
					)),
				grid.ColWidthPerc(15,
					grid.Widget(
						app.elements.currentStatus,
						container.BorderTitle("Status"),
						container.Border(linestyle.Light),
						container.BorderColor(cell.ColorGreen),
					)),
				grid.ColWidthPerc(55,
					grid.Widget(
						app.elements.currentSource,
						container.BorderTitle("Source"),
						container.Border(linestyle.Light),
						container.BorderColor(cell.ColorNavy),
					),
				)),
			grid.RowHeightPerc(87,
				grid.Widget(
					app.elements.resultsBoard,
					container.Border(linestyle.Light),
					container.BorderTitle("Search results"),
					container.BorderColor(cell.ColorAqua),
				),
			),
			grid.RowHeightPerc(3,
				grid.Widget(
					app.elements.searchOptions,
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
			switch { // keyboard.KeyEsc
			case lo.Contains([]keyboard.Key{keyboard.KeyCtrlC}, k.Key):
				cancel()
			case k.Key == keyboard.KeyCtrlLsqBracket:
				app.previousTab()
				app.updateWidgetTexts()

			case k.Key == keyboard.KeyCtrlRsqBracket:
				app.nextTab()
				app.updateWidgetTexts()
			}
		}),
		termdash.RedrawInterval(500 * time.Millisecond),
	}

	return termdash.Run(ctx, t, c, termOptions...)
}
