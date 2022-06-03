package relay

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

//DataSet : DataSet is data type used for filtering data.
type DataSet struct {
	bMatched bool
	CurRec   interface{}
}

//IsColumnNameValid : To validate Column Name
func IsColumnNameValid(Column *string, val reflect.Value) bool {
	ColumnUp := strings.ToUpper(*Column)
	var ColumnFound = false
	for j := 0; j < val.NumField(); j++ {
		if ColumnUp == strings.ToUpper(val.Type().Field(j).Name) {
			*Column = val.Type().Field(j).Name
			ColumnFound = true
			break
		}
	}
	return ColumnFound
}

//GetProfilesSubQueries :  Function to Get Sub Queries with ProfileIDs
func GetProfilesSubQueries(filterConditions string) ([]string, error) {
	subQuery := GetSubQueries(filterConditions, OPAND, OPOR)

	if strings.Contains(filterConditions, ProfileKeyCheck) && len(subQuery) > 1 {
		subQueryLen := len(subQuery)
		//ProfileID filter needs to be at the end
		if !strings.Contains(subQuery[subQueryLen-1], ProfileKeyCheck) {
			if strings.Contains(subQuery[0], ProfileKeyCheck) {
				return nil, errors.New("Profile ID Filter Key cannot be first argument!!!")
			} else {
				for i := range subQuery {
					if strings.Contains(subQuery[i], ProfileKeyCheck) {
						subQuery[subQueryLen-1], subQuery[i] = subQuery[i], subQuery[subQueryLen-1]
						break
					}
				}
			}
		}
	}
	return subQuery, nil
}

//FilterProfiles : Filters Profile data on given filterConditions.
func FilterProfiles(filterConditions string, sortCond string, val reflect.Value, arraySlice []interface{}) ([]interface{}, int, error) {
	var isProfileIDFound = false
	arraySliceRS := make(map[int]DataSet)
	var arraySliceRet []interface{}
	var arraySliceRetZeroScore []interface{}
	var arraySliceRetProNil []interface{}
	var arraySliceRetNotFound []interface{}

	for ind := range arraySlice {
		arraySliceRS[ind] = DataSet{bMatched: false, CurRec: arraySlice[ind]}
	}

	subQuery, err := GetProfilesSubQueries(filterConditions)
	if err != nil {
		return nil, 0, err
	}

	for i := range subQuery {
		var LOP, Key, Column, Op, CurQuery string
		CurQuery = subQuery[i]
		if CurQuery[:1] == OPAND || CurQuery[:1] == OPOR {
			LOP = CurQuery[:1]
			CurQuery = CurQuery[1:]
		} else {
			LOP = OPBLANK
		}
		Key, Column, Op = GetQueryDetails(CurQuery)
		if Key == "" {
			return nil, 0, errors.New("Filter Key not found!!!")
		}
		if Column == "" {
			return nil, 0, errors.New("Filter Column not found!!!")
		}

		var ColumnSlice []string
		if strings.Contains(Column, ".") {
			ColumnSlice = strings.Split(Column, ".")
			Column = ColumnSlice[0]
		}

		if !IsColumnNameValid(&Column, val) {
			return nil, 0, errors.New("Filter [" + Column + "] No such column exist!!!")
		}

		for k := 0; k < len(arraySliceRS); k++ {
			bConsiderMe := false
			if LOP == OPBLANK || LOP == OPOR {
				bConsiderMe = !arraySliceRS[k].bMatched
			} else if LOP == OPAND {
				bConsiderMe = arraySliceRS[k].bMatched
			}

			if bConsiderMe {
				bMatched := false
				CurField := reflect.Indirect(reflect.ValueOf(arraySliceRS[k].CurRec)).FieldByName(Column)

				if isPrimitive(CurField) {
					if bMatched, err = processPrimitive(CurField, Key, Op); err != nil {
						return nil, 0, err
					}

				} else if CurField.Kind() == reflect.Slice {
					bMatched = true
					if reflect.Indirect(reflect.ValueOf(arraySliceRS[k].CurRec)).FieldByName(Column).IsNil() {
						arraySliceRetProNil = append(arraySliceRetProNil, arraySliceRS[k].CurRec)
					} else {

						sliceLen := CurField.Len()
						if !IsColumnNameValid(&ColumnSlice[1], CurField.Index(0)) {
							return nil, 0, errors.New("Filter [" + ColumnSlice[1] + "] No such column exist!!!")
						}
						for lenslice := 0; lenslice < sliceLen; lenslice++ {
							CurFieldVal := CurField.Index(lenslice).FieldByName(ColumnSlice[1])
							if CurFieldVal.String() == Key {
								if lenslice != 0 {
									x, y := CurField.Index(lenslice).Interface(), CurField.Index(0).Interface()
									CurField.Index(lenslice).Set(reflect.ValueOf(y))
									CurField.Index(0).Set(reflect.ValueOf(x))
								}

								isProfileIDFound = true
								if strings.Contains(sortCond, "RPTSCORE") {
									ScoreVal := CurField.Index(0).FieldByName("Score")
									if ScoreVal.Int() == 0 {
										arraySliceRetZeroScore = append(arraySliceRetZeroScore, arraySliceRS[k].CurRec)
									} else {
										arraySliceRet = append(arraySliceRet, arraySliceRS[k].CurRec)
									}
								} else {
									arraySliceRet = append(arraySliceRet, arraySliceRS[k].CurRec)
								}

								break
							} else if lenslice+1 == sliceLen {
								arraySliceRetNotFound = append(arraySliceRetNotFound, arraySliceRS[k].CurRec)
							}
						}
					}
				}
				arraySliceRS[k] = DataSet{bMatched: bMatched, CurRec: arraySliceRS[k].CurRec}
			}
		}
	}
	var arraySliceRetLen = 0
	arraySlice = make([]interface{}, 0)
	if !isProfileIDFound {
		for l := 0; l < len(arraySliceRS); l++ {
			if arraySliceRS[l].bMatched {
				arraySlice = append(arraySlice, arraySliceRS[l].CurRec)
			}
		}
	} else if isProfileIDFound {
		arraySliceRetLen = len(arraySliceRet)
		arraySlice = append(arraySlice, arraySliceRet...)
		arraySlice = append(arraySlice, arraySliceRetZeroScore...)
		arraySlice = append(arraySlice, arraySliceRetNotFound...)
		arraySlice = append(arraySlice, arraySliceRetProNil...)
	}

	return arraySlice, arraySliceRetLen, nil
}

//Filter : Filters data on given filterConditions.
func Filter(filterConditions string, val reflect.Value, arraySlice []interface{}) ([]interface{}, error) {
	var err error
	arraySliceRS := make(map[int]DataSet)
	for ind := range arraySlice {
		arraySliceRS[ind] = DataSet{bMatched: false, CurRec: arraySlice[ind]}
	}

	filterConditions = CleanConditions(filterConditions)
	arraySliceRS, err = markRecords(filterConditions, val, arraySliceRS)
	if err != nil {
		return nil, err
	}

	arraySlice = make([]interface{}, 0)
	for l := 0; l < len(arraySliceRS); l++ {
		if arraySliceRS[l].bMatched {
			arraySlice = append(arraySlice, arraySliceRS[l].CurRec)
		}
	}
	return arraySlice, nil
}

//markRecords : mark records to be filtered based on filterConditions.
func markRecords(filterConditions string, val reflect.Value, arraySlice map[int]DataSet) (map[int]DataSet, error) {
	subQuery := GetSubQueries(filterConditions, OPAND+STARTBRACE, OPOR+STARTBRACE)
	for i := range subQuery {
		var LOP, Key, Column, Op, CurQuery string
		CurQuery = subQuery[i]
		if CurQuery[:1] == OPAND || CurQuery[:1] == OPOR {
			LOP = CurQuery[:1]
			CurQuery = CurQuery[1:]
		} else {
			LOP = OPBLANK
		}

		Key, Column, Op = GetQueryDetails(CurQuery)
		if Key == "" {
			return nil, errors.New("Filter Key not found!!!")
		}
		if Column == "" {
			return nil, errors.New("Filter Column not found!!!")
		}

		var ColumnStruct []string
		if strings.Contains(Column, ".") {
			ColumnStruct = strings.Split(Column, ".")
			Column = ColumnStruct[0]
		}

		ColumnUp := strings.ToUpper(Column)
		var ColumnFound = false
		for j := 0; j < val.NumField(); j++ {
			if ColumnUp == strings.ToUpper(val.Type().Field(j).Name) {
				Column = val.Type().Field(j).Name
				ColumnFound = true
			}
		}

		if !ColumnFound {
			return nil, errors.New("Filter [" + Column + "] No such column exist!!!")
		}

		for k := 0; k < len(arraySlice); k++ {
			bConsiderMe := false
			if LOP == OPBLANK || LOP == OPOR {
				bConsiderMe = !arraySlice[k].bMatched
			} else if LOP == OPAND {
				bConsiderMe = arraySlice[k].bMatched
			}

			var err error
			if bConsiderMe {
				bMatched := false
				CurField := reflect.Indirect(reflect.ValueOf(arraySlice[k].CurRec)).FieldByName(Column)

				if isPrimitive(CurField) {
					if bMatched, err = processPrimitive(CurField, Key, Op); err != nil {
						return nil, err
					}

				} else if CurField.Kind() == reflect.Struct {
					//to avoid out of bound errors we need to check this
					if len(ColumnStruct) < 2 {
						return nil, errors.New("Filter applied to struct column but inner field not specified")
					}
					innerField := CurField.FieldByName(ColumnStruct[1])
					if !innerField.IsValid() {
						return nil, errors.New("Filter [" + ColumnStruct[1] + "] No such column exist!!!")
					}
					if isPrimitive(innerField) {
						if bMatched, err = processPrimitive(innerField, Key, Op); err != nil {
							return nil, err
						}
					} else {
						return nil, errors.New("Filter [" + ColumnStruct[1] + "] Column is not primitive type")
					}

				} else if CurField.Kind() == reflect.Slice {
					if len(ColumnStruct) < 2 {
						continue
					}
					fieldName := ColumnStruct[1]

					var fieldValueExists bool
					for i := 0; i < CurField.Len(); i++ {
						element := CurField.Index(i)

						if !element.IsValid() {
							return nil, fmt.Errorf("Filter [%s] No such column exist!!!", fieldName)
						}

						if element.FieldByName("ID").String() == fieldName {
							fieldValueExists = true
							if isPrimitive(element.FieldByName("ID")) {
								if bMatched, err = processPrimitive(element.FieldByName("Value"), Key, Op); err != nil {
									return nil, err
								}
							} else {
								return nil, fmt.Errorf("Filter [%s] Column is not primitive type", fieldName)
							}
						}
					}
					if !fieldValueExists {
						return nil, fmt.Errorf("Filter [%s] No such column exist!!!", fieldName)
					}
				}
				arraySlice[k] = DataSet{bMatched: bMatched, CurRec: arraySlice[k].CurRec}
			}
		}
	}
	return arraySlice, nil
}

func isPrimitive(CurField reflect.Value) bool {
	switch CurField.Kind() {
	case reflect.Bool, reflect.String, reflect.Float32, reflect.Float64, reflect.Int, reflect.Int32, reflect.Int64:
		return true
	}
	return false
}

func processPrimitive(CurField reflect.Value, Key string, Op string) (bMatched bool, err error) {
	switch CurField.Kind() {
	case reflect.Bool:
		b, _ := strconv.ParseBool(Key)
		if Op == "" || Op == "==" {
			if CurField.Bool() == b {
				bMatched = true
			}
		} else if Op == "!" {
			if CurField.Bool() != b {
				bMatched = true
			}
		} else {
			return bMatched, errors.New("Filter Invalid Operator [" + Op + "] Applied!!!")
		}
	case reflect.String:
		if Key == NILQUERY {
			Key = ""
		}
		if Op == "" {
			CurStr := strings.ToLower(CurField.String())
			CurKey := strings.ToLower(Key)
			if strings.Contains(CurStr, CurKey) {
				bMatched = true
			}
		} else if Op == "==" {
			if strings.Contains(CurField.String(), Key) {
				bMatched = true
			}
		} else if Op == "===" {
			if CurField.String() == Key {
				bMatched = true
			}
		} else if Op == "!" {
			if CurField.String() != Key {
				bMatched = true
			}
		} else if Op == ">" {
			if CurField.String() > Key {
				bMatched = true
			}
		} else if Op == ">=" {
			if CurField.String() >= Key {
				bMatched = true
			}
		} else if Op == "<" {
			if CurField.String() < Key {
				bMatched = true
			}
		} else if Op == "<=" {
			if CurField.String() <= Key {
				bMatched = true
			}
		} else {
			return bMatched, errors.New("Filter Invalid Operator [" + Op + "] Applied!!!")
		}
	case reflect.Float32, reflect.Float64:
		f, _ := strconv.ParseFloat(Key, 64)
		if Op == "" || Op == "==" {
			if CurField.Float() == f {
				bMatched = true
			}
		} else if Op == "!" {
			if CurField.Float() != f {
				bMatched = true
			}
		} else if Op == ">" {
			if CurField.Float() > f {
				bMatched = true
			}
		} else if Op == ">=" {
			if CurField.Float() >= f {
				bMatched = true
			}
		} else if Op == "<" {
			if CurField.Float() < f {
				bMatched = true
			}
		} else if Op == "<=" {
			if CurField.Float() <= f {
				bMatched = true
			}
		} else {
			return bMatched, errors.New("Filter Invalid Operator [" + Op + "] Applied!!!")
		}
	case reflect.Int, reflect.Int32, reflect.Int64:
		i, _ := strconv.ParseInt(Key, 10, 64)
		if Op == "" || Op == "==" {
			if CurField.Int() == i {
				bMatched = true
			}
		} else if Op == "!" {
			if CurField.Int() != i {
				bMatched = true
			}
		} else if Op == ">" {
			if CurField.Int() > i {
				bMatched = true
			}
		} else if Op == ">=" {
			if CurField.Int() >= i {
				bMatched = true
			}
		} else if Op == "<" {
			if CurField.Int() < i {
				bMatched = true
			}
		} else if Op == "<=" {
			if CurField.Int() <= i {
				bMatched = true
			}
		} else {
			return bMatched, errors.New("Filter Invalid Operator [" + Op + "] Applied!!!")
		}

	default:
		return bMatched, errors.New("Field is not primitive type")
	}

	return bMatched, nil
}

//PriorityFilter : Filters data based on priority of given complex filterConditions.
func PriorityFilter(filterConditions string, val reflect.Value, arraySlice []interface{}) ([]interface{}, error) {
	var err error
	arraySliceRS := make(map[int]DataSet)
	for ind := range arraySlice {
		arraySliceRS[ind] = DataSet{bMatched: false, CurRec: arraySlice[ind]}
	}

	filterConditions = CleanConditions(filterConditions)
	arraySliceRS, err = ResolveFilterConditions(filterConditions, val, arraySliceRS)
	if err != nil {
		return nil, err
	}

	arraySlice = make([]interface{}, 0)
	for l := 0; l < len(arraySliceRS); l++ {
		if arraySliceRS[l].bMatched {
			arraySlice = append(arraySlice, arraySliceRS[l].CurRec)
		}
	}
	return arraySlice, nil
}

//ResolveFilterConditions : Resolve Filter Conditions based on priority
func ResolveFilterConditions(filter string, val reflect.Value, dataMap map[int]DataSet) (map[int]DataSet, error) {
	var err error
	dataMapRet := make(map[int]DataSet)
	var bExtractForward, bIsComplexFilter, bImbalance bool
	var pendingFilter, curCondition, curOperator string
	if strings.HasPrefix(filter, STARTBRACE+STARTBRACE) && strings.HasSuffix(filter, ENDBRACE+ENDBRACE) {
		iBalancePos, iLength := GetBalancePosition(filter)
		if iBalancePos == iLength-1 {
			filter = filter[1 : len(filter)-1]
		} else {
			bImbalance = true
			curCondition = filter[:iBalancePos+1]
			curOperator = filter[iBalancePos+1 : iBalancePos+2]
			pendingFilter = filter[iBalancePos+2:]
		}
	}

	if bImbalance {
		dataMapCopy := make(map[int]DataSet)
		for k, v := range dataMap {
			dataMapCopy[k] = v
		}

		dataMap, err = ResolveFilterConditions(pendingFilter, val, dataMap)
		if err != nil {
			return nil, err
		}
		dataMapCopy, err = ResolveFilterConditions(curCondition, val, dataMapCopy)
		if err != nil {
			return nil, err
		}
		dataMapRet, err = MergeFilterResults(dataMap, dataMapCopy, curOperator)
		if err != nil {
			return nil, err
		}
	} else {
		if strings.HasPrefix(filter, STARTBRACE) && strings.HasSuffix(filter, ENDBRACE+ENDBRACE) {
			bExtractForward = true
			bIsComplexFilter = true
		} else if strings.HasPrefix(filter, STARTBRACE+STARTBRACE) && strings.HasSuffix(filter, ENDBRACE) {
			bExtractForward = false
			bIsComplexFilter = true
		} else if strings.Contains(filter, STARTBRACE+STARTBRACE) {
			bExtractForward = true
			bIsComplexFilter = true
		} else if strings.Contains(filter, ENDBRACE+ENDBRACE) {
			bExtractForward = false
			bIsComplexFilter = true
		} else {
			bIsComplexFilter = false
		}

		if bIsComplexFilter {
			pendingFilter, curCondition, curOperator = ExtractConditions(filter, bExtractForward)
			dataMapCopy := make(map[int]DataSet)
			for k, v := range dataMap {
				dataMapCopy[k] = v
			}

			dataMap, err = ResolveFilterConditions(pendingFilter, val, dataMap)
			if err != nil {
				return nil, err
			}
			dataMapCopy, err = markRecords(curCondition, val, dataMapCopy)
			if err != nil {
				return nil, err
			}
			dataMapRet, err = MergeFilterResults(dataMap, dataMapCopy, curOperator)
			if err != nil {
				return nil, err
			}
		} else {
			curCondition = filter
			dataMapRet, err = markRecords(curCondition, val, dataMap)
			if err != nil {
				return nil, err
			}
		}
	}
	return dataMapRet, nil
}

//ExtractConditions : Extracts simplest conditions from complex conditions.
func ExtractConditions(sFilter string, bExtractForward bool) (pendingFilter string, curCondition string, curOperator string) {
	var i, iAnd, iOr int
	i = -1
	if bExtractForward {
		iAnd = strings.Index(sFilter, OPAND)
		iOr = strings.Index(sFilter, OPOR)
	} else {
		iAnd = strings.LastIndex(sFilter, OPAND)
		iOr = strings.LastIndex(sFilter, OPOR)
	}

	if iAnd > -1 && iOr > -1 {
		i = iOr
		if bExtractForward && (iAnd < iOr) {
			i = iAnd
		} else if !bExtractForward && (iAnd > iOr) {
			i = iAnd
		}
	} else if iAnd > -1 {
		i = iAnd
	} else if iOr > -1 {
		i = iOr
	}

	if bExtractForward {
		curCondition = sFilter[0:i]
		curOperator = sFilter[i : i+1]
		pendingFilter = sFilter[i+1:]
	} else {
		curCondition = sFilter[i+1:]
		curOperator = sFilter[i : i+1]
		pendingFilter = sFilter[:i]
	}
	return pendingFilter, curCondition, curOperator
}

//MergeFilterResults : Merges Filter Conditions results
func MergeFilterResults(dataMapLeft, dataMapRight map[int]DataSet, curOperation string) (map[int]DataSet, error) {
	if len(dataMapLeft) != len(dataMapRight) {
		return nil, errors.New("MergeFilterResults Length Check Failed!!!")
	}

	if curOperation != OPOR && curOperation != OPAND {
		return nil, errors.New("Filter Invalid Operator [" + curOperation + "] Applied!!!")
	}

	for i := 0; i < len(dataMapLeft); i++ {
		bMatched := false
		if curOperation == OPOR {
			bMatched = dataMapLeft[i].bMatched || dataMapRight[i].bMatched
		} else if curOperation == OPAND {
			bMatched = dataMapLeft[i].bMatched && dataMapRight[i].bMatched
		}
		dataMapLeft[i] = DataSet{bMatched: bMatched, CurRec: dataMapLeft[i].CurRec}
	}
	return dataMapLeft, nil
}

//GetBalancePosition : Gets balancing position for condition.
func GetBalancePosition(sFilter string) (iBalancePos int, iLength int) {
	var iStart, iEnd int
	iLength = len(sFilter)
	for i := 0; i < iLength; i++ {
		if sFilter[i:i+1] == STARTBRACE {
			iStart++
		} else if sFilter[i:i+1] == ENDBRACE {
			iEnd++
		}

		if (iStart == iEnd) && (iStart != 0) {
			iBalancePos = i
			break
		}
	}

	return iBalancePos, iLength
}
