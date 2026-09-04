// Chapter 2 — Types: describing objects; required vs optional fields (page 47-48)

// Good practice: prefix type aliases with 'T' so they are not confused with a
// class or a React component.
type TUser = {
  username: string
  email: string
  name: string
  age: number
  website: string
  active: boolean
}

const user: TUser = {
  username: 'czantany',
  email: 'carlos@milkzoft.com',
  name: 'Carlos Santana',
  age: 33,
  website: 'http://www.js.education',
  active: true
}
// Omitting or passing an invalid value for `age` would raise:
// "Figure 2.2: Age is missing in type User but is required"
// Suppose we insert this data with Sequelize:
// models.User.create({ ...user })

// Optional property: prefix with '?' so the node may be omitted.
type TUserOptional = {
  username: string
  email: string
  name: string
  age?: number // optional
  website: string
  active: boolean
}
