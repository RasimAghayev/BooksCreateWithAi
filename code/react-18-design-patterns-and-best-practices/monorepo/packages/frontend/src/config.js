const { is } = require('@web-creator/utils');
const blankPageConfig = require('./sites/blank-page/config').config;
const sanPanchoConfig = require('./sites/san-pancho/config').config;
const { Site } = require('./types/config');

const isProduction = process.env.NODE_ENV === 'production';
const isLocal = process.env.LOCAL === 'true';
const isLocalProduction = isProduction && isLocal;

const getSiteConfig = (site) => {
  switch (site) {
    case Site.SanPancho:
      return sanPanchoConfig;
    default:
      return blankPageConfig;
  }
};

const buildConfig = () => {
  let site = process.env.SITE || 'blank-page';

  if (typeof window !== 'undefined') {
    const props = window.__NEXT_DATA__ && window.__NEXT_DATA__.props;
    if (props && props.site) {
      site = props.site;
    }
  } else if (!site) {
    throw 'You must specify a site (E.g. SITE=san-pancho npm run dev)';
  }

  const siteConfig = getSiteConfig(site);

  const config = {
    ...siteConfig,
    api: {
      uri: isProduction && !isLocalProduction
        ? `https://${siteConfig.domainName}/graphql`
        : 'http://localhost:4000/graphql',
    },
    site,
    homeUrl: `https://${siteConfig.domainName}`,
    hostname: isProduction && !isLocalProduction ? siteConfig.domainName : 'localhost',
    mode: isProduction ? 'production' : 'development',
  };

  return config;
};

const Config = buildConfig();
module.exports = { Config, getSiteConfig, buildConfig };
