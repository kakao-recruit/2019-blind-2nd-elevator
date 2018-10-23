package main


const (
	htmlBody = `
	<!doctype html>
	<link rel=stylesheet href=/viewer/main.css>
	<p>Timestamp: <input id=ts type=number min=0 value=0 autofocus> / <span id=maxts>0</span>
	<p> AveWait: %f, AveTravel: %f, AveTotal: %f, LastTs: %d, Status: %s
	<main>
	<div class=floors>
	<h4>F</h4>
	<ul><li>25<li>24<li>23<li>22<li>21<li>20<li>19<li>18<li>17<li>16<li>15<li>14<li>13<li>12<li>11<li>10<li>9<li>8<li>7<li>6<li>5<li>4<li>3<li>2<li>1</ul>
	</div>
	<div class=ev id=ev0>
	<h4>0</h4>
	<ul><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li></ul>
	<p>?
	</div>
	<div class=ev id=ev1>
	<h4>1</h4>
	<ul><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li></ul>
	<p>?
	</div>
	<div class=ev id=ev2>
	<h4>2</h4>
	<ul><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li></ul>
	<p>?
	</div>
	<div class=ev id=ev3>
	<h4>3</h4>
	<ul><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li></ul>
	<p>?
	</div>
	<div class=calls>
	<h4>Calls</h4>
	<ul border=1 id=calls><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li><li></ul>
	</div>
	</main>
	<pre hidden>%s</pre>
	<script src="https://code.jquery.com/jquery-3.3.1.min.js"></script>
	<script src="https://unpkg.com/axios/dist/axios.min.js"></script>
	<script src="/viewer/main.js"></script>`

)
