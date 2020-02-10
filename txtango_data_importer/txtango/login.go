//Package txtango manages Transics TX-TANGO connection
package txtango

import (
	"os"
	"time"
)

var loginTemplate = `
{{ define "login" }}
<Login>
    <DateTime>{{.DateTime}}</DateTime>
    <Version>{{.Version}}</Version>
    <Dispatcher>{{.Dispatcher}}</Dispatcher>
    <Password>{{.Password}}</Password>
    <SystemNr>{{.SystemNr}}</SystemNr>
    <ApplicationName />
    <ApplicationVersion />
    <PcName />
    <Integrator>{{.Integrator}}</Integrator>
    <Language>{{.Language}}</Language>
</Login>
{{ end }}
`

//Login to login in tx-tango
//login happens in each request
type Login struct {
	Date       time.Time
	Version    int
	Dispatcher string
	Password   string
	SystemNr   string
	Integrator string
	Language   string
}

//Authenticate build authentication bloc
func Authenticate() (*Login, error) {
	var login Login

	//fill in login credentials
	login.Dispatcher = os.Getenv("TX_USERNAME")
	login.Password = os.Getenv("TX_PASSWORD")
	login.Integrator = os.Getenv("TX_INTEGRATOR")
	login.SystemNr = os.Getenv("TX_SYSTEM_NR")
	login.Language = "EN"

	return &login, nil
}
