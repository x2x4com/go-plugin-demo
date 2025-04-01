package dynamic_plugin_shared

type DynamicPluginInterface interface {
	Invoke(method string, args, options []interface{}) (interface{}, error)
	Help(method string) (string, error)
	Version() string
}
