const TEST_DATA = {
  schedule: {
    group: "А-10",
    days: [
      {
        day: "Понедельник",
        lessons: [
          { time: "9:00-10:30", subject: "Программирование", room: "301", teacher: "Иванов И.И." },
          { time: "10:45-12:15", subject: "Математика", room: "215", teacher: "Петрова А.С." }
        ]
      },
      {
        day: "Вторник",
        lessons: [
          { time: "9:00-10:30", subject: "Физика", room: "101", teacher: "Сидоров В.П." }
        ]
      },
      {
        day: "Среда",
        lessons: [
          { time: "9:00-12:15", subject: "Математика", room: "220", teacher: "Петрова А.С." },
          { time: "13:00-14:30", subject: "Английский язык", room: "410", teacher: "Козлова Е.В." }
        ]
      },
      {
        day: "Четверг",
        lessons: [
          { time: "10:45-12:15", subject: "Программирование", room: "301", teacher: "Иванов И.И." }
        ]
      },
      {
        day: "Пятница",
        lessons: [
          { time: "9:00-10:30", subject: "История", room: "205", teacher: "Николаев П.Д." },
          { time: "10:45-12:15", subject: "Физкультура", room: "102", teacher: "Смирнов А.В." }
        ]
      }
    ]
  },
  subjects: [
    { name: "Программирование", teacher: "Иванов И.И." },
    { name: "Математика", teacher: "Петрова А.С." },
    { name: "Физика", teacher: "Сидоров В.П." },
    { name: "Английский язык", teacher: "Козлова Е.В." },
    { name: "История", teacher: "Николаев П.Д." },
    { name: "Физкультура", teacher: "Смирнов А.В." }
  ]
};

document.addEventListener('DOMContentLoaded', function() {
  checkAuth();

  const currentPage = window.location.pathname.split('/').pop();
  
  if (currentPage === 'index.html' || currentPage === '') {
    loadSchedule();
  }
  
  if (currentPage === 'profile.html') {
    loadProfile();
  }
});

function checkAuth() {
  const user = localStorage.getItem('user');
  const currentPage = window.location.pathname.split('/').pop();
  
  if (currentPage === 'login.html' && user) {
    window.location.href = 'index.html';
  }
  
  if (currentPage !== 'login.html' && currentPage !== 'register.html' && !user) {
    window.location.href = 'login.html';
  }
}

function login() {
  const email = document.getElementById('email').value;
  const password = document.getElementById('password').value;
  
  if (email && password) {
    const user = {
      name: "Иванов Иван Иванович",
      email: email,
      group: "A-10"
    };
    
    localStorage.setItem('user', JSON.stringify(user));
    window.location.href = 'index.html';
  } else {
    alert("Введите email и пароль");
  }
}

function logout() {
  localStorage.removeItem('user');
  window.location.href = 'login.html';
}

function loadSchedule() {
  const scheduleData = TEST_DATA.schedule;
  
  const title = document.getElementById('pageTitle');
  if (title) {
    title.textContent = `Расписание группы ${scheduleData.group}`;
  }
  
  const scheduleContainer = document.getElementById('scheduleContainer');
  if (scheduleContainer) {
    scheduleContainer.innerHTML = '';
    
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
      
      scheduleContainer.appendChild(dayElement);
    });
  }
}

function loadProfile() {
  const userStr = localStorage.getItem('user');
  let user = {
    name: "Иванов Иван Иванович",
    email: "test@example.com",
    group: "А-10"
  };
  
  if (userStr) {
    try {
      user = JSON.parse(userStr);
    } catch (e) {
      console.log('Ошибка загрузки пользователя, используем тестовые данные');
    }
  }

  const profileInfo = document.getElementById('profileInfo');
  if (profileInfo) {
    profileInfo.innerHTML = `
      <p><strong>ФИО:</strong> ${user.name}</p>
      <p><strong>Группа:</strong> ${user.group}</p>
      <p><strong>Email:</strong> ${user.email}</p>
    `;
  }
  
  const subjectsContainer = document.getElementById('subjectsContainer');
  if (subjectsContainer) {
    subjectsContainer.innerHTML = '';
    
    TEST_DATA.subjects.forEach(subject => {
      const subjectCard = document.createElement('div');
      subjectCard.className = 'subject-card';
      subjectCard.innerHTML = `
        <h3>${subject.name}</h3>
        <p><strong>Преподаватель:</strong> ${subject.teacher}</p>
      `;
      subjectsContainer.appendChild(subjectCard);
    });
  }
}