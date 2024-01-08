import {
  LitElement,
  html,
  css,
} from "https://cdn.jsdelivr.net/gh/lit/dist@3.1.0/core/lit-core.min.js";

export class PersonaSwitcherComponent extends LitElement {
  static styles = css`
    :host {
      position: fixed;
      bottom: 0;
      left: 0;
      display: flex;
      justify-content: center;
      width: 100vw;
    }

    p {
      margin: 0;
      padding: 0;
    }
    a {
      color: unset;
      text-decoration: none;
    }
    #root {
      background-color: #fff;
      border: 1px solid #f5f5f5;
      padding: 0.75rem 1.25rem 0.75rem 1.25rem;
      margin-bottom: 1rem;
      border-radius: 0.75rem;
      -webkit-box-shadow: 0px 0px 5px 0px rgba(0, 0, 0, 0.25);
      -moz-box-shadow: 0px 0px 5px 0px rgba(0, 0, 0, 0.25);
      box-shadow: 0px 0px 5px 0px rgba(0, 0, 0, 0.25);
      transition: 200ms;
    }

    #root:hover {
      transition: 200ms;
      color: white;
      background-color: #476ade;
    }
  `;

  static properties = {};

  constructor() {
    super();
  }

  render() {
    return html`<a href="/launchpad"
      ><section id="root">
        <p>Switch persona</p>
      </section></a
    >`;
  }
}

customElements.define("persona-switcher", PersonaSwitcherComponent);

let switcher = document.createElement("persona-switcher");
document.body.appendChild(switcher);
