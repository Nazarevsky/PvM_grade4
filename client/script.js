// Адреса серверу
const apiUrl = 'http://localhost:8000';

function sendData() {
    // Поле зі введеним повідомленням
    let message = document.getElementById("message").value
    // Поле з відповіддю
    let result = document.getElementById("result")

    // Структура запиту
    const msg = {
        message: message,
    };

    // POST запит на ендпоїнт
    fetch(`${apiUrl}/replace_round_brackets`, {
        method: "POST",
        body: JSON.stringify(msg), // Серіалізація запиту у json формат
    })
        // Проміс для отримання відповіді
        .then(response => {
            // Обробка неуспішного запиту
            if (!response.ok){
                alert('Помилка! Статус ' + response.statusText)
                return
            }

            // Обробка успішного запиту (статус - 200)
            response.json().then(data => {
                result.innerHTML = 'Відповідь: ' + data.message_processed
            })
        })
}
