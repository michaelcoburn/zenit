package mem

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
	"github.com/swapbyt3s/zenit/plugins/lists/loader"
	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
)

type OSMEM struct {}

func (l *OSMEM) Collect() {
	if ! config.File.OS.Alerts.MEM.Enable {
		log.Info("Require to enable OS MEM in config file.")
		return
	}

	var metrics = metrics.Load()
	var message string = ""
	var value = metrics.FetchOne("zenit_os", "name", "mem")
	var percentage = uint64(common.InterfaceToFloat64(value))

//	if percentage == -1 {
//		return
//	}

	message += fmt.Sprintf("*Memory:* %d%%\n", percentage)

	alerts.Load().Register(
		"mem",
		"MEM",
		config.File.OS.Alerts.MEM.Duration,
		config.File.OS.Alerts.MEM.Warning,
		config.File.OS.Alerts.MEM.Critical,
		percentage,
		message,
	)
}

func init() {
	loader.Add("AlertOSMEM", func() loader.Plugin { return &OSMEM{} })
}
