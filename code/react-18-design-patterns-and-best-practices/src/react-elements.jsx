// Chapter 1 — How React elements work (pages 35-38)
import ReactDOM from 'react-dom/client'

// JSX for an element (Title is a component; h1 is a DOM node)
const element = (
  <Title color="red">
    <h1>Hello, H1!</h1>
  </Title>
)

// JSX is compiled to plain JS objects via React.createElement():
// {
//   type: Title,           // function -> component; string -> DOM node
//   props: {
//     color: 'red',
//     children: {
//       type: 'h1',
//       props: { children: 'Hello, H1!' }
//     }
//   }
// }

// Rendering a component: mixing logic + templating, with inline styles (#CSSinJS)
const divStyle = {
  color: 'white',
  backgroundImage: `url(${imgUrl})`,
  WebkitTransition: 'all', // capital 'W'
  msTransition: 'all' // 'ms' is the only lowercase vendor prefix
}

return (
  <button style={{ color: 'red' }} onClick={handleClick}>
    Click me!
  </button>
)

ReactDOM.render(<div style={divStyle}>Hello World!</div>, mountNode)

// Declarative toggle equivalent to the imperative DOM version above:
//   <Toggle on />   // on (green)
//   <Toggle />     // off (gray)
