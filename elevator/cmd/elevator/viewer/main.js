'use strict';

const timeline = JSON.parse($('pre').text());

timeline.forEach(({elevators}, i, a) => {
   if (i > 0) {
       elevators.forEach(({passengers}, j) => {
           passengers.sort((a, b) => a.to - b.to);
           const b = a[i - 1].elevators[j].passengers;
           for (const p of passengers) {
               if (!b.some(q => q.id === p.id)) {
                   p.entered = true;
               }
           }
       });
   }
   if (i + 1 < a.length) {
       elevators.forEach(({passengers}, j) => {
           passengers.sort((a, b) => a.to - b.to);
           const b = a[i + 1].elevators[j].passengers;
           for (const p of passengers) {
               if (!b.some(q => q.id === p.id)) {
                   p.exited = true;
               }
           }
       });
   }
});

const commandText = {
    S: 'STOP',
    O: 'OPEN',
    E: 'ENTER',
    X: 'EXIT',
    C: 'CLOSE',
    U: 'UP',
    D: 'DOWN',
};

const stateText = {
    S: 'STOPPED',
    O: 'OPENED',
    U: 'UPWARD',
    D: 'DOWNWARD',
};

function show(ts) {
    $('#ts').val(ts);
    $('button').removeClass('active').eq(ts).addClass('active');

    console.log(timeline[ts]);
    if (!timeline[ts]) {
        return;
    }

    timeline[ts].elevators.forEach((ev, i) => {
        const $li = $('#ev' + i).find('p').text((commandText[ev.command] || ev.command) + ' : ' + (stateText[ev.state] || ev.state)).end().find('li');
        $li.removeClass('target state-S state-O state-U state-D command-E command-X').empty()
            .eq(25 - ev.floor).addClass('state-' + ev.state).addClass('command-' + ev.command)
            .append(ev.passengers.map(c => '<small>(' + c.id + ')</small>' + (c.exited ? '<strong>' + c.to + '</strong>' : c.entered ? '<em>' + c.to + '</em>' : c.to)).join(' '));
            console.log(ev.passengers);
        for (const {to} of ev.passengers) {
            $li.eq(25 - to).addClass('target');
        }
        switch (ev.state) {
            case 'U':
                $li.eq(25 - ev.floor + 1).text('↑');
                break;
            case 'D':
                $li.eq(25 - ev.floor - 1).text('↓');
                break;
        }
    });

    const $calls = $('#calls li').empty();
    for (const c of timeline[ts].calls) {
        $calls.eq(25 - c.from).append('<small>(' + c.id + ')</small>' + c.to + ' ');
    }
}

$('#ts').prop('max', timeline.length - 1).change(e => {
    show(e.currentTarget.value);
});

$('#maxts').text(timeline.length - 1);

show(0);
