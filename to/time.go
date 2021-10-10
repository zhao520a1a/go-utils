/**
 * author: suchenghao10349
 */
package to

import "time"

const defaultTimeFormat = "2006-01-02 15:04:05"

type timeOption struct {
	timestamp int64
	nsec      int64
	format    string
}

func Nsec(nsec int64) func(option *timeOption) {
	return func(option *timeOption) {
		option.nsec = nsec
	}
}

func Format(format string) func(option *timeOption) {
	return func(option *timeOption) {
		option.format = format
	}
}

// timestamp to string default format:2006-01-02 15:04:05
func TimestampString(timestamp int64, params ...func(option *timeOption)) string {
	args := &timeOption{
		timestamp: timestamp,
		nsec:      0,
		format:    defaultTimeFormat,
	}

	for _, param := range params {
		param(args)
	}

	return time.Unix(args.timestamp, args.nsec).Format(args.format)
}
