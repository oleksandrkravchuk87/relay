package relay

import (
	"fmt"
)

type ConnectionCursor string

type PageInfo struct {
	StartCursor     int  `json:"startCursor"`
	EndCursor       int  `json:"endCursor"`
	HasPreviousPage bool `json:"hasPreviousPage"`
	HasNextPage     bool `json:"hasNextPage"`
	TotalCount      int  `json:"totalCount"`
}

type Connection struct {
	Edges      []*Edge  `json:"edges"`
	PageInfo   PageInfo `json:"pageInfo"`
	StaticInfo string   `json:"staticInfo"`
}

func NewConnection() *Connection {
	return &Connection{
		Edges:    []*Edge{},
		PageInfo: PageInfo{},
	}
}

type Edge struct {
	Node   interface{} `json:"node"`
	Cursor int         `json:"cursor"`
}

// Use NewConnectionArguments() to properly initialize default values
type ConnectionArguments struct {
	Before int              `json:"before"`
	After  int              `json:"after"`
	First  int              `json:"first"` // -1 for undefined, 0 would return zero results
	Last   int              `json:"last"`  //  -1 for undefined, 0 would return zero results
	Sort   ConnectionCursor `json:"sort"`
	Filter ConnectionCursor `json:"filter"`
}
type ConnectionArgumentsConfig struct {
	Before int `json:"before"`
	After  int `json:"after"`

	// use pointers for `First` and `Last` fields
	// so constructor would know when to use default values
	First *int `json:"first"`
	Last  *int `json:"last"`

	Sort   ConnectionCursor `json:"sort"`
	Filter ConnectionCursor `json:"filter"`
}

func NewConnectionArguments(filters map[string]interface{}) ConnectionArguments {
	conn := ConnectionArguments{
		First:  -1,
		Last:   -1,
		Before: -1,
		After:  -1,
		Sort:   "",
		Filter: "",
	}
	if filters != nil {
		if first, ok := filters["first"]; ok {
			if first, ok := first.(int); ok {
				conn.First = first
			}
		}
		if last, ok := filters["last"]; ok {
			if last, ok := last.(int); ok {
				conn.Last = last
			}
		}
		if before, ok := filters["before"]; ok {
			if before, ok := before.(int); ok {
				conn.Before = before
			}
		}
		if after, ok := filters["after"]; ok {
			if after, ok := after.(int); ok {
				conn.After = after
			}
		}
		if sort, ok := filters["sort"]; ok {
			conn.Sort = ConnectionCursor(fmt.Sprintf("%v", sort))
		}
		if filter, ok := filters["filter"]; ok {
			conn.Filter = ConnectionCursor(fmt.Sprintf("%v", filter))
		}
	}
	return conn
}
