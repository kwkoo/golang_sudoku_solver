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
					color: gray;
				}
				#error {
					visibility: hidden;
				}
			</style>
		</head>
		<body onload="initPage()">
			<h1>Sudoku Solver</h1>
			<div id="main">
				<div style="padding: 10px;">
					<input id="solveButton" type="button" value="Solve Puzzle" onclick="solvePuzzle()"/>
				</div>
				<div id="grid" class="grid">
				</div>
			</div>
			<div id="error">
				<h2>Error</h2>
				<pre id="errormessage"></pre>
			</div>
			<script type="text/javascript">
				var reply;

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
								populateGrid();
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

var tmp;
				function solvePuzzle() {
					document.getElementById("solveButton").disabled=true;
					console.log("Initiating call to server");
					websocket = new WebSocket("ws://" + window.location.host + "/solve/" + reply.puzzle);
					websocket.binaryType = 'arraybuffer';

					websocket.onerror = function(evt) {
						console.log(evt);
						showError("Websocket error");
					}
					websocket.onopen = function(evt) {
						console.log("open:");
						console.log(evt);
					}
					websocket.onmessage = function (evt) {
						if (evt.type != "message") { return; }
						var data = new Uint8Array(evt.data);

						if (data.length != 2) { return; }
						var index = data[0];
						var value = data[1];
						if ((index > 255) || (value > 9)) { return; }
						setCell(index, value);
					}
				}
	
				function initPage() {
					document.getElementById("solveButton").disabled=true;
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

				function setCell(index, value) {
					var cell = document.getElementById("cell" + index);
					if (value == "0") { value = ""; }
					cell.innerText = value;
				}
	
				function populateGrid() {
					var s;
					for (var i=0; i<81; i++) {
						s = reply.puzzle.charAt(i);
						if (s != 0) { document.getElementById("cell" + i).className = "cell static"; }
						setCell(i, reply.puzzle.charAt(i));
					}
					document.getElementById("solveButton").disabled=false;
				}

			</script>
		</body>
	</html>`)
}
