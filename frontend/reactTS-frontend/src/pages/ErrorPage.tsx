import React from 'react';
import { useRouteError, Link } from 'react-router-dom';

const ErrorPage: React.FC = () => {
  const error: any = useRouteError();

  return (
    <div style={{ textAlign: 'center', marginTop: '50px' }}>
      <h1>Oops!</h1>
      <p>Sorry, an unexpected error has occurred.</p>
      <p>
        <i>{error.statusText || error.message}</i>
      </p>
      <Link to="/">Go back to Home Page</Link>
    </div>
  );
};

export default ErrorPage;