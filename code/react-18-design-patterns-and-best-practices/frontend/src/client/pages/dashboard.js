const { UserProvider } = require('../contexts/user');
const DashboardLayout = require('../components/dashboard/DashboardLayout');

const DashboardPage = () => (
  <UserProvider>
    <DashboardLayout />
  </UserProvider>
);

module.exports = DashboardPage;
