const LOGIN_URL = '/user';
const WS_URL = 'ws://' + document.location.host + '/match';
let user = null;
//TODO https://developer.mozilla.org/en-US/docs/Web/API/window/requestAnimationFrame
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
    const canvas = document.getElementById('game');
    if (!canvas.getContext) {
        console.error('Canvas not supported by browser')
        return;
    }
    const ctx = canvas.getContext('2d');
    const game = new GameCanvas(canvas, ctx);
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
        game.clear();
        game.updateState(message);
        game.draw();
    });
    document.addEventListener('keydown', (event) => {
        let input = 0;
        switch (event.key) {
            case 'ArrowUp': input = -1;
                break;
            case 'ArrowDown': input = 1;
                break;
        }
        if (input != 0) {
            game.clear();
            game.moveLeftPaddle(input)
            game.draw();
            const input_message = { input: game.state.LeftPaddle.Pos.Y };
            console.log(input_message);
            websocket.send(JSON.stringify(input_message));
        }
    });
}
class GameCanvas {
    state;
    canvas;
    ctx;
    constructor(canvas, ctx) {
        this.canvas = canvas;
        this.ctx = ctx;
    }
    moveLeftPaddle(y) {
        if (this.state.LeftPaddle.Pos !== undefined) {
            this.state.LeftPaddle.Pos.Y += y;
        }
    }
    updateState(state) {
        this.state = state;
    }
    draw() {
        const scale = this.canvas.height / 100;
        this.ctx.beginPath();
        this.ctx.setLineDash([4 * scale, 4 * scale]);
        this.ctx.moveTo(this.canvas.width / 2, 0);
        this.ctx.lineTo(this.canvas.width / 2, this.canvas.height);
        this.ctx.stroke();

        this.drawScore(this.state.LeftScore, this.state.RightScore);

        this.drawRectangle(this.state.LeftPaddle);
        this.drawRectangle(this.state.RightPaddle);
        this.drawRectangle(this.state.Ball);
    }

    clear() {
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
    }

    drawRectangle(rect) {
        const scaleX = this.canvas.width / 200;
        const scaleY = this.canvas.height / 100
        const x = rect.Pos.X * scaleX;
        const y = rect.Pos.Y * scaleY;
        const width = rect.Width * scaleX;
        const height = rect.Height * scaleY;
        this.ctx.fillRect(x, y, width, height);
    }

    drawScore(left, right) {
        const scaleX = this.canvas.width / 200;
        const scaleY = this.canvas.height / 100
        const size = scaleY * 14;
        this.ctx.font = size + 'px sans-serif';
        const leftX = 46 * scaleX;
        const y = 14 * scaleY;
        const rightX = 146 * scaleX;
        this.ctx.fillText(left, leftX, y);
        this.ctx.fillText(right, rightX, y);
    }
}

// window.onload = function () {
//     const canvas = document.getElementById('game-field');
//     if (!canvas.getContext) {
//         console.error('Canvas not supported by browser')
//         return;
//     }
//     const ctx = canvas.getContext('2d');
//     const game = new Game(canvas, ctx);

//     game.updateState(state);
//     game.draw();
//     document.addEventListener('keydown', (event) => {
//         let input = 0;
//         switch (event.key) {
//             case 'ArrowUp': input = -1;
//                 break;
//             case 'ArrowDown': input = 1;
//                 break;
//         }
//         if (input != 0) {
//             game.clear();
//             state.LeftPaddle.Pos.Y += input;
//             game.updateState(state);
//             game.draw();
//         }
//     });
// }


// const state = {
//     "LeftPaddle": {
//         "Pos": {
//             "X": 5,
//             "Y": 43
//         },
//         "Width": 4,
//         "Height": 14
//     },
//     "LeftScore": 0,
//     "RightPaddle": {
//         "Pos": {
//             "X": 191,
//             "Y": 43
//         },
//         "Width": 4,
//         "Height": 14
//     },
//     "RightScore": 0,
//     "Ball": {
//         "Pos": {
//             "X": 98,
//             "Y": 48
//         },
//         "Width": 4,
//         "Height": 4
//     },
//     "BallDir": {
//         "X": 1,
//         "Y": -1
//     },
//     "Time": 7095537900
// };