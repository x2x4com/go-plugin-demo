package main

import (
	"errors"
	"fmt"
	dynamic_plugin_shared "go-plugin-demo/src/internal/plugin/shared"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

var ExportFuncMap = map[string]dynamic_plugin_shared.DynamicFunc{
	"AddDays": {
		Name:       "AddDays",
		Call:       AddDays,
		Help:       "Adds days to a given date.",
		HasArgs:    true,
		HasOptions: false,
	},
	"Format": {
		Name:       "Format",
		Call:       Format,
		Help:       "Formats a date to a specified layout.",
		HasArgs:    true,
		HasOptions: false,
	},
	"Parse": {
		Name:       "Parse",
		Call:       Parse,
		Help:       "Parses a date string into a time.Time object.",
		HasArgs:    true,
		HasOptions: false,
	},
	"Between": {
		Name:       "Between",
		Call:       Between,
		Help:       "Calculates the number of days between two dates.",
		HasArgs:    true,
		HasOptions: false,
	},
}

type DataUtilsImpl struct {
	logger hclog.Logger
}

func (ds *DataUtilsImpl) Invoke(method string, args, options []interface{}) (interface{}, error) {
	// 判断method是否在ExportFuncMap中, 如果存在返回nil和error，方法不存在
	if f, ok := ExportFuncMap[method]; ok {
		if f.HasArgs && f.HasOptions {
			return f.Call(args, options)
		} else if f.HasArgs {
			return f.Call(args, []interface{}{})
		} else {
			return f.Call([]interface{}{}, []interface{}{})
		}
	} else {
		msg := fmt.Sprintf("method %s not found", method)
		ds.logger.Error(msg)
		return nil, errors.New(msg)
	}
}

func (ds *DataUtilsImpl) Help(method string) (string, error) {
	if f, ok := ExportFuncMap[method]; ok {
		return f.GetFuncHelp(), nil
	} else {
		msg := fmt.Sprintf("method %s not found", method)
		ds.logger.Error(msg)
		return "", errors.New(msg)
	}
}

func (ds *DataUtilsImpl) Version() string {
	return "1.0.0"
}

// AddDays 日期加减
// func AddDays(date time.Time, days int) time.Time {
func AddDays(args, options []interface{}) (interface{}, error) {
	//return date.AddDate(0, 0, days)
	if len(args) != 2 {
		return nil, errors.New("AddDays requires exactly 2 arguments: date and days")
	}
	//fmt.Printf("%v", args...)
	date, ok1 := time.Parse(time.RFC3339, args[0].(string))
	if ok1 != nil {
		return nil, ok1
	}
	days, ok2 := args[1].(int)
	if !ok2 {
		return nil, errors.New("AddDays requires arguments of type time.Time and int")
	}
	return date.AddDate(0, 0, days).Format(time.RFC3339), nil
}

// Format 日期格式化
// func Format(date time.Time, layout string) string {
func Format(args, options []interface{}) (interface{}, error) {
	//return date.Format(layout)
	if len(args) != 2 {
		return nil, errors.New("Format requires exactly 2 arguments: date and layout")
	}
	date, ok1 := time.Parse(time.RFC3339, args[0].(string))
	if ok1 != nil {
		return nil, ok1
	}
	layout, ok2 := args[1].(string)
	if !ok2 {
		return nil, errors.New("Format requires arguments of type time.Time and string")
	}
	return date.Format(layout), nil
}

// Parse 日期解析
// func Parse(dateStr, layout string) (time.Time, error) {
func Parse(args, options []interface{}) (interface{}, error) {
	//return time.Parse(layout, dateStr)
	if len(args) != 2 {
		return nil, errors.New("Parse requires exactly 2 arguments: dateStr and layout")
	}
	dateStr, ok1 := args[0].(string)
	layout, ok2 := args[1].(string)
	if !ok1 || !ok2 {
		return nil, errors.New("Parse requires arguments of type string and string")
	}
	res, err := time.Parse(layout, dateStr)
	if err != nil {
		return nil, err
	}
	return res.Format(time.RFC3339), nil
}

// Between 计算日期差值
// func Between(start, end time.Time) int {
func Between(args, options []interface{}) (interface{}, error) {
	//duration := end.Sub(start)
	//return int(duration.Hours() / 24)
	if len(args) != 2 {
		return nil, errors.New("Between requires exactly 2 arguments: start and end")
	}
	start, ok1 := time.Parse(time.RFC3339, args[0].(string))
	end, ok2 := time.Parse(time.RFC3339, args[1].(string))
	if ok1 != nil || ok2 != nil {
		return nil, errors.New("Between requires arguments of type time.Time and time.Time")
	}
	duration := end.Sub(start)
	return int(duration.Hours() / 24), nil
}

func main() {
	handshake := dynamic_plugin_shared.GenHandShakeConfig("calculator")

	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	ds := &DataUtilsImpl{logger: logger}

	var pluginMap = map[string]plugin.Plugin{
		"date_utils": &dynamic_plugin_shared.DynamicPlugin{Impl: ds},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshake,
		Plugins:         pluginMap,
	})
}
