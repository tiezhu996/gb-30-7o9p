import axios from 'axios'
import { message } from 'antd'
import { useAuthStore } from '@/stores/authStore'

const request = axios.create({
  baseURL: '/api/v1',
  timeout: 15000,
})

request.interceptors.request.use((config) => {
  const token = localStorage.getItem('gbadopt_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

request.interceptors.response.use(
  (res) => {
    const body = res.data
    if (body && typeof body.code === 'number' && body.code !== 0) {
      message.error(body.message || '请求失败')
      return Promise.reject(new Error(body.message))
    }
    return body?.data
  },
  (err) => {
    const status = err.response?.status
    const msg = err.response?.data?.message || '网络异常'
    if (status === 401) {
      localStorage.removeItem('gbadopt_token')
      useAuthStore.getState().setUser(null)
    }
    message.error(msg)
    return Promise.reject(err)
  },
)

export default request
