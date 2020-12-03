import logo from './logo.svg';
import './App.css';
import { useState, useEffect } from 'react';

function App() {
  const [word, setWord] = useState('Bye World');
  const getRoot = async () => {
    let result = await fetch('/api/');
    let data = await result.text();
    setWord(data);
  };
  useEffect(()=>{
    getRoot();
  },[])

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.js</code> and save to reload or just don't bother.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          {word}
        </a>
      </header>
    </div>
  );
}

export default App;
