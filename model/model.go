package model

type ConfigSettings struct{
	Servers []string	`jsom:"servers"`
	RefreshTime int `json:"refresh_time"`
	TimeDuration int64 `json:"time_duration"`
}