# Cheat Sheet — Understanding GraphQL with a Real Project

## GraphQL + Apollo

### Apollo Server qurma

**Nə edir:** GraphQL server yaradır, type-defs ilə schema.

**Kod:**

```ts
import { ApolloServer, gql } from 'apollo-server';
const typeDefs = gql`
  type User { id: ID!, name: String!, email: String! }
  type Query { users: [User!]! }
`;
const server = new ApolloServer({ typeDefs, resolvers });
server.listen().then(({ url }) => console.log(url));
```

**Mənbə:** Chapter 13, page 281-353

### Apollo Client qurma

**Nə edir:** GraphQL client React-ə inject olunur.

**Kod:**

```ts
import { ApolloClient, InMemoryCache, ApolloProvider } from '@apollo/client';
const client = new ApolloClient({
  uri: 'http://localhost:4000/graphql',
  cache: new InMemoryCache(),
});
<ApolloProvider client={client}><App /></ApolloProvider>
```

**Mənbə:** Chapter 13, page 281-353

### useQuery

**Nə edir:** Sorğu: data, loading, error.

**Kod:**

```tsx
import { gql, useQuery } from '@apollo/client';
const GET_USERS = gql`query { users { id name } }`;
function Users() {
  const { loading, error, data } = useQuery(GET_USERS);
  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error</p>;
  return <ul>{data.users.map(u => <li key={u.id}>{u.name}</li>)}</ul>;
}
```

**Mənbə:** Chapter 13, page 281-353

### useMutation

**Nə edir:** Mutation çağırışı, variables ilə arqumentlər.

**Kod:**

```tsx
const LOGIN = gql`
  mutation Login($email: String!, $password: String!) {
    login(email: $email, password: $password) { token }
  }
`;
const [login, { data, loading, error }] = useMutation(LOGIN);
await login({ variables: { email, password } });
```

**Mənbə:** Chapter 13, page 281-353

### JWT auth

**Nə edir:** Token yaratma. `jwt.verify` ilə doğrulama.

**Kod:**

```ts
import jwt from 'jsonwebtoken';
const token = jwt.sign({ userId: user.id }, SECRET, { expiresIn: '1h' });
```

**Mənbə:** Chapter 13, page 281-353
