package app

import "github.com/samber/lo"

func (a App) getDefaultSearchTab() *searchTabConfig {
	config, _ := lo.Find(a.getSearchTabs(), func(item searchTabConfig) bool {
		return item.Tab == HOME_MANAGER_OPTIONS
	})

	return &config
}
