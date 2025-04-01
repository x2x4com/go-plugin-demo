package dynamic_plugin_shared

import "github.com/hashicorp/go-plugin"

func GenHandShakeConfig(pluginName string) plugin.HandshakeConfig {
	return plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "DYNAMIC_PLUGIN_" + pluginName,
		MagicCookieValue: pluginName,
	}
}
