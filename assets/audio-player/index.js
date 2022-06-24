import EventEmitter from 'events';

/**
 * A wrapper class for the AudioContext API
 */
export default class Index extends EventEmitter {
  /**
   * @constructor
   * @param {number} replayTimeout - wait time in ms until a new song is started
   */
  constructor(replayTimeout = 150) {
    super();

    this.replayTimeout = replayTimeout;
    this.timeout = null;

    this.context = new AudioContext();
    this.enabled = this.context.state !== 'suspended';
    this.playing = false;
    this.source = null;
  }

  /**
   * Fetches an audio file and injects it into the AudioContext
   * @param {string} url
   * @return {Promise<AudioBuffer>}
   * @private
   */
  _playPreviewAudio(url) {
    return fetch(url)
        .then((response) => response.arrayBuffer())
        .then((arrayBuffer) => this.context.decodeAudioData(arrayBuffer))
        .then((audioBuffer) => {
          this.stop();

          this.source = this.context.createBufferSource();
          this.source.connect(this.context.destination);
          this.source.buffer = audioBuffer;
          this.source.start();
        });
  }

  /**
   * Starts playing an audio file
   * @param {string} url
   * @return {Promise<void>}
   */
  async play(url) {
    if (this.timeout) {
      clearTimeout(this.timeout);
    }

    if (!this.isEnabled()) {
      return;
    }

    await new Promise((resolve) => {
      this.timeout = setTimeout(() => {
        resolve();
      }, this.replayTimeout);
    });

    await this._playPreviewAudio(url);
    this.playing = true;
  }

  /**
   * Stops playing audio
   */
  stop() {
    if (this.timeout) {
      clearTimeout(this.timeout);
    }
    if (this.source) {
      this.source.stop();
      this.source.currentTime = 0;
    }

    this.playing = false;
  }

  /**
   * @return {boolean}
   */
  isPlaying() {
    return this.playing;
  }

  /**
   * @return {boolean}
   */
  isEnabled() {
    return this.enabled;
  }

  /**
   * Enables audio playback
   * @return {Promise<void>}
   */
  async enable() {
    await this.context.resume();

    this.enabled = true;
    this.emit('change', true);
  }

  /**
   * Disables audio playback
   * @return {Promise<void>}
   */
  async disable() {
    this.stop();

    await this.context.suspend();

    this.enabled = false;
    this.playing = false;
    this.emit('change', false);
  }
}
