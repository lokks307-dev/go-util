package mt

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var RegExpTelephone *regexp.Regexp
var RegExpMobilePhoneNum *regexp.Regexp
var RegExpNonNumeric *regexp.Regexp
var RegExpEnglishName *regexp.Regexp
var RegExpWhitespaces *regexp.Regexp
var RegExpThreeFourDigits *regexp.Regexp
var RegExpYmd *regexp.Regexp
var RegExpHms *regexp.Regexp
var RegExpRRN *regexp.Regexp
var RegExpTimeZone *regexp.Regexp

var KoLoc *time.Location
var LocalLoc *time.Location

func init() {
	RegExpTelephone = regexp.MustCompile(`^(02|0[3456][1234]|01[016789])-{0,1}(\d{3,4})-{0,1}(\d{4})$`)
	RegExpMobilePhoneNum = regexp.MustCompile(`^(01[016789])-{0,1}(\d{3,4})-{0,1}(\d{4})$`)
	RegExpNonNumeric = regexp.MustCompile(`([^0-9.]+)`)
	RegExpEnglishName = regexp.MustCompile(`(^[a-zA-Z\s]+)`)
	RegExpWhitespaces = regexp.MustCompile(`\s+`)
	RegExpThreeFourDigits = regexp.MustCompile(`[0-9]{3,4}`)
	RegExpYmd = regexp.MustCompile(`(\d{4})-{0,1}(0[1-9]|10|11|12)-{0,1}([0-2][1-9]|30|31)`)
	RegExpHms = regexp.MustCompile(`([0-1][0-9]|2[0-3]):{0,1}(0[0-9]|[1-5][0-9]):{0,1}(0[0-9]|[1-5][0-9])`)
	RegExpRRN = regexp.MustCompile(`(\d{6})-{0,1}(\d{7})`)
	RegExpTimeZone = regexp.MustCompile(`(\+{0,1}|-)([0-9]|0[0-9]|1[0-3]):{0,1}(00|15|30|45)`)
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

func IsThreeFourDigits(v string) bool {
	return RegExpThreeFourDigits.MatchString(v)
}

func MakePath(base string, dd ...string) string {
	format := strings.Repeat("/%s", len(dd)+1)
	args := make([]interface{}, 0)
	args = append(args, url.QueryEscape(base))
	for idx := range dd {
		args = append(args, url.QueryEscape(dd[idx]))
	}
	return fmt.Sprintf(format, args...)
}

func MakePathForCareEase(base string, dd ...string) string {
	format := strings.Repeat("/%s", len(dd)+1)
	args := make([]interface{}, 0)
	args = append(args, base)
	for idx := range dd {
		args = append(args, dd[idx])
	}
	return fmt.Sprintf(format, args...)
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

func WithInTime(b, t, m int64) bool {
	return AbsInt64(b-t) <= m
}

func NoneEmptyMatch(s string, c string) bool {
	if s == "" || c == "" {
		return false
	}

	return (s == c)
}

func GetUpdatedString(olds, news string) string {
	newstr := strings.TrimSpace(news)
	if newstr != "" && newstr != olds {
		return newstr
	}

	return olds
}

func IsOptionSet(opt []bool) bool {
	return len(opt) > 0 && opt[0]
}

func B64toHex(v string) string {
	p, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return ""
	}
	h := hex.EncodeToString(p)
	return h
}

func IsSameNonEmptyString(v1, v2 string) bool {
	if v1 == "" && v2 == "" {
		return false
	}

	return v1 == v2
}

func FindIndexStringSlice(slice []string, item string) int {
	for i := range slice {
		if slice[i] == item {
			return i
		}
	}
	return -1
}

func FindIndexInt64Slice(slice []int64, item int64) int {
	for i := range slice {
		if slice[i] == item {
			return i
		}
	}
	return -1
}

func FiltOutInt64Slice(slice []int64, item int64) []int64 {
	nslice := make([]int64, 0)
	for i := range slice {
		if slice[i] != item {
			nslice = append(nslice, slice[i])
		}
	}

	return nslice
}
