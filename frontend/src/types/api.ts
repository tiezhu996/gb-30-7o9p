export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

export interface PageData<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

export interface UserInfo {
  id: number
  username: string
  email: string
  nickname: string
  avatar: string
  phone: string
  role: 'user' | 'org' | 'admin'
  created_at: string
}

export interface AdoptionApplication {
  id: number
  user_id: number
  pet_id: number
  org_id: number
  questionnaire: string
  status: string
  created_at: string
  updated_at: string
}

export interface VisitReview {
  id: number
  application_id: number
  user_id: number
  org_id: number
  scheduled_days: number
  due_date: string
  status: 'pending' | 'submitted' | 'overdue'
  photos: string
  note: string
  created_at: string
}

export interface CommunityPost {
  id: number
  user_id: number
  org_id: number
  title: string
  content: string
  images: string
  post_type: 'story' | 'lost_notice'
  like_count: number
  comment_count: number
  status: string
  created_at: string
}

export interface PostComment {
  id: number
  post_id: number
  user_id: number
  content: string
  created_at: string
}

export interface Donation {
  id: number
  user_id: number
  org_id: number
  amount: number
  transaction_id: string
  status: string
  created_at: string
}

export interface DonationUsage {
  id: number
  org_id: number
  donation_id: number
  amount: number
  usage_desc: string
  proof_url: string
  created_at: string
}

export interface Favorite {
  id: number
  user_id: number
  target_type: string
  target_id: number
  created_at: string
}
