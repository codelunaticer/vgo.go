package utilsLogic

import "time"

// FormatShangHaiTime 将时间格式化成上海时间
func FormatShangHaiTime(mtime string) (time.Time, error) {
	shanghaiZone, _ := time.LoadLocation("Asia/Shanghai")
	shangHaiTime, err := time.ParseInLocation("2006-01-02 15:04:05", mtime, shanghaiZone)
	if err != nil {
		return time.Time{}, err
	}
	return shangHaiTime, nil
}

type GetCurrentDayTimeStruct struct {
	StartTime string
	EndTime   string
}

// GetCurrentDayTime 获取当天开始时间跟结束时间
func GetCurrentDayTime() (GetCurrentDayTimeStruct, error) {
	now := time.Now()
	// 00:00:00
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// 23:59:59
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
	return GetCurrentDayTimeStruct{
		StartTime: startOfDay.Format("2006-01-02 15:04:05"),
		EndTime:   endOfDay.Format("2006-01-02 15:04:05"),
	}, nil
}

// MContains 切片中是否存在某个子元素
func MContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
