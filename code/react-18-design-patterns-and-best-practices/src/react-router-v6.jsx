// React Router v6 (data router) example: routes, params, loaders, navigation state
// From Roldan, "React 18 Design Patterns and Best Practices" Ch 9 (pp. 205-220)

import {
  createBrowserRouter,
  createRoutesFromElements,
  Route,
  Link,
  Outlet,
  RouterProvider,
  useLoaderData,
  useNavigation,
  useParams,
} from 'react-router-dom';

// --- Page components ---
const Home = () => (
  <div className="Home"><h1>Home</h1></div>
);
const About = () => (
  <div className="About"><h1>About</h1></div>
);
const Error404 = () => (
  <div className="Error404"><h1>Error404</h1></div>
);

// Root layout: holds the nav menu and renders child routes via <Outlet />
const Root = () => (
  <>
    <ul className="menu">
      <li><Link to="/">Home</Link></li>
      <li><Link to="/about">About</Link></li>
      <li><Link to="/pokemons">Pokemons</Link></li>
    </ul>
    <div><Outlet /></div>
  </>
);

// --- Route params example: /contacts and /contacts/:contactId ---
const contacts = [
  { id: 1, name: 'Carlos Santana', email: 'carlos.santana@dev.education', phone: '415-307-3112' },
  { id: 2, name: 'John Smith', email: 'john.smith@dev.education', phone: '223-344-5122' },
  { id: 3, name: 'Alexis Nelson', email: 'alexis.nelson@dev.education', phone: '664-291-4477' },
];

const Contacts = () => {
  const { contactId = '0' } = useParams();
  const selectedContact =
    Number(contactId) > 0
      ? contacts.find(c => c.id === Number(contactId))
      : undefined;

  if (selectedContact) {
    return (
      <>
        <h2>{selectedContact.name}</h2>
        <p>{selectedContact.email}</p>
        <p>{selectedContact.phone}</p>
      </>
    );
  }

  return (
    <div className="Contacts">
      <h1>Contacts</h1>
      <ul>
        {contacts.map(c => (
          <li key={c.id}>
            <Link to={`/contacts/${c.id}`}>{c.name}</Link>
          </li>
        ))}
      </ul>
    </div>
  );
};

// --- Data router + loader (v6.4+) ---
const imgUrl =
  'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/';

const Pokemons = () => {
  const pokemons = useLoaderData();
  const navigation = useNavigation();

  if (navigation.state === 'loading') return <h1>Loading...</h1>;

  return (
    <div className="Home">
      <h1>Pokemons</h1>
      {pokemons.map((pokemon, index) => (
        <div key={pokemon.name}>
          <h2>{index + 1} {pokemon.name}</h2>
          <img
            src={`${imgUrl}/${pokemon.url.split('/').slice(-2, -1)}.png`}
            alt={pokemon.name}
          />
          <p>
            <a href={pokemon.url} target="_blank" rel="noreferrer">
              {pokemon.url}
            </a>
          </p>
        </div>
      ))}
    </div>
  );
};

const dataLoader = async () => {
  const response = await fetch('https://pokeapi.co/api/v2/pokemon?limit=151');
  const data = await response.json();
  return data.results;
};

// Data router definition (v6.4+ createBrowserRouter API)
const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path="/" element={<Root />}>
      <Route index element={<Home />} />
      <Route path="/about" element={<About />} />
      <Route path="/contacts" element={<Contacts />} />
      <Route path="/contacts/:contactId" element={<Contacts />} />
      <Route path="/pokemons" element={<Pokemons />} loader={dataLoader} />
      <Route path="*" element={<Error404 />} />
    </Route>
  )
);

const App = () => (
  <div className="App">
    <RouterProvider router={router} />
  </div>
);

export default App;
