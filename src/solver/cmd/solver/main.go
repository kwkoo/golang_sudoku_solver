package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"solver"
)

const staticHTML = `<!DOCTYPE html>
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
			#main {
				visibility: hidden;
			}
			#loading {
				visibility: visible;
			}
			#error {
				visibility: hidden;
			}
		</style>
	</head>
	<body onload="initPage()">
		<h1>Sudoku Solver</h1>
		<div id="loading">
			<center>
				<img src="data:image/gif;base64,R0lGODlhIAAgAPMAAP///wAAAMbGxoSEhLa2tpqamjY2NlZWVtjY2OTk5Ly8vB4eHgQEBAAAAAAAAAAAACH/C05FVFNDQVBFMi4wAwEAAAAh/hpDcmVhdGVkIHdpdGggYWpheGxvYWQuaW5mbwAh+QQJCgAAACwAAAAAIAAgAAAE5xDISWlhperN52JLhSSdRgwVo1ICQZRUsiwHpTJT4iowNS8vyW2icCF6k8HMMBkCEDskxTBDAZwuAkkqIfxIQyhBQBFvAQSDITM5VDW6XNE4KagNh6Bgwe60smQUB3d4Rz1ZBApnFASDd0hihh12BkE9kjAJVlycXIg7CQIFA6SlnJ87paqbSKiKoqusnbMdmDC2tXQlkUhziYtyWTxIfy6BE8WJt5YJvpJivxNaGmLHT0VnOgSYf0dZXS7APdpB309RnHOG5gDqXGLDaC457D1zZ/V/nmOM82XiHRLYKhKP1oZmADdEAAAh+QQJCgAAACwAAAAAIAAgAAAE6hDISWlZpOrNp1lGNRSdRpDUolIGw5RUYhhHukqFu8DsrEyqnWThGvAmhVlteBvojpTDDBUEIFwMFBRAmBkSgOrBFZogCASwBDEY/CZSg7GSE0gSCjQBMVG023xWBhklAnoEdhQEfyNqMIcKjhRsjEdnezB+A4k8gTwJhFuiW4dokXiloUepBAp5qaKpp6+Ho7aWW54wl7obvEe0kRuoplCGepwSx2jJvqHEmGt6whJpGpfJCHmOoNHKaHx61WiSR92E4lbFoq+B6QDtuetcaBPnW6+O7wDHpIiK9SaVK5GgV543tzjgGcghAgAh+QQJCgAAACwAAAAAIAAgAAAE7hDISSkxpOrN5zFHNWRdhSiVoVLHspRUMoyUakyEe8PTPCATW9A14E0UvuAKMNAZKYUZCiBMuBakSQKG8G2FzUWox2AUtAQFcBKlVQoLgQReZhQlCIJesQXI5B0CBnUMOxMCenoCfTCEWBsJColTMANldx15BGs8B5wlCZ9Po6OJkwmRpnqkqnuSrayqfKmqpLajoiW5HJq7FL1Gr2mMMcKUMIiJgIemy7xZtJsTmsM4xHiKv5KMCXqfyUCJEonXPN2rAOIAmsfB3uPoAK++G+w48edZPK+M6hLJpQg484enXIdQFSS1u6UhksENEQAAIfkECQoAAAAsAAAAACAAIAAABOcQyEmpGKLqzWcZRVUQnZYg1aBSh2GUVEIQ2aQOE+G+cD4ntpWkZQj1JIiZIogDFFyHI0UxQwFugMSOFIPJftfVAEoZLBbcLEFhlQiqGp1Vd140AUklUN3eCA51C1EWMzMCezCBBmkxVIVHBWd3HHl9JQOIJSdSnJ0TDKChCwUJjoWMPaGqDKannasMo6WnM562R5YluZRwur0wpgqZE7NKUm+FNRPIhjBJxKZteWuIBMN4zRMIVIhffcgojwCF117i4nlLnY5ztRLsnOk+aV+oJY7V7m76PdkS4trKcdg0Zc0tTcKkRAAAIfkECQoAAAAsAAAAACAAIAAABO4QyEkpKqjqzScpRaVkXZWQEximw1BSCUEIlDohrft6cpKCk5xid5MNJTaAIkekKGQkWyKHkvhKsR7ARmitkAYDYRIbUQRQjWBwJRzChi9CRlBcY1UN4g0/VNB0AlcvcAYHRyZPdEQFYV8ccwR5HWxEJ02YmRMLnJ1xCYp0Y5idpQuhopmmC2KgojKasUQDk5BNAwwMOh2RtRq5uQuPZKGIJQIGwAwGf6I0JXMpC8C7kXWDBINFMxS4DKMAWVWAGYsAdNqW5uaRxkSKJOZKaU3tPOBZ4DuK2LATgJhkPJMgTwKCdFjyPHEnKxFCDhEAACH5BAkKAAAALAAAAAAgACAAAATzEMhJaVKp6s2nIkolIJ2WkBShpkVRWqqQrhLSEu9MZJKK9y1ZrqYK9WiClmvoUaF8gIQSNeF1Er4MNFn4SRSDARWroAIETg1iVwuHjYB1kYc1mwruwXKC9gmsJXliGxc+XiUCby9ydh1sOSdMkpMTBpaXBzsfhoc5l58Gm5yToAaZhaOUqjkDgCWNHAULCwOLaTmzswadEqggQwgHuQsHIoZCHQMMQgQGubVEcxOPFAcMDAYUA85eWARmfSRQCdcMe0zeP1AAygwLlJtPNAAL19DARdPzBOWSm1brJBi45soRAWQAAkrQIykShQ9wVhHCwCQCACH5BAkKAAAALAAAAAAgACAAAATrEMhJaVKp6s2nIkqFZF2VIBWhUsJaTokqUCoBq+E71SRQeyqUToLA7VxF0JDyIQh/MVVPMt1ECZlfcjZJ9mIKoaTl1MRIl5o4CUKXOwmyrCInCKqcWtvadL2SYhyASyNDJ0uIiRMDjI0Fd30/iI2UA5GSS5UDj2l6NoqgOgN4gksEBgYFf0FDqKgHnyZ9OX8HrgYHdHpcHQULXAS2qKpENRg7eAMLC7kTBaixUYFkKAzWAAnLC7FLVxLWDBLKCwaKTULgEwbLA4hJtOkSBNqITT3xEgfLpBtzE/jiuL04RGEBgwWhShRgQExHBAAh+QQJCgAAACwAAAAAIAAgAAAE7xDISWlSqerNpyJKhWRdlSAVoVLCWk6JKlAqAavhO9UkUHsqlE6CwO1cRdCQ8iEIfzFVTzLdRAmZX3I2SfZiCqGk5dTESJeaOAlClzsJsqwiJwiqnFrb2nS9kmIcgEsjQydLiIlHehhpejaIjzh9eomSjZR+ipslWIRLAgMDOR2DOqKogTB9pCUJBagDBXR6XB0EBkIIsaRsGGMMAxoDBgYHTKJiUYEGDAzHC9EACcUGkIgFzgwZ0QsSBcXHiQvOwgDdEwfFs0sDzt4S6BK4xYjkDOzn0unFeBzOBijIm1Dgmg5YFQwsCMjp1oJ8LyIAACH5BAkKAAAALAAAAAAgACAAAATwEMhJaVKp6s2nIkqFZF2VIBWhUsJaTokqUCoBq+E71SRQeyqUToLA7VxF0JDyIQh/MVVPMt1ECZlfcjZJ9mIKoaTl1MRIl5o4CUKXOwmyrCInCKqcWtvadL2SYhyASyNDJ0uIiUd6GGl6NoiPOH16iZKNlH6KmyWFOggHhEEvAwwMA0N9GBsEC6amhnVcEwavDAazGwIDaH1ipaYLBUTCGgQDA8NdHz0FpqgTBwsLqAbWAAnIA4FWKdMLGdYGEgraigbT0OITBcg5QwPT4xLrROZL6AuQAPUS7bxLpoWidY0JtxLHKhwwMJBTHgPKdEQAACH5BAkKAAAALAAAAAAgACAAAATrEMhJaVKp6s2nIkqFZF2VIBWhUsJaTokqUCoBq+E71SRQeyqUToLA7VxF0JDyIQh/MVVPMt1ECZlfcjZJ9mIKoaTl1MRIl5o4CUKXOwmyrCInCKqcWtvadL2SYhyASyNDJ0uIiUd6GAULDJCRiXo1CpGXDJOUjY+Yip9DhToJA4RBLwMLCwVDfRgbBAaqqoZ1XBMHswsHtxtFaH1iqaoGNgAIxRpbFAgfPQSqpbgGBqUD1wBXeCYp1AYZ19JJOYgH1KwA4UBvQwXUBxPqVD9L3sbp2BNk2xvvFPJd+MFCN6HAAIKgNggY0KtEBAAh+QQJCgAAACwAAAAAIAAgAAAE6BDISWlSqerNpyJKhWRdlSAVoVLCWk6JKlAqAavhO9UkUHsqlE6CwO1cRdCQ8iEIfzFVTzLdRAmZX3I2SfYIDMaAFdTESJeaEDAIMxYFqrOUaNW4E4ObYcCXaiBVEgULe0NJaxxtYksjh2NLkZISgDgJhHthkpU4mW6blRiYmZOlh4JWkDqILwUGBnE6TYEbCgevr0N1gH4At7gHiRpFaLNrrq8HNgAJA70AWxQIH1+vsYMDAzZQPC9VCNkDWUhGkuE5PxJNwiUK4UfLzOlD4WvzAHaoG9nxPi5d+jYUqfAhhykOFwJWiAAAIfkECQoAAAAsAAAAACAAIAAABPAQyElpUqnqzaciSoVkXVUMFaFSwlpOCcMYlErAavhOMnNLNo8KsZsMZItJEIDIFSkLGQoQTNhIsFehRww2CQLKF0tYGKYSg+ygsZIuNqJksKgbfgIGepNo2cIUB3V1B3IvNiBYNQaDSTtfhhx0CwVPI0UJe0+bm4g5VgcGoqOcnjmjqDSdnhgEoamcsZuXO1aWQy8KAwOAuTYYGwi7w5h+Kr0SJ8MFihpNbx+4Erq7BYBuzsdiH1jCAzoSfl0rVirNbRXlBBlLX+BP0XJLAPGzTkAuAOqb0WT5AH7OcdCm5B8TgRwSRKIHQtaLCwg1RAAAOwAAAAAAAAAAAA==" /><span style="font-size: 3em;">&nbsp;Loading...</span>
			</center>
		</div>
		<div id="error">
			<h2>Error</h2>
			<pre id="errormessage"></pre>
		</div>
		<div id="main">
			<input id="show" type="checkbox" onclick="populateGrid()"/>Show Answers
			<div id="grid" class="grid">
			</div>
		</div>
		<script type="text/javascript">
			//var reply = {"question":"009060000040010000050700320890400070000507000002009180400000002005000760060200400", "answer":"729365841346812957158794326893421675614587293572639184481976532235148769967253418"}
			var reply;

			function loadData() {
				var xmlhttp = new XMLHttpRequest();
				xmlhttp.onreadystatechange = function() {
				  if (this.readyState == 4) {
					if (this.status == 200) {
						try {
							var loading = document.getElementById("loading");
							loading.style.visibility="hidden";
							loading.style.display="none";

							reply = JSON.parse(this.responseText);
							
							if (reply.error != null) {
								document.getElementById("errormessage").innerText = reply.error;
								document.getElementById("error").visibility="visible";
								return;
							}
							document.getElementById("error").style.display="none";
							document.getElementById("main").style.visibility="visible";
							populateGrid();
							return;
					  } catch (e) {
						  document.getElementById("errormessage").innerText = "Got an unexpected response while fetching JSON: " + e;
						  document.getElementById("error").style.visibility="visible";
						  document.getElementById("loading").style.visibility="hidden";
						  return;
					  }
					} else {
						  document.getElementById("errormessage").innerText = "Got an unexpected response while fetching JSON: " + this.status;
						  document.getElementById("error").style.visibility="visible";
						  document.getElementById("loading").style.visibility="hidden";
						  return;
					}
				  }
				}
				xmlhttp.open("GET", "/puzzle", true);
				xmlhttp.send();
			}

			function initPage() {
				document.getElementById("main").style.visibility="hidden";
				document.getElementById("error").style.visibility="hidden";
				document.getElementById("loading").style.visibility="visible";
				buildGrid();
				loadData();
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
						cell.className = "cell";
						cell.id = 'cell' + index;
						box.appendChild(cell);
					}
					grid.appendChild(box);
				}
			}

			function populateGrid() {
				var content;
				var s;
				var cell;
				console.log(reply);
				if (document.getElementById("show").checked == true) {
					content = reply.answer;
				} else {
					content = reply.question;
				}
				for (var i=0; i<81; i++) {
					cell = document.getElementById("cell" + i);
					s = content.charAt(i);
					if (s == "0") { s = ""; }
					cell.innerText = s;
				}
			}
		</script>
	</body>
</html>`

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/puzzle" {
		w.Header().Set("Content-Type", "application/json")
		grid, err := solver.LoadPuzzle()
		question := grid.Clone()
		if err != nil {
			outputError(w, err)
			return
		}
		if !grid.Solve() {
			outputError(w, errors.New("could not solve puzzle"))
			return
		}

		fmt.Fprintf(w, "{\"question\":\"%s\", \"answer\":\"%s\"}", question, grid)
		return
	}

	fmt.Fprint(w, staticHTML)
}

func outputError(w http.ResponseWriter, err error) {
	fmt.Fprint(w, `{"error":`)
	output, _ := json.Marshal(err)
	buffer := bytes.NewBufferString(string(output))
	buffer.WriteTo(w)
	fmt.Fprintln(w, "}")
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
