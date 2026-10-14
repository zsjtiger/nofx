interface IconProps {
  width?: number
  height?: number
  className?: string
}

// Returns the icon for an AI model
export const getModelIcon = (modelType: string, props: IconProps = {}) => {
  // Supports full ID or type name
  const type = modelType.includes('_') ? modelType.split('_').pop() : modelType

  let iconPath: string | null = null

  switch (type) {
    case 'deepseek':
      iconPath = '/icons/deepseek.svg'
      break
    case 'qwen':
      iconPath = '/icons/qwen.svg'
      break
    case 'claude':
      iconPath = '/icons/claude.svg'
      break
    case 'kimi':
      iconPath = '/icons/kimi.svg'
      break
    case 'gemini':
      iconPath = '/icons/gemini.svg'
      break
    case 'grok':
      iconPath = '/icons/grok.svg'
      break
    case 'openai':
      iconPath = '/icons/openai.svg'
      break
    case 'minimax':
      iconPath = '/icons/minimax.svg'
      break
    default:
      return null
  }

  return (
    <img
      src={iconPath}
      alt={`${type} icon`}
      width={props.width || 24}
      height={props.height || 24}
      className={props.className}
    />
  )
}
