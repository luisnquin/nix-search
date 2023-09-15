package gui

import "github.com/samber/lo"

func (g GUI) getCurrentBranchIndex() int {
	_, index, found := lo.FindIndexOf(g.tabs.search.ChannelIDs, func(channelId string) bool {
		return channelId == g.tabs.search.CurrentChannelID
	})
	if found {
		return index
	}

	return 0
}

func (g *GUI) nextChannel() {
	index := g.getCurrentBranchIndex()

	if len(g.tabs.search.ChannelIDs) == 0 {
		return
	} else if index+1 == len(g.tabs.search.ChannelIDs) {
		index = 0
	} else {
		index++
	}

	g.tabs.search.CurrentChannelID = g.tabs.search.ChannelIDs[index]
	g.updateCurrentChannelID()
}
