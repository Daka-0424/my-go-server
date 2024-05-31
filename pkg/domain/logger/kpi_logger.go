package logger

type KpiLogName string

const (
	KpiLogInstall KpiLogName = "kpi/Install"
	KpiLogLogin   KpiLogName = "kpi/Login"
	KpiLogPayment KpiLogName = "kpi/Payment"
	KpiLogConsume KpiLogName = "kpi/Consume"
	KpiLogEarn    KpiLogName = "kpi/Earn"
)

type IKpi interface {
	// LogEvent logs an event to the logger
	LogEvent(logName KpiLogName, date map[string]interface{})
	// Flush flushes the logger
	Flush() error
}
