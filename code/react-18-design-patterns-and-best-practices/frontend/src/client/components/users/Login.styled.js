const styled = require('styled-components');

const StyledLogin = styled.div`
  .wrapper {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 2rem;
  }
  .alert {
    background: #fee;
    color: #c00;
    padding: 0.5rem;
    margin-bottom: 1rem;
    border-radius: 4px;
  }
  .form {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    width: 100%;
    max-width: 320px;
  }
  .form p {
    margin: 0;
  }
  .form input {
    padding: 0.5rem;
    border: 1px solid #ccc;
    border-radius: 4px;
    font-size: 14px;
  }
  .actions {
    margin-top: 0.5rem;
  }
  .actions button {
    padding: 0.5rem 1rem;
    background: #007bff;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
  }
  .actions button:hover {
    background: #0056b3;
  }
`;

module.exports = StyledLogin;
