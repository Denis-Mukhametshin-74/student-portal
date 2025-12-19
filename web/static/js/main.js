import { loginWithGoogle, login, register, logout } from "./auth.js";

document.addEventListener('DOMContentLoaded', async function() {
    const loginBtn = document.querySelector('.login-button');
    const registerBtn = document.querySelector('.register-button');
    const logoutBtn = document.querySelector('.logout-button');

    const googleLoginBtn = document.getElementById('googleLogin');
    if (googleLoginBtn) googleLoginBtn.addEventListener('click', loginWithGoogle);

    if (loginBtn) loginBtn.addEventListener('click', login);
    if (registerBtn) registerBtn.addEventListener('click', register);
    if (logoutBtn) logoutBtn.addEventListener('click', logout);

    await checkAuth();
    await checkOAuthTokenInHash();

    if (window.location.pathname.endsWith('index.html') || window.location.pathname.endsWith('/')) {
        await loadSchedule();
    }
    if (window.location.pathname.endsWith('profile.html')) {
        await loadProfile();
    }
});

async function checkAuth() {
    const token = localStorage.getItem('token');
    const currentPage = window.location.pathname.split('/').pop();
    
    if ((currentPage === 'login.html' || currentPage === 'register.html') && token) {
        window.location.href = 'index.html';
        return;
    }
    
    if (currentPage !== 'login.html' && currentPage !== 'register.html' && !token) {
        window.location.href = 'login.html';
        return;
    }
}

async function checkOAuthTokenInHash() {
    const hash = window.location.hash.substring(1);
    if (!hash) return;
    
    const params = new URLSearchParams(hash);
    const encodedToken = params.get('token');
    
    if (encodedToken) {
        try {
            const tokenData = JSON.parse(atob(encodedToken));
            
            localStorage.setItem('token', tokenData.token);
            localStorage.setItem('student', JSON.stringify(tokenData.student));
            
            window.location.hash = '';
            window.location.href = 'index.html';
            
        } catch (error) {
            console.error('Error parsing OAuth token:', error);
            alert('Ошибка обработки токена авторизации');
        }
    }
    
    const error = params.get('error');
    if (error) {
        alert('Ошибка авторизации: ' + decodeURIComponent(error));
        window.location.hash = '';
    }
}

async function loadSchedule() {
    try {
        const token = localStorage.getItem('token');
        const response = await fetch('/api/schedule', {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (response.ok) {
            const scheduleData = await response.json();
            renderSchedule(scheduleData);
        } else {
            console.log('Ошибка получения данных о расписании');
            document.getElementById('pageTitle').innerHTML =
            '<p class="error">Ошибка загрузки расписания</p>';
            document.getElementById('scheduleContainer').innerHTML =
            '<p class="error">Ошибка загрузки расписания</p>';
        }
    } catch (error) {
        console.error('Failed to load schedule:', error);
        document.getElementById('pageTitle').innerHTML =
        '<p class="error">Ошибка соединения с сервером</p>';
        document.getElementById('scheduleContainer').innerHTML =
        '<p class="error">Ошибка соединения с сервером</p>';
    }
}

async function loadProfile() {
    try {
        const token = localStorage.getItem('token');
        const response = await fetch('/api/profile', {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (response.ok) {
            const profile = await response.json();
            renderProfile(profile);
        } else {
            console.log('Ошибка получения данных о профиле');
            document.getElementById('profileInfo').innerHTML =
            '<p class="error">Ошибка загрузки профиля</p>';
            document.getElementById('subjectsContainer').innerHTML =
            '<p class="error">Ошибка загрузки предметов</p>';
        }
    } catch (error) {
        console.error('Failed to load profile:', error);
        document.getElementById('profileInfo').innerHTML =
        '<p class="error">Ошибка соединения с сервером</p>';
        document.getElementById('subjectsContainer').innerHTML =
        '<p class="error">Ошибка соединения с сервером</p>';
    }
}

function renderSchedule(scheduleData) {
    const title = document.getElementById('pageTitle');
    const container = document.getElementById('scheduleContainer');
    
    if ((title && scheduleData.group == "Не указана") || (title && !scheduleData.group)) {
        title.innerHTML = '<p>Нет данных о группе</p>';
    } else if (title && scheduleData.group) {
        title.textContent = `Расписание группы ${scheduleData.group}`;
    }
    
    if (container && scheduleData.days) {
        container.innerHTML = '';
        
        scheduleData.days.forEach(day => {
            const dayElement = document.createElement('div');
            dayElement.className = 'day';
            
            const dayTitle = document.createElement('h3');
            dayTitle.textContent = day.day;
            dayElement.appendChild(dayTitle);
            
            day.lessons.forEach(lesson => {
                const lessonElement = document.createElement('div');
                lessonElement.className = 'lesson';
                lessonElement.innerHTML = `
                    <span class="time">${lesson.time}</span>
                    <span class="subject">${lesson.subject}</span>
                    <span class="room">${lesson.room}</span>
                    <span class="teacher">${lesson.teacher}</span>
                `;
                dayElement.appendChild(lessonElement);
            });
            
            container.appendChild(dayElement);
        });
    } else if (container) {
        container.innerHTML = '<p>Нет данных о расписании</p>';
    }
}

function renderProfile(profile) {
    const profileContainer = document.getElementById('profileInfo');
    if (profileContainer && profile.student) {
        const student = profile.student;
        profileContainer.innerHTML = `
            <p><strong>ФИО:</strong> ${student.name}</p>
            <p><strong>Группа:</strong> ${student.group}</p>
            <p><strong>Email:</strong> ${student.email}</p>
        `;
    } else if (profileContainer) {
        profileContainer.innerHTML = '<p>Нет данных о студенте</p>';
    }
    
    const subjectsContainer = document.getElementById('subjectsContainer');
    if (subjectsContainer && profile.subjects && profile.subjects.length > 0) {
        subjectsContainer.innerHTML = profile.subjects.map(subject => `
            <div class="subject-card">
                <h3>${subject.name}</h3>
                <p><strong>Преподаватель:</strong> ${subject.teacher}</p>
            </div>
        `).join('');
    } else if (subjectsContainer) {
        subjectsContainer.innerHTML = '<p>Нет данных о предметах</p>';
    }
}