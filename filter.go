package relay

import (
	"errors"
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

	arraySliceRS := make(map[int]DataSet)
	for ind := range arraySlice {
		arraySliceRS[ind] = DataSet{bMatched: false, CurRec: arraySlice[ind]}
	}

	filterConditions = CleanConditions(filterConditions)
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
		for k := 0; k < len(arraySliceRS); k++ {
			bConsiderMe := false
			if LOP == OPBLANK || LOP == OPOR {
				bConsiderMe = !arraySliceRS[k].bMatched
			} else if LOP == OPAND {
				bConsiderMe = arraySliceRS[k].bMatched
			}

			var err error

			if bConsiderMe {
				bMatched := false
				CurField := reflect.Indirect(reflect.ValueOf(arraySliceRS[k].CurRec)).FieldByName(Column)

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

				}

				arraySliceRS[k] = DataSet{bMatched: bMatched, CurRec: arraySliceRS[k].CurRec}
			}
		}
	}

	arraySlice = make([]interface{}, 0)
	for l := 0; l < len(arraySliceRS); l++ {
		if arraySliceRS[l].bMatched {
			arraySlice = append(arraySlice, arraySliceRS[l].CurRec)
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
