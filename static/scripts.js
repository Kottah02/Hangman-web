// static/script.js

document.addEventListener("DOMContentLoaded", function() {
    const guessForm = document.getElementById("guessForm");

    guessForm.addEventListener("submit", function(event) {
        const guessInput = guessForm.querySelector("input[name='guess']");
        const guess = guessInput.value.trim();

        // Vérifie que l'entrée contient seulement des lettres
        if (!/^[a-zA-Z]+$/.test(guess)) {
            alert("Veuillez entrer une lettre valide.");
            event.preventDefault();
        }
    });
});
