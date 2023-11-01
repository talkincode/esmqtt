/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package timeutil

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/spf13/cast"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	Datetime14Layout      = "20060102150405"
	Datetime8Layout       = "20060102"
	Datetime6Layout       = "200601"
	YYYYMMDDHHMMSS_LAYOUT = "2006-01-02 15:04:05"
	YYYYMMDDHHMM_LAYOUT   = "2006-01-02 15:04"
	YYYYMMDD_LAYOUT       = "2006-01-02"
)

var (
	ShangHaiLOC, _ = time.LoadLocation("Asia/Shanghai")
	EmptyTime, _   = time.Parse("2006-01-02 15:04:05 Z0700 MST", "1979-11-30 00:00:00 +0000 GMT")
)

// yyyy-MM-dd hh:mm:ss 年-月-日 时:分:秒
func FmtDatetimeString(t time.Time) string {
	return t.Format(YYYYMMDDHHMMSS_LAYOUT)
}

// yyyy-MM-dd hh:mm 年-月-日 时:分
func FmtDatetimeMString(t time.Time) string {
	return t.Format(YYYYMMDDHHMM_LAYOUT)
}

// yy-MM-dd 年-月-日
func FmtDateString(t time.Time) string {
	return t.Format(YYYYMMDD_LAYOUT)
}

// yyyyMMddhhmmss 年月日时分秒
func FmtDatetime14String(t time.Time) string {
	return t.Format(Datetime14Layout)
}

// yyyyMMdd 年月日
func FmtDatetime8String(t time.Time) string {
	return t.Format(Datetime8Layout)
}

// yyyyMM  年月
func FmtDatetime6String(t time.Time) string {
	return t.Format(Datetime6Layout)
}

func ComputeEndTime(times int, unit string) time.Time {
	ctime := time.Now()
	switch unit {
	case "second":
		return ctime.Add(time.Second * time.Duration(times))
	case "minute":
		return ctime.Add(time.Minute * time.Duration(times))
	case "hour":
		return ctime.Add(time.Hour * time.Duration(times))
	case "day":
		return ctime.Add(time.Hour * 24 * time.Duration(times))
	case "week":
		return ctime.Add(time.Hour * 24 * 7 * time.Duration(times))
	case "month":
		return ctime.Add(time.Hour * 24 * 30 * time.Duration(times))
	case "year":
		return ctime.Add(time.Hour * 24 * 365 * time.Duration(times))
	default:
		return ctime
	}
}

type LocalTime time.Time

func ParseLocalTime(value string) (*LocalTime, error) {
	ts, err := dateparse.ParseAny(value)
	if err != nil {
		return nil, err
	}
	lt := LocalTime(ts)
	return &lt, nil
}

func (t *LocalTime) IsZero() bool {
	return time.Time(*t).IsZero()
}

func NewLocalTime() LocalTime {
	return LocalTime(time.Now())
}

func Now() LocalTime {
	return LocalTime(time.Now())
}

func (t *LocalTime) UnmarshalParam(src string) error {
	ts, err := dateparse.ParseAny(strings.Replace(src, "\"", "", -1))
	if err != nil {
		return err
	}
	*t = LocalTime(ts)
	return nil
}

func (t *LocalTime) MarshalParam() string {
	lt := time.Time(*t)
	return lt.Format(time.DateTime)
}

func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	now, err := dateparse.ParseAny(strings.Replace(string(data), "\"", "", -1))
	if err != nil {
		return
	}
	*t = LocalTime(now)
	return
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Format(time.DateTime))
}

func (t *LocalTime) Format(layout string) string {
	lt := time.Time(*t)
	return lt.Format(layout)
}

func (t *LocalTime) Value() time.Time {
	return time.Time(*t)
}

// MarshalMsgpack 实现了 Msgpack 编码接口
func (t LocalTime) MarshalMsgpack() ([]byte, error) {
	return msgpack.Marshal(t.Value())
}

// UnmarshalMsgpack 实现了 Msgpack 解码接口
func (t *LocalTime) UnmarshalMsgpack(data []byte) error {
	var tt time.Time
	err := msgpack.Unmarshal(data, &tt)
	if err != nil {
		return err
	}
	*t = LocalTime(tt)
	return nil
}

func ParseAny(timestr string) (LocalTime, error) {
	ts, err := dateparse.ParseLocal(strings.Replace(timestr, "\"", "", -1))
	if err != nil {
		return LocalTime(time.Time{}), err
	}
	return LocalTime(ts), nil
}

func ParseTimeDesc(timestr string) string {
	switch {
	case strings.HasPrefix(timestr, "now-") && strings.HasSuffix(timestr, "hour"):
		v := cast.ToInt(timestr[4 : len(timestr)-4])
		return time.Now().Add(time.Hour * time.Duration(v*-1)).Format(time.RFC3339)
	case strings.HasPrefix(timestr, "now-") && strings.HasSuffix(timestr, "min"):
		v := cast.ToInt(timestr[4 : len(timestr)-3])
		return time.Now().Add(time.Minute * time.Duration(v*-1)).Format(time.RFC3339)
	case strings.HasPrefix(timestr, "now-") && strings.HasSuffix(timestr, "sec"):
		v := cast.ToInt(timestr[4 : len(timestr)-3])
		return time.Now().Add(time.Second * time.Duration(v*-1)).Format(time.RFC3339)
	case strings.HasPrefix(timestr, "now-") && strings.HasSuffix(timestr, "day"):
		v := cast.ToInt(timestr[4 : len(timestr)-3])
		return time.Now().Add(time.Hour * 24 * time.Duration(v*-1)).Format(time.RFC3339)
	case timestr == "now":
		return time.Now().Format(time.RFC3339)
	default:
		return time.Now().Format(time.RFC3339)
	}
}

func TimeBucket(t time.Time, interval time.Duration) time.Time {
	return t.Truncate(interval)
}

func SplitDaterange(start, end string) ([]time.Time, error) {
	const layout = "2006-01-02 15:04:05" // Go 的特定时间格式
	startTime, err := time.Parse(layout, start)
	if err != nil {
		return nil, err
	}

	endTime, err := time.Parse(layout, end)
	if err != nil {
		return nil, err
	}

	var dates []time.Time
	for d := startTime; d.Before(endTime) || d.Equal(endTime); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d)
	}

	return dates, nil
}
