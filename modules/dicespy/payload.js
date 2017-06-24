function getPlayers() {
    var players = {};
    $('.by').each((i,r) => players[$(r).parent().attr('data-playerid')] = $(r).text().slice(0, -1));
    console.log(players);
    $.post('http://127.0.0.1:1323/players', JSON.stringify(players));
}

function injectReporter() {
    let oldParse = JSON.parse;
    JSON.parse = function newParse(d, r) {
        let data = oldParse(d, r);
        setTimeout(function() {
            if (data.d && data.d.b && data.d.b.d && data.d.b.d.type == 'rollresult') {
                getPlayers();
                $.post('http://127.0.0.1:1323/roll', JSON.stringify(data.d.b));
            }}, 0);
        return data;
    };
}

getPlayers();
injectReporter();
console.info('DiceSpy injected.');
let c = $('#textchat .content');
let m = $('<div class="message private whisper"></div>');
m.append($('<div class="spacer"></div>'));
m.append($('<p>DiceSpy injected.</p>'));
c.append(m);
