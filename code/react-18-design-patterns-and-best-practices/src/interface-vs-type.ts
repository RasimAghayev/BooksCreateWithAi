// Chapter 2 — Interfaces vs types: extending & implementing (page 49-50)

// Good practice: prefix interfaces with 'I'.
interface IUser {
  username: string
  email: string
  name: string
  age?: number
  website: string
  active: boolean
}

// Extending an interface uses `extends`.
interface IWork {
  company: string
  position: string
}
interface IPerson extends IWork {
  name: string
  age: number
}

// Extending a type uses the intersection operator `&`.
type TWork = {
  company: string
  position: string
}
type TPerson = TWork & {
  name: string
  age: number
}

// An interface can also be extended into a type (type uses &).
type TPerson2 = IWork & {
  name: string
  age: number
}

// Rule: a type alias can be extended with & only if it does NOT name a union
// type. A class may `implements` an interface or a non-union type alias.
