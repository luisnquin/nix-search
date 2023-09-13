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

type App struct {
	nixClient *nix_search.Client

	searchInput  *textinput.TextInput
	resultsBoard *text.Text
}

func New(config *config.Config) (App, error) {
	app := App{
		nixClient: nix_search.NewClient(config),
	}

	var err error

	app.resultsBoard, err = app.getResultsBoard()
	if err != nil {
		return App{}, fmt.Errorf("unable to create results board: %w", err)
	}

	app.searchInput, err = app.getSearchTextInput()
	if err != nil {
		return App{}, fmt.Errorf("unable to create search input: %w", err)
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
					a.searchInput,
					container.Border(linestyle.Light),
					container.AlignHorizontal(align.HorizontalCenter),
					container.Focused(),
					container.PaddingLeft(1),
				),
			),
			grid.RowHeightPerc(90,
				grid.Widget(
					a.resultsBoard,
					container.Border(linestyle.Light),
					container.BorderTitle("Search results"),
					container.AlignHorizontal(align.HorizontalCenter),
					container.AlignHorizontal(align.Horizontal(align.VerticalBottom)),
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
		textinput.Label("What do you want to search: ", cell.FgColor(cell.ColorAqua)),
		textinput.Border(linestyle.None),
		textinput.PlaceHolder("Enter any text"),
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
	return text.New()
}
