import Player from './player';

/**
 * Updates the player button icons inside the track list.
 * @param {boolean} isPlaying
 */
function updateState(isPlaying) {
  document.querySelectorAll('.player-button').forEach((button) => {
    if (isPlaying) {
      button.classList.add('player-button--playing');
    } else {
      button.classList.remove('player-button--playing');
    }
  });
}

/**
 * Initializes audio playback on table list item hover
 */
export default function() {
  const player = new Player();
  const tracks = document.querySelectorAll('#tracks tbody tr');

  player.on('change', (isPlaying) => {
    updateState(isPlaying);
  });
  updateState(player.isEnabled());

  document.querySelectorAll('.player-button').forEach(function(button) {
    button.addEventListener('click', async () => {
      if (player.isPlaying()) {
        player.disable();
      } else {
        const url = button.closest('tr').dataset.previewUrl;

        await player.enable();
        player.play(url);
      }
    });
  });

  tracks.forEach(function(track) {
    track.addEventListener('mouseenter', function() {
      if (!track.dataset.previewUrl) {
        return;
      }

      player.play(track.dataset.previewUrl);
    });

    track.addEventListener('mouseleave', function() {
      player.stop();
    });
  });
}
