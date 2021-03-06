// This parse OLD Style
// https://dev.mysql.com/doc/refman/5.5/en/audit-log-file-formats.html
// TODO: Move this package to inputs/parsers/mysqlauditlog

package audit

import (
	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/common/sql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/outputs/clickhouse"
)

func Start() {
	if config.File.MySQL.Inputs.AuditLog.Enable {
		if config.File.General.Debug {
			log.Debug("Load MySQL AuditLog")
			log.Debug("Read MySQL AuditLog: " + config.File.MySQL.Inputs.AuditLog.LogPath)
		}

		if !clickhouse.Check() {
			log.Error("AuditLog require active connection to ClickHouse.")
		}

		if config.File.MySQL.Inputs.AuditLog.Format == "xml-old" {
			channel_tail := make(chan string)
			channel_parser := make(chan map[string]string)
			channel_data := make(chan map[string]string)

			event := &clickhouse.Event{
				Type:    "AuditLog",
				Schema:  "zenit",
				Table:   "mysql_audit_log",
				Size:    config.File.MySQL.Inputs.AuditLog.BufferSize,
				Timeout: config.File.MySQL.Inputs.AuditLog.BufferTimeOut,
				Wildcard: map[string]string{
					"_time":          "'%s'",
					"command_class":  "'%s'",
					"connection_id":  "%s",
					"db":             "'%s'",
					"host":           "'%s'",
					"host_ip":        "IPv4StringToNum('%s')",
					"host_name":      "'%s'",
					"ip":             "'%s'",
					"name":           "'%s'",
					"os_login":       "'%s'",
					"os_user":        "'%s'",
					"priv_user":      "'%s'",
					"proxy_user":     "'%s'",
					"record":         "'%s'",
					"sqltext":        "'%s'",
					"sqltext_digest": "'%s'",
					"status":         "%s",
					"user":           "'%s'",
				},
				Values: []map[string]string{},
			}

			go common.Tail(config.File.MySQL.Inputs.AuditLog.LogPath, channel_tail)
			go Parser(config.File.MySQL.Inputs.AuditLog.LogPath, channel_tail, channel_parser)
			go clickhouse.Send(event, channel_data)

			go func() {
				for channel_event := range channel_parser {
					channel_data <- channel_event
				}
			}()
		}
	}
}

func Parser(path string, tail <-chan string, parser chan<- map[string]string) {
	var buffer []string

	go func() {
		defer close(parser)

		for line := range tail {
			if line == "<AUDIT_RECORD" && len(buffer) > 0 {
				result := map[string]string{
					"_time":          "",
					"command_class":  "",
					"connection_id":  "",
					"db":             "",
					"host":           "",
					"host_ip":        "",
					"host_name":      "",
					"ip":             "",
					"name":           "",
					"os_login":       "",
					"os_user":        "",
					"priv_user":      "",
					"proxy_user":     "",
					"record":         "",
					"sqltext":        "",
					"sqltext_digest": "",
					"status":         "",
					"user":           "",
				}

				for _, kv := range buffer {
					key, value := common.SplitKeyAndValue(&kv)
					result[key] = common.Trim(&value)
				}

				buffer = buffer[:0]

				if common.KeyInMap("user", result) {
					result["user"] = mysql.ClearUser(result["user"])
				}

				if common.KeyInMap("sqltext", result) {
					result["sqltext_digest"] = sql.Digest(result["sqltext"])
				}

				// Convert timestamp ISO 8601 (UTC) to RFC 3339:
				result["_time"] = common.ToDateTime(result["timestamp"], "2006-01-02T15:04:05 UTC")
				result["host_ip"] = config.IPAddress
				result["host_name"] = config.File.General.Hostname
				result["sqltext"] = common.Escape(result["sqltext"])
				result["sqltext_digest"] = common.Escape(result["sqltext_digest"])

				// Remove unnused key:
				delete(result, "timestamp")
				delete(result, "record")

				parser <- result
			} else if line != "/>" {
				buffer = append(buffer, line)
			}
		}
	}()
}
