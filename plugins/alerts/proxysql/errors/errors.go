package errors

import (
	"fmt"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
	"github.com/swapbyt3s/zenit/plugins/lists/loader"
	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
)

type Server struct {
	Host   string
	Group  string
	Errors int
}

type ProxySQLErrors struct {}

func (l *ProxySQLErrors) Collect() {
	var m = metrics.Load()
	var s Server

	for _, m := range *m {
		if m.Key == "zenit_proxysql_connection_pool" {
			for _, metricTag := range m.Tags {
				if metricTag.Name == "host" {
					s.Host = metricTag.Value
				} else if metricTag.Name == "group" {
					s.Group = metricTag.Value
				}
			}

			for _, value := range m.Values.([]metrics.Value) {
				if value.Key == "errors" {
					if v, ok := value.Value.(uint); ok {
						s.Errors = int(v)
						break
					}
				}
			}

			// Build one message with details for notification:
			var message = fmt.Sprintf("*Server:* %s\n*Group:* %s\n*Error:* %d\n", s.Host, s.Group, s.Errors)

			// fmt.Printf("ProxySQL Error Message: %s\n", message)

			// Register new check and update last status:
			alerts.Load().Register(
				"proxysql_pool_errors_" + s.Host + s.Group,
				"ProxySQL Connection Pool Errors",
				config.File.ProxySQL.Alerts.Errors.Duration,
				config.File.ProxySQL.Alerts.Errors.Warning,
				config.File.ProxySQL.Alerts.Errors.Critical,
				s.Errors,
				message,
			)

		}
	}
}

func init() {
	loader.Add("AlertProxySQLErrors", func() loader.Plugin { return &ProxySQLErrors{} })
}
