package algo

import "testing"

func TestIsCentrilized(t *testing.T) {

	cases := []struct {
		name string

		head []int
		x    int
		y    int
		z    int

		want bool
	}{
		{
			name: "in centre",
			head: []int{90, 90, 30},
			x:    180,
			y:    180,
			z:    60,
			want: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := isCentralized(tc.head, tc.x, tc.y, tc.z)
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
