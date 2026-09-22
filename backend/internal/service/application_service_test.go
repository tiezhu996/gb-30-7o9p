package service

import "testing"

func TestEndReasonText(t *testing.T) {
	cases := map[string]string{
		"rejected":         "审核未通过",
		"final_rejected":   "预留后最终拒绝",
		"holder_abandoned": "获选人放弃",
		"org_cancelled":    "机构取消预留",
		"adopted_by_other": "名额已被其他领养人获得",
		"final_adopted":    "最终领养确认",
		"":                 "",
		"unknown_reason":   "",
	}
	for in, want := range cases {
		if got := EndReasonText(in); got != want {
			t.Errorf("EndReasonText(%q) = %q, want %q", in, got, want)
		}
	}
}
