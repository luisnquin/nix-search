package gui

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
	GUI struct {
		// The nix client used to perform searches.
		nixClient *nix_search.Client
		// The configuration for the application.
		config *config.Config

		// The GUI widgets.
		widgets widgets
		// The GUI Tabs.
		tabs *tabs
	}

	// The GUI tabs.
	tabs struct {
		// The current search tab.
		search *searchTabConfig
	}

	// The GUI widgets.
	widgets struct {
		searchInput      *textinput.TextInput
		resultsBoard     *text.Text
		currentChannelId *text.Text
		currentStatus    *text.Text
		currentSource    *text.Text
		currentLabel     *text.Text
	}
)

func New(config *config.Config) (GUI, error) {
	gui := GUI{
		nixClient: nix_search.NewClient(config),
		config:    config,
	}

	gui.tabs = &tabs{
		search: gui.getDefaultSearchTab(),
	}

	if err := gui.initWidgets(); err != nil {
		return GUI{}, err
	}

	return gui, nil
}

func (g GUI) Run(ctx context.Context) error {
	return g.run(ctx)
}

func (g GUI) run(ctx context.Context) error {
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
					g.widgets.searchInput,
					container.Border(linestyle.Light),
					container.AlignHorizontal(align.HorizontalCenter),
					container.Focused(),
					container.PaddingLeft(1),
				),
			),
			grid.RowHeightPerc(6,
				grid.ColWidthPerc(25,
					grid.Widget(
						g.widgets.currentLabel,
						container.BorderTitle("Current tab"),
						container.Border(linestyle.Light),
						container.BorderColor(cell.ColorMagenta),
					)),
				grid.ColWidthPerc(20,
					grid.Widget(
						g.widgets.currentChannelId,
						container.BorderTitle("Channel ID"),
						container.Border(linestyle.Light),
						container.BorderColor(cell.ColorNavy),
					)),
				grid.ColWidthPerc(20,
					grid.Widget(
						g.widgets.currentStatus,
						container.BorderTitle("Status"),
						container.Border(linestyle.Light),
						container.BorderColor(cell.ColorGreen),
					)),
				grid.ColWidthPerc(35,
					grid.Widget(
						g.widgets.currentSource,
						container.BorderTitle("Source"),
						container.Border(linestyle.Light),
						container.BorderColor(cell.ColorFuchsia),
					)),
			),
			grid.RowHeightPerc(90,
				grid.Widget(
					g.widgets.resultsBoard,
					container.Border(linestyle.Light),
					container.BorderTitle("Search results"),
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
				g.previousTab()

			case k.Key == keyboard.KeyCtrlRsqBracket:
				g.nextTab()

			case k.Key == keyboard.KeyCtrlSpace:
				g.nextChannel()

			case k.Key == keyboard.KeyCtrlQ:
				g.clearSearchInput()
			}
		}),
		termdash.RedrawInterval(350 * time.Millisecond),
	}

	return termdash.Run(ctx, t, c, termOptions...)
}
