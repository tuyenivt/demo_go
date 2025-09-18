import API from "../services/API.js";
import { MovieItemComponent } from "./MovieItem.js";

export default class HomePage extends HTMLElement {
  async render() {
    const topMovies = await API.getTopMovies();
    renderMoviesInList(topMovies, this.querySelector("#top-10 ul"));

    function renderMoviesInList(movies, ul) {
      ul.innerHTML = "";
      movies.forEach((movie) => {
        const li = document.createElement("li");
        li.appendChild(new MovieItemComponent(movie));
        ul.appendChild(li);
      });
    }
  }

  connectedCallback() {
    const template = document.getElementById("template-home");
    const content = template.content.cloneNode(true);
    this.appendChild(content);

    this.render();
  }
}
customElements.define("home-page", HomePage);
