import { useState } from "react";

function App() {
  const [goUsers, setGoUsers] = useState([]);
  const [pythonUsers, setPythonUsers] = useState([]);
  const [goName, setGoName] = useState("");
  const [pythonName, setPythonName] = useState("");

  const backendGo = process.env.REACT_APP_BACKEND_GO;
  const backendPython = process.env.REACT_APP_BACKEND_PYTHON;

  async function createGoUser() {
    await fetch(backendGo + "/users", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name: goName })
    });
    setGoName("");
    fetchGoUsers();
  }

  async function fetchGoUsers() {
    const res = await fetch(backendGo + "/users");
    setGoUsers(await res.json());
  }

  async function createPythonUser() {
    await fetch(backendPython + "/users", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name: pythonName })
    });
    setPythonName("");
    fetchPythonUsers();
  }

  async function fetchPythonUsers() {
    const res = await fetch(backendPython + "/users");
    setPythonUsers(await res.json());
  }

  return (
    <div style={{ padding: 20 }}>
      <h1>Full Stack App</h1>
      <div style={{ marginBottom: 40 }}>
        <h2>Go Backend Users</h2>
        <input value={goName} onChange={e => setGoName(e.target.value)} placeholder="Name" />
        <button onClick={createGoUser}>Create User</button>
        <button onClick={fetchGoUsers}>Fetch Users</button>
        <ul>{goUsers.map(u => <li key={u.id}>{u.name}</li>)}</ul>
      </div>
      <div>
        <h2>Python Backend Users</h2>
        <input value={pythonName} onChange={e => setPythonName(e.target.value)} placeholder="Name" />
        <button onClick={createPythonUser}>Create User</button>
        <button onClick={fetchPythonUsers}>Fetch Users</button>
        <ul>{pythonUsers.map(u => <li key={u.id}>{u.name}</li>)}</ul>
      </div>
    </div>
  );
}

export default App;
