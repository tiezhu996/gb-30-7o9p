package util

import (
	"fmt"
	"time"
)

// Shared formatters: date, money, status/role/species text live together.

// FormatDate renders YYYY-MM-DD.
func FormatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

// FormatDateTime renders YYYY-MM-DD HH:mm.
func FormatDateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04")
}

// FormatMoney renders a decimal amount with two places.
func FormatMoney(amount float64) string {
	return fmt.Sprintf("¥%.2f", amount)
}

// AppStatusText maps application status to Chinese text.
func AppStatusText(s string) string {
	switch s {
	case "submitted":
		return "已提交"
	case "org_review":
		return "机构审核中"
	case "communicating":
		return "沟通中"
	case "confirmed":
		return "已确认"
	case "offline_interview":
		return "线下面签"
	case "approved":
		return "已通过"
	case "rejected":
		return "已拒绝"
	default:
		return "未知"
	}
}

// RoleText maps a role to Chinese text.
func RoleText(r string) string {
	switch r {
	case "user":
		return "领养人"
	case "org":
		return "救助机构"
	case "admin":
		return "管理员"
	default:
		return "未知"
	}
}

// PetSpeciesText maps a species code to Chinese text.
func PetSpeciesText(s string) string {
	switch s {
	case "dog":
		return "犬"
	case "cat":
		return "猫"
	case "rabbit":
		return "兔"
	default:
		return "其他"
	}
}

// PetStatusText maps a pet status to Chinese text.
func PetStatusText(s string) string {
	switch s {
	case "available":
		return "可领养"
	case "pending":
		return "申请中"
	case "adopted":
		return "已领养"
	default:
		return "未知"
	}
}

// OrgStatusText maps an org status to Chinese text.
func OrgStatusText(s string) string {
	switch s {
	case "pending":
		return "待审核"
	case "approved":
		return "已认证"
	case "rejected":
		return "已驳回"
	default:
		return "未知"
	}
}
