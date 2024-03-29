import {
  LitElement,
  css,
  html,
} from "https://cdn.jsdelivr.net/gh/lit/dist@3.1.0/core/lit-core.min.js";

export class UserTabComponent extends LitElement {
  static styles = css`
    p {
      margin: 0;
    }
    div.user {
      --primary: #2f3235;
      --secondary: #575b60;
      display: flex;
      border-bottom: 2px solid var(--primary);
      padding-bottom: 1rem;
      padding-top: 1rem;
      align-items: center;
      justify-content: space-between;
    }

    .user #about {
      display: flex;
      gap: 0.65rem;
      align-items: center;
      margin-bottom: 0.2rem;
    }

    .user #about #username {
      font-size: 1.15rem;
      font-weight: 600;
    }

    .user #about #role {
      font-size: 1rem;
      color: var(--secondary);
    }

    .user #info {
      font-size: 0.85rem;
      display: flex;
      gap: 0.65rem;
      color: var(--secondary);
    }

    .user #info p:not(:last-of-type)::after {
      margin-left: 0.65rem;
    }

    div#actions {
      display: flex;
      gap: 1rem;
    }

    div#actions button {
      font-size: 0.75rem;
      padding: 0.5rem;
      border: 0;
      transition: 200ms;
      border-radius: 0.25rem;
      --color: white;
      --baseText: #000;
      --hoverText: #000;

      color: var(--baseText);
      border: 2px solid var(--baseText);
      background-color: white;
    }

    div#actions button.waiting,
    div#actions button:hover {
      transition: 200ms;
      color: var(--hoverText);
      background-color: var(--baseText);
    }

    div#actions button#details {
      --baseText: #a7afba;
      --hoverText: #fff;
    }

    div#actions button#delete.waiting,
    div#actions button#delete {
      --baseText: #d13c32;
      --hoverText: #fff;
    }

    div#actions button.waiting {
      filter: brightness(0.8);
      transition: 200ms;
    }
  `;

  static properties = {
    user: { type: String },
    userObj: { type: Object },
    isDeleting: { type: Boolean },
  };

  constructor() {
    super();
    this.userObj = {};
    this.isDeleting = false;
  }

  delete() {
    if (this.isDeleting) return;

    this.isDeleting = true;
    const payload = {
      username: this.userObj["username"],
    };
    const options = {
      method: "DELETE",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    };

    fetch("/api/users/delete", options)
      .then((response) => response.status)
      .then((response) => {
        if (response !== 200) {
          this.isDeleting = false;
          alert("Something went wrong, please reload and try again.");
          return;
        }

        setTimeout(() => {
          this.parentNode.removeChild(this);
        }, 1000);
      })
      .catch((err) => console.error(err));
  }

  firstUpdated() {
    if (this.user == undefined) return;

    this.userObj = JSON.parse(this.user);
  }

  capitalize(string) {
    return string.charAt(0).toUpperCase() + string.slice(1);
  }

  render() {
    return html`<div class="user">
      <div>
        <div id="about">
          <p id="username">${this.userObj["username"] || "Unknown"}</p>
          <p id="role">${this.capitalize(this.userObj["role"] || "unknown")}</p>
        </div>
        <div id="info">
          <p>Last Active: ${this.userObj["lastactive"] || "Unknown"}</p>
          <p>Active Sessions: ${this.userObj["sessions"] || "0"}</p>
        </div>
      </div>

      <div id="actions">
        <button id="details">Details</button>
        <button
          id="delete"
          class="${this.isDeleting ? "waiting" : ""}"
          @click=${this.delete}
        >
          ${this.isDeleting ? "Deleting..." : "Delete"}
        </button>
      </div>
    </div>`;
  }
}

customElements.define("user-tab", UserTabComponent);
