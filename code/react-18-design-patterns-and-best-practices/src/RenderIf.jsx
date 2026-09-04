// Conditional rendering: ternary vs let+if (and helper functions for complex conditions)
// RenderIf reusable conditional component
// From Roldan, "React 18 Design Patterns and Best Practices" Ch 3 (pp. 71-73)

import React, { FC, ReactElement } from 'react';

// Verbose let+if+else pattern (avoid for complex cases)
export const VerboseConditional = ({ isLoggedIn }) => {
  let button;
  if (isLoggedIn) {
    button = <LogoutButton />;
  } else {
    button = <LoginButton />;
  }
  return <div>{button}</div>;
};

// Preferred: ternary inline
export const TernaryConditional = ({ isLoggedIn }) => (
  <div>
    {isLoggedIn ? <LogoutButton /> : <LoginButton />}
  </div>
);

// Complex inline condition (hurts readability)
export const InlineComplex = ({ dataIsReady, isAdmin, userHasPermissions }) => (
  <div>
    {dataIsReady && (isAdmin || userHasPermissions) && <SecretData />}
  </div>
);

// Better: extracted helper function with a descriptive name
export const WithHelper = ({ dataIsReady, isAdmin, userHasPermissions }) => {
  const canShowSecretData = () =>
    dataIsReady && (isAdmin || userHasPermissions);

  return (
    <div>
      {canShowSecretData() && <SecretData />}
    </div>
  );
};

// Computed-property helper (isolated + testable)
export const PriceComponent = ({ currency, value }) => {
  const getPrice = () => `${currency}${value}`;
  return <div>{getPrice()}</div>;
};

// Reusable RenderIf component
interface RenderIfProps {
  children: ReactElement | string;
  isTrue?: boolean;
  isFalse?: boolean;
}

export const RenderIf: FC<RenderIfProps> = ({ children, isTrue, isFalse }) => {
  if (isTrue === true) {
    return <>{children}</>;
  }
  if (isFalse === false) {
    return <>{children}</>;
  }
  return null;
};

// List rendering with map (always supply a unique key)
export const UserList = ({ users }) => (
  <ul>
    {users.map(user => (
      <li key={user.id}>{user.name}</li>
    ))}
  </ul>
);

export default RenderIf;

function LogoutButton() { return <button>Logout</button>; }
function LoginButton() { return <button>Login</button>; }
function SecretData() { return <div>Secret</div>; }
