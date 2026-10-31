async function testLargeContextCombo() {
  const url = "http://127.0.0.1:20128/v1/chat/completions";
  const apiKey = "sk-4bc1e878b525f20a-qo8rt7-7ca55e91";

  // Generate ~360k characters of context history (>100k tokens)
  const bigHistory = "This is a detailed historical trace of long coding session. ".repeat(6000);
  console.log(`Generated context length: ${bigHistory.length} chars (~${Math.ceil(bigHistory.length / 3.5)} tokens)`);

  const payload = {
    model: "code-pro",
    messages: [
      { role: "user", content: `Context history:\n${bigHistory}\n\nPlease respond with just: CONTEXT_RECEIVED_OK` }
    ],
    stream: false,
    max_tokens: 20
  };

  console.log("Sending large-context request to combo:", payload.model);
  const startTime = Date.now();
  const res = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${apiKey}`,
      "X-9Router-Token-Saver": "off"
    },
    body: JSON.stringify(payload)
  });

  console.log(`HTTP STATUS: ${res.status} ${res.statusText} (${Date.now() - startTime}ms)`);
  const json = await res.json();
  console.log("RESPONSE:", JSON.stringify(json.choices?.[0]?.message?.content || json, null, 2));
}

testLargeContextCombo();
