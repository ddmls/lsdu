package human

import "testing"

func Test_constants(t *testing.T) {
	tests := []int64{KiB, MiB, GiB, TiB, PiB}
	size := int64(1024)
	for _, test := range tests {
		if test != size {
			t.Errorf("%v != %v", test, size)
		}
		size *= 1024
	}
}

func Test_Humanize(t *testing.T) {
	type args struct {
		size int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"zero", args{int64(0)}, "0b"},
		{"50 bytes", args{int64(50)}, "50b"},
		{"1023 bytes", args{int64(1023)}, "1023b"},
		{"1 KiB", args{KiB}, "1.0K"},
		{"1 KiB + 1b", args{KiB + 1}, "1.0K"},
		{"1 KiB - 1b", args{KiB - 1}, "1023b"},
		{"2 MiB", args{2 * MiB}, "2.0M"},
		{"2 MiB - 1b", args{2*MiB - 1}, "2.0M"},
		{"2 MiB + 10b", args{2*MiB + 10}, "2.0M"},
		{"667 MiB + 1030b", args{667*MiB + 1030}, "667.0M"},
		{"4 GiB + 10Kib", args{4*GiB + 10*KiB}, "4.0G"},
		{"2 TiB - 1", args{2*TiB - 1}, "2.0T"},
		{"1 PiB", args{PiB}, "1.0P"},
		{"1.4 MiB", args{14 * MiB / 10}, "1.4M"},
		{"2.86 GiB", args{286 * GiB / 100}, "2.9G"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Humanize(tt.args.size); got != tt.want {
				t.Errorf("humanize() = %v, want %v", got, tt.want)
			}
		})
	}
}
