package gui

import (
	"context"
	"fmt"
	"log"
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
)

type (
	GUI struct {
		// The nix client used to perform searches.
		nixClient *nix_search.Client
		// The configuration for the application.
		config *config.Config

		// The tcell terminal.
		terminal *tcell.Terminal
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

// Prepare the the GUI components and returns it.
func New(config *config.Config) (GUI, error) {
	terminal, err := tcell.New()
	if err != nil {
		return GUI{}, fmt.Errorf("unable to create tcell terminal: %w", err)
	}

	gui := GUI{
		nixClient: nix_search.NewClient(config),
		terminal:  terminal,
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

// Run the GUI of the program.
func (g GUI) Run(ctx context.Context) error {
	defer g.handleProgramPanic()
	defer g.terminal.Close()

	if err := g.run(ctx); err != nil {
		return err
	}

	return nil
}

func (g GUI) run(ctx context.Context) error {
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

	c, err := container.New(g.terminal, gridOpts...)
	if err != nil {
		return err
	}

	termOptions := []termdash.Option{
		termdash.KeyboardSubscriber(func(k *terminalapi.Keyboard) {
			switch { // keyboard.KeyEsc
			case k.Key == keyboard.KeyCtrlLsqBracket:
				g.previousTab()

			case k.Key == keyboard.KeyCtrlRsqBracket:
				g.nextTab()

			case k.Key == keyboard.KeyCtrlSpace:
				g.nextChannel()

			case k.Key == keyboard.KeyCtrlQ:
				g.clearSearchInput()

			case k.Key == keyboard.KeyCtrlC:
				cancel()
			}
		}),
		termdash.RedrawInterval(350 * time.Millisecond),
	}

	return termdash.Run(ctx, g.terminal, c, termOptions...)
}

func (g GUI) handleProgramPanic() {
	if err := recover(); err != nil {
		g.terminal.Close()
		log.Fatal(err)
	}
}
