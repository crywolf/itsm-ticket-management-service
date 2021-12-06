package repository

import (
	"reflect"
	"testing"
)

func TestNewPagination(t *testing.T) {
	type args struct {
		collectionLength int
		page             uint
		perPage          uint
	}

	tests := []struct {
		name string
		args args
		want *Pagination
	}{
		{
			"empty collection",
			args{0, 1, 10},
			&Pagination{
				Total:             0,
				Size:              0,
				Page:              1,
				Prev:              0,
				Next:              0,
				First:             1,
				Last:              1,
				FirstElementIndex: 0,
				LastElementIndex:  0,
			},
		},
		{
			"only one half-page",
			args{5, 1, 10},
			&Pagination{
				Total:             5,
				Size:              5,
				Page:              1,
				Prev:              0,
				Next:              0,
				First:             1,
				Last:              1,
				FirstElementIndex: 0,
				LastElementIndex:  4,
			},
		},
		{
			"one full page",
			args{10, 1, 10},
			&Pagination{
				Total:             10,
				Size:              10,
				Page:              1,
				Prev:              0,
				Next:              0,
				First:             1,
				Last:              1,
				FirstElementIndex: 0,
				LastElementIndex:  9,
			},
		},
		{
			"second of two full pages",
			args{20, 2, 10},
			&Pagination{
				Total:             20,
				Size:              10,
				Page:              2,
				Prev:              1,
				Next:              0,
				First:             1,
				Last:              2,
				FirstElementIndex: 10,
				LastElementIndex:  19,
			},
		},
		{
			"second page with 1 item",
			args{11, 2, 10},
			&Pagination{
				Total:             11,
				Size:              1,
				Page:              2,
				Prev:              1,
				Next:              0,
				First:             1,
				Last:              2,
				FirstElementIndex: 10,
				LastElementIndex:  10,
			},
		},
		{
			"first page from 2 pages",
			args{11, 1, 10},
			&Pagination{
				Total:             11,
				Size:              10,
				Page:              1,
				Prev:              0,
				Next:              2,
				First:             1,
				Last:              2,
				FirstElementIndex: 0,
				LastElementIndex:  9,
			},
		},
		{
			"first page from 2 pages, one more item",
			args{12, 1, 10},
			&Pagination{
				Total:             12,
				Size:              10,
				Page:              1,
				Prev:              0,
				Next:              2,
				First:             1,
				Last:              2,
				FirstElementIndex: 0,
				LastElementIndex:  9,
			},
		},
		{
			"first page from 3 pages",
			args{22, 1, 10},
			&Pagination{
				Total:             22,
				Size:              10,
				Page:              1,
				Prev:              0,
				Next:              2,
				First:             1,
				Last:              3,
				FirstElementIndex: 0,
				LastElementIndex:  9,
			},
		},
		{
			"second page from 3 pages",
			args{22, 2, 10},
			&Pagination{
				Total:             22,
				Size:              10,
				Page:              2,
				Prev:              1,
				Next:              3,
				First:             1,
				Last:              3,
				FirstElementIndex: 10,
				LastElementIndex:  19,
			},
		},
		{
			"third page from 3 pages",
			args{22, 3, 10},
			&Pagination{
				Total:             22,
				Size:              2,
				Page:              3,
				Prev:              2,
				Next:              0,
				First:             1,
				Last:              3,
				FirstElementIndex: 20,
				LastElementIndex:  21,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPagination(tt.args.collectionLength, tt.args.page, tt.args.perPage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPagination() = %v, want %v", got, tt.want)
			}
		})
	}
}
