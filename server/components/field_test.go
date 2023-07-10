package components

import (
	"image"
	"testing"

	"gocv.io/x/gocv"
)

func TestGetInnerCorners(t *testing.T) {
	type args struct {
		markers [4]Marker
	}
	tests := []struct {
		name string
		args args
		want [4]image.Point
	}{
		{
			name: "test",
			args: args{
				markers: [4]Marker{
					{
						ID: 0,
						Corners: []gocv.Point2f{
							{
								X: 0,
								Y: 3,
							}, {
								X: 1,
								Y: 2,
							}, {
								X: 2,
								Y: 3,
							}, {
								X: 1,
								Y: 4,
							},
						},
					},
					{
						ID: 1,
						Corners: []gocv.Point2f{
							{
								X: 3,
								Y: 0,
							}, {
								X: 4,
								Y: 1,
							}, {
								X: 3,
								Y: 2,
							}, {
								X: 2,
								Y: 1,
							},
						},
					},
					{
						ID: 2,
						Corners: []gocv.Point2f{
							{
								X: 6,
								Y: 3,
							}, {
								X: 5,
								Y: 4,
							}, {
								X: 4,
								Y: 3,
							}, {
								X: 5,
								Y: 2,
							},
						},
					},
					{
						ID: 3,
						Corners: []gocv.Point2f{
							{
								X: 3,
								Y: 4,
							}, {
								X: 4,
								Y: 5,
							}, {
								X: 3,
								Y: 6,
							}, {
								X: 2,
								Y: 5,
							},
						},
					},
				},
			},
			want: [4]image.Point{
				{
					X: 2,
					Y: 3,
				},
				{
					X: 3,
					Y: 2,
				}, {
					X: 4,
					Y: 3,
				},
				{
					X: 3,
					Y: 4,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getInnerCorners(tt.args.markers)
			t.Log(got)
			if got != tt.want {
				t.Errorf("GetInnerCorners() = %v, want %v", got, tt.want)
			}
		})
	}
}
