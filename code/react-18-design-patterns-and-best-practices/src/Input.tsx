// Chapter 2 — Namespaced styled-components consumer (page 56)
import React, { ComponentPropsWithoutRef, FC } from 'react'
import { CSS } from './Input.styled'

export interface Props extends ComponentPropsWithoutRef<'input'> {
  error?: boolean
}

const Input: FC<Props> = ({
  type = 'text',
  error = false,
  value = '',
  disabled = false,
  ...restProps
}) => (
  <CSS.InputWrapper style={error ? { border: '1px solid red' } : {}}>
    <CSS.InputBase type={type} value={value} disabled={disabled} {...restProps} />
  </CSS.InputWrapper>
)
