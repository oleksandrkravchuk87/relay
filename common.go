package relay

import (
	"strings"
)

//OPBLANK : OPBLANK is no operation specified.
const OPBLANK = ""

//OPAND : OPAND is and condition to join multile filer criteria.
const OPAND = "&"

//OPOR : OPOR is or condition to join multile filer criteria.
const OPOR = "|"

//SEPENTRY : SEPENTRY is entry seperator.
const SEPENTRY = ","

//SEPDATA : SEPDATA is data seperator.
const SEPDATA = ":"

//INDKEY : INDKEY is key identifier.
const INDKEY = "key"

//INDCOL : INDCOL is column identifier.
const INDCOL = "column"

//INDOP : INDCOP is operation identifier.
const INDOP = "op"

//ORDASC : ORDASC is ascending Order.
const ORDASC = "ASC"

//ORDDESC : ORDDESC is descending Order.
const ORDDESC = "DESC"

//NILQUERY : NILQUERY is blank query.
const NILQUERY = "NIL"

//SPACE : SPACE.
const SPACE = " "

//STARTBRACE : STARTBRACE is starting brace.
const STARTBRACE = "{"

//ENDBRACE : ENDBRACE is ending brace.
const ENDBRACE = "}"

//KeyCheck : Profile ID needs column
const ProfileKeyCheck = "profileScoreData.profileID"


//CleanConditions : Function to Clean Sort Filter Conditions
func CleanConditions(str string) string {
	var sRet string
	bEscape := false
	for i := 0; i < len(str); i++ {
		if str[i:i+1] == SEPDATA {
			bEscape = true
		} else if str[i:i+1] == SEPENTRY || str[i:i+1] == ENDBRACE {
			bEscape = false
		}
 
		if str[i:i+1] != SPACE || bEscape {
			sRet = sRet + str[i:i+1]
		}
	}

	return sRet
}

//GetSubQueries : Function to Get Sub Queries
func GetSubQueries(str, s1, s2 string) []string {
	s := str
	var Op string
	var ret []string
	breakLoop := false

	for i := 0; ; i++ {
		if s[:1] == OPAND || s[:1] == OPOR {
			Op = s[:1]
			s = s[1:]
		}

		i1 := strings.Index(s, s1)
		i2 := strings.Index(s, s2)

		var curInd int
		curInd = -1
		if i1 > -1 && i2 > -1 {
			if i1 < i2 {
				curInd = i1
			} else {
				curInd = i2
			}
		} else if i1 > -1 {
			curInd = i1
		} else if i2 > -1 {
			curInd = i2
		}

		if curInd == -1 {
			curInd = len(s)
			breakLoop = true
		}
		ret = append(ret, Op+s[:curInd])
		s = s[curInd:]

		if breakLoop {
			break
		}
	}

	return ret
}

//GetQueryDetails : Function to Get Query Details
func GetQueryDetails(Query string) (Key string, Column string, Op string) {
	bOPFound := false
	var FSQ3 []string
	var FSQ1 []string
	var FSQ2 []string

	FQR := Query[1 : len(Query)-1]
	FQ := strings.Split(string(FQR), SEPENTRY)

	if len(strings.Split(FQ[0], SEPDATA)) > 2 {
		FSQ1 = strings.SplitN(FQ[0], SEPDATA, 2)
	} else {
		FSQ1 = strings.Split(FQ[0], SEPDATA)
	}

	if len(strings.Split(FQ[1], SEPDATA)) > 2 {
		FSQ2 = strings.SplitN(FQ[1], SEPDATA, 2)
	} else {
		FSQ2 = strings.Split(FQ[1], SEPDATA)
	}

	if len(FQ) > 2 {
		bOPFound = true
		FSQ3 = strings.Split(FQ[2], SEPDATA)
	}

	if strings.ToLower(FSQ1[0]) == INDKEY {
		Key = FSQ1[1]
	}
	if strings.ToLower(FSQ2[0]) == INDKEY {
		Key = FSQ2[1]
	}
	if bOPFound {
		if strings.ToLower(FSQ3[0]) == INDKEY {
			Key = FSQ3[1]
		}
	}

	if strings.ToLower(FSQ1[0]) == INDCOL {
		Column = FSQ1[1]
	}
	if strings.ToLower(FSQ2[0]) == INDCOL {
		Column = FSQ2[1]
	}
	if bOPFound {
		if strings.ToLower(FSQ3[0]) == INDCOL {
			Column = FSQ3[1]
		}
	}

	if bOPFound {
		if strings.ToLower(FSQ1[0]) == INDOP {
			Op = FSQ1[1]
		}
		if strings.ToLower(FSQ2[0]) == INDOP {
			Op = FSQ2[1]
		}
		if strings.ToLower(FSQ3[0]) == INDOP {
			Op = FSQ3[1]
		}
	}
	return
}

//StringLessOp : Function to check Strings less operation
func StringLessOp(strFirst, strSecond string) bool {
	lenFirst := len(strFirst)
	lenSecond := len(strSecond)

	minLength := lenFirst
	if lenFirst > lenSecond {
		minLength = lenSecond
	}

	for i := 0; i < minLength; i++ {
		cFirst := strFirst[i]
		cSecond := strSecond[i]
		if cFirst > 96 && cFirst < 123 {
			cFirst = cFirst ^ 32
		}
		if cSecond > 96 && cSecond < 123 {
			cSecond = cSecond ^ 32
		}

		if cFirst != cSecond {
			return cFirst < cSecond
		}
	}

	if lenFirst < lenSecond {
		return true
	}
	return false
}
