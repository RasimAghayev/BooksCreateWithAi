const { useContext } = require('react');
const { UserContext } = require('../../contexts/user');
const Dashboard = require('./Dashboard');

const DashboardLayout = () => {
  const { connectedUser } = useContext(UserContext);
  if (connectedUser) {
    return <Dashboard connectedUser={connectedUser} />;
  }
  return <div />;
};

module.exports = DashboardLayout;
