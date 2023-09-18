package gui

import (
	"context"
	"fmt"
	"time"

	"github.com/luisnquin/nix-search/internal/config"
	"github.com/luisnquin/nix-search/internal/log"
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
		logger log.Logger

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

const REDRAW_INTERVAL = 350 * time.Millisecond

// Prepare the GUI components and returns it.
func New(logger log.Logger, config *config.Config, nixClient *nix_search.Client) (GUI, error) {
	terminal, err := tcell.New()
	if err != nil {
		logger.Err(err).Msg("unable to create tcell terminal")

		return GUI{}, fmt.Errorf("unable to create tcell terminal: %w", err)
	}

	logger.Trace().Msg("initializing program components...")

	gui := GUI{
		nixClient: nixClient,
		terminal:  terminal,
		config:    config,
		logger:    logger,
	}

	logger.Trace().Msg("initializing GUI widgets...")

	gui.tabs = &tabs{
		search: gui.getSelectedOrDefaultTab(),
	}

	logger.Trace().Msgf("current search tab: %s", gui.tabs.search.Name)

	if err := gui.initWidgets(); err != nil {
		logger.Err(err).Msg("an error was detected while initializing GUI widgets")

		return GUI{}, err
	}

	logger.Trace().Msg("initialized, providing a GUI instance ready to be used...")

	return gui, nil
}

// Run the GUI of the program.
func (g GUI) Run(ctx context.Context) error {
	defer g.handleProgramPanic()
	defer func() {
		g.logger.Debug().Msg("closing tcell terminal...")
		g.terminal.Close()
		g.logger.Debug().Msg("tcell terminal has been closed")
	}()

	tsize := g.terminal.Size()

	g.logger.Trace().
		Int("tx", tsize.X).Int("ty", tsize.Y).
		Bool("¿context is nil?", ctx == nil).Msg("starting the program GUI...")

	if err := g.run(ctx); err != nil {
		g.logger.Err(err).Str("current tab", g.tabs.search.Name.String()).
			Msg("error detected while running GUI...")

		return err
	}

	g.logger.Trace().Msg("user exited GUI without errors(i guess)")

	return nil
}

func (g GUI) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	builder := grid.New()

	//nolint:gomnd
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
			switch {
			case k.Key == keyboard.KeyCtrlLsqBracket:
				g.logger.Trace().Msg("going to previous tab...")
				g.previousTab()

			case k.Key == keyboard.KeyCtrlRsqBracket:
				g.logger.Trace().Msg("going to next tab...")
				g.nextTab()

			case k.Key == keyboard.KeyCtrlSpace:
				g.logger.Trace().Msg("going to next channel...")
				g.nextChannel()

			case lo.Contains([]keyboard.Key{keyboard.KeyCtrlQ, keyboard.KeyBackspace}, k.Key):
				g.logger.Trace().Msg("going to clear search input...")
				g.clearSearchInput()

			case k.Key == keyboard.KeyCtrlC:
				g.logger.Trace().Msg("closing application(cancelling root context)...")
				cancel()
			}
		}),
		termdash.RedrawInterval(REDRAW_INTERVAL),
	}

	return termdash.Run(ctx, g.terminal, c, termOptions...)
}

func (g GUI) handleProgramPanic() {
	if err := recover(); err != nil {
		g.terminal.Close()
		g.logger.Fatal().Any("error", err).Msg("catched during program panic")
	}
}
