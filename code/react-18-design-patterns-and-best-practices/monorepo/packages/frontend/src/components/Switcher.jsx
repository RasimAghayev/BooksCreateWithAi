const React = require('react');

const Switcher = ({ routerParams, props = {}, sitePages }) => {
  const {
    page,
    section = 'index',
    subSection = '',
    urlParams,
    queryParams = {},
  } = routerParams;

  const extraProps = {
    queryParams,
    router: { section, subSection },
    urlParams,
  };

  let PageToRender;
  let sectionPages = {};

  if (sitePages[page]) {
    sectionPages = sitePages[page];
    PageToRender = sectionPages.index;

    if (sectionPages[section] && sectionPages[section][subSection]) {
      PageToRender = sectionPages[section][subSection];
    } else if (section !== 'index') {
      PageToRender = sectionPages[section].index;
    }
  } else {
    const ErrorPage = require('@/components/ErrorPage');
    PageToRender = ErrorPage;
  }

  return React.createElement(PageToRender, { ...props, ...extraProps });
};

module.exports = Switcher;
