package status

var Registered = "REGISTERED"
var Invalid = "INVALID"
var Processing = "PROCESSING"
var Processed = "PROCESSED"

// AccrualStatusesOrderStatusesMap
// соответствие статусов пришедших с сервиса акруал с внутренними статусами заказа
var AccrualStatusesOrderStatusesMap = map[string]int{
	Registered: NewID,
	Invalid:    InvalidID,
	Processing: ProcessingID,
	Processed:  ProcessedID,
}
