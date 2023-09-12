package config

import "github.com/samber/lo"

type (
	NixConfig struct {
		DefaultChannel string       `json:"default_channel"`
		Channels       []NixChannel `json:"channels"`
		Sources        NixSources   `json:"sources"`
	}

	NixChannel struct {
		ID     string `json:"id"`
		Branch string `json:"branch"`
		JobSet string `json:"jobset"`
		Status string `json:"status"`
	}

	NixSources struct {
		HomeManagerOptions HomeManagerOptionsConfig `json:"home_manager_options"`
		ElasticSearch      ElasticSearchConfig      `json:"elastic_search"`
	}

	ElasticSearchConfig struct {
		URL            string `json:"url"`
		Username       string `json:"username"`
		Password       string `json:"password"`
		MappingVersion string `json:"mapping_version"`
	}

	HomeManagerOptionsConfig struct {
		URL string `json:"url"`
	}
)

func (nc NixConfig) FindChannel(channelId string) (NixChannel, bool) {
	return lo.Find(nc.Channels, func(channel NixChannel) bool {
		return channel.ID == channelId
	})
}

func (nc NixConfig) FindChannelWithStatus(channelStatus string) (NixChannel, bool) {
	return lo.Find(nc.Channels, func(channel NixChannel) bool {
		return channel.Status == channelStatus
	})
}
