import { LitElement, css, html } from 'https://cdn.jsdelivr.net/gh/lit/dist@2/core/lit-core.min.js';

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

    div#actions button:hover {
        transition: 200ms;
        color: var(--hoverText);
        background-color: var(--baseText);
    }

    div#actions button#details {
        --baseText: #a7afba;
        --hoverText: #FFF;
    }

    div#actions button#delete {
        --baseText: #d13c32;
        --hoverText: #FFF;
    }
    `;

  static properties = {
    user: { type: String },
    userObj: { type: Object }
  };

  constructor() {
    super()
    this.userObj = {}
  }

  firstUpdated() {
    if (this.user == undefined) return;

    this.userObj = JSON.parse(this.user)
  }

  render() {
    return html`<div class="user">
              <div>
                  <div id="about">
                      <p id="username">${this.userObj["username"] || "Unknown"}</p>
                      <p id="role">${this.userObj["role"] || "Unknown"}</p>
                  </div>
                  <div id="info">
                      <p>Last Active: ${this.userObj["lastactive"] || "Unknown"}</p>
                      <p>Active Sessions: ${this.userObj["sessions"] || "Unknown"}</p>
                  </div>
              </div>

              <div id="actions">
                  <button id="details">Details</button>
                  <button id="delete">Delete</button>
              </div>
          </div>`;
  }
}

customElements.define('user-tab', UserTabComponent)
