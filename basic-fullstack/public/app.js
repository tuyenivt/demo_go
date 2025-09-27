window.app = {
  API,
  Router,
  showError: (
    message = "There was an error loading the page",
    goToHome = true
  ) => {
    document.querySelector("#alert-modal").showModal();
    document.querySelector("#alert-modal p").textContents = message;
    if (goToHome) app.Router.go("/");
    return;
  },
  closeError: () => {
    document.getElementById("alert-modal").close();
  },
  search: (event) => {
    event.preventDefault();
    const keywords = document.querySelector("input[type=search]").value;
    if (keywords.length > 1) {
      app.Router.go(`/movies?q=${keywords}`);
    }
  },
  searchOrderChange: (order) => {
    const urlParams = new URLSearchParams(window.location.search);
    const q = urlParams.get("q");
    const genre = urlParams.get("genre") ?? "";
    app.Router.go(`/movies?q=${q}&order=${order}&genre=${genre}`);
  },
  searchFilterChange: (genre) => {
    const urlParams = new URLSearchParams(window.location.search);
    const q = urlParams.get("q");
    const order = urlParams.get("order") ?? "";
    app.Router.go(`/movies?q=${q}&order=${order}&genre=${genre}`);
  },
  register: async (event) => {
    event.preventDefault();
    let errors = [];
    const name = document.getElementById("register-name").value;
    const email = document.getElementById("register-email").value;
    const password = document.getElementById("register-password").value;
    const passwordConfirm = document.getElementById(
      "register-password-confirm"
    ).value;

    if (name.length < 4) errors.push("Enter your complete name");
    if (email.length < 8) errors.push("Enter your complete email");
    if (password.length < 6) errors.push("Enter a password with 6 characters");
    if (password != passwordConfirm) errors.push("Passwords don't match");
    if (errors.length == 0) {
      const response = await API.register(name, email, password);
      if (response.success) {
        app.Store.jwt = response.jwt;
        app.Router.go("/account/");
      } else {
        app.showError(response.message, false);
      }
    } else {
      app.showError(errors.join(". "), false);
    }
  },
  login: async (event) => {
    event.preventDefault();
    let errors = [];
    const email = document.getElementById("login-email").value;
    const password = document.getElementById("login-password").value;

    if (email.length < 8) errors.push("Enter your complete email");
    if (password.length < 6) errors.push("Enter a password with 6 characters");
    if (errors.length == 0) {
      const response = await API.authenticate(email, password);
      if (response.success) {
        app.Store.jwt = response.jwt;
        app.Router.go("/account/");
      } else {
        app.showError(response.message, false);
      }
    } else {
      app.showError(errors.join(". "), false);
    }
  },
  saveToCollection: async (movie_id, collection) => {
    if (app.Store.loggedIn) {
      try {
        const response = await API.saveToCollection(movie_id, collection);
        if (response.success) {
          switch (collection) {
            case "favorite":
              app.Router.go("/account/favorites");
              break;
            case "watchlist":
              app.Router.go("/account/watchlist");
          }
        } else {
          app.showError("We couldn't save the movie.");
        }
      } catch (e) {
        console.log(e);
      }
    } else {
      app.Router.go("/account/");
    }
  },
};

window.addEventListener("DOMContentLoaded", () => {
  app.Router.init();
});
