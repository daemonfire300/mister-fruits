package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductSet_In(t *testing.T) {
	cases := []struct {
		Name   string
		In     bool
		Set    ProductSet
		Groups []*ProductGroup
	}{
		{
			"One Banana is in the Banana Set",
			true,
			ProductSet{
				Products: []Product{
					*Banana,
				},
			},
			[]*ProductGroup{
				{Banana, 2},
			},
		},
		{
			"No Banana is in the Banana Set",
			false,
			ProductSet{
				Products: []Product{
					*Banana,
				},
			},
			[]*ProductGroup{},
		},
		{
			"Pear and Banana are in the Banana and Pear Set",
			true,
			ProductSet{
				Products: []Product{
					*Banana,
					*Pear,
				},
			},
			[]*ProductGroup{
				{Banana, 2},
				{Pear, 4},
			},
		},
		{
			"Pear and Banana and Apple are in the Banana and Pear Set",
			false,
			ProductSet{
				Products: []Product{
					*Banana,
					*Pear,
					*Apple,
				},
			},
			[]*ProductGroup{
				{Banana, 2},
				{Pear, 4},
			},
		},
		{
			"Pear and Banana are in the Banana and Pear and Apple Set",
			true,
			ProductSet{
				Products: []Product{
					*Banana,
					*Pear,
				},
			},
			[]*ProductGroup{
				{Banana, 2},
				{Pear, 4},
				{Apple, 1},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.In, tc.Set.In(tc.Groups))
		})
	}
}

func TestProductSet_ExactlyIn(t *testing.T) {
	cases := []struct {
		Name   string
		In     bool
		Set    ProductSet
		Groups []*ProductGroup
	}{
		{
			"One Banana is in the Banana Set",
			true,
			ProductSet{
				Products: []Product{
					*Banana,
				},
			},
			[]*ProductGroup{
				{Banana, 2},
			},
		},
		{
			"No Banana is in the Banana Set",
			false,
			ProductSet{
				Products: []Product{
					*Banana,
				},
			},
			[]*ProductGroup{},
		},
		{
			"Pear and Banana are in the Banana and Pear Set",
			true,
			ProductSet{
				Products: []Product{
					*Banana,
					*Pear,
				},
			},
			[]*ProductGroup{
				{Banana, 2},
				{Pear, 4},
			},
		},
		{
			"Pear and Banana and Apple are in the Banana and Pear Set",
			false,
			ProductSet{
				Products: []Product{
					*Banana,
					*Pear,
					*Apple,
				},
			},
			[]*ProductGroup{
				{Banana, 2},
				{Pear, 4},
			},
		},
		{
			"Pear and Banana are in the Banana and Pear and Apple Set",
			false,
			ProductSet{
				Products: []Product{
					*Banana,
					*Pear,
				},
			},
			[]*ProductGroup{
				{Banana, 2},
				{Pear, 4},
				{Apple, 1},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.In, tc.Set.ExactlyIn(tc.Groups))
		})
	}
}
