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
        let encodedMessage = ev.data.replaceAll('"', '').replaceAll('\n', '')
        const message = JSON.parse(atob(encodedMessage))
        console.log(message);
        // if (isGame(message)) {
        // document.getElementById('left-paddle').style.top = message.Left.Y;
        // document.getElementById('left-paddle').style.left = message.Left.X;
        // document.getElementById('right-paddle').style.top = message.Right.Y;
        // document.getElementById('right-paddle').style.left = message.Right.X;

        // document.getElementById('ball').style.top = message.Ball.Y;
        // document.getElementById('ball').style.left = message.Ball.X;
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
            const input_message = { input: input };
            console.log(input_message);
            websocket.send(JSON.stringify(input_message));
        }
    });
}

function isGame(game) {
    return game.left !== undefined && game.right !== undefined && ball != undefined;
}

function drawGame() {
    const canvas = document.getElementById('game-field');
    if (canvas.getContext) {
        const ctx = canvas.getContext('2d');
        return new Game(canvas, ctx);
    } else {
        console.error('Canvas not supported by browser')
    }
}
window.onload = function () {
    const game = drawGame();
    const scale = game.canvas.width / 100;

    state.Ball.Width *= scale;
    state.Ball.Height *= scale;
    state.Ball.Pos.X *= scale;
    state.Ball.Pos.Y *= scale;

    state.LeftPaddle.Width *= scale;
    state.LeftPaddle.Height *= scale;
    state.LeftPaddle.Pos.X *= scale;
    state.LeftPaddle.Pos.Y *= scale;

    state.RightPaddle.Width *= scale;
    state.RightPaddle.Height *= scale;
    state.RightPaddle.Pos.X *= scale;
    state.RightPaddle.Pos.Y *= scale;

    game.updateState(state);
    game.draw();
}

class Game {
    state;
    canvas;
    ctx;
    constructor(canvas, ctx) {
        this.canvas = canvas;
        this.ctx = ctx;
    }
    updateState(state) {
        this.state = state;
    }
    draw() {
        const leftPaddle = state.LeftPaddle;
        this.ctx.fillRect(leftPaddle.Pos.X, leftPaddle.Pos.Y, leftPaddle.Width, leftPaddle.Height);

        const rightPaddle = state.RightPaddle;
        this.ctx.fillRect(rightPaddle.Pos.X, rightPaddle.Pos.Y, rightPaddle.Width, rightPaddle.Height);

        const ball = state.Ball;
        this.ctx.fillRect(ball.Pos.X, ball.Pos.Y, ball.Width, ball.Height);

    }
    clear() {
        const leftPaddle = state.LeftPaddle;
        this.ctx.clearRect(leftPaddle.Pos.X, leftPaddle.Pos.Y, leftPaddle.Width, leftPaddle.Height);

        const rightPaddle = state.RightPaddle;
        this.ctx.clearRect(rightPaddle.Pos.X, rightPaddle.Pos.Y, rightPaddle.Width, rightPaddle.Height);

        const ball = state.Ball;
        this.ctx.clearRect(ball.Pos.X, ball.Pos.Y, ball.Width, ball.Height);
    }
}
const state = {
    "LeftPaddle": {
        "Pos": {
            "X": 5,
            "Y": 50
        },
        "Width": 1,
        "Height": 8
    },
    "LeftScore": 0,
    "RightPaddle": {
        "Pos": {
            "X": 95,
            "Y": 50
        },
        "Width": 1,
        "Height": 8
    },
    "RightScore": 0,
    "Ball": {
        "Pos": {
            "X": 50,
            "Y": 50
        },
        "Width": 1,
        "Height": 1
    },
    "BallDir": {
        "X": 1,
        "Y": -1
    },
    "Time": 7095537900
};