package mt

import (
	"strconv"
	"strings"
)

const (
	SEX_MALE   = 1
	SEX_FEMALE = 2
	SEX_OTHERS = 3
)

func GetFullBirthdate(yymmdd string, rc int) string {
	if len(yymmdd) > 6 {
		return strings.ReplaceAll(yymmdd, "-", "")
	}

	if rc == 0 || rc == 9 {
		return "18" + yymmdd
	}

	if rc == 1 || rc == 2 || rc == 5 || rc == 6 {
		return "19" + yymmdd
	}

	if rc == 3 || rc == 4 || rc == 7 || rc == 8 {
		return "20" + yymmdd
	}

	return "18000101"
}

func GetSex(rc int) int {
	if rc%2 == 1 {
		return SEX_MALE
	}

	if rc%2 == 0 {
		return SEX_FEMALE
	}

	return SEX_OTHERS
}

func TrimRRN(input string) string {
	if input == "" || len(input) < 13 {
		return ""
	}

	return strings.ReplaceAll(input, "-", "")
}

func GetSexFromRRN(input string) int {
	if len(input) < 13 {
		return SEX_OTHERS
	}

	if xCode, err := strconv.Atoi(string(input[6])); err != nil {
		if xCode, err = strconv.Atoi(string(input[7])); err == nil {
			return GetSex(xCode)
		}
	} else {
		return GetSex(xCode)
	}

	return SEX_OTHERS
}

func GetFullBirthDateFromRRN(input string) string {
	if len(input) < 13 {
		return ""
	}
	var xCode int
	var err error

	if xCode, err = strconv.Atoi(string(input[6])); err != nil {
		if xCode, err = strconv.Atoi(string(input[7])); err != nil {
			xCode = 1
		}
	}

	return GetFullBirthdate(input[0:6], xCode)
}

func GetPartialRRN(birthYmd, sex, foreginerYn string) string {

	var yy, mm, dd string

	found := RegExpYmd.FindAllStringSubmatch(birthYmd, -1)
	if len(found) > 0 && len(found[0]) >= 4 {
		yy = found[0][1]
		mm = found[0][2]
		dd = found[0][3]
	} else {
		return "000000-0000000"
	}

	RRNQ := make(map[string]string)
	RRNQ["19NM"] = "1"
	RRNQ["19NF"] = "2"
	RRNQ["20NM"] = "3"
	RRNQ["20NF"] = "4"
	RRNQ["19YM"] = "5"
	RRNQ["19YF"] = "6"
	RRNQ["20YM"] = "7"
	RRNQ["20YF"] = "8"

	key := strings.ToUpper(yy[0:2] + foreginerYn + sex)

	if qv, ok := RRNQ[key]; ok {
		return yy[2:4] + mm + dd + "-" + qv + "000000"
	}

	return yy[2:4] + mm + dd + "-0000000"
}

func GetNakedRRN(rrn string) string {
	found := RegExpRRN.FindAllStringSubmatch(rrn, -1)
	if len(found) > 0 && len(found[0]) == 3 {
		return found[0][1] + found[0][2]
	}

	return "000000-0000000"
}

func GetMaskedRRN(rrn string) string {
	found := RegExpRRN.FindAllStringSubmatch(rrn, -1)
	if len(found) > 0 && len(found[0]) == 3 {
		return found[0][1] + "-" + (found[0][2])[0:1] + "******"
	}

	return "******-*******"
}

func GetCuttedRRN(rrn string) string {
	found := RegExpRRN.FindAllStringSubmatch(rrn, -1)
	if len(found) > 0 && len(found[0]) == 3 {
		return found[0][1] + "-" + (found[0][2])[0:1]
	}

	return "******-*"
}
