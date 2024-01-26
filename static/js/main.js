const LOGIN_URL = '/user';
const WS_URL = 'ws://' + document.location.host + '/match';
let user = null;

window.onload = function () {
    const user = localStorage.getItem('user');
    if (user) {
        document.getElementById('login-name').value = user.name;
        document.getElementById('login').hidden = false;
    }
}


function register() {
    const name = document.getElementById('name').value
    fetch(LOGIN_URL + '?name=' + name, { method: 'POST' }).then(response => {
        if (!response.ok) {
            console.log(response);
        } else {
            response.json().then((id) => {
                console.log(id);
                document.getElementById('register').hidden = true;
                user = { id: id, name: name };
                localStorage.setItem('user', user);
                startMatch(id);
            });
        }
    }).catch(err => {
        console.error(err);
    })
}

function startMatch() {
    if (!user) {
        console.error('User is null');
        return;
    }
    const websocket = new WebSocket(WS_URL + '?id=' + user.id);
    websocket.onopen = (ev => {
        console.log(ev);
        document.getElementById('game').hidden = false;
    }
    );
    websocket.onclose = (ev => console.log(ev))
    websocket.onmessage = (ev => {
        let encodedMessage = ev.data.replaceAll('"','').replaceAll('\n','')
        const message = JSON.parse(atob(encodedMessage))
        // if (isGame(message)) {
            document.getElementById('left-paddle').style.top = message.Left.Y;
            document.getElementById('left-paddle').style.left = message.Left.X;
            document.getElementById('right-paddle').style.top = message.Right.Y;
            document.getElementById('right-paddle').style.left = message.Right.X;

            document.getElementById('ball').style.top = message.Ball.Y;
            document.getElementById('ball').style.left = message.Ball.X;
        // }
    });
    document.addEventListener('keydown', (event) => {
        let input = 0;
        switch (event.key) {
            case 'ArrowUp': input = 1
                break;
            case 'ArrowDown': input = -1
                break;
        }
        if (input != 0) {
            const input_message = { input: input};
            console.log(input_message);
            websocket.send(JSON.stringify(input_message));
        }
    });
}

function isGame(game) {
    return game.left !== undefined && game.right !== undefined && ball != undefined;
}