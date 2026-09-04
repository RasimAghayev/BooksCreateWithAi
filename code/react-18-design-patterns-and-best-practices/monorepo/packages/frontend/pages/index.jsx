const React = require('react');

const Page = ({ siteTitle }) => {
  const SwitcherPage = require(`~/sites/${require('@/config').Config.site}/switcher`);
  const getRouterParams = require(`~/sites/${require('@/config').Config.site}/server/routerParams`);
  const routerParams = getRouterParams({});
  return React.createElement(SwitcherPage.default, { routerParams, siteTitle });
};

module.exports = Page;
