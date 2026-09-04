const React = require('react');
const { useRouter } = require('next/router');

const Page = ({ siteTitle, serverData }) => {
  const router = useRouter();
  const Config = require('@/config').Config;
  const SwitcherPage = require(`~/sites/${Config.site}/switcher`);
  const getRouterParams = require(`~/sites/${Config.site}/server/routerParams`);
  const routerParams = getRouterParams(router.query);
  return React.createElement(SwitcherPage.default, {
    routerParams,
    siteTitle,
    props: { serverData },
  });
};

module.exports = Page;
