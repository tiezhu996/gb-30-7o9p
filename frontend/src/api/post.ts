import request from '@/utils/request'
import type { CommunityPost, PostComment } from '@/types/api'
import type { PageData } from '@/types/api'

export function listPosts(params: { page?: number; page_size?: number; post_type?: string; keyword?: string }) {
  return request.get<never, PageData<CommunityPost>>('/posts', { params })
}

export function getPost(id: number | string) {
  return request.get<never, CommunityPost>(`/posts/${id}`)
}

export function createPost(payload: { title: string; content: string; images?: string; post_type?: string }) {
  return request.post<never, CommunityPost>('/posts', payload)
}

export function likePost(id: number) {
  return request.put<never, CommunityPost>(`/posts/${id}/like`)
}

export function listComments(postId: number) {
  return request.get<never, PostComment[]>(`/posts/${postId}/comments`)
}

export function createComment(postId: number, content: string) {
  return request.post<never, PostComment>(`/posts/${postId}/comments`, { content })
}

export function deleteComment(id: number) {
  return request.delete<never, { deleted: boolean }>(`/comments/${id}`)
}
