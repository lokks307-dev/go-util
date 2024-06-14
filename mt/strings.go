package mt

import (
	"strconv"
	"strings"

	"github.com/volatiletech/null/v8"
)

func SubStr(v string, length int) string {
	if len(v) > length {
		return v[0:length]
	}
	return v
}

// SubStrEx("1234567890",-2,0) = "90"

func SubStrEx(a string, begin int, length int) string {
	ra := []rune(a)
	if len(ra) == 0 {
		return ""
	}

	if begin < 0 {
		begin = len(ra) + begin
	}
	if begin < 0 {
		begin = 0
	}

	if begin > len(ra) {
		return ""
	}

	end := begin + length
	if length <= 0 {
		end = len(ra)
	} else {
		if end > len(ra) {
			end = len(ra)
		}
	}

	return string(ra[begin:end])
}

func GetNakedTelephone(untTel string) string {
	untTel = strings.ReplaceAll(untTel, " ", "")
	untTel = strings.ReplaceAll(untTel, "-", "")
	untTel = strings.ReplaceAll(untTel, "(", "")
	untTel = strings.ReplaceAll(untTel, ")", "")
	untTel = strings.ReplaceAll(untTel, "+821", "01") // +821X => 01X
	return string(RegExpNonNumeric.ReplaceAll([]byte(untTel), []byte{}))
}

func GetHypenedTelephone(unTel string) string {
	return GetDividerTelephone(unTel, "-")
}

func GetSpacedTelephone(unTel string) string {
	return GetDividerTelephone(unTel, " ")
}

func GetDottedTelephone(unTel string) string {
	return GetDividerTelephone(unTel, ".")
}

func GetDividerTelephone(unTel string, divider string) string {
	if RegExpTelephone == nil {
		return unTel
	}

	ntel := GetNakedTelephone(unTel) // to remove international code and brackets

	found := RegExpTelephone.FindAllStringSubmatch(ntel, -1)
	if len(found) > 0 && len(found[0]) >= 4 {
		return found[0][1] + divider + found[0][2] + divider + found[0][3]
	}

	return unTel
}

func GetNakedName(untName string) string {
	trimName := strings.TrimSpace(untName)
	if IsEnglishName(trimName) {
		return RegExpWhitespaces.ReplaceAllString(trimName, " ")
	} else {
		return strings.ReplaceAll(trimName, " ", "")
	}
}

func IsSameName(aName, bName string) bool {
	aNameTrim := strings.ReplaceAll(aName, " ", "")
	bNameTrim := strings.ReplaceAll(bName, " ", "")

	return strings.EqualFold(aNameTrim, GetNakedName(bNameTrim))
}

func UTF8toEUCKR(s string) []byte {
	return ToCP949([]byte(s))
}

func EscapeSingle(s string) string {
	return strings.ReplaceAll(s, "'", "\\'")
}

func EscapeDouble(s string) string {
	return strings.ReplaceAll(s, "\"", "\\\"")
}

func SortOrderCompString(s1, s2 string) bool {
	s1rune := []rune(s1)
	s2rune := []rune(s2)

	lenToInspect := len(s1rune)
	if len(s2rune) < lenToInspect {
		lenToInspect = len(s2rune)
	}

	for k := 0; k < lenToInspect; k++ {
		if s1rune[k] == s2rune[k] {
			continue
		}

		if s1rune[k] > s2rune[k] {
			return true
		}

		return false
	}

	return len(s1rune) > len(s2rune)
}

func IsEnglishName(s string) bool {
	return RegExpEnglishName.Match([]byte(s))
}

func IsMobilePhoneNum(v string) bool {
	return RegExpMobilePhoneNum.Match([]byte(v))
}

func WrapCrmIpAddrHttp(ipAddr string) string {
	if ipAddr == "" {
		return ""
	}

	return "http://" + ipAddr + ":8080" // Fixed Port num
}

func WrapCrmIpAddrSocketNoScheme(ipAddr string) string {
	if ipAddr == "" {
		return ""
	}

	return ipAddr + ":8081"
}

func ToInt(s string, dd ...int) int {
	var defv int
	if len(dd) > 0 {
		defv = dd[0]
	}
	if ival, err := strconv.ParseInt(s, 10, 64); err != nil {
		return defv
	} else {
		return int(ival)
	}
}

func ToInt32(s string, dd ...int32) int32 {
	var defv int32
	if len(dd) > 0 {
		defv = dd[0]
	}
	if ival, err := strconv.ParseInt(s, 10, 64); err != nil {
		return defv
	} else {
		return int32(ival)
	}
}

func ToInt64(s string, dd ...int64) int64 {
	var defv int64
	if len(dd) > 0 {
		defv = dd[0]
	}
	if ival, err := strconv.ParseInt(s, 10, 64); err != nil {
		return defv
	} else {
		return ival
	}
}

func IsEmptyStr(v interface{}) bool {
	switch t := v.(type) {
	case *string:
		if t != nil && len(*t) > 0 {
			return false
		}
	case string:
		if t != "" {
			return false
		}
	case null.String:
		if t.Valid && t.String != "" {
			return false
		}
	default:
		return false
	}

	return true
}

func IsContain(dd string, vv ...string) bool {
	for _, v := range vv {
		if strings.Contains(dd, v) {
			return true
		}
	}

	return false
}

func ToStr(s interface{}) string {
	switch t := s.(type) {
	case string:
		return t
	case *string:
		if t == nil {
			return ""
		}
		return *t
	case null.String:
		if t.Valid {
			return t.String
		}
		return ""
	}
	return ""
}

func TrueOr[T any](t bool, a T, b T) T {
	if t {
		return a
	}
	return b
}

func NoneEmptyMatch(s string, c string) bool {
	if s == "" || c == "" {
		return false
	}

	return (s == c)
}

func IsSameNonEmptyString(v1, v2 string) bool {
	if v1 == "" && v2 == "" {
		return false
	}

	return v1 == v2
}

func GetUpdatedString(olds, news string) string {
	newstr := strings.TrimSpace(news)
	if newstr != "" && newstr != olds {
		return newstr
	}

	return olds
}

func Unescape(raw string) string {
	return strings.Replace(raw, `\\u`, `\u`, -1)
}

func UnescapeUnicodeCharactersInJSON(raw string) string {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(raw), `\\u`, `\u`, -1))
	if err != nil {
		return ""
	}
	return str
}

func EmptyOr(aa ...interface{}) string {
	r := ""

	for _, a := range aa {
		switch t := a.(type) {
		case null.String:
			if t.Valid {
				r = t.String
			}
		case string:
			r = t
		case *string:
			if t != nil {
				r = *t
			}
		}

		if r != "" {
			return r
		}

	}

	return ""
}

func StringToInt64Slice(v []string) []int64 {
	vv := make([]int64, 0)

	for i := range v {
		if conv, err := strconv.ParseInt(v[i], 10, 64); err == nil && conv > 0 {
			vv = append(vv, conv)
		}
	}

	return vv
}

func ReplaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}
