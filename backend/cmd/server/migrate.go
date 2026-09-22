package main

import (
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/gbadopt/gbadopt/internal/constants"
	"github.com/gbadopt/gbadopt/internal/model"
)

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Organization{},
		&model.Pet{},
		&model.AdoptionApplication{},
		&model.VisitReview{},
		&model.CommunityPost{},
		&model.PostComment{},
		&model.Donation{},
		&model.DonationUsage{},
		&model.Favorite{},
	)
}

func seed(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.User{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	logger := slog.Default()

	adminHash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	userHash, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	orgHash, _ := bcrypt.GenerateFromPassword([]byte("org123"), bcrypt.DefaultCost)
	admin := &model.User{Username: "admin", Email: "admin@gbadopt.local", PasswordHash: string(adminHash), Nickname: "平台管理员", Role: "admin"}
	user := &model.User{Username: "adopter", Email: "adopter@gbadopt.local", PasswordHash: string(userHash), Nickname: "爱心领养人", Role: "user"}
	orgUser := &model.User{Username: "shelter", Email: "shelter@gbadopt.local", PasswordHash: string(orgHash), Nickname: "暖窝救助站", Role: "org"}
	if err := db.Create(admin).Error; err != nil {
		return err
	}
	if err := db.Create(user).Error; err != nil {
		return err
	}
	if err := db.Create(orgUser).Error; err != nil {
		return err
	}

	org := &model.Organization{
		UserID: orgUser.ID, Name: "暖窝动物救助站", CertType: "registered",
		Status: constants.OrgStatusApproved, Contact: "13800000000", City: "上海",
		Description: "致力于流浪猫狗救助与领养的专业机构。",
	}
	org2 := &model.Organization{
		UserID: orgUser.ID, Name: "城市伴侣宠物收容所", CertType: "registered",
		Status: constants.OrgStatusApproved, Contact: "13900000000", City: "北京",
		Description: "提供宠物收容、医疗与领养服务。",
	}
	if err := db.Create(org).Error; err != nil {
		return err
	}
	if err := db.Create(org2).Error; err != nil {
		return err
	}

	pets := []model.Pet{
		{OrgID: org.ID, Name: "旺财", Species: "dog", Breed: "中华田园犬", Age: 2, Gender: "male", Size: "medium", City: "上海", Description: "性格温顺忠诚，已绝育疫苗齐全。", Personality: "亲人活泼", HealthStatus: "健康", Neutered: true, Vaccinated: true, ImageURLs: `["https://images.unsplash.com/photo-1543466835-00a7907e9de1?w=600"]`, Status: "available"},
		{OrgID: org.ID, Name: "雪球", Species: "cat", Breed: "英短", Age: 1, Gender: "female", Size: "small", City: "上海", Description: "安静粘人的小猫咪，已驱虫。", Personality: "温顺", HealthStatus: "健康", Neutered: true, Vaccinated: true, ImageURLs: `["https://images.unsplash.com/photo-1514888286974-6c03e2ca1dba?w=600"]`, Status: "available"},
		{OrgID: org2.ID, Name: "跳跳", Species: "rabbit", Breed: "垂耳兔", Age: 1, Gender: "male", Size: "small", City: "北京", Description: "活泼好动的垂耳兔，喜欢胡萝卜。", Personality: "活泼", HealthStatus: "健康", Neutered: false, Vaccinated: false, ImageURLs: `["https://images.unsplash.com/photo-1585110396000-c9ffd4e4b308?w=600"]`, Status: "available"},
		{OrgID: org2.ID, Name: "豆豆", Species: "dog", Breed: "柯基", Age: 3, Gender: "male", Size: "small", City: "北京", Description: "短腿萌宠，粘人爱撒娇。", Personality: "粘人", HealthStatus: "健康", Neutered: true, Vaccinated: true, ImageURLs: `["https://images.unsplash.com/photo-1529778873920-4da4926a72c2?w=600"]`, Status: "available"},
	}
	if err := db.Create(&pets).Error; err != nil {
		return err
	}

	posts := []model.CommunityPost{
		{UserID: orgUser.ID, OrgID: org.ID, Title: "旺财的救助故事：从流浪到新生", Content: "旺财是在街角被发现的流浪狗，经过治疗与照顾，如今已经健康活泼，等待有缘家庭领养。", PostType: "story", LikeCount: 32, CommentCount: 6},
		{UserID: user.ID, Title: "寻主公告：走失的橘猫", Content: "昨天在小区附近捡到一只橘猫，脖子上有红色项圈，请失主联系。", PostType: "lost_notice", LikeCount: 12, CommentCount: 3},
	}
	if err := db.Create(&posts).Error; err != nil {
		return err
	}

	comments := []model.PostComment{
		{PostID: posts[0].ID, UserID: user.ID, Content: "太暖心了，希望旺财早日找到新家！"},
		{PostID: posts[1].ID, UserID: orgUser.ID, Content: "已转发，希望猫咪早日回家。"},
	}
	if err := db.Create(&comments).Error; err != nil {
		return err
	}

	donations := []model.Donation{
		{UserID: user.ID, OrgID: org.ID, Amount: 100, TransactionID: "sandbox_1001", Status: "success"},
		{UserID: user.ID, OrgID: org2.ID, Amount: 200, TransactionID: "sandbox_1002", Status: "success"},
	}
	if err := db.Create(&donations).Error; err != nil {
		return err
	}

	usages := []model.DonationUsage{
		{OrgID: org.ID, DonationID: donations[0].ID, Amount: 100, UsageDesc: "采购猫粮与驱虫药", ProofURL: ""},
	}
	if err := db.Create(&usages).Error; err != nil {
		return err
	}

	apps := []model.AdoptionApplication{
		{UserID: user.ID, PetID: pets[0].ID, OrgID: org.ID, Questionnaire: `{"has_yard":false,"pet_experience":"有养狗经验"}`, Status: "submitted"},
	}
	if err := db.Create(&apps).Error; err != nil {
		return err
	}

	reviews := []model.VisitReview{
		{ApplicationID: apps[0].ID, UserID: user.ID, OrgID: org.ID, ScheduledDays: 30, DueDate: time.Now().AddDate(0, 0, 30), Status: "pending"},
	}
	if err := db.Create(&reviews).Error; err != nil {
		return err
	}

	logger.Info("gbadopt seed data created", "users", 3, "orgs", 2, "pets", len(pets), "posts", len(posts), "apps", len(apps))
	return nil
}
