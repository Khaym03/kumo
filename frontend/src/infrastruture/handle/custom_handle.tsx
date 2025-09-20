import { Handle, type HandleProps } from '@xyflow/react'
import { memo } from 'react'

// const styles = {
//     width: 8,
//     height: 8,
//     backgroundColor: 'var(--secondary)',
//     border: '1px solid var(--muted-foreground)',
// }

const CustomHandle = memo((props: HandleProps) => {
  return <Handle className='!w-2 !h-2 !bg-muted !border-muted-foreground' {...props} />
})

export default CustomHandle
