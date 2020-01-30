package relay

import (
	"sort"
)

//SortBy : SortBy is function type of less function
type SortBy func(p1, p2 interface{}) bool

//Sort : Sort is a function on functiontype
func (sortBy SortBy) Sort(obj []interface{}) {
	ps := &DataSorter{
		obj:    obj,
		sortBy: sortBy,
	}
	sort.Sort(ps)
}

//DataSorter : DataSorter combines SortBy function and data to sort
type DataSorter struct {
	obj    []interface{}
	sortBy func(p1, p2 interface{}) bool
}

//Len : Len gives length of data to sort
func (s *DataSorter) Len() int {
	return len(s.obj)
}

//Swap : Swap is function to swap elements
func (s *DataSorter) Swap(i, j int) {
	s.obj[i], s.obj[j] = s.obj[j], s.obj[i]
}

//Less : Less will be called by calling SortBy closure in the sorter
func (s *DataSorter) Less(i, j int) bool {
	return s.sortBy(s.obj[i], s.obj[j])
}
