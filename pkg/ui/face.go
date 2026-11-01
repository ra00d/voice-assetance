package ui

import (
	"fmt"

	webview "github.com/webview/webview_go"
)

// Face handles the digital face UI.
type Face struct {
	w webview.WebView
}

// NewFace creates a new Face instance.
func NewFace() *Face {
	w := webview.New(true)
	w.SetTitle("Kind Echo")
	w.SetSize(800, 600, webview.HintNone)

	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<style>
			body {
				background-color: #1a1a1a;
				color: white;
				font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
				margin: 0;
				overflow: hidden;
				height: 100vh;
				width: 100vw;
			}

			/* Navigation */
			.nav-trigger {
				position: absolute;
				top: 20px;
				right: 20px;
				z-index: 100;
				opacity: 0.3;
				transition: opacity 0.3s;
				cursor: pointer;
				font-size: 24px;
			}
			.nav-trigger:hover { opacity: 1; }

			.screen {
				position: absolute;
				top: 0;
				left: 0;
				width: 100%;
				height: 100%;
				display: none; /* Hidden by default */
				flex-direction: column;
				justify-content: center;
				align-items: center;
				background-color: #1a1a1a;
			}
			.screen.active { display: flex; }

			/* Face Screen */
			.face-container {
				position: relative;
				width: 400px;
				height: 400px;
				display: flex;
				flex-direction: column;
				justify-content: center;
				align-items: center;
			}
			.eyes-container {
				display: flex;
				gap: 60px;
				margin-bottom: 40px;
			}
			.eye {
				width: 80px;
				height: 80px;
				background: radial-gradient(circle at 30% 30%, #4facfe, #00f2fe);
				border-radius: 50%;
				box-shadow: 0 0 20px rgba(0, 242, 254, 0.4);
				transition: all 0.3s cubic-bezier(0.4, 0.0, 0.2, 1);
			}
			.mouth {
				width: 60px;
				height: 8px;
				background-color: rgba(255, 255, 255, 0.8);
				border-radius: 10px;
				transition: all 0.3s cubic-bezier(0.4, 0.0, 0.2, 1);
			}
			
			/* Expressions */
			.joy .eye { 
				height: 40px; 
				border-radius: 50% 50% 0 0;
				transform: translateY(10px);
			}
			.joy .mouth { 
				width: 80px;
				height: 40px; 
				background-color: transparent;
				border-bottom: 8px solid rgba(255, 255, 255, 0.9);
				border-radius: 0 0 40px 40px;
				margin-top: -20px;
			}
			
			.sadness .eye { 
				height: 60px; 
				transform: rotate(-20deg); 
			}
			.sadness .eye:nth-child(2) { 
				transform: rotate(20deg); 
			}
			.sadness .mouth { 
				width: 60px;
				height: 30px; 
				background-color: transparent;
				border-top: 6px solid rgba(255, 255, 255, 0.8);
				border-radius: 30px 30px 0 0;
				margin-top: 10px;
			}

			.interest .eye {
				height: 90px;
				width: 90px;
			}
			.interest .mouth {
				width: 30px;
				height: 30px;
				border-radius: 50%;
				border: 4px solid rgba(255, 255, 255, 0.8);
				background-color: transparent;
			}

			/* Animations */
			@keyframes pulse {
				0% { transform: scale(1); opacity: 1; box-shadow: 0 0 20px rgba(0, 242, 254, 0.4); }
				50% { transform: scale(1.1); opacity: 0.8; box-shadow: 0 0 40px rgba(0, 242, 254, 0.6); }
				100% { transform: scale(1); opacity: 1; box-shadow: 0 0 20px rgba(0, 242, 254, 0.4); }
			}

			@keyframes talk {
				0% { height: 8px; border-radius: 10px; }
				50% { height: 20px; border-radius: 20px; }
				100% { height: 8px; border-radius: 10px; }
			}

			.listening .eye {
				animation: pulse 2s infinite ease-in-out;
			}

			.talking .mouth {
				animation: talk 0.3s infinite alternate;
			}

			.status {
				position: absolute;
				bottom: 20px;
				font-size: 18px;
				color: rgba(255, 255, 255, 0.5);
				text-transform: uppercase;
				letter-spacing: 2px;
			}

			/* Settings & Testing Screens */
			h2 { border-bottom: 1px solid #444; padding-bottom: 10px; width: 80%; text-align: center; }
			.control-group {
				margin: 15px 0;
				display: flex;
				flex-direction: column;
				align-items: center;
				width: 100%;
			}
			label { margin-bottom: 5px; color: #aaa; }
			input[type="range"], select { width: 60%; padding: 5px; background: #333; color: white; border: 1px solid #555; border-radius: 4px; }
			button {
				margin: 5px;
				padding: 10px 20px;
				background: #00aaff;
				color: white;
				border: none;
				border-radius: 5px;
				cursor: pointer;
				font-size: 14px;
				transition: background 0.2s;
				width: 200px;
			}
			button:hover { background: #0088cc; }
			button.secondary { background: #444; }
			button.secondary:hover { background: #555; }
			
			.recording .eye {
				background: radial-gradient(circle at 30% 30%, #ff4b4b, #ff0000);
				box-shadow: 0 0 20px rgba(255, 0, 0, 0.4);
				animation: pulse 0.5s infinite ease-in-out;
			}

			.status {
				position: absolute;
				bottom: 20px;
				font-size: 18px;
				color: rgba(255, 255, 255, 0.5);
				text-transform: uppercase;
				letter-spacing: 2px;
			}
			
			.back-btn {
				position: absolute;
				top: 20px;
				left: 20px;
				width: auto;
			}
		</style>
	</head>
	<body>
		<!-- Navigation Trigger -->
		<div class="nav-trigger" onclick="showScreen('settings')">⚙️</div>

		<!-- Face Screen -->
		<div id="screen-face" class="screen active">
			<div id="face" class="face-container">
				<div class="eyes-container">
					<div class="eye left"></div>
					<div class="eye right"></div>
				</div>
				<div class="mouth"></div>
			</div>
			<div id="status" class="status"></div>
		</div>

		<!-- Settings Screen -->
		<div id="screen-settings" class="screen">
			<button class="back-btn secondary" onclick="showScreen('face')">← Back</button>
			<h2>Settings</h2>
			<div class="control-group">
				<label>Microphone Sensitivity</label>
				<input type="range" min="0" max="100" value="50">
			</div>
			<div class="control-group">
				<label>Volume</label>
				<input type="range" min="0" max="100" value="80">
			</div>
			<div class="control-group">
				<label>Voice Model</label>
				<select>
					<option>Piper - Lessac (Medium)</option>
					<option>Piper - Ryan (Low)</option>
				</select>
			</div>
			<button onclick="showScreen('testing')">Go to Testing >></button>
		</div>

		<!-- Testing Screen -->
		<div id="screen-testing" class="screen">
			<button class="back-btn secondary" onclick="showScreen('settings')">← Back</button>
			<h2>Testing & Diagnostics</h2>
			
			<div class="control-group">
				<label>Expressions</label>
				<button onclick="testExpression('joy')">Joy</button>
				<button onclick="testExpression('sadness')">Sadness</button>
				<button onclick="testExpression('interest')">Interest</button>
				<button onclick="testExpression('')">Neutral</button>
			</div>

			<div class="control-group">
				<label>States</label>
				<button onclick="testExpression('listening')">Simulate Listening</button>
				<button onclick="testExpression('recording')">Simulate Recording</button>
				<button onclick="testExpression('talking')">Simulate Talking</button>
				<button onclick="testExpression('processing')">Simulate Thinking</button>
			</div>
		</div>

		<script>
			function showScreen(id) {
				document.querySelectorAll('.screen').forEach(s => s.classList.remove('active'));
				document.getElementById('screen-' + id).classList.add('active');
			}

			function setExpression(expr) {
				// Clear previous states but keep base face class
				const face = document.getElementById('face');
				face.className = 'face-container'; // Reset to base class
				if (expr) {
					face.classList.add(expr);
				}
				
				// Update status text
				const status = document.getElementById('status');
				if (expr === 'listening') {
					status.innerText = 'Listening...';
				} else if (expr === 'recording') {
					status.innerText = 'Recording...';
				} else if (expr === 'talking') {
					status.innerText = 'Speaking...';
				} else if (expr === 'processing') {
					status.innerText = 'Thinking...';
				} else {
					status.innerText = '';
				}
			}

			function testExpression(expr) {
				setExpression(expr);
				// If we wanted to notify Go, we'd call a bound function here
				// window.goCallback(expr);
			}
		</script>
	</body>
	</html>
	`
	w.SetHtml(html)

	return &Face{w: w}
}

// Run starts the UI loop. This blocks.
func (f *Face) Run() {
	f.w.Run()
}

// SetExpression updates the face expression.
// Valid expressions: "", "joy", "sadness"
func (f *Face) SetExpression(expr string) {
	f.w.Dispatch(func() {
		f.w.Eval(fmt.Sprintf("setExpression('%s')", expr))
	})
}

// Close closes the UI.
func (f *Face) Close() {
	f.w.Destroy()
}
