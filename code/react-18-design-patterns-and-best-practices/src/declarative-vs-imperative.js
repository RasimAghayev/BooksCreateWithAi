// Chapter 1 — Declarative vs Imperative (page 33-34)

// ---- toUpperCase: describe WHAT (declarative) vs HOW (imperative) ----
// Imperative: loop through input and push uppercase values into a new array
const toUpperCaseImperative = input => {
  const output = []
  for (let i = 0; i < input.length; i++) {
    output.push(input[i].toUpperCase())
  }
  return output
}

// Declarative: let Array.map return a new array; no mutation of variables
const toUpperCaseDeclarative = input => input.map(value => value.toUpperCase())

// ---- toggle button ----
// Imperative DOM: manually query the node and toggle CSS classes
const toggleButton = document.querySelector('#toggle')
toggleButton.addEventListener('click', () => {
  if (toggleButton.classList.contains('on')) {
    toggleButton.classList.remove('on')
    toggleButton.classList.add('off')
  } else {
    toggleButton.classList.remove('off')
    toggleButton.classList.add('on')
  }
})

// See react-elements.jsx for the declarative React equivalent (<Toggle on />)
