// React Hooks migration example: class (Issues) → functional with useState/useEffect
// From Roldan, "React 18 Design Patterns and Best Practices" Ch 8 (pp. 164-169)

const Issue = {
  number: 0,
  title: '',
  state: '',
};

// Class component version (legacy) - equivalent functional Hooks version below
// class Issues extends Component<Props, State> {
//   constructor(props) {
//     super(props);
//     this.state = { issues: [] };
//   }
//   componentDidMount() {
//     axios.get(url).then(res => this.setState({ issues: res.data }));
//   }
//   render() { ... }
// }

// Functional component with Hooks (migrated)
// - useState replaces constructor state + this.setState
// - useEffect(() => {...}, []) replaces componentDidMount
import { FC, useState, useEffect } from 'react';
import axios from 'axios';

export const Issues: FC = () => {
  const [issues, setIssues] = useState([]);

  // Empty dependency array = run once on mount (componentDidMount equivalent)
  // NOTE: fires after paint (deferred); use useLayoutEffect if you must mutate DOM before paint.
  useEffect(() => {
    axios
      .get('https://api.github.com/repos/ContentPI/ContentPI/issues')
      .then((response) => {
        setIssues(response.data);
      });
  }, []);

  return (
    <>
      <h1>ContentPI Issues</h1>
      {issues.map((issue) => (
        <p key={issue.title}>
          <strong>#{issue.number}</strong>{' '}
          <a
            href={`https://github.com/ContentPI/ContentPI/issues/${issue.number}`}
            target="_blank"
            rel="noopener noreferrer"
          >
            {issue.title}
          </a>{' '}
          {issue.state}
        </p>
      ))}
    </>
  );
};
