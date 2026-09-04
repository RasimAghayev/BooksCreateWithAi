// Chapter 2 — Implementing interfaces/types, interface merging, template-literal unions (pages 51-56)

interface IWork {
  company: string
  position: string
}

type TWork = {
  company: string
  position: string
}

// A type alias that names a UNION cannot be implemented by a class
// (Figure 2.3: class can only implement an object type / intersection of object
// types with statically known members).
type TWork2 = {
  company: string
  position: string
} | {
  name: string
  age: number
}

class Person implements IWork {
  name: 'Carlos'
  age: 35
}
class Person2 implements TWork {
  name: 'Cristina'
  age: 34
}
class Person3 implements TWork2 {
  company: 'Google'
  position: 'Senior Software Engineer'
}

// Interfaces MERGE: redeclaring IUser adds `country`.
interface IUser {
  username: string
  email: string
  name: string
  age?: number
  website: string
  active: boolean
}
interface IUser {
  country: string
}

const user: IUser = {
  username: 'czantany',
  email: 'carlos@milkzoft.com',
  name: 'Carlos Santana',
  country: 'Mexico',
  age: 35,
  website: 'http://www.js.education',
  active: true
}

// Template literal / string-literal union types: compile-time value safety.
type Theme = 'light' | 'dark'
