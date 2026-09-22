import { Empty } from 'antd'
import { parseImages } from '@/constants/pet'

export default function ImageGallery({ imageUrls }: { imageUrls: string }) {
  const images = parseImages(imageUrls)
  if (!images.length) return <Empty description="暂无图片" />
  return (
    <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap' }}>
      {images.map((img, i) => (
        <img key={i} src={img} alt="" style={{ width: 200, height: 150, objectFit: 'cover', borderRadius: 8 }} />
      ))}
    </div>
  )
}
