let currentUser = null;
let isRegisterMode = false;

function toggleAuthMode() {
    isRegisterMode = !isRegisterMode;
    document.getElementById('auth-title').innerText = isRegisterMode ? 'Регистрация' : 'Вход в Helpdesk';
    document.getElementById('auth-btn').innerText = isRegisterMode ? 'Зарегистрироваться' : 'Войти';
    document.getElementById('toggle-auth').innerText = isRegisterMode ? 'Уже есть аккаунт? Войти' : 'Зарегистрироваться';
    document.getElementById('role-group').classList.toggle('hidden', !isRegisterMode);
}

async function submitAuth() {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    const role = document.getElementById('role').value;

    const endpoint = isRegisterMode ? '/api/register' : '/api/login';
    const payload = isRegisterMode ? { username, password, role } : { username, password };

    const res = await fetch(endpoint, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
    });

    const data = await res.json();
    if (res.ok && !isRegisterMode) {
        currentUser = data.user;
        showApp();
    } else if (res.ok && isRegisterMode) {
        alert('Успешная регистрация! Теперь войдите.');
        toggleAuthMode();
    } else {
        alert(data.message || 'Ошибка');
    }
}

function showApp() {
    document.getElementById('auth-box').classList.add('hidden');
    document.getElementById('app-box').classList.remove('hidden');
    document.getElementById('user-display').innerText = `${currentUser.username} [${currentUser.role}]`;

    const createSection = document.getElementById('create-ticket-section');
    if (currentUser.role === 'admin') {
        createSection.classList.add('hidden');
    } else {
        createSection.classList.remove('hidden');
    }

    loadTickets();
}

async function loadTickets() {
    const res = await fetch(`/api/tickets?role=${currentUser.role}&user_id=${currentUser.id}`);
    const tickets = await res.json();
    const list = document.getElementById('tickets-list');
    list.innerHTML = '';

    if (!tickets || tickets.length === 0) {
        list.innerHTML = '<p>Заявок пока нет</p>';
        return;
    }

    tickets.forEach(t => {
        const div = document.createElement('div');
        div.className = 'ticket';
        
        let statusText = t.status === 'done' ? 'Завершена' : (t.status === 'in_progress' ? 'В работе' : 'Новая');
        let statusClass = t.status === 'done' ? 'status-done' : 'status-new';

        let adminControls = '';
        if (currentUser.role === 'admin') {
            adminControls = `
                <div style="margin-top: 10px; padding-top: 10px; border-top: 1px dashed #ccc;">
                    <label><b>Ответ администратора:</b></label>
                    <textarea id="resp-${t.id}" placeholder="Напишите ответ..." style="margin-top: 5px; margin-bottom: 5px;">${t.admin_response || ''}</textarea>
                    
                    <div style="display: flex; gap: 10px;">
                        <select id="status-${t.id}">
                            <option value="new" ${t.status === 'new' ? 'selected' : ''}>Новая</option>
                            <option value="in_progress" ${t.status === 'in_progress' ? 'selected' : ''}>В работе</option>
                            <option value="done" ${t.status === 'done' ? 'selected' : ''}>Завершена</option>
                        </select>
                        <button onclick="updateTicket(${t.id})" style="width: auto;">Сохранить ответ</button>
                    </div>
                </div>
            `;
        } else if (t.admin_response) {
            adminControls = `
                <div style="margin-top: 10px; padding: 10px; background: #e9ecef; border-radius: 4px;">
                    <b>Ответ администратора:</b>
                    <p style="margin: 5px 0 0 0;">${t.admin_response}</p>
                    </div>
            `;
        }

        div.innerHTML = `
            <div class="ticket-header">
                <span>#${t.id} ${t.title} (${t.username})</span>
                <span class="${statusClass}">[${statusText}]</span>
            </div>
            <p>${t.description}</p>
            ${adminControls}
        `;
        list.appendChild(div);
    });
}

async function updateTicket(ticketId) {
    const status = document.getElementById(`status-${ticketId}`).value;
    const adminResponse = document.getElementById(`resp-${ticketId}`).value;

    const res = await fetch(`/api/tickets/status?id=${ticketId}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ status: status, admin_response: adminResponse })
    });

    if (res.ok) {
        alert('Заявка обновлена');
        loadTickets();
    } else {
        alert('Ошибка обновления');
    }
}

async function createTicket() {
    const title = document.getElementById('ticket-title').value;
    const description = document.getElementById('ticket-desc').value;

    const res = await fetch('/api/tickets', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title, description, user_id: currentUser.id })
    });

    if (res.ok) {
        document.getElementById('ticket-title').value = '';
        document.getElementById('ticket-desc').value = '';
        loadTickets();
    }
}

function logout() {
    currentUser = null;
    document.getElementById('app-box').classList.add('hidden');
    document.getElementById('auth-box').classList.remove('hidden');
}