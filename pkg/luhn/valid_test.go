package luhn

import "testing"

func TestValid(t *testing.T) {
	type args struct {
		number int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "4950 2174 6090 5278",
			args: args{number: 4950217460905278},
			want: true,
		},
		{
			name: "9278923470",
			args: args{
				number: 9278923470,
			},
			want: true,
		},
		{
			name: "12345678903",
			args: args{
				number: 12345678903,
			},
			want: true,
		},
		{
			name: "zero",
			args: args{number: 0},
			want: true,
		},
		{
			name: "invalid 1",
			args: args{
				number: 1,
			},
			want: false,
		},
		{
			name: "invalid card number",
			args: args{
				number: 495021746090521,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Valid(tt.args.number); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
