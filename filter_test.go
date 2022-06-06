package relay

import (
	"reflect"
	"testing"
)

type data struct {
	SomeID        string             `json:"someId"`
	OutOfSupport  bool               `json:"outOfSupport"`
	DynamicFields []dynamicFieldType `json:"dynamicField"`
	SliceField    []string           `json:"sliceField"`
}

type dynamicFieldType struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type dataNotPrimitiveValue struct {
	SomeID        string                              `json:"someId"`
	OutOfSupport  bool                                `json:"outOfSupport"`
	DynamicFields []DynamicFieldTypeNotPrimitiveValue `json:"dynamicField"`
	SliceField    []string                            `json:"sliceField"`
}

type DynamicFieldTypeNotPrimitiveValue struct {
	ID    string   `json:"id"`
	Value []string `json:"value"`
}

func Test_markRecords(t *testing.T) {
	type args struct {
		filterConditions string
		val              reflect.Value
		arraySlice       map[int]DataSet
	}
	tests := []struct {
		name    string
		args    args
		want    map[int]DataSet
		wantErr bool
	}{
		{
			name: "marked and not marked, dynamic ID not match, multiple filters",
			args: args{
				filterConditions: "{key:someValue,column:dynamicFields.someID}&{key:5,column:SomeID}",
				val: func() reflect.Value {
					return reflect.Indirect(reflect.ValueOf(&data{}))
				}(),
				arraySlice: func() map[int]DataSet {
					arraySlice := []interface{}{
						data{
							SomeID:       "5",
							OutOfSupport: true,
							DynamicFields: []dynamicFieldType{
								{
									ID:    "someID",
									Value: "someValue",
								},
							},
						},
						data{
							SomeID:       "5",
							OutOfSupport: true,
							DynamicFields: []dynamicFieldType{
								{
									ID:    "someID_2",
									Value: "someValue",
								},
							},
						},
					}
					arraySliceRS := make(map[int]DataSet)
					for ind := range arraySlice {
						arraySliceRS[ind] = DataSet{bMatched: false, CurRec: arraySlice[ind]}
					}
					return arraySliceRS
				}(),
			},
			want: map[int]DataSet{
				0: {
					bMatched: true,
					CurRec: data{
						SomeID:       "5",
						OutOfSupport: true,
						DynamicFields: []dynamicFieldType{
							{
								ID:    "someID",
								Value: "someValue",
							},
						},
					},
				},
				1: {
					bMatched: false,
					CurRec: data{
						SomeID:       "5",
						OutOfSupport: true,
						DynamicFields: []dynamicFieldType{
							{
								ID:    "someID_2",
								Value: "someValue",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "marked and not marked, dynamic Value not match, multiple filters",
			args: args{
				filterConditions: "{key:someValue,column:dynamicFields.someID,op:===}&{key:5,column:SomeID}",
				val: func() reflect.Value {
					return reflect.Indirect(reflect.ValueOf(&data{}))
				}(),
				arraySlice: func() map[int]DataSet {
					arraySlice := []interface{}{
						data{
							SomeID:       "5",
							OutOfSupport: true,
							DynamicFields: []dynamicFieldType{
								{
									ID:    "someID",
									Value: "someValue",
								},
							},
						},
						data{
							SomeID:       "5",
							OutOfSupport: true,
							DynamicFields: []dynamicFieldType{
								{
									ID:    "someID",
									Value: "someValue_2",
								},
							},
						},
					}
					arraySliceRS := make(map[int]DataSet)
					for ind := range arraySlice {
						arraySliceRS[ind] = DataSet{bMatched: false, CurRec: arraySlice[ind]}
					}
					return arraySliceRS
				}(),
			},
			want: map[int]DataSet{
				0: {
					bMatched: true,
					CurRec: data{
						SomeID:       "5",
						OutOfSupport: true,
						DynamicFields: []dynamicFieldType{
							{
								ID:    "someID",
								Value: "someValue",
							},
						},
					},
				},
				1: {
					bMatched: false,
					CurRec: data{
						SomeID:       "5",
						OutOfSupport: true,
						DynamicFields: []dynamicFieldType{
							{
								ID:    "someID",
								Value: "someValue_2",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "marked, multiple filters",
			args: args{
				filterConditions: "{key:someValue,column:dynamicFields.someID}&{key:5,column:SomeID}&{key:true,column:outOfSupport,op:==}",
				val: func() reflect.Value {
					return reflect.Indirect(reflect.ValueOf(&data{}))
				}(),
				arraySlice: func() map[int]DataSet {
					arraySlice := []interface{}{
						data{
							SomeID:       "5",
							OutOfSupport: true,
							DynamicFields: []dynamicFieldType{
								{
									ID:    "someID",
									Value: "someValue",
								},
							},
						},
					}
					arraySliceRS := make(map[int]DataSet)
					for ind := range arraySlice {
						arraySliceRS[ind] = DataSet{bMatched: false, CurRec: arraySlice[ind]}
					}
					return arraySliceRS
				}(),
			},
			want: map[int]DataSet{
				0: {
					bMatched: true,
					CurRec: data{
						SomeID:       "5",
						OutOfSupport: true,
						DynamicFields: []dynamicFieldType{
							{
								ID:    "someID",
								Value: "someValue",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "marked, primitive type",
			args: args{
				filterConditions: "{key:5,column:SomeID}",
				val: func() reflect.Value {
					return reflect.Indirect(reflect.ValueOf(&data{}))
				}(),
				arraySlice: func() map[int]DataSet {
					arraySlice := []interface{}{
						data{
							SomeID:       "5",
							OutOfSupport: true,
							DynamicFields: []dynamicFieldType{
								{
									ID:    "someID",
									Value: "someValue",
								},
							},
						},
					}
					arraySliceRS := make(map[int]DataSet)
					for ind := range arraySlice {
						arraySliceRS[ind] = DataSet{bMatched: false, CurRec: arraySlice[ind]}
					}
					return arraySliceRS
				}(),
			},
			want: map[int]DataSet{
				0: {
					bMatched: true,
					CurRec: data{
						SomeID:       "5",
						OutOfSupport: true,
						DynamicFields: []dynamicFieldType{
							{
								ID:    "someID",
								Value: "someValue",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "not marked, primitive not match",
			args: args{
				filterConditions: "{key:someValue,column:dynamicFields.someID}&{key:4,column:SomeID}",
				val: func() reflect.Value {
					return reflect.Indirect(reflect.ValueOf(&data{}))
				}(),
				arraySlice: func() map[int]DataSet {
					arraySlice := []interface{}{
						data{
							SomeID:       "5",
							OutOfSupport: true,
							DynamicFields: []dynamicFieldType{
								{
									ID:    "someID",
									Value: "someValue",
								},
							},
						},
					}
					arraySliceRS := make(map[int]DataSet)
					for ind := range arraySlice {
						arraySliceRS[ind] = DataSet{bMatched: false, CurRec: arraySlice[ind]}
					}
					return arraySliceRS
				}(),
			},
			want: map[int]DataSet{
				0: {
					bMatched: false,
					CurRec: data{
						SomeID:       "5",
						OutOfSupport: true,
						DynamicFields: []dynamicFieldType{
							{
								ID:    "someID",
								Value: "someValue",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "not marked, dynamicFields ID not match",
			args: args{
				filterConditions: "{key:someValue1,column:dynamicFields.someID}&{key:5,column:SomeID}",
				val: func() reflect.Value {
					return reflect.Indirect(reflect.ValueOf(&data{}))
				}(),
				arraySlice: func() map[int]DataSet {
					arraySlice := []interface{}{
						data{
							SomeID:       "5",
							OutOfSupport: true,
							DynamicFields: []dynamicFieldType{
								{
									ID:    "someID",
									Value: "someValue",
								},
							},
						},
					}
					arraySliceRS := make(map[int]DataSet)
					for ind := range arraySlice {
						arraySliceRS[ind] = DataSet{bMatched: false, CurRec: arraySlice[ind]}
					}
					return arraySliceRS
				}(),
			},
			want: map[int]DataSet{
				0: {
					bMatched: false,
					CurRec: data{
						SomeID:       "5",
						OutOfSupport: true,
						DynamicFields: []dynamicFieldType{
							{
								ID:    "someID",
								Value: "someValue",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "marked, dynamicFields values match, contains",
			args: args{
				filterConditions: "{key:someValue,column:dynamicFields.someID}&{key:5,column:SomeID}",
				val: func() reflect.Value {
					return reflect.Indirect(reflect.ValueOf(&data{}))
				}(),
				arraySlice: func() map[int]DataSet {
					arraySlice := []interface{}{
						data{
							SomeID:       "5",
							OutOfSupport: true,
							DynamicFields: []dynamicFieldType{
								{
									ID:    "someID",
									Value: "someValue1",
								},
							},
						},
					}
					arraySliceRS := make(map[int]DataSet)
					for ind := range arraySlice {
						arraySliceRS[ind] = DataSet{bMatched: false, CurRec: arraySlice[ind]}
					}
					return arraySliceRS
				}(),
			},
			want: map[int]DataSet{
				0: {
					bMatched: true,
					CurRec: data{
						SomeID:       "5",
						OutOfSupport: true,
						DynamicFields: []dynamicFieldType{
							{
								ID:    "someID",
								Value: "someValue1",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "not marked, dynamicFields values not match, ===",
			args: args{
				filterConditions: "{key:someValue,column:dynamicFields.someID,op:===}&{key:5,column:SomeID}",
				val: func() reflect.Value {
					return reflect.Indirect(reflect.ValueOf(&data{}))
				}(),
				arraySlice: func() map[int]DataSet {
					arraySlice := []interface{}{
						data{
							SomeID:       "5",
							OutOfSupport: true,
							DynamicFields: []dynamicFieldType{
								{
									ID:    "someID",
									Value: "someValue1",
								},
							},
						},
					}
					arraySliceRS := make(map[int]DataSet)
					for ind := range arraySlice {
						arraySliceRS[ind] = DataSet{bMatched: false, CurRec: arraySlice[ind]}
					}
					return arraySliceRS
				}(),
			},
			want: map[int]DataSet{
				0: {
					bMatched: false,
					CurRec: data{
						SomeID:       "5",
						OutOfSupport: true,
						DynamicFields: []dynamicFieldType{
							{
								ID:    "someID",
								Value: "someValue1",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "field not found, not marked",
			args: args{
				filterConditions: "{key:someValue,column:dynamicFields.someID_1}&{key:5,column:SomeID}",
				val: func() reflect.Value {
					return reflect.Indirect(reflect.ValueOf(&data{}))
				}(),
				arraySlice: func() map[int]DataSet {
					arraySlice := []interface{}{
						data{
							SomeID:       "5",
							OutOfSupport: true,
							DynamicFields: []dynamicFieldType{
								{
									ID:    "someID",
									Value: "someValue1",
								},
							},
						},
					}
					arraySliceRS := make(map[int]DataSet)
					for ind := range arraySlice {
						arraySliceRS[ind] = DataSet{bMatched: false, CurRec: arraySlice[ind]}
					}
					return arraySliceRS
				}(),
			},
			want: map[int]DataSet{
				0: {
					bMatched: false,
					CurRec: data{
						SomeID:       "5",
						OutOfSupport: true,
						DynamicFields: []dynamicFieldType{
							{
								ID:    "someID",
								Value: "someValue1",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error, dynamic field name not provided",
			args: args{
				filterConditions: "{key:someValue,column:dynamicFields.}&{key:5,column:SomeID}",
				val: func() reflect.Value {
					return reflect.Indirect(reflect.ValueOf(&data{}))
				}(),
				arraySlice: func() map[int]DataSet {
					arraySlice := []interface{}{
						data{
							SomeID:       "5",
							OutOfSupport: true,
							DynamicFields: []dynamicFieldType{
								{
									ID:    "someID",
									Value: "someValue1",
								},
							},
							SliceField: []string{"some"},
						},
					}
					arraySliceRS := make(map[int]DataSet)
					for ind := range arraySlice {
						arraySliceRS[ind] = DataSet{bMatched: false, CurRec: arraySlice[ind]}
					}
					return arraySliceRS
				}(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not marked, no field, skipping item",
			args: args{
				filterConditions: "{key:someValue,column:dynamicFields}",
				val: func() reflect.Value {
					return reflect.Indirect(reflect.ValueOf(&data{}))
				}(),
				arraySlice: func() map[int]DataSet {
					arraySlice := []interface{}{
						data{
							SomeID:       "5",
							OutOfSupport: true,
							DynamicFields: []dynamicFieldType{
								{
									ID:    "someID",
									Value: "someValue",
								},
							},
							SliceField: []string{"some"},
						},
					}
					arraySliceRS := make(map[int]DataSet)
					for ind := range arraySlice {
						arraySliceRS[ind] = DataSet{bMatched: false, CurRec: arraySlice[ind]}
					}
					return arraySliceRS
				}(),
			},
			want: map[int]DataSet{
				0: {
					bMatched: false,
					CurRec: data{
						SomeID:       "5",
						OutOfSupport: true,
						DynamicFields: []dynamicFieldType{
							{
								ID:    "someID",
								Value: "someValue",
							},
						},
						SliceField: []string{"some"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "not marked, missing id",
			args: args{
				filterConditions: "{key:someValue,column:dynamicFields.someID}",
				val: func() reflect.Value {
					return reflect.Indirect(reflect.ValueOf(&data{}))
				}(),
				arraySlice: func() map[int]DataSet {
					arraySlice := []interface{}{
						data{
							SomeID:       "5",
							OutOfSupport: true,
							DynamicFields: []dynamicFieldType{
								{
									Value: "someValue",
								},
							},
							SliceField: []string{"some"},
						},
					}
					arraySliceRS := make(map[int]DataSet)
					for ind := range arraySlice {
						arraySliceRS[ind] = DataSet{bMatched: false, CurRec: arraySlice[ind]}
					}
					return arraySliceRS
				}(),
			},
			want: map[int]DataSet{
				0: {
					bMatched: false,
					CurRec: data{
						SomeID:       "5",
						OutOfSupport: true,
						DynamicFields: []dynamicFieldType{
							{
								Value: "someValue",
							},
						},
						SliceField: []string{"some"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "not marked, no value",
			args: args{
				filterConditions: "{key:someValue,column:dynamicFields.someID}",
				val: func() reflect.Value {
					return reflect.Indirect(reflect.ValueOf(&data{}))
				}(),
				arraySlice: func() map[int]DataSet {
					arraySlice := []interface{}{
						data{
							SomeID:       "5",
							OutOfSupport: true,
							DynamicFields: []dynamicFieldType{
								{
									ID: "someID",
								},
							},
							SliceField: []string{"some"},
						},
					}
					arraySliceRS := make(map[int]DataSet)
					for ind := range arraySlice {
						arraySliceRS[ind] = DataSet{bMatched: false, CurRec: arraySlice[ind]}
					}
					return arraySliceRS
				}(),
			},
			want: map[int]DataSet{
				0: {
					bMatched: false,
					CurRec: data{
						SomeID:       "5",
						OutOfSupport: true,
						DynamicFields: []dynamicFieldType{
							{
								ID: "someID",
							},
						},
						SliceField: []string{"some"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error, not primitive value",
			args: args{
				filterConditions: "{key:someValue,column:dynamicFields.someID}&{key:5,column:SomeID}",
				val: func() reflect.Value {
					return reflect.Indirect(reflect.ValueOf(&data{}))
				}(),
				arraySlice: func() map[int]DataSet {
					arraySlice := []interface{}{
						dataNotPrimitiveValue{
							SomeID:       "5",
							OutOfSupport: true,
							DynamicFields: []DynamicFieldTypeNotPrimitiveValue{
								{
									ID:    "someID",
									Value: []string{"someValue"},
								},
							},
						},
						data{
							SomeID:       "5",
							OutOfSupport: true,
							DynamicFields: []dynamicFieldType{
								{
									ID:    "someID_2",
									Value: "someValue",
								},
							},
						},
					}
					arraySliceRS := make(map[int]DataSet)
					for ind := range arraySlice {
						arraySliceRS[ind] = DataSet{bMatched: false, CurRec: arraySlice[ind]}
					}
					return arraySliceRS
				}(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := markRecords(tt.args.filterConditions, tt.args.val, tt.args.arraySlice)
			if (err != nil) != tt.wantErr {
				t.Errorf("markRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("markRecords() got = %v, want %v", got, tt.want)
			}
		})
	}
}
