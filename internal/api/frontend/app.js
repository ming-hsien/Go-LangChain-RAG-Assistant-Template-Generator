async function sendMessage() {
    const input = document.getElementById('user-input');
    const container = document.getElementById('chat-container');
    const question = input.value.trim();
    if (!question) return;

    // Add user message
    appendMessage('USER', question);
    input.value = '';

    try {
        const response = await fetch('/v1/query', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ question })
        });
        const data = await response.json();
        if (data.error) {
            appendMessage('AI', 'Error: ' + data.error);
        } else {
            appendMessage('AI', data.answer);
        }
    } catch (error) {
        appendMessage('AI', 'Error connecting to server: ' + error.message);
    }
}

function appendMessage(sender, text) {
    const container = document.getElementById('chat-container');
    const msgDiv = document.createElement('div');
    msgDiv.className = 'flex items-start';
    const isAi = sender === 'AI';

    msgDiv.innerHTML = `
        <div class="w-10 h-10 ${isAi ? 'bg-indigo-600 text-white shadow-lg' : 'bg-white border border-slate-200 text-slate-400'} rounded-2xl flex items-center justify-center flex-shrink-0 font-bold text-[10px] z-10">${sender}</div>
        <div class="ml-4 max-w-[85%]">
            <div class="glass p-5 rounded-3xl rounded-tl-none text-slate-800 font-medium leading-relaxed shadow-sm text-sm">
                ${text.replace(/\n/g, '<br>')}
            </div>
        </div>
    `;
    container.appendChild(msgDiv);
    container.scrollTop = container.scrollHeight;
}

async function loadTools() {
    try {
        const response = await fetch('/v1/tools');
        const data = await response.json();
        const list = document.getElementById('tools-list');
        list.innerHTML = data.tools.map(t => `
            <div class="p-3 bg-white/60 rounded-xl border border-slate-100 text-[11px] shadow-sm">
                <div class="font-bold text-indigo-600">${t.name}</div>
                <div class="text-slate-500 mt-1 leading-snug">${t.description}</div>
            </div>
        `).join('');
    } catch (e) {
        console.error("Failed to load tools", e);
    }
}

async function reindex() {
    try {
        const response = await fetch('/v1/index', { method: 'POST' });
        const data = await response.json();
        alert(data.message);
    } catch (e) {
        alert("Reindex request failed: " + e.message);
    }
}

async function loadConfig() {
    try {
        const response = await fetch('/v1/config');
        const data = await response.json();
        if (data.app_name) {
            document.getElementById('head-title').innerText = data.app_name;
            document.getElementById('app-header').innerText = data.app_name;
            // Update welcome message text if it exists
            const welcome = document.getElementById('welcome-message');
            if (welcome) {
                welcome.innerHTML = welcome.innerHTML.replace('AI 助手', `${data.app_name} 助手`);
            }
        }
    } catch (e) {
        console.error("Failed to load config", e);
    }
}

// Initial Load
loadConfig();
loadTools();
