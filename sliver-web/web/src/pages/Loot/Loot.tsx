import React from 'react'
import { Card, Empty } from 'antd'

export const Loot: React.FC = () => {
  return (
    <div>
      <h1>Loot</h1>
      <Card>
        <Empty description="No loot collected yet">
        </Empty>
      </Card>
    </div>
  )
}
