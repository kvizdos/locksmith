import { LitElement, html, css } from 'https://cdn.jsdelivr.net/gh/lit/dist@2/core/lit-core.min.js';

export class UserListComponent extends LitElement {
  static styles = css`
    section {
      --primary: #2f3235;
      --secondary: #575b60;
    }
    `;

  static properties = {
    users: { type: Array }
  };

  constructor() {
    super()
    console.log("Loading users..")
  }

  render() {
    return html`<section>
          <slot>No user tab present.</slot>
      </section>`;
  }
}

customElements.define('user-list', UserListComponent)
