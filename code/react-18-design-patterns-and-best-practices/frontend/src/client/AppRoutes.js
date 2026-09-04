const { BrowserRouter as Router, Routes, Route } = require('react-router-dom');
const HomePage = require('./pages/home');
const DashboardPage = require('./pages/dashboard');
const LoginPage = require('./pages/login');
const Error404 = require('./pages/error404');

const AppRoutes = () => (
  <Router>
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="/dashboard" element={<DashboardPage />} />
      <Route path="/login" element={<LoginPage />} />
      <Route element={<Error404 />} />
    </Routes>
  </Router>
);

module.exports = AppRoutes;
