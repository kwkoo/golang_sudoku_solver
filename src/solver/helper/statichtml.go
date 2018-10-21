package helper

import (
	"fmt"
	"io"
)

// StaticHTML returns a hardcoded HTML page for the application.
func StaticHTML(w io.Writer) {
	fmt.Fprint(w, `<!DOCTYPE html>
	<html>
		<head>
			<title>Sudoku Solver</title>
			<style>
				.grid {
					display: grid;
					grid-template-columns: repeat(3, 180px);
					grid-template-rows: repeat(3, 180px);
				}
				.box {
					border: 1px solid rgba(0, 0, 0, 0.8);
					display: grid;
					grid-template-columns: repeat(3, 60px);
					grid-template-rows: repeat(3, 60px);
				}
				.cell {
					border: 1px solid rgba(0, 0, 0, 0.3);
					font-size: 40px;
					text-align: center;
					justify-content: center;
					vertical-align: bottom;
					padding: 5px;
				}
				.static {
					color: black;
				}
				.dynamic {
					color: blue;
				}
				#error {
					visibility: hidden;
				}
				#keypad {
					visibility: hidden;
				}
				.highlighted {
					background-color: lightgray;
				}
			</style>
		</head>
		<body onload="initPage()">
			<h1>Sudoku Solver</h1>
			<div id="main">
				<div style="padding: 10px;">
					<input id="enterButton" type="button" value="Enter Puzzle" onclick="prepForManualEntry()"/>
					&nbsp;
					<input id="solveButton" type="button" value="Solve Puzzle" onclick="solvePuzzle()"/>
					&nbsp;
					<input id="delayRange" type="range" min="0" max="10" value="1" onchange="sendDelay()"/>
				</div>
				<div id="grid" class="grid">
				</div>
				<div id="keypad">
					<input type="button" value="1" onclick="manualSet(1)"/>
					<input type="button" value="2" onclick="manualSet(2)"/>
					<input type="button" value="3" onclick="manualSet(3)"/>
					<input type="button" value="4" onclick="manualSet(4)"/>
					<input type="button" value="5" onclick="manualSet(5)"/>
					<input type="button" value="6" onclick="manualSet(6)"/>
					<input type="button" value="7" onclick="manualSet(7)"/>
					<input type="button" value="8" onclick="manualSet(8)"/>
					<input type="button" value="9" onclick="manualSet(9)"/>
					<input type="button" value="✕" onclick="manualSet(0)"/>
				</div>
			</div>
			<div id="error">
				<h2>Error</h2>
				<pre id="errormessage"></pre>
			</div>
			<script type="text/javascript">
				var globalSocket = null;

				function getDelay() {
					return document.getElementById("delayRange").value;
				}

				function sendDelay() {
					if (globalSocket != null) {
						var arrBuf = new ArrayBuffer(1);
						var view = new Uint8Array(arrBuf);
						view[0] = getDelay();
						globalSocket.send(arrBuf);
					}
				}

				function showError(message) {
					document.getElementById("errormessage").innerText = message;
					document.getElementById("error").style.display="block";
					document.getElementById("error").style.visibility="visible";
				}
	
				function loadPuzzle() {
					var xmlhttp = new XMLHttpRequest();
					xmlhttp.onreadystatechange = function() {
					  if (this.readyState == 4) {
						if (this.status == 200) {
							try {
								reply = JSON.parse(this.responseText);
								
								if (reply.error != null) {
									showError(reply.error);
									return;
								}
								populateGrid(reply.puzzle);
								return;
						  } catch (e) {
							  showError("Got an unexpected response while fetching JSON: " + e);
							  return;
						  }
						} else {
							  showError("Got an unexpected response while fetching JSON: " + this.status);
							  return;
						}
					  }
					}
					xmlhttp.open("GET", "/puzzle", true);
					xmlhttp.send();
				}

				function solvePuzzle() {
					document.getElementById("enterButton").disabled=true;
					document.getElementById("solveButton").disabled=true;
					document.getElementById("keypad").style.visibility="hidden";
					var websocket = new WebSocket("ws://" + window.location.host + "/solve/" + getGridState());
					websocket.binaryType = 'arraybuffer';

					websocket.onerror = function(evt) {
						console.log(evt);
						showError("Websocket error");
					}
					websocket.onopen = function(evt) {
						globalSocket = websocket;
						sendDelay();
					}
					websocket.onmessage = function (evt) {
						if (evt.type != "message") { return; }
						var data = new Uint8Array(evt.data);

						var len = data.length
						if ((len < 2) || (len%2 == 1)) { return; }
						for (var i=0; i<len; i+=2) {
							var index = data[i];
							var value = data[i+1];
							if ((index < 81) && (value < 10)) {
								setCell(index, value);
							}
						}
					}
					websocket.onclose = function (evt) {
						globalSocket = null;
					}
				}
	
				function initPage() {
					document.getElementById("solveButton").disabled=true;
					document.getElementById("enterButton").disabled=true;
					document.getElementById("error").style.visibility="hidden";
					buildGrid();
					loadPuzzle();
				}
	
				function buildGrid() {
					var b = 1;
					var boxStart;
					var s;
					var grid = document.getElementById("grid");
					for (var b=0; b<9; b++) {
						var box = document.createElement("div");
						box.className = "box";
						box.id = 'box' + b;
						boxStart = parseInt(b/3)*27 + (b%3)*3;
						for (var i=0; i<9; i++) {
							var index = boxStart + parseInt(i/3)*6 + i;
							var cell = document.createElement("div");
							cell.className = "cell dynamic";
							cell.id = 'cell' + index;
							box.appendChild(cell);
						}
						grid.appendChild(box);
					}
				}

				function prepForManualEntry() {
					document.getElementById("enterButton").disabled=true;
					document.getElementById("keypad").style.display="block";
					document.getElementById("keypad").style.visibility="visible";
					resetGrid();
				}

				function manualSet(value) {
					var hl = document.getElementsByClassName("highlighted");
					if (hl.length < 1) { return; }
					var cell = hl[0];
					if (value == 0) {
						cell.className = "cell dynamic highlighted";
						cell.innerText = "";
					} else {
						cell.className = "cell static highlighted";
						cell.innerText = value;
					}
				}

				function resetGrid() {
					for (var i=0; i<81; i++) {
						var cell = document.getElementById("cell" + i);
						cell.className = "cell dynamic";
						cell.innerText = "";
						cell.addEventListener("click", highlightCell);
					}
				}

				function highlightCell(evt) {
					var hl = document.getElementsByClassName("highlighted");
					for (var i=0; i<hl.length; i++) {
						hl[i].className = hl[i].className.substring(0, hl[i].className.lastIndexOf(" "));
					}
					evt.target.className += " highlighted";
				}

				function getGridState() {
					var state = "";
					for (var i=0; i<81; i++) {
						var value = document.getElementById("cell" + i).innerText;
						if (value == "") { value = "0"; }
						state = state + value;
					}
					return state;
				}

				function setCell(index, value) {
					var cell = document.getElementById("cell" + index);
					if (value == "0") { value = ""; }
					cell.innerText = value;
				}
	
				function populateGrid(state) {
					var s;
					for (var i=0; i<81; i++) {
						s = state.charAt(i);
						if (s != 0) { document.getElementById("cell" + i).className = "cell static"; }
						setCell(i, state.charAt(i));
					}
					document.getElementById("solveButton").disabled=false;
					document.getElementById("enterButton").disabled=false;
				}

			</script>
		</body>
	</html>`)
}
