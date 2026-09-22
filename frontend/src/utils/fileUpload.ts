import { message } from 'antd'
import axios from 'axios'

export async function uploadFile(file: File): Promise<string> {
  const token = localStorage.getItem('gbadopt_token')
  const form = new FormData()
  form.append('file', file)
  const res = await axios.post('/api/v1/uploads', form, {
    headers: { Authorization: `Bearer ${token}`, 'Content-Type': 'multipart/form-data' },
  })
  message.success('上传成功')
  return res.data.data.url
}
