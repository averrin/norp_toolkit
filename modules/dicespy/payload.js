alert("Script injected");
function getPlayers() {
    var players = {};
    // $('.by').each((i,r) => players[$(r).parent().attr('data-playerid')] = $(r).text().slice(0, -1));
    var elements = document.querySelectorAll('.by');
    Array.prototype.forEach.call(elements, function(el, i){
        players[el.parentElement.getAttribute('data-playerid')] = el.innerText.slice(0, -1);
    });
    alert(JSON.stringify(players));

    var request = new XMLHttpRequest();
    request.open('POST', 'http://127.0.0.1:1323/players', true);
    request.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
    request.send(JSON.stringify(players));
}

function injectReporter() {
    let oldParse = JSON.parse;
    JSON.parse = function newParse(d, r) {
        let data = oldParse(d, r);
        setTimeout(function() {
            if (data.d && data.d.b && data.d.b.d && data.d.b.d.type == 'rollresult') {
                getPlayers();

                var request = new XMLHttpRequest();
                request.open('POST', 'http://127.0.0.1:1323/roll', true);
                request.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
                alert(JSON.stringify(data.d.b));
                request.onreadystatechange = function() { // (3)
                    if (request.readyState != 4) return;

                    if (request.status != 200) {
                        alert(request.status + ': ' + request.statusText);
                    } else {
                        alert(request.responseText);
                    }
                };
                request.send(JSON.stringify(data.d.b));
            }}, 0);
        return data;
    };
}

getPlayers();
injectReporter();