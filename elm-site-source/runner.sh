while true; do elm make src/Main.elm --output main.js ; ln -sf $(pwd)/main.js /tmp/main.js ; sleep 2; clear; done
