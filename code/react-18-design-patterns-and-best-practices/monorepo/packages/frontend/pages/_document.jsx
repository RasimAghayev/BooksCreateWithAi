const React = require('react');
const Document = require('next/document').default;
const { Head, Html, Main, NextScript } = require('next/document');
const { ServerStyleSheet } = require('styled-components');
const Config = require('@/config');

class MyDocument extends Document {
  static async getInitialProps(ctx) {
    const sheet = new ServerStyleSheet();
    const originalRenderPage = ctx.renderPage;
    try {
      ctx.renderPage = () =>
        originalRenderPage({
          enhanceApp: (App) => (props) =>
            sheet.collectStyles(
              React.createElement('body', { className: 'theme--light' },
                React.createElement(App, { ...props, title: Config.siteTitle })
              )
            ),
        });
      const initialProps = await Document.getInitialProps(ctx);
      return {
        initialProps,
        styles: (
          React.createElement(React.Fragment, null,
            initialProps.styles,
            sheet.getStyleElement()
          )
        ),
      };
    } finally {
      sheet.seal();
    }
  }

  render() {
    return (
      <Html>
        <Head>
          <link rel="icon" type="image/x-icon" href="/images/favicon.png" />
        </Head>
        <Main />
        <NextScript />
      </Html>
    );
  }
}

module.exports = MyDocument;
