import React from 'react'

interface IconProps {
  width?: number
  height?: number
  className?: string
}

// Fallback icon: a tile with the first letter of the account type. Only the
// Alpaca paper and CN paper accounts are supported, so no bundled logos are used.
const FallbackIcon: React.FC<IconProps & { label: string }> = ({
  width = 24,
  height = 24,
  className,
  label,
}) => (
  <div
    className={className}
    style={{
      width,
      height,
      borderRadius: 6,
      background: '#E8E2D5',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      fontSize: Math.max(10, (width || 24) * 0.4),
      fontWeight: 'bold',
      color: '#1A1813',
      flexShrink: 0,
    }}
  >
    {label[0]?.toUpperCase() || '?'}
  </div>
)

// Returns the icon for an exchange account type (e.g. "alpaca_paper", "cn_paper")
export const getExchangeIcon = (
  exchangeType: string,
  props: IconProps = {}
) => {
  const label = exchangeType.toLowerCase()
  return (
    <FallbackIcon
      width={props.width || 24}
      height={props.height || 24}
      className={props.className}
      label={label}
    />
  )
}
