package gui

import "github.com/samber/lo"

func (app App) getCurrentBranchIndex() int {
	_, index, found := lo.FindIndexOf(app.tabs.search.ChannelIDs, func(channelId string) bool {
		return channelId == app.tabs.search.CurrentChannelID
	})
	if found {
		return index
	}

	return 0
}

// func (app *App) previousChannel() {
// 	index := app.getCurrentBranchIndex()

// 	if index-1 >= 0 {
// 		app.tabs.search.CurrentChannelID = app.tabs.search.ChannelIDs[index-1]
// 		app.updateCurrentChannelID()
// 	}
// }

// func (app *App) nextChannel() {
// 	index := app.getCurrentBranchIndex()

// 	if index+1 < len(app.tabs.search.ChannelIDs) {
// 		app.tabs.search.CurrentChannelID = app.tabs.search.ChannelIDs[index+1]
// 		app.updateCurrentChannelID()
// 	}
// }

func (app *App) nextChannel() {
	index := app.getCurrentBranchIndex()

	if len(app.tabs.search.ChannelIDs) == 0 {
		return
	} else if index+1 == len(app.tabs.search.ChannelIDs) {
		index = 0
	} else {
		index++
	}

	app.tabs.search.CurrentChannelID = app.tabs.search.ChannelIDs[index]
	app.updateCurrentChannelID()
}
