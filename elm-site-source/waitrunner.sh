while inotifywait -e close_write src/Main.elm; do clear; elm make src/Main.elm --output main.js ; ln -sf $(pwd)/main.js /tmp/main.js ; done
