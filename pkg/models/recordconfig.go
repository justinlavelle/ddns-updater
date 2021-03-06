package models

import (
	"sync"
	"time"
)

// RecordConfigType contains all the information to update and display a DNS record
type RecordConfigType struct { // internal
	Settings SettingsType  // fixed
	Status   statusType    // changes for each update
	History  historyType   // past information
	M        *sync.RWMutex // TODO inherit
}

// NewRecordConfig returns a new recordConfig with settings
func NewRecordConfig(settings SettingsType, IPs []string, tSuccess time.Time) RecordConfigType {
	return RecordConfigType{
		Settings: settings,
		History: historyType{
			IPs:      IPs,
			TSuccess: tSuccess,
		},
		M: new(sync.RWMutex),
	}
}

func (conf *RecordConfigType) String() string {
	conf.M.RLock()
	defer conf.M.RUnlock()
	return conf.Settings.String() + ": " + conf.Status.string() + "; " + conf.History.string()
}

func (conf RecordConfigType) toHTML() HTMLRow {
	row := HTMLRow{
		Domain:   conf.Settings.getHTMLDomain(),
		Host:     conf.Settings.Host,
		Provider: conf.Settings.getHTMLProvider(),
		IPMethod: conf.Settings.getHTMLIPMethod(),
	}
	if conf.Status.Code == UPTODATE {
		conf.Status.Message = "No IP change for " + durationString(conf.History.TSuccess)
	}
	row.Status = conf.Status.toHTML()
	if len(conf.History.IPs) > 0 {
		row.IP = "<a href=\"https://ipinfo.io/" + conf.History.IPs[0] + "\">" + conf.History.IPs[0] + "</a>"
	} else {
		row.IP = "N/A"
	}
	if len(conf.History.IPs) > 1 {
		row.IPs = conf.History.IPs[1:]
		for i := range row.IPs {
			if i == len(row.IPs)-1 {
				break
			}
			row.IPs[i] += ", "
		}
	} else {
		row.IPs = []string{"N/A"}
	}
	return row
}
