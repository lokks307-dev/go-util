package mt

import (
	"regexp"
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
}

func IsThreeFourDigits(v string) bool {
	return RegExpThreeFourDigits.MatchString(v)
}
