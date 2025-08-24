const API_BASE = "http://localhost:3000/v1"; // change if needed

function appendOutput(msg) {
  const out = document.getElementById("output");
  out.textContent += msg + "\n";
  out.scrollTop = out.scrollHeight;
}

document.getElementById("startJob").onclick = async () => {
  const res = await fetch(`${API_BASE}/job`, { method: "POST" });
  const data = await res.json();
  appendOutput("Started job: " + data.id);
  document.getElementById("jobIdInput").value = data.id;
};

document.getElementById("getJob").onclick = async () => {
  const id = document.getElementById("jobIdInput").value;
  if (!id) return alert("Enter Job ID");
  const res = await fetch(`${API_BASE}/job/${id}`);
  const data = await res.json();
  appendOutput("Metadata: " + JSON.stringify(data));
};

document.getElementById("streamOutput").onclick = () => {
  const id = document.getElementById("jobIdInput").value;
  if (!id) return alert("Enter Job ID");
  const ws = new WebSocket(`ws://localhost:3000/v1/job/${id}/output`);
  ws.onmessage = (event) => appendOutput(event.data);
  ws.onopen = () => appendOutput("📡 Connected to output stream...");
  ws.onclose = () => appendOutput("❌ Stream closed");
};

document.getElementById("stopJob").onclick = async () => {
  const id = document.getElementById("jobIdInput").value;
  if (!id) return alert("Enter Job ID");
  await fetch(`${API_BASE}/job/${id}`, { method: "DELETE" });
  appendOutput("Stopped job: " + id);
};
