var socket = new WebSocket("ws://127.0.0.1:1323/ws");
socket.onmessage = function() {
    window.location.reload();
};
