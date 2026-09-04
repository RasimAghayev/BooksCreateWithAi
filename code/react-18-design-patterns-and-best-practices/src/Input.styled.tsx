// Chapter 2 — Namespaces: themed styled-components (page 55)
import styled from 'styled-components'

export namespace CSS {
  export const InputWrapper = styled.div`
    padding: 10px;
    margin: 0;
    background: white;
    width: 250px;
  `
  export const InputBase = styled.input`
    width: 100%;
    background: transparent;
    border: none;
    font-size: 14px;
  `
}
