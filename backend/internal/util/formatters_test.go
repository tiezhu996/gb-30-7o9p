package util

import "testing"

func TestFormatMoney(t *testing.T) {
	if got := FormatMoney(123.5); got != "¥123.50" {
		t.Errorf("FormatMoney = %s", got)
	}
}

func TestAppStatusText(t *testing.T) {
	cases := map[string]string{
		"submitted": "已提交", "approved": "已通过", "rejected": "已拒绝", "x": "未知",
	}
	for in, want := range cases {
		if got := AppStatusText(in); got != want {
			t.Errorf("AppStatusText(%s) = %s, want %s", in, got, want)
		}
	}
}

func TestPetSpeciesText(t *testing.T) {
	if got := PetSpeciesText("cat"); got != "猫" {
		t.Errorf("PetSpeciesText = %s", got)
	}
}
