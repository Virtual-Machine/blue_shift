class Canvas {
	constructor(width, height, bgWidth, bgHeight, scale){
		this.width = width
		this.height = height
		this.backgroundWidth = bgWidth
		this.backgroundHeight = bgHeight
		this.dragWindowFlag = false
		this.dragObjectFlag = false
		this.deltaWX = 15
		this.deltaWY = 15
		this.pos1X
		this.pos1Y
		this.pos2X
		this.pos2Y
		this.buffer = 15
		this.scale = scale
		this.grid = new Grid(this.backgroundWidth, this.backgroundHeight, this.scale)
		this.state = new State("John")

		// MARKER Client -> Canvas init
		this.layers = {
			layer1: document.getElementById('BackgroundLayer'),
			layer2: document.getElementById('ItemLayer'),
			layer3: document.getElementById('CharacterLayer'),
			layer4: document.getElementById('ClickLayer')
		}

		this.backgroundLayer = this.layers.layer1.getContext('2d')
		this.itemLayer = this.layers.layer2.getContext('2d')
		this.characterLayer = this.layers.layer3.getContext('2d')
		this.clickLayer = this.layers.layer4.getContext('2d')

		this.turnOffSmoothing()

		this.drawGrid()
	}

	turnOffSmoothing(){
		this.backgroundLayer.imageSmoothingEnabled = false
		this.itemLayer.imageSmoothingEnabled = false
		this.characterLayer.imageSmoothingEnabled = false
		this.clickLayer.imageSmoothingEnabled = false
	}

	clearBackground(){
		this.backgroundLayer.clearRect(0, 0, this.width, this.height)
	}

	clearItemLayer(){
		this.itemLayer.clearRect(0, 0, this.width, this.height)
	}

	clearCharacterLayer(){
		this.characterLayer.clearRect(0, 0, this.width, this.height)
	}

	clearClickLayer(){
		this.clickLayer.clearRect(0, 0, this.width, this.height)
	}

	clearAll(){
		this.clearBackground()
		this.clearItemLayer()
		this.clearCharacterLayer()
		this.clearClickLayer()
	}

	drawPath(layer, x, y, width, height, color){
		this[layer].beginPath()
		this[layer].rect(x + this.deltaWX, y + this.deltaWY, width, height)
		this[layer].lineWidth = 1
		this[layer].strokeStyle = color
		this[layer].stroke()
	}

	drawRect(layer, x, y, width, height){
		switch(layer){
		case 'background':
			this.drawPath('backgroundLayer', x, y, width, height, "#004400")
			break
		case 'item':
			this.drawPath('itemLayer', x, y, width, height, "#443300")
			break
		case 'character':
			this.drawPath('characterLayer', x, y, width, height, "#440000")
			break
		}
	}

	drawGrid(){
		this.clearAll()
		for(var i in this.grid.cells){
			let xPixels = this.grid.cells[i].getXPixels(this.scale)
			let yPixels = this.grid.cells[i].getYPixels(this.scale)
			this.drawRect('background', xPixels, yPixels, this.scale, this.scale)
		}
	}

	setBindings(){
		var self = this
		// MARKER Client -> Canvas click and keyboard handlers
		this.layers.layer4.addEventListener('mousedown', function(e){
			if(self.state.getActive() == window.activeClient){
				var pos = self.getClickedCell(e)
				if (pos.x < 0 || pos.y < 0 || pos.x >= 60 || pos.y >= 40){
					appendMessage("Click was out of bounds")	
				} else {
					window.sConn.send(JSON.stringify({type: "Click",x: pos.x, y: pos.y}))
				}
			} else {
				appendMessage("It is currently " + self.state.getActive() + "'s turn")
			}
		})

		document.addEventListener('keydown', function(e){
			if(e.which == 37){
				self.deltaWX += 45
				if (self.deltaWX > 15){self.deltaWX = 15}
				self.drawGrid()
			}
			if(e.which == 38){
				self.deltaWY += 45
				if (self.deltaWY > 15){self.deltaWY = 15}
				self.drawGrid()
			}
			if(e.which == 39){
				self.deltaWX -= 45
				if (self.deltaWX < -2855){self.deltaWX = -2855}
				self.drawGrid()
			}
			if(e.which == 40){
				self.deltaWY -= 45
				if (self.deltaWY < -2040){self.deltaWY = -2040}
				self.drawGrid()
			}
			var posX = (Math.abs(self.deltaWX - 15) / 2860) * 185
			var posY = (Math.abs(self.deltaWY - 15) / 2055) * 131.667
			$minimapCursor.style.left = (posX + "px")
			$minimapCursor.style.top = (posY + "px")
		})
	}

	getClickedCell(e){
		var rect = this.layers.layer4.getBoundingClientRect()
		var pos = {
			x: e.clientX - rect.left,
			y: e.clientY - rect.top
		}
		pos.x -= this.deltaWX
		pos.y -= this.deltaWY
		pos.x /= 64
		pos.y /= 64
		pos.x = Math.floor(pos.x)
		pos.y = Math.floor(pos.y)
		return pos
	}
}

window.canvas = new Canvas(1000, 535, 3840, 2560, 64)