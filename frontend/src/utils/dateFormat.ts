import dayjs from 'dayjs'

export function formatDate(input?: string | number | Date): string {
  if (!input) return '-'
  const d = dayjs(input)
  return d.isValid() ? d.format('YYYY-MM-DD') : '-'
}

export function formatDateTime(input?: string | number | Date): string {
  if (!input) return '-'
  const d = dayjs(input)
  return d.isValid() ? d.format('YYYY-MM-DD HH:mm') : '-'
}
