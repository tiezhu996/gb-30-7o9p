-- gbadopt 宠物领养平台 初始化脚本（PostgreSQL）
-- 由 postgres 官方镜像 /docker-entrypoint-initdb.d 首次启动时自动执行。

CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR(64) NOT NULL,
  email VARCHAR(128) NOT NULL,
  CONSTRAINT uni_users_username UNIQUE (username),
  CONSTRAINT uni_users_email UNIQUE (email),
  password_hash VARCHAR(255) NOT NULL,
  nickname VARCHAR(64),
  avatar VARCHAR(255),
  phone VARCHAR(32),
  role VARCHAR(16) NOT NULL DEFAULT 'user',
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS organizations (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  name VARCHAR(128) NOT NULL,
  license_url VARCHAR(255),
  cert_type VARCHAR(32),
  status VARCHAR(16) NOT NULL DEFAULT 'pending',
  contact VARCHAR(64),
  city VARCHAR(64),
  description TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT fk_org_user FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS pets (
  id BIGSERIAL PRIMARY KEY,
  org_id BIGINT NOT NULL,
  name VARCHAR(64) NOT NULL,
  species VARCHAR(16) NOT NULL,
  breed VARCHAR(64),
  age INT DEFAULT 0,
  gender VARCHAR(8),
  size VARCHAR(16),
  city VARCHAR(64),
  description TEXT,
  personality VARCHAR(255),
  health_status VARCHAR(128),
  neutered BOOLEAN DEFAULT FALSE,
  vaccinated BOOLEAN DEFAULT FALSE,
  image_urls JSONB DEFAULT '[]',
  status VARCHAR(16) NOT NULL DEFAULT 'available',
  created_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT fk_pet_org FOREIGN KEY (org_id) REFERENCES organizations(id)
);

CREATE TABLE IF NOT EXISTS adoption_applications (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  pet_id BIGINT NOT NULL,
  org_id BIGINT NOT NULL,
  questionnaire JSONB DEFAULT '{}',
  status VARCHAR(32) NOT NULL DEFAULT 'submitted',
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT fk_app_user FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT fk_app_pet FOREIGN KEY (pet_id) REFERENCES pets(id)
);

CREATE TABLE IF NOT EXISTS visit_reviews (
  id BIGSERIAL PRIMARY KEY,
  application_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  org_id BIGINT,
  scheduled_days INT DEFAULT 30,
  due_date DATE,
  status VARCHAR(16) NOT NULL DEFAULT 'pending',
  photos JSONB DEFAULT '[]',
  note TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT fk_review_app FOREIGN KEY (application_id) REFERENCES adoption_applications(id)
);

CREATE TABLE IF NOT EXISTS community_posts (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  org_id BIGINT,
  title VARCHAR(255) NOT NULL,
  content TEXT,
  images JSONB DEFAULT '[]',
  post_type VARCHAR(16) NOT NULL DEFAULT 'story',
  like_count INT DEFAULT 0,
  comment_count INT DEFAULT 0,
  status VARCHAR(16) NOT NULL DEFAULT 'published',
  created_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT fk_post_user FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS post_comments (
  id BIGSERIAL PRIMARY KEY,
  post_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  content VARCHAR(1000) NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT fk_comment_post FOREIGN KEY (post_id) REFERENCES community_posts(id)
);

CREATE TABLE IF NOT EXISTS donations (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  org_id BIGINT NOT NULL,
  amount NUMERIC(12,2) NOT NULL,
  transaction_id VARCHAR(64),
  status VARCHAR(16) NOT NULL DEFAULT 'pending',
  created_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT fk_donation_user FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT fk_donation_org FOREIGN KEY (org_id) REFERENCES organizations(id)
);

CREATE TABLE IF NOT EXISTS donation_usages (
  id BIGSERIAL PRIMARY KEY,
  org_id BIGINT NOT NULL,
  donation_id BIGINT,
  amount NUMERIC(12,2) NOT NULL,
  usage_desc VARCHAR(512) NOT NULL,
  proof_url VARCHAR(255),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT fk_usage_org FOREIGN KEY (org_id) REFERENCES organizations(id)
);

CREATE TABLE IF NOT EXISTS favorites (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  target_type VARCHAR(16) NOT NULL,
  target_id BIGINT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT idx_fav_user_target UNIQUE (user_id, target_type, target_id)
);

-- 种子数据
INSERT INTO users (username, email, password_hash, nickname, role) VALUES
  ('admin', 'admin@gbadopt.local', '$2a$10$TUjlusy5GigbkU945MoMauH6AXrUFCuj4tMR6EFdXtbKHi5kyigiO', '平台管理员', 'admin'),
  ('adopter', 'adopter@gbadopt.local', '$2a$10$3KG3E2lKg2QAZs4dnzmfqePQJYGFKD9HEwG6hvBThDMlW1pZYb46a', '爱心领养人', 'user'),
  ('shelter', 'shelter@gbadopt.local', '$2a$10$eDmmXgXm9Y8byb0hzYBBC.mVwXOqTtuxWzY0bub/0nRWGYgYLpE1O', '暖窝救助站', 'org');

INSERT INTO organizations (user_id, name, cert_type, status, contact, city, description) VALUES
  (3, '暖窝动物救助站', 'registered', 'approved', '13800000000', '上海', '致力于流浪猫狗救助与领养的专业机构。'),
  (3, '城市伴侣宠物收容所', 'registered', 'approved', '13900000000', '北京', '提供宠物收容、医疗与领养服务。');

INSERT INTO pets (org_id, name, species, breed, age, gender, size, city, description, personality, health_status, neutered, vaccinated, image_urls, status) VALUES
  (1, '旺财', 'dog', '中华田园犬', 2, 'male', 'medium', '上海', '性格温顺忠诚，已绝育疫苗齐全。', '亲人活泼', '健康', TRUE, TRUE, '["https://images.unsplash.com/photo-1543466835-00a7907e9de1?w=600"]', 'available'),
  (1, '雪球', 'cat', '英短', 1, 'female', 'small', '上海', '安静粘人的小猫咪，已驱虫。', '温顺', '健康', TRUE, TRUE, '["https://images.unsplash.com/photo-1514888286974-6c03e2ca1dba?w=600"]', 'available'),
  (2, '跳跳', 'rabbit', '垂耳兔', 1, 'male', 'small', '北京', '活泼好动的垂耳兔，喜欢胡萝卜。', '活泼', '健康', FALSE, FALSE, '["https://images.unsplash.com/photo-1585110396000-c9ffd4e4b308?w=600"]', 'available'),
  (2, '豆豆', 'dog', '柯基', 3, 'male', 'small', '北京', '短腿萌宠，粘人爱撒娇。', '粘人', '健康', TRUE, TRUE, '["https://images.unsplash.com/photo-1529778873920-4da4926a72c2?w=600"]', 'available');

INSERT INTO community_posts (user_id, org_id, title, content, post_type, like_count, comment_count) VALUES
  (3, 1, '旺财的救助故事：从流浪到新生', '旺财是在街角被发现的流浪狗，经过治疗与照顾，如今已经健康活泼，等待有缘家庭领养。', 'story', 32, 6),
  (2, NULL, '寻主公告：走失的橘猫', '昨天在小区附近捡到一只橘猫，脖子上有红色项圈，请失主联系。', 'lost_notice', 12, 3);

INSERT INTO post_comments (post_id, user_id, content) VALUES
  (1, 2, '太暖心了，希望旺财早日找到新家！'),
  (2, 3, '已转发，希望猫咪早日回家。');

INSERT INTO donations (user_id, org_id, amount, transaction_id, status) VALUES
  (2, 1, 100.00, 'sandbox_1001', 'success'),
  (2, 2, 200.00, 'sandbox_1002', 'success');

INSERT INTO donation_usages (org_id, donation_id, amount, usage_desc) VALUES
  (1, 1, 100.00, '采购猫粮与驱虫药');

INSERT INTO adoption_applications (user_id, pet_id, org_id, questionnaire, status) VALUES
  (2, 1, 1, '{"has_yard":false,"pet_experience":"有养狗经验"}', 'submitted');

INSERT INTO visit_reviews (application_id, user_id, org_id, scheduled_days, due_date, status) VALUES
  (1, 2, 1, 30, CURRENT_DATE + 30, 'pending');
