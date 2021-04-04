import EventEmitter from 'events'

export default class Player extends EventEmitter {
    constructor(replayTimeout = 300) {
        super();

        this.replayTimeout = replayTimeout
        this.timeout = null;

        this.context = new AudioContext();
        this.autoPlayEnabled = this.context.state !== "suspended";
        this.playing = false;
        this.source = null;
    }

    playPreviewAudio(url) {
        return fetch(url)
            .then((response) => response.arrayBuffer())
            .then((arrayBuffer) => this.context.decodeAudioData(arrayBuffer))
            .then((audioBuffer) => {
                this.stop()

                this.source = this.context.createBufferSource();
                this.source.connect(this.context.destination);
                this.source.buffer = audioBuffer;
                this.source.start();
            });
    }

    async play(url) {
        if (this.timeout) {
            clearTimeout(this.timeout)
        }

        if (!this.isEnabled()) {
            return
        }

        await this.tryToGetAutoplay();

        await this.playPreviewAudio(url)
        await new Promise((resolve) => {
            this.timeout = setTimeout(() => {
                resolve()
            }, this.replayTimeout)
        });

        await this.context.resume()
        this.playing = true;
        this.autoPlayEnabled = true;
    }

    stop() {
        if (this.timeout) {
            clearTimeout(this.timeout);
        }
        if (this.source) {
            this.source.stop();
            this.source.currentTime = 0;
        }

        return this.context.suspend()
            .then((_) => {
                this.playing = false;
            })
    }

    isPlaying() {
        return this.playing
    }

    async tryToGetAutoplay() {
        try {
            const p = await this.context.resume()
            this.autoPlayEnabled = true
        } catch (e) {
            this.autoPlayEnabled = false
        }
    }

    isEnabled() {
        return this.autoPlayEnabled
    }

    enable() {
        this.autoPlayEnabled = true;
        this.emit('change', true)
    }

    disable() {
        this.autoPlayEnabled = false;
        this.stop()
        this.emit('change', false)
    }
}