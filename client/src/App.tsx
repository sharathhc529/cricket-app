import { useEffect, useState } from 'react'
import './App.css'

function App() {
  const [data, setData] = useState<[]>([]);
  const API_BASE = "https://cricket-app-si7p.onrender.com";
  // const API_BASE = "http://localhost:8080";

  useEffect(() => {
    // fetch("http://localhost:8080/entries")
    fetch(`${API_BASE}/entries`)
      .then(res => res.json())
      .then(json => {
        console.log("Sheet data:", json.data);
        setData(json.data);
      })
      .catch(err => console.error("Fetch error:", err));
  }, []);

  const [message, setMessage] = useState('');

const handleSubmit = async () => {
  if (!message.trim()) return alert("Please enter something before submitting.");

  try {
    const res = await fetch(`${API_BASE}/submit`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ value: message }),
    });

    const result = await res.json();
    console.log("Server response:", result);
    setMessage('');
  } catch (err) {
    console.error("Error submitting:", err);
  }
};



  return (
    <>
      <div>
        <div>
          <h1>Google Sheet Data</h1>
          <table border={1}>
            <tbody>
              {data.map((row, i) => (
                <tr key={i}>
                  {(row as string[]).map((cell, j) => (
                    <td key={j}>{cell}</td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>

          <input
            value={message}
            onChange={e => setMessage(e.target.value)}
            placeholder="Type something..."
          />
          <button onClick={handleSubmit}>Submit</button>

        </div>

      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </>
  )
}

export default App
