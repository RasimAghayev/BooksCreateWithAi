// withInnerWidth HOC: reads window.innerWidth, passes as prop, listens to resize
// From Roldan, "React 18 Design Patterns and Best Practices" Ch 4 (pp. 98-99)

import { useEffect, useState } from 'react';

const withInnerWidth = Component => props => {
  const [innerWidth, setInnerWidth] = useState(window.innerWidth);

  useEffect(() => {
    const handleResize = () => setInnerWidth(window.innerWidth);
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  return <Component {...props} innerWidth={innerWidth} />;
};

export default withInnerWidth;
