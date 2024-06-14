package mt

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var KoLoc *time.Location
var LocalLoc *time.Location

func init() {
	KoLoc = time.FixedZone("", 9*60*60)

	var err error
	LocalLoc, err = time.LoadLocation("Local")
	if err != nil {
		LocalLoc = time.UTC
	}
}
func SetLocalLocToKoLoc() {
	LocalLoc = KoLoc
}

func GetTime(ts int64) time.Time {
	return time.Unix(ts, 0)
}

func GetNowYYYYMMDD(loc *time.Location) string {
	if loc == nil {
		loc = LocalLoc
	}
	return time.Now().In(loc).Format("20060102")
}

func GetNowYYYYMMDDAsInt64(loc *time.Location) int64 {
	if loc == nil {
		loc = LocalLoc
	}

	t := GetNowYYYYMMDD(loc)
	if ival, err := strconv.ParseInt(t, 10, 64); err != nil {
		return 0
	} else {
		return ival
	}
}

func TimeToDayInt(t time.Time, loc *time.Location) int {
	if loc == nil {
		loc = LocalLoc
	}
	s := t.In(loc).Format("20060102")
	return ToInt(s)
}

func TimeToMonthInt(t time.Time, loc *time.Location) int {
	if loc == nil {
		loc = LocalLoc
	}
	s := t.In(loc).Format("200601")
	return ToInt(s)
}

func UnixToDayInt(unixt int64, loc *time.Location) int {
	if loc == nil {
		loc = LocalLoc
	}
	return TimeToDayInt(time.Unix(unixt, 0), loc)
}

func UnixToMinInt(unixt int64, loc *time.Location) int64 {
	if loc == nil {
		loc = LocalLoc
	}
	return TimeToMinInt(time.Unix(unixt, 0), loc)
}

func UnixToMonthInt(unixt int64, loc *time.Location) int {
	if loc == nil {
		loc = LocalLoc
	}
	return TimeToMonthInt(time.Unix(unixt, 0), loc)
}

func DayIntToTime[T int | int32 | int64](tdint T, loc *time.Location) (time.Time, error) {
	if loc == nil {
		loc = LocalLoc
	}

	return time.ParseInLocation("20060102", strconv.Itoa(int(tdint)), loc)
}

func MinIntToTime(mdint int64, loc *time.Location) (time.Time, error) {
	if loc == nil {
		loc = LocalLoc
	}
	return time.ParseInLocation("200601021504", strconv.FormatInt(mdint, 10), loc)
}

func TimeToMinInt(t time.Time, loc *time.Location) int64 {
	if loc == nil {
		loc = LocalLoc
	}

	mint, err := strconv.ParseInt(t.Format("200601021504"), 10, 64)
	if err != nil {
		return 0
	}
	return mint
}

func DayIntStrToTime(tdint string, loc *time.Location) (time.Time, error) {
	if loc == nil {
		loc = LocalLoc
	}
	return time.ParseInLocation("20060102", tdint, loc)
}

func TimeToHour(t time.Time, loc *time.Location) int {
	if loc == nil {
		loc = LocalLoc
	}
	return t.In(loc).Hour()
}

func TimeToDayMinute(t time.Time, loc *time.Location) int {
	if loc == nil {
		loc = LocalLoc
	}
	return t.In(loc).Hour()*60 + t.In(loc).Minute()
}

func GetNowHHMM(loc *time.Location) string {
	if loc == nil {
		loc = LocalLoc
	}
	return time.Now().In(loc).Format("1504")
}

func IsSameDate(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.In(KoLoc).Date()
	y2, m2, d2 := date2.In(KoLoc).Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func IsToday(date1 time.Time) bool {
	now := time.Now().In(KoLoc)
	return IsSameDate(date1, now)
}

func BeginOfDay(t time.Time, loc *time.Location) time.Time {
	if loc == nil {
		loc = LocalLoc
	}

	year, month, day := t.In(loc).Date()
	return time.Date(year, month, day, 0, 0, 0, 0, loc)
}

func EndOfDay(t time.Time, loc *time.Location) time.Time {
	if loc == nil {
		loc = LocalLoc
	}

	year, month, day := t.In(loc).Date()
	return time.Date(year, month, day, 23, 59, 59, 0, loc)
}

func BeginOfToday(loc *time.Location) time.Time {
	if loc == nil {
		loc = LocalLoc
	}

	return BeginOfDay(time.Now(), loc)
}

func EndOfToday(loc *time.Location) time.Time {
	if loc == nil {
		loc = LocalLoc
	}

	return EndOfDay(time.Now(), loc)
}

func BeginOfDayUnix(ts int64, loc *time.Location) int64 {
	if loc == nil {
		loc = LocalLoc
	}

	ttime := time.Unix(ts, 0)
	return BeginOfDay(ttime, loc).Unix()
}

func EndOfDayUnix(ts int64, loc *time.Location) int64 {
	if loc == nil {
		loc = LocalLoc
	}

	ttime := time.Unix(ts, 0)
	return EndOfDay(ttime, loc).Unix()
}

func GetNakedYmd(ymd string) string {
	found := RegExpYmd.FindAllStringSubmatch(ymd, -1)
	if len(found) > 0 && len(found[0]) >= 4 {
		return found[0][1] + found[0][2] + found[0][3]
	}

	return ""
}

func HmsTo(hhmmss string) (int, int, int) {
	if len(hhmmss) < 6 {
		hhmmss += ":00"
	}

	found := RegExpHms.FindAllStringSubmatch(hhmmss, -1)
	if len(found) > 0 && len(found[0]) >= 4 {
		return ToInt(found[0][1]), ToInt(found[0][2]), ToInt(found[0][3])
	}

	return 0, 0, 0
}

func HmsToSec(hhmmss string) int {
	h, m, s := HmsTo(hhmmss)
	return h*60*60 + m*60 + s
}

func SecToHms[T int | int32 | int64 | uint | uint32 | uint64](sec T) string {
	h, m, s := SecToHMS(sec)

	if s == 0 {
		return fmt.Sprintf("%02d:%02d", h, m)
	}

	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func SecToHMS[T int | int32 | int64 | uint | uint32 | uint64](s T) (int, int, int) {
	h, hq := Divide(s, 3600)
	m, s := Divide(hq, 60)
	return int(h % 24), int(m), int(s)
}

func IntSliceToHmsSlice(hm []int) []string {
	ret := make([]string, 0)
	for _, s := range hm {
		ret = append(ret, SecToHms(s))
	}

	return ret
}

func HmsSliceToIntSlice(ss []string) []int {
	ret := make([]int, 0)
	for _, s := range ss {
		ret = append(ret, HmsToSec(s))
	}

	return ret
}

func IsTimeZone(tz string) bool {
	return RegExpTimeZone.Match([]byte(tz))
}

func TimeZoneToLocation(tz string, def ...*time.Location) *time.Location {
	found := RegExpTimeZone.FindAllStringSubmatch(tz, -1)
	if len(found) > 0 && len(found[0]) >= 4 {
		offset := ToInt(found[0][2])*60*60 + ToInt(found[0][3])*60
		if found[0][1] == "-" {
			offset = -1 * offset
		}

		return time.FixedZone(tz, offset)
	}

	if len(def) >= 0 {
		return def[0]
	}

	return time.UTC
}

func LocationToTimeZone(loc *time.Location) string {
	if loc == nil {
		return "+00:00"
	}
	return time.Now().In(loc).Format("-07:00")
}

func WeekDayIntToSlice(w int) []int {
	s := make([]int, 0)
	for i := 0; i < 7; i++ {
		if w/Pow(10, 7-i-1) > 0 {
			s = append(s, i)
		}
		w = w % Pow(10, 7-i-1)
	}
	return s
}

func TimeSliceToInt64Slice(tt []time.Time) []int64 {
	ts := make([]int64, 0)
	for _, t := range tt {
		ts = append(ts, t.Unix())
	}

	return ts
}

func ToInterfaceSlice[T any](tt []T) []interface{} {
	si := make([]interface{}, 0)
	for _, t := range tt {
		si = append(si, t)
	}

	return si
}

func GetDaysBetween(after, before int64) int {
	days := 1

	tgap := AbsInt64(before - after)
	if tgap >= ONE_DAY_SEC {
		days_, q := Divide(tgap, ONE_DAY_SEC)
		if q != 0 {
			days_++
		}

		days = int(days_)
	}

	return days
}

func TimeDayRangeFunc(a, b time.Time, x func(time.Time)) {
	for td := a; td.Before(b); td = td.AddDate(0, 0, 1) {
		x(td)
	}
}

func WithInTime(b, t, m int64) bool {
	return AbsInt64(b-t) <= m
}

func SlotNoToTimeString(slotNo int) string {
	hourInt := int(slotNo / 2)
	minInt := 0
	if slotNo%2 != 0 {
		minInt = 30
	}

	return fmt.Sprintf("%02d:%02d", hourInt, minInt)
}

func TimeStringToSlotNo(hhmm string) int {
	splits := strings.Split(hhmm, ":")
	if len(splits) != 2 {
		return 0
	}

	hourInt, err := strconv.Atoi(splits[0])
	if err != nil {
		return 0
	}

	carry := 0
	if splits[1] == "30" {
		carry = 1
	}

	return hourInt*2 + carry
}

func TimestampToSlotNo(stamp int64, loc *time.Location) int {
	t := time.Unix(stamp, 0).In(loc)
	if t.Minute() >= 30 {
		return t.Hour()*2 + 1
	}
	return t.Hour() * 2
}

func GetBirthdateFromAge(age int, onTime time.Time) time.Time {
	return onTime.AddDate(-age, 0, 0)
}

func GetDoBIntFromAge(age int, onTime time.Time) int {
	t := GetBirthdateFromAge(age, onTime)
	return TimeToDayInt(t, onTime.Location())
}

func GetMinutesBetweenHHmmInt(start, end int) int {
	diff := end - start
	if diff < 0 {
		return 0
	}

	return (diff/100)*60 + (diff % 100)
}

func DateTimeIntToTime[T int | int32 | int64](dateInt, HHmmInt T, loc *time.Location) (time.Time, error) {
	dateConv, err := DayIntToTime(dateInt, loc)
	if err != nil {
		return dateConv, err
	}

	hour := time.Duration(HHmmInt / 100)
	minute := time.Duration(HHmmInt % 100)

	return dateConv.Add(hour*time.Hour + minute*time.Minute), nil
}

func TimeToTimeInt(t time.Time, loc *time.Location) int {
	if loc == nil {
		loc = LocalLoc
	}

	return t.In(loc).Hour()*100 + t.In(loc).Minute()
}
