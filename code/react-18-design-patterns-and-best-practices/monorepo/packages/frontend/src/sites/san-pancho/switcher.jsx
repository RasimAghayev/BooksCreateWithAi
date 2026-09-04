const React = require('react');
const dynamic = require('next/dynamic');
const Switcher = require('@/components/Switcher');

const dynamicPages = {
  index: {
    index: dynamic(() => import('./pages/index')),
  },
  login: {
    index: dynamic(() => import('./pages/login')),
  },
  dashboard: {
    index: dynamic(() => import('./pages/dashboard/index')),
  },
};

const SanPanchoSwitcher = ({ routerParams, siteTitle, props }) => (
  React.createElement(Switcher, {
    routerParams,
    siteTitle,
    props,
    dynamicPages,
  })
);

module.exports = SanPanchoSwitcher;
