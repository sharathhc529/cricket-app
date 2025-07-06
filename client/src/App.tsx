import { useEffect, useState } from 'react'
import './App.css'


function App() {

const [data, setData] = useState<any[][]>([]);
// const API_BASE = "https://cricket-app-si7p.onrender.com";
const API_BASE = "http://localhost:8080";
  // Unified handler for both single and multi input
  // Helper to fetch and update sheet data
  const fetchSheetData = async () => {
    try {
      const res = await fetch(`${API_BASE}/entries`);
      const json = await res.json();
      setData(json.data);
    } catch (err) {
      console.error("Fetch error:", err);
    }
  };

  // Unified handler for both single and multi input
  const handleUnifiedSubmit = async (payload: { value?: string; player?: string; score?: string | number }) => {
    try {
      // Convert score to number if present and not already a number
      let fixedPayload = { ...payload };
      if (fixedPayload.score !== undefined && typeof fixedPayload.score === 'string') {
        const numScore = Number(fixedPayload.score);
        if (!isNaN(numScore)) {
          fixedPayload.score = numScore;
        } else {
          delete fixedPayload.score;
        }
      }

      const res = await fetch(`${API_BASE}/submit`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(fixedPayload),
      });
      const result = await res.json();
      if (result.status === 'success') {
        await fetchSheetData();
        alert('Submitted successfully!');
      } else {
        alert('Submission failed. Please try again.');
      }
    } catch (err) {
      alert('Submission failed. Please try again.');
      console.error('Error submitting:', err);
    }
  };
  // For single input field
  const [message, setMessage] = useState('');

  // Handler for single input field
  const handleSingleSubmit = async () => {
    if (!message.trim()) return alert("Please enter something before submitting.");
    await handleUnifiedSubmit({ value: message });
    setMessage('');
  };

  useEffect(() => {
    fetchSheetData();
  }, []);

  // React state for form fields
  const [player, setPlayer] = useState('');
  const [score, setScore] = useState('');

  // Handler for the player/score form
  const handleScoreSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!player.trim() || !score.trim()) {
      alert('Please enter both player and score.');
      return;
    }
    await handleUnifiedSubmit({ player, score });
    setPlayer('');
    setScore('');
  };

  return (
    <>
      <div>
        <div>
          <h1>Google Sheet Data</h1>
          <table border={1}>
            <tbody>
              {Array.isArray(data) && data.length > 0 ? (
                data.map((row, i) => (
                  <tr key={i}>
                    {Array.isArray(row) ? row.map((cell, j) => (
                      <td key={j}>{cell}</td>
                    )) : null}
                  </tr>
                ))
              ) : (
                <tr><td colSpan={2}>No data</td></tr>
              )}
            </tbody>
          </table>

          {/* Single input field logic */}
          <div style={{ marginBottom: 16 }}>
            <input
              value={message}
              onChange={e => setMessage(e.target.value)}
              placeholder="Type something..."
            />
            <button onClick={handleSingleSubmit} style={{ marginLeft: 8 }}>Submit</button>
          </div>

          {/* Multi-input form logic */}
          <form onSubmit={handleScoreSubmit}>
            <input
              name="player"
              placeholder="Player Name"
              value={player}
              onChange={e => setPlayer(e.target.value)}
              required
            />
            <input
              name="score"
              type="number"
              placeholder="Score"
              value={score}
              onChange={e => setScore(e.target.value)}
              required
            />
            <button type="submit">Submit</button>
          </form>

        </div>

      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </>
  )
}

export default App
