const users = document.getElementById('users');
const tBody = document.getElementById('rankingBody');

window.addEventListener('load', (event) => {
    users.click();
});

users.addEventListener('click', (event) => {
    const response = JSON.parse(users.textContent);
    console.log(users.textContent)

    if (response.data && response.data != null) {
        if (response.data?.length >= 1) {
            for(let i = 0; i < response.data?.length; i++) {
                let tr = document.createElement('tr');
                let id = document.createElement('td');
                let name = document.createElement('td');
                let score = document.createElement('td');
    
                id.textContent = i + 1;
                name.textContent = response.data[i].name;
                score.textContent = response.data[i].score;
    
                tr.appendChild(id);
                tr.appendChild(name);
                tr.appendChild(score);
                tBody.appendChild(tr);
            }
        }
    }

    users.style.display = 'none';
})