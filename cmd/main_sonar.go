package cmd

// SonarHTTPCheck represents sonar check record
// Example:
//
//	{
//	   "id": 81994,
//	   "name": "msi-prod-pdfeditor-backup",
//	   "host": "18.180.166.79",
//	   "port": 443,
//	   "protocolType": "HTTPS",
//	   "ipVersion": "IPV4",
//	   "fqdn": "msi-pdfeditor.surfly.jp",
//	   "path": "/blackbox/HealthCheck",
//	   "searchString": "",
//	   "connectionTimeout": 5,
//	   "expectedStatusCode": 200,
//	   "userAgent": "IE",
//	   "note": "",
//	   "runTraceroute": "DISABLED",
//	   "scheduleInterval": "NONE",
//	   "sslPolicy": "IGNORE",
//	   "userId": 300003895,
//	   "interval": "ONEMINUTE",
//	   "monitorIntervalPolicy": "PARALLEL",
//	   "checkSites": [
//	     15
//	   ],
//	   "notificationGroups": [
//	     41680,
//	     41675
//	   ],
//	   "scheduleId": 0,
//	   "notificationReportTimeout": 1440,
//	   "verificationPolicy": "SIMPLE"
//	 }
type SonarHTTPCheck struct {
	Name                      string `json:"name"`
	Host                      string `json:"host"`
	Port                      int    `json:"port"`
	ProtocolType              string `json:"protocolType"`
	IPVersion                 string `json:"ipVersion"`
	FQDN                      string `json:"fqdn"`
	Path                      string `json:"path"`
	SearchString              string `json:"searchString"`
	ConnectionTimeout         int    `json:"connectionTimeout"`
	ExpectedStatusCode        int    `json:"expectedStatusCode"`
	UserAgent                 string `json:"userAgent"`
	Note                      string `json:"note"`
	RunTraceroute             string `json:"runTraceroute"`
	ScheduleInterval          string `json:"scheduleInterval"`
	SSLPolicy                 string `json:"sslPolicy"`
	UserID                    int    `json:"userId"`
	Interval                  string `json:"interval"`
	MonitorIntervalPolicy     string `json:"monitorIntervalPolicy"`
	CheckSites                []int  `json:"checkSites"`
	NotificationGroups        []int  `json:"notificationGroups"`
	ScheduleID                int    `json:"scheduleId"`
	NotificationReportTimeout int    `json:"notificationReportTimeout"`
	VerificationPolicy        string `json:"verificationPolicy"`
}
