class Sprite {
	constructor(tag, img, clipX, clipY, size){
		this.tag = tag
		this.img = img
		this.clipX = clipX
		this.clipY = clipY
		this.size = size
	}
}

class Coordinate {
	constructor(x, y){
		this.x = x
		this.y = y
	}

	distanceTo(x, y){
		return (Math.abs(x - this.x) + Math.abs(y - this.y))
	}

	getXPixels(scale){
		return this.x * scale
	}

	getYPixels(scale){
		return this.y * scale
	}
}

class Grid {
	constructor(width, height, scale){
		let cellsWide = width / scale
		let cellsHigh = height / scale
		this.scale = scale
		this.cells = {}
		for(var i = 0; i < cellsWide; i++){
			for(var k = 0; k < cellsHigh; k++){
				this.cells[i + "-" + k] = new Coordinate(i, k)
			}
		}
	}
}

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

		let layers = {
			layer1: document.getElementById('BackgroundLayer'),
			layer2: document.getElementById('ItemLayer'),
			layer3: document.getElementById('CharacterLayer'),
			layer4: document.getElementById('ClickLayer')
		}

		layers.layer4.addEventListener('mousedown', function(e){
			if(e.region){
				// TODO Client should check if it is the active player
				let cursorPosition = getMousePos(layers.layer4, e)
				cursorPosition = translatePosition(cursorPosition, this.deltaWX, this.deltaWY)
				window.sConn.send(JSON.stringify({type: "Click",x: cursorPosition.x, y: cursorPosition.y}))
			} else {
				if(!this.dragWindowFlag && !this.dragObjectFlag){
					this.dragWindowFlag = true
					let cursorPosition = getMousePos(layers.layer4, e)
					this.pos1X = cursorPosition.x
					this.pos1Y = cursorPosition.y
				}
			}
		}.bind(this))

		layers.layer4.addEventListener('mousemove', function(e){
			if(this.dragWindowFlag){
				this.clearAll()
				let cursorPosition = getMousePos(layers.layer4, e)
				this.pos2X = cursorPosition.x
				this.pos2Y = cursorPosition.y
				this.deltaWX += this.pos2X - this.pos1X
				this.deltaWY += this.pos2Y - this.pos1Y
				if(this.deltaWX > this.buffer) { 
					this.deltaWX = this.buffer 
				}
				if(this.deltaWY > this.buffer) { 
					this.deltaWY = this.buffer 
				}
				if(this.deltaWX < this.backgroundWidth * -1 + this.width - this.buffer) { 
					this.deltaWX = this.backgroundWidth * -1 + this.width - this.buffer 
				}
				if(this.deltaWY < this.backgroundHeight * -1 + this.height - this.buffer) { 
					this.deltaWY = this.backgroundHeight * -1 + this.height - this.buffer 
				}
				this.pos1X = this.pos2X
				this.pos1Y = this.pos2Y
				this.drawBoxes()
			}
		}.bind(this))

		layers.layer4.addEventListener('mouseleave', function(e){
			this.dragWindowFlag = false
			this.dragObjectFlag = false
		}.bind(this))

		layers.layer4.addEventListener('mouseup', function(e){
			this.dragWindowFlag = false
			this.dragObjectFlag = false
		}.bind(this))

		this.backgroundLayer = layers.layer1.getContext('2d')
		this.itemLayer = layers.layer2.getContext('2d')
		this.characterLayer = layers.layer3.getContext('2d')
		this.clickLayer = layers.layer4.getContext('2d')

		this.turnOffSmoothing()

		this.drawBoxes()
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

	addHitPath(type, x, y, width, height){
		this.clickLayer.beginPath()
		this.clickLayer.rect(x + this.deltaWX, y + this.deltaWY, width, height)
		this.clickLayer.addHitRegion({id:type})
	}

	drawRect(layer, x, y, width, height){
		switch(layer){
		case 'background':
			this.drawPath('backgroundLayer', x, y, width, height, "#004400")
			break
		case 'item':
			this.drawPath('itemLayer', x, y, width, height, "#443300")
			this.addHitPath('item', x, y, width, height)
			break
		case 'character':
			this.drawPath('characterLayer', x, y, width, height, "#440000")
			this.addHitPath('character', x, y, width, height)
			break
		}
	}

	drawSprite(layer, sprite, position){
		let tag = sprite.tag
		let img = sprite.img
		let clipX = sprite.clipX
		let clipY = sprite.clipY
		let size = sprite.size
		let destX = position.x + this.deltaWX
		let destY = position.y + this.deltaWY
		switch(layer){
		case 'background':
			this.backgroundLayer.drawImage(img, clipX, clipY, size, size, destX, destY, this.scale, this.scale)
			break
		case 'item':
			this.itemLayer.drawImage(img, clipX, clipY, size, size, destX, destY, this.scale, this.scale)
			this.clickLayer.beginPath()
			this.clickLayer.rect(destX, destY, this.scale, this.scale)
			this.clickLayer.addHitRegion({id:'item'})
			break
		case 'character':
			this.characterLayer.drawImage(img, clipX, clipY, size, size, destX, destY, this.scale, this.scale)
			this.clickLayer.beginPath()
			this.clickLayer.rect(destX, destY, this.scale, this.scale)
			this.clickLayer.addHitRegion({id:'character'})
			break
		}
	}

	drawBoxes(){
		this.drawGrid()
	}

	drawGrid(){
		for(var i in this.grid.cells){
			let xPixels = this.grid.cells[i].getXPixels(this.scale)
			let yPixels = this.grid.cells[i].getYPixels(this.scale)
			this.drawRect('background', xPixels, yPixels, this.scale, this.scale)
		}
	}
}

function getMousePos(canvas, evt) {
	var rect = canvas.getBoundingClientRect()
	return {
		x: evt.clientX - rect.left,
		y: evt.clientY - rect.top
	}
}

function translatePosition(cursorPosition, deltaWX, deltaWY){
	cursorPosition.x -= deltaWX
	cursorPosition.y -= deltaWY
	cursorPosition.x /= 64
	cursorPosition.y /= 64
	cursorPosition.x = Math.floor(cursorPosition.x)
	cursorPosition.y = Math.floor(cursorPosition.y)
	console.log(cursorPosition)
	return cursorPosition
}

window.canvas = new Canvas(1260, 675, 3840, 2560, 64)