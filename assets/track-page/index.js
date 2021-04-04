import Player from './player'

export default function () {
    const player = new Player();

    player.on('change', isPlaying => {
        if (isPlaying) {
            document.querySelectorAll('.player-button').forEach((button) => {
                button.classList.add('player-button--playing');
            });
        } else {
            document.querySelectorAll('.player-button').forEach((button) => {
                button.classList.remove('player-button--playing');
            });
        }
    })

    const tracks = document.querySelectorAll('#tracks tbody tr');

    if (!player.isEnabled()) {
        document.querySelectorAll('.player-button').forEach((button) => {
            button.classList.remove('player-button--playing');
        });
    } else {
        document.querySelectorAll('.player-button').forEach((button) => {
            button.classList.add('player-button--playing');
        });
    }

    document.querySelectorAll('.player-button').forEach(function (button) {
        button.addEventListener('click', function () {
            if (player.isPlaying()) {
                player.disable()
            } else {
                const url = button.closest('tr').dataset.previewUrl;

                player.enable()
                player.play(url)
            }
        });
    });

    tracks.forEach(function (track) {
        track.addEventListener('mouseenter', function () {
            if (!track.dataset.previewUrl) {
                return;
            }

            player.play(track.dataset.previewUrl)
        });
        track.addEventListener('mouseleave', function () {
            player.stop()
        });
    });
}