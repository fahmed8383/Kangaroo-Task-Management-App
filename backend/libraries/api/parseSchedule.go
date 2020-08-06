package api

import "encoding/json"

// ScheduleInfo holds the schedule and sorting method info received from the frontend in a struct
type ScheduleInfo struct {
	Schedule      json.RawMessage `json:"schedule"`
	SortingMethod string          `json:"sort"`
}

// ParseScheduleInfo unmarshalls the byte schedule data into a golang struct
func ParseScheduleInfo(data []byte) (ScheduleInfo, error) {
	var info ScheduleInfo
	err := json.Unmarshal(data, &info)
	return info, err
}

// SetScheduleResponse repacks schedule and sorting method into ScheduleInfo struct and marshalls them into bytes
func SetScheduleResponse(schedule json.RawMessage, sorting string) ([]byte, error) {
	response := ScheduleInfo{schedule, sorting}
	res, err := json.Marshal(response)
	return res, err
}
