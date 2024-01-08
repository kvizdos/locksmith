import {
  LitElement,
  html,
  css,
} from "https://cdn.jsdelivr.net/gh/lit/dist@3.1.0/core/lit-core.min.js";

export class UserListComponent extends LitElement {
  static styles = css`
    section {
      --primary: #2f3235;
      --secondary: #575b60;
    }
    slot {
      display: none;
    }
  `;

  static properties = {
    users: { type: Array },
  };

  fetchUsers() {
    const options = {
      method: "GET",
      headers: { "Content-Type": "application/json" },
    };

    fetch("/api/users/list", options)
      .then((response) => response.json())
      .then((response) => {
        this.users = response;
        this.renderUsers();
      })
      .catch((err) => console.error(err));
  }

  constructor() {
    super();
    console.log("Loading users..");
    this.fetchUsers();
  }

  renderUsers() {
    const slotElement = this.shadowRoot.querySelector("slot");
    const assignedNodes = slotElement.assignedNodes();
    const firstNode = assignedNodes[0];

    console.log(firstNode.tagName);

    if (!(firstNode instanceof HTMLElement)) {
      console.log("No user tab present.");
      return;
    }

    const container = this.shadowRoot.querySelector("section#root");

    for (let user of this.users) {
      const node = document.createElement(firstNode.tagName);
      node.setAttribute("user", JSON.stringify(user));
      container.appendChild(node);
    }
  }

  render() {
    return html`<section id="root">
      <slot name="tab">No user tab present.</slot>
    </section>`;
  }
}

customElements.define("user-list", UserListComponent);
