const React = require('react');

const ErrorPage = ({ siteTitle }) => (
  <div className="error-page">
    <h1>Error 404</h1>
    <p>Page not found for {siteTitle}</p>
  </div>
);

module.exports = ErrorPage;
