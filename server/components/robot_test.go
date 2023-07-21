package components

import (
	"testing"

	"gocv.io/x/gocv"
)

func TestRobot_Position(t *testing.T) {
	type fields struct {
		angle      float32
		ids        [2]int
		dist       float32
		getMarkers func() map[int]Marker
	}
	tests := []struct {
		name    string
		fields  fields
		want    float32
		want1   float32
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				angle: 0,
				ids:   [2]int{0, 1},
				dist:  1,
				getMarkers: func() map[int]Marker {
					return map[int]Marker{
						0: {
							ID: 0,
							Corners: []gocv.Point2f{
								{X: 0, Y: 1},
								{X: 1, Y: 1},
								{X: 1, Y: 0},
								{X: 0, Y: 0},
							},
						},
					}
				},
			},
			want:    1.5,
			want1:   0.5,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Robot{
				angle:      tt.fields.angle,
				ids:        tt.fields.ids,
				dist:       tt.fields.dist,
				getMarkers: tt.fields.getMarkers,
			}
			got, got1, err := r.Position()
			if (err != nil) != tt.wantErr {
				t.Errorf("Robot.Position() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Robot.Position() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Robot.Position() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
