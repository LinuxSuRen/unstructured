package pkg

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNestedField(t *testing.T) {
	type args struct {
		obj    map[string]interface{}
		fields []string
	}
	tests := []struct {
		name    string
		args    args
		wantVal interface{}
		want    bool
		wantErr bool
	}{{
		name: "parameter is nil",
		args: args{obj: map[string]interface{}{
			"name": nil,
		}, fields: []string{"name", "good"}},
		want:    false,
		wantErr: false,
	}, {
		name: "type is map[string]string",
		args: args{
			obj: map[string]interface{}{
				"name": "rick",
			},
			fields: []string{"name"},
		},
		wantVal: "rick",
		want:    true,
		wantErr: false,
	}, {
		name: "complex map structure, correct case",
		args: args{
			obj: map[string]interface{}{
				"jenkins": map[string]interface{}{
					"name": "jenkins",
				},
			},
			fields: []string{"jenkins", "name"},
		},
		wantVal: "jenkins",
		want:    true,
		wantErr: false,
	}, {
		name: "complex map structure, invalid fields",
		args: args{
			obj: map[string]interface{}{
				"jenkins": map[string]interface{}{
					"name": "jenkins",
				},
			},
			fields: []string{"jenkins", "node"},
		},
		want:    false,
		wantErr: false,
	}, {
		name: "complex map structure, invalid fields",
		args: args{
			obj: map[string]interface{}{
				"jenkins": map[string]interface{}{
					"name": "jenkins",
				},
			},
			fields: []string{"jenkins", "name", "node"},
		},
		want:    false,
		wantErr: true,
	}, {
		name: "get item from a nested array, correct case",
		args: args{
			obj: map[string]interface{}{
				"jenkins": map[string]interface{}{
					"clouds": []map[string]interface{}{{
						"name": "one",
					}, {
						"name": "two",
					}},
				},
			},
			fields: []string{"jenkins", "clouds[1]", "name"},
		},
		wantVal: "two",
		want:    true,
		wantErr: false,
	}, {
		name: "get item from a nested array, correct case",
		args: args{
			obj: map[string]interface{}{
				"jenkins": map[string]interface{}{
					"clouds": []interface{}{
						"one", "two",
					},
				},
			},
			fields: []string{"jenkins", "clouds[1]"},
		},
		wantVal: "two",
		want:    true,
		wantErr: false,
	}, {
		name: "get item from a nested array, correct case",
		args: args{
			obj: map[string]interface{}{
				"jenkins": map[string]interface{}{
					"clouds": []string{
						"one", "two",
					},
				},
			},
			fields: []string{"jenkins", "clouds[1]"},
		},
		wantVal: "two",
		want:    true,
		wantErr: false,
	}, {
		name: "get item from a nested array, invalid index",
		args: args{
			obj: map[string]interface{}{
				"jenkins": map[string]interface{}{
					"clouds": []map[string]interface{}{{
						"name": "one",
					}, {
						"name": "two",
					}},
				},
			},
			fields: []string{"jenkins", "clouds[9]", "name"},
		},
		wantVal: nil,
		want:    false,
		wantErr: true,
	}, {
		name: "",
		args: args{
			obj: map[string]interface{}{
				"jenkins": map[string]interface{}{
					"clouds": []map[string]interface{}{
						{
							"kubernetes": map[string]interface{}{
								"name": "kubernetes",
							},
						},
					},
				},
			},
			fields: []string{"jenkins", "clouds[0]", "kubernetes", "name"},
		},
		wantVal: "kubernetes",
		want:    true,
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := NestedField(tt.args.obj, tt.args.fields...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NestedField() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantVal) {
				t.Errorf("NestedField() got = %v, wantVal %v", got, tt.wantVal)
			}
			if got1 != tt.want {
				t.Errorf("NestedField() got1 = %v, wantVal %v", got1, tt.want)
			}
		})
	}
}

func Test_getIndex(t *testing.T) {
	type args struct {
		word string
	}
	tests := []struct {
		name string
		args args
		want int
	}{{
		name: "case, [1]",
		args: args{
			word: "[1]",
		},
		want: 1,
	}, {
		name: "case, name[1]",
		args: args{
			word: "name[1]",
		},
		want: 1,
	}, {
		name: "case, name[100]",
		args: args{
			word: "name[100]",
		},
		want: 100,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getIndex(tt.args.word); got != tt.want {
				t.Errorf("getIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isIndexedField(t *testing.T) {
	type args struct {
		field string
	}
	tests := []struct {
		name        string
		args        args
		wantMatched bool
	}{{
		name:        "normal",
		args:        args{field: "name[1]"},
		wantMatched: true,
	}, {
		name:        "invalid",
		args:        args{field: "good"},
		wantMatched: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMatched := isIndexedField(tt.args.field); gotMatched != tt.wantMatched {
				t.Errorf("isIndexedField() = %v, want %v", gotMatched, tt.wantMatched)
			}
		})
	}
}

func Test_removeIndexFromField(t *testing.T) {
	type args struct {
		field string
	}
	tests := []struct {
		name string
		args args
		want string
	}{{
		name: "normal",
		args: args{field: "name[1]"},
		want: "name",
	}, {
		name: "is not a indexed field",
		args: args{field: "name"},
		want: "name",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeIndexFromField(tt.args.field); got != tt.want {
				t.Errorf("removeIndexFromField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetNestedField(t *testing.T) {
	type args struct {
		obj    map[string]interface{}
		target interface{}
		fields []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		verify  func(t *testing.T, obj map[string]interface{})
	}{{
		name: "a plain map, update the exist key-value",
		args: args{
			obj: map[string]interface{}{
				"name": "rick",
			},
			target: "good",
			fields: []string{"name"},
		},
		wantErr: false,
		verify: func(t *testing.T, obj map[string]interface{}) {
			assert.Equal(t, "good", obj["name"])
		},
	}, {
		name: "a plain map, put a non-exit key-value",
		args: args{
			obj: map[string]interface{}{
				"name": "rick",
			},
			target: "good",
			fields: []string{"luck"},
		},
		wantErr: false,
		verify: func(t *testing.T, obj map[string]interface{}) {
			assert.Equal(t, "good", obj["luck"])
			assert.Equal(t, "rick", obj["name"])
		},
	}, {
		name: "update a two non-exist levels map key-value",
		args: args{
			obj: map[string]interface{}{
				"name": "rick",
			},
			target: "good",
			fields: []string{"job", "luck"},
		},
		wantErr: false,
		verify: func(t *testing.T, obj map[string]interface{}) {
			assert.Equal(t, map[string]interface{}{
				"luck": "good",
			}, obj["job"])
		},
	}, {
		name: "update a two levels map key-value",
		args: args{
			obj: map[string]interface{}{
				"name": map[string]interface{}{
					"job": "back",
				},
			},
			target: "good",
			fields: []string{"name", "job"},
		},
		wantErr: false,
		verify: func(t *testing.T, obj map[string]interface{}) {
			assert.Equal(t, map[string]interface{}{
				"job": "good",
			}, obj["name"])
		},
	}, {
		name: "have a nested array",
		args: args{
			obj: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{
						"name": "rick",
					},
				},
			},
			target: "good",
			fields: []string{"items[0]", "name"},
		},
		wantErr: false,
		verify: func(t *testing.T, obj map[string]interface{}) {
			assert.Equal(t, []interface{}{
				map[string]interface{}{
					"name": "good",
				},
			}, obj["items"])
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetNestedField(tt.args.obj, tt.args.target, tt.args.fields...); (err != nil) != tt.wantErr {
				t.Errorf("SetNestedField() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.verify != nil {
				tt.verify(t, tt.args.obj)
			}
		})
	}
}

func TestNestedFieldAsString(t *testing.T) {
	type args struct {
		obj    map[string]interface{}
		fields []string
	}
	tests := []struct {
		name       string
		args       args
		wantStrVal string
		wantOk     bool
		wantErr    assert.ErrorAssertionFunc
	}{{
		name: "normal",
		args: args{
			obj: map[string]interface{}{
				"name": "rick",
			},
			fields: []string{"name"},
		},
		wantStrVal: "rick",
		wantOk:     true,
		wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
			assert.Nil(t, err)
			return true
		},
	}, {
		name: "non string",
		args: args{
			obj: map[string]interface{}{
				"age": 12,
			},
			fields: []string{"age"},
		},
		wantStrVal: "",
		wantOk:     false,
		wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
			assert.NotNil(t, err)
			return true
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStrVal, gotOk, err := NestedFieldAsString(tt.args.obj, tt.args.fields...)
			if !tt.wantErr(t, err, fmt.Sprintf("NestedFieldAsString(%v, %v)", tt.args.obj, tt.args.fields)) {
				return
			}
			assert.Equalf(t, tt.wantStrVal, gotStrVal, "NestedFieldAsString(%v, %v)", tt.args.obj, tt.args.fields)
			assert.Equalf(t, tt.wantOk, gotOk, "NestedFieldAsString(%v, %v)", tt.args.obj, tt.args.fields)
		})
	}
}

func TestNestedFieldAsInt(t *testing.T) {
	type args struct {
		obj    map[string]interface{}
		fields []string
	}
	tests := []struct {
		name       string
		args       args
		wantIntVal int
		wantOk     bool
		wantErr    assert.ErrorAssertionFunc
	}{{
		name: "normal",
		args: args{
			obj: map[string]interface{}{
				"age": 12,
			},
			fields: []string{"age"},
		},
		wantIntVal: 12,
		wantOk:     true,
		wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
			assert.Nil(t, err)
			return true
		},
	}, {
		name: "non int type",
		args: args{
			obj: map[string]interface{}{
				"age": "12",
			},
			fields: []string{"age"},
		},
		wantIntVal: 0,
		wantOk:     false,
		wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
			assert.NotNil(t, err)
			return true
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIntVal, gotOk, err := NestedFieldAsInt(tt.args.obj, tt.args.fields...)
			if !tt.wantErr(t, err, fmt.Sprintf("NestedFieldAsInt(%v, %v)", tt.args.obj, tt.args.fields)) {
				return
			}
			assert.Equalf(t, tt.wantIntVal, gotIntVal, "NestedFieldAsInt(%v, %v)", tt.args.obj, tt.args.fields)
			assert.Equalf(t, tt.wantOk, gotOk, "NestedFieldAsInt(%v, %v)", tt.args.obj, tt.args.fields)
		})
	}
}
