import {Component, ElementRef, OnDestroy, ViewChild} from '@angular/core';
import {IonContent, IonHeader, IonLabel, IonTitle, IonToolbar} from '@ionic/angular/standalone';
import {AudioAnalyzer} from "./audio-analyzer";


@Component({
  selector: 'app-speech',
  templateUrl: './speech.page.html',
  styleUrl: './speech.page.scss',
  imports: [IonHeader, IonToolbar, IonTitle, IonContent, IonLabel]
})
export class SpeechPage implements OnDestroy {

  static readonly SAMPLING_RATE = 16_000;

  @ViewChild('canvas', {static: true}) canvas!: ElementRef;
  modelReady = false;
  private ctx!: CanvasRenderingContext2D;
  private snakeSize = 10;
  private width = 300;
  private height = 350;
  private score = 0;
  private snake: { x: number, y: number }[] = [];
  private food!: { x: number, y: number };
  private direction: string | null = null;
  private gameloop: number | null = null;
  private listenLoop: number | null = null;
  private worker!: Worker;
  private recorder: MediaRecorder | null = null;
  private audioContext = new AudioContext({sampleRate: SpeechPage.SAMPLING_RATE});
  private audioAnalyzer = new AudioAnalyzer({sampleRate: this.audioContext.sampleRate});

  constructor() {
    this.worker = new Worker(new URL('../app.worker', import.meta.url));
    this.initListener();
    this.worker.postMessage({type: 'load'});
  }

  initListener(): void {
    const onMessageReceived = (e: MessageEvent) => {
      switch (e.data.status) {
        case 'ready':
          this.modelReady = true;
          this.startListen();
          break;
        case 'complete':
          this.handleTranscriptionResult(e.data.output[0]);
          break;
      }
    };

    this.worker.addEventListener('message', onMessageReceived);
  }

  async startListen(): Promise<void> {
    this.ctx = this.canvas.nativeElement.getContext('2d');
    if (!navigator.mediaDevices.getUserMedia) {
      console.error("getUserMedia not supported on your browser!");
      return
    }
    const stream = await navigator.mediaDevices.getUserMedia({
      audio: {
        sampleRate: SpeechPage.SAMPLING_RATE,
        channelCount: 1,
        echoCancellation: true,
        noiseSuppression: true
      }
    });

    this.recorder = new MediaRecorder(stream);
    this.recorder.onstart = () => {
    }
    this.recorder.ondataavailable = (e) => {
      const blob = new Blob([e.data], {type: this.recorder!.mimeType});
      const fileReader = new FileReader();
      fileReader.onloadend = async () => {
        try {
          const arrayBuffer = fileReader.result;

          const decoded = await this.audioContext.decodeAudioData(arrayBuffer as ArrayBuffer);
          const channelData = decoded.getChannelData(0);

          const analysis = this.audioAnalyzer.analyzeAudioData(channelData);
          if (!analysis.hasSpeech) {
            return;
          }

          this.worker.postMessage({type: 'generate', data: {audio: channelData, language: 'english'}});
        } catch (e) {
          console.error('Error decoding audio data:', e);
        }
      }
      fileReader.readAsArrayBuffer(blob);
    };

    this.recorder.onstop = () => {
    };


    this.listenLoop = window.setInterval(() => {
      this.recorder!.stop();
      this.recorder!.start();
    }, 800);

  }

  ngOnDestroy(): void {
    this.stopGame();
    if (this.listenLoop) {
      clearInterval(this.listenLoop);
      this.listenLoop = null
    }
  }

  handleTranscriptionResult(word: string): void {
    word = word.toLowerCase();
    if (word.includes('go')) {
      this.handleGo();
    } else if (word.includes('stop')) {
      this.handleStop();
    } else if (word.includes('left')) {
      this.handleLeft();
    } else if (word.includes('right')) {
      this.handleRight();
    } else if (word.includes('up')) {
      this.handleUp();
    } else if (word.includes('down')) {
      this.handleDown();
    }
  }

  handleGo(): void {
    this.stopGame();
    this.direction = 'down';
    this.initSnake();
    this.createFood();
    this.gameloop = window.setInterval(this.paint.bind(this), 80);
  }

  handleStop(): void {
    this.stopGame();
  }

  handleLeft(): void {
    if (this.direction !== 'right') {
      this.direction = 'left';
    }
  }

  handleRight(): void {
    if (this.direction !== 'left') {
      this.direction = 'right';
    }
  }

  handleUp(): void {
    if (this.direction !== 'down') {
      this.direction = 'up';
    }
  }

  handleDown(): void {
    if (this.direction !== 'up') {
      this.direction = 'down';
    }
  }

  stopGame(): void {
    if (this.gameloop) {
      clearInterval(this.gameloop);
      this.gameloop = null;
    }
  }

  drawSnake(x: number, y: number): void {
    this.ctx.fillStyle = 'green';
    this.ctx.fillRect(x * this.snakeSize, y * this.snakeSize, this.snakeSize, this.snakeSize);
    this.ctx.strokeStyle = 'darkgreen';
    this.ctx.strokeRect(x * this.snakeSize, y * this.snakeSize, this.snakeSize, this.snakeSize);
  }

  drawFood(x: number, y: number): void {
    this.ctx.fillStyle = 'yellow';
    this.ctx.fillRect(x * this.snakeSize, y * this.snakeSize, this.snakeSize, this.snakeSize);
    this.ctx.fillStyle = 'red';
    this.ctx.fillRect(x * this.snakeSize + 1, y * this.snakeSize + 1, this.snakeSize - 2, this.snakeSize - 2);
  }

  drawScore(): void {
    const scoreText = 'Score: ' + this.score;
    this.ctx.fillStyle = 'blue';
    this.ctx.fillText(scoreText, 145, this.height - 5);
  }

  initSnake(): void {
    const length = 4;
    this.snake = [];
    for (let i = length - 1; i >= 0; i--) {
      this.snake.push({x: i, y: 0});
    }
  }

  paint(): void {
    this.ctx.fillStyle = 'lightgrey';
    this.ctx.fillRect(0, 0, this.width, this.height);
    this.ctx.strokeStyle = 'black';
    this.ctx.strokeRect(0, 0, this.width, this.height);

    let snakeX = this.snake[0].x;
    let snakeY = this.snake[0].y;

    if (this.direction === 'right') {
      snakeX++;
    } else if (this.direction === 'left') {
      snakeX--;
    } else if (this.direction === 'up') {
      snakeY--;
    } else if (this.direction === 'down') {
      snakeY++;
    }

    if (snakeX === -1) {
      snakeX = this.width / this.snakeSize;
    } else if (snakeY === -1) {
      snakeY = this.height / this.snakeSize;
    } else if (snakeY === this.height / this.snakeSize) {
      snakeY = 0;
    } else if (snakeX === this.width / this.snakeSize) {
      snakeX = 0;
    }

    let tail: { x: number, y: number } | undefined;
    if (snakeX === this.food.x && snakeY === this.food.y) {
      tail = {x: snakeX, y: snakeY};
      this.score++;

      this.createFood();
    } else {
      tail = this.snake.pop();
      if (tail) {
        tail.x = snakeX;
        tail.y = snakeY;
      }
    }

    if (tail) {
      this.snake.unshift(tail);
    }

    for (const sn of this.snake) {
      this.drawSnake(sn.x, sn.y);
    }

    this.drawFood(this.food.x, this.food.y);
    this.drawScore();
  }

  createFood(): void {
    this.food = {
      x: Math.floor((Math.random() * 30) + 1),
      y: Math.floor((Math.random() * 30) + 1)
    };

    for (let i = 0; i < this.snake.length; i++) {
      const snakeX = this.snake[i].x;
      const snakeY = this.snake[i].y;

      if (this.food.x === snakeX && this.food.y === snakeY || this.food.y === snakeY && this.food.x === snakeX) {
        this.food.x = Math.floor((Math.random() * 30) + 1);
        this.food.y = Math.floor((Math.random() * 30) + 1);
      }
    }
  }


}
