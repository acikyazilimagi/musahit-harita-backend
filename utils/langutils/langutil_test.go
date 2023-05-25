package langutils

import "testing"

func Test_convertTurkishCharsToEnglish(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "convert turkish characters to english",
			args: args{
				str: "İstanbul",
			},
			want: "Istanbul",
		},
		{
			name: "convert turkish characters to english",
			args: args{
				str: "Şanlıurfa",
			},
			want: "Sanliurfa",
		},
		{
			name: "convert turkish characters to english",
			args: args{
				str: "Çanakkale",
			},
			want: "Canakkale",
		},
		{
			name: "convert turkish characters to english",
			args: args{
				str: "İzmir",
			},
			want: "Izmir",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertTurkishCharsToEnglish(tt.args.str); got != tt.want {
				t.Errorf("convertTurkishCharsToEnglish() = %v, want %v", got, tt.want)
			}
		})
	}
}
