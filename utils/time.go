package utils

import "time"

func GetJktTime() (res *time.Time, err error) {
	tz, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return nil, err
	}

	now := time.Now().In(tz)

	return &now, nil
}
