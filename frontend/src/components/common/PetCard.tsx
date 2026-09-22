import { Card, Tag } from 'antd'
import { useNavigate } from 'react-router-dom'
import { PetSpeciesMap, PetStatusMap, parseImages, type Pet } from '@/constants/pet'

export default function PetCard({ pet }: { pet: Pet }) {
  const navigate = useNavigate()
  const cover = parseImages(pet.image_urls)[0]
  const status = PetStatusMap[pet.status]
  return (
    <Card
      hoverable
      className="pet-card"
      onClick={() => navigate(`/pets/${pet.id}`)}
      cover={<img alt={pet.name} src={cover} style={{ height: 180, objectFit: 'cover' }} />}
    >
      <Card.Meta
        title={
          <span>
            {pet.name} <Tag color="blue">{PetSpeciesMap[pet.species]}</Tag>
            <Tag color={status.color}>{status.text}</Tag>
          </span>
        }
        description={`${pet.breed || '未知品种'} · ${pet.age}岁 · ${pet.city || '未知城市'}`}
      />
    </Card>
  )
}
