package main

import (
	"context"
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

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	appConfig := config.Load()
	nixClient := nix_search.NewClient(appConfig)

	t, err := tcell.New()
	if err != nil {
		panic(err)
	}

	defer t.Close()

	resultsBoard, err := text.New()
	if err != nil {
		panic(err)
	}

	input, err := textinput.New(
		textinput.Label("What do you want to search:", cell.FgColor(cell.ColorNumber(33))),
		// textinput.MaxWidthCells(80),
		textinput.Border(linestyle.Light),
		textinput.PlaceHolder("Enter any text"),
		textinput.ExclusiveKeyboardOnFocus(),
		textinput.OnSubmit(func(text string) error {
			options, err := nixClient.SearchHomeManagerOptions(ctx, strings.TrimSpace(text))
			if err != nil {
				return err
			}

			resultsBoard.Reset()

			results := strings.Join(lo.Map(options, func(opt *nix.HomeManagerOption, _ int) string {
				return opt.String()
			}), " ")

			resultsBoard.Write(results)

			return nil
		}))
	if err != nil {
		panic(err)
	}

	builder := grid.New()

	builder.Add(
		grid.RowHeightPerc(99,
			grid.RowHeightFixed(20,
				grid.Widget(
					input,
					container.AlignHorizontal(align.HorizontalCenter),
					container.AlignVertical(align.VerticalBottom),
					container.MarginBottom(1),
					container.Focused(),
				),
			),
			grid.RowHeightPerc(50,
				grid.Widget(
					resultsBoard,
					container.Border(linestyle.Light),
					container.BorderTitle("Search results"),
					container.AlignHorizontal(align.HorizontalCenter),
					container.AlignHorizontal(align.Horizontal(align.VerticalBottom)),
				),
			),
		),
	)

	gridOpts, err := builder.Build()
	if err != nil {
		panic(err)
	}

	c, err := container.New(t, gridOpts...)
	if err != nil {
		panic(err)
	}

	termOptions := []termdash.Option{
		termdash.KeyboardSubscriber(func(k *terminalapi.Keyboard) {
			if lo.Contains([]keyboard.Key{keyboard.KeyEsc, keyboard.KeyCtrlC}, k.Key) {
				cancel()
			}
		}),
		termdash.RedrawInterval(500 * time.Millisecond),
	}

	if err := termdash.Run(ctx, t, c, termOptions...); err != nil {
		panic(err)
	}
}
