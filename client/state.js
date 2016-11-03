var $warnText = document.getElementById('warnText')
var $loginButton = document.getElementById('loginSubmit')
var $userName = document.getElementById('uName')
var $userPassword = document.getElementById('uPassword')
var $chatInput = document.getElementById('chatInput')
var $history = document.getElementById('history')
var $userList = document.getElementById('users')
var $chatDisplay = document.getElementById('chatDisplay')
var $minimapCursor = document.getElementById('minimapCursor')
var $minimapMap = document.getElementById('minimapMap')
var $adminPanel = document.getElementById('adminPanel')


class State {
	constructor(firstPlayer){
		this.active = firstPlayer
		this.started = false
	}

	getActive(){
		if (this.started){
			return this.active
		} else {
			return "NOBODY"
		}
	}

	setActive(player){
		this.active = player 
	}

	start(){
		this.started = true
	}
}