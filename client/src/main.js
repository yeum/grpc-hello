import Phaser from 'phaser';

const SPEED = 140;
let player, cursors;

const config = {
    type: Phaser.AUTO,
    parent: 'game',
    width: 512,
    height: 320,
    pixelArt: true,
    physics: { default: 'arcade' },
    scene: {
        preload,
        create,
        update,
    }
};

function preload() {
    // 경로와 파일명은 본인 에셋에 맞춰 수정
    // 1) 타일 이미지
    this.load.image('tiles', 'assets/Tilemap/tilemap.png');
    // 2) Tiled에서 Export한 맵 JSON
    this.load.tilemapTiledJSON('map', 'assets/maps/sampleMap.json');
    // 3) Load player images
    let fileNum = Math.floor(Math.random() * 15) + 1;
    let fileName = 'assets/characters/character_' + String(fileNum).padStart(2, "0") + '.png';
    this.load.image('player', fileName);

    // 로딩 완료 로그(디버그)
    this.load.once('complete', () => console.log('[load] assets loaded'));
}

function create() {
    // Tiled에서 타일셋 이름(“Name”)이 ‘tilesetName’이라면 동일하게 넣어야 함
    const map = this.make.tilemap({ key: 'map' });

    // 디버그: 타일셋/레이어 이름 확인
    console.log('[map] tilesets =', map.tilesets.map(t => t.name));
    console.log('[map] layers   =', map.layers.map(l => l.name));

    // ── 타일셋 이름 자동 사용 (첫 번째 타일셋)
    const tsName = map.tilesets[0]?.name;
    if (!tsName) {
        console.error('No tileset found in map. Check your .tmj/.tsx');
        return;
    }

    // ★ 중요: 아래 'tilesetNameInTiled'는 Tiled의 타일셋 “Name”과 같아야 함
    // 예) Tiled 타일셋 Name이 "kenney_tiles"면 아래도 "kenney_tiles"
    const tileset = map.addTilesetImage('tileset', 'tiles');
    if (!tileset) {
        console.error(`addTilesetImage failed. tsName="${tsName}" 이미지 키='tiles' 확인 필요`);
        return;
    }

    // 레이어 이름도 Tiled의 레이어 이름과 일치해야 함 (예: 'Ground' 또는 'Tile Layer 1')
    map.createLayer('Dungeon', tileset, 0, 0);
    map.createLayer('Objects', tileset, 0, 0);
    map.createLayer('Carts', tileset, 0, 0);

    // Create Player
    const spawnX = map.widthInPixels / 2;
    const spawnY = map.heightInPixels / 2;
    player = this.physics.add.sprite(spawnX, spawnY, 'player');
    player.setCollideWorldBounds(true);

    // camera is following character
    this.cameras.main.setBounds(0, 0, map.widthInPixels, map.heightInPixels);
    this.cameras.main.startFollow(player, true, 0.1, 0.1);

    // 입력
    cursors = this.input.keyboard.createCursorKeys();
}

function update() {
    if (!player) return;
    const body = player.body;
    body.setVelocity(0);

    if (cursors.left.isDown)    body.setVelocityX(-SPEED);
    if (cursors.right.isDown)    body.setVelocityX(SPEED);
    if (cursors.up.isDown)    body.setVelocityY(-SPEED);
    if (cursors.down.isDown)    body.setVelocityY(SPEED);

    // 대각선 속도 보정
    body.velocity.normalize().scale(SPEED);
}

new Phaser.Game(config);