const API_BASE = "/v1"; // Use relative URL for load balancer

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
  
  // Use relative URL for WebSocket connection
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const ws = new WebSocket(`${protocol}//${window.location.host}/v1/job/${id}/output`);
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
