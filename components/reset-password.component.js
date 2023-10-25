import { LitElement, html, css } from 'https://cdn.jsdelivr.net/gh/lit/dist@2/core/lit-core.min.js';

export class ResetPasswordFormComponent extends LitElement {
  static styles = css`
    div#root {
      --error: #c43f33;
      display: flex;
      flex-direction: column;
      gap: 1rem;
    }

    div.input {
        display: flex;
        flex-direction: column;
        gap: 0.35rem;
    }

    p {
      max-width: 28ch;
      margin: 0;
    }

    div.input label {
      padding-left: 0.5rem;
      color: #878d9c;
    }

    div.input input {
      font-size: 1rem;
      padding: 0.5rem;
      border-radius: 0.25rem;
      border: 1px solid #c9ccd4;
    }

    div.input.error input {
      border: 2px solid var(--error);
    }

    p#error {
      color: var(--error);
      padding: 0;
      margin: 0;
      text-align: center;
    }

    p#success {
      color: #33c466;
      padding: 0;
      margin: 0;
      text-align: center;
    }

    button {
      margin-top: 0.5rem;
      border: 0;
      color: white;
      display: flex;
      gap: 0.65rem;
      background-color: var(--color);
      border-radius: 0.25rem;
      padding: 0.65rem 1rem 0.65rem 1rem;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      transition: 200ms;
      font-size: 1rem;
    }

    a {
      color: var(--color);
      text-decoration: none;
    }
    `;

  static properties = {
    username: { type: String },
    password1: { type: String },
    password2: { type: String },
    backgroundColor: { type: String },
    stage: { type: Number },
    emailAsUsername: { type: Boolean },
    isSending: { type: Boolean },
  };

  constructor() {
    super();
    this.username = ""
    this.backgroundColor = "#565b66"
    this.stage = 0
    this.emailAsUsername = false
    this.isSending = false
    this.password1 = ""
    this.password2 = ""
  }

  async sendReset() {
    if (this.isSending) { return }

    if (this.username.length == 0) {
      alert("Please enter the required field.")
      return
    }

    this.isSending = true

    const options = {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
    };

    const response = await fetch(`/api/reset-password?username=${this.username}`, options)
    const status = response.status

    if (status != 200) {
      alert("Something went wrong. Please reload and try again.")
      return
    }

    this.stage = 1;
  }

  async resetPassword() {
    if (this.isSending) { return }

    console.log(this.password1, this.password2)

    if (this.password1.length == 0 || this.password2.length == 0) {
      alert("Please confirm both password fields are filled.")
      return
    }

    if (this.password1 != this.password2) {
      alert("Passwords do not match. Please try again.")
      return
    }

    this.isSending = true

    const options = {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        "password": this.password1,
      })
    };

    const response = await fetch(`/api/reset-password`, options)
    const status = response.status

    if (status != 200) {
      alert("Something went wrong. Please reload and try again.")
      return
    }

    this.stage = 3
  }

  updateUsername(e) {
    this.username = e.srcElement.value;
  }

  updatePassword1(e) {
    this.password1 = e.srcElement.value;
  }

  updatePassword2(e) {
    this.password2 = e.srcElement.value;
  }

  render() {
    return html`<div id="root">
      ${this.stage == 0 ? html`
      <div class="input">
        <label for="username">${this.emailAsUsername ? "Email" : "Username"}</label>
        <input id="username" type="text" placeholder="${this.emailAsUsername ? "Email" : "Username"}" autocorrect="off" autocapitalize="off" value="${this.username}" @input="${this.updateUsername}" />
      </div>
      <button style="--color: ${this.backgroundColor};" @click="${this.sendReset}">Send Reset Link</button>
      ` : this.stage == 1 ? html`
        <p>If an account was found with the provided information, you will be receiving an email momentarily with instructions on how to continue.</p>
        ` : this.stage == 2 ? html`
          <div class="input">
            <label for="password1">Password</label>
            <input id="password1" type="password" placeholder="Password" autocorrect="off" autocapitalize="off" value="${this.password1}" @input="${this.updatePassword1}" />
          </div>
          <div class="input">
            <label for="password2">Re-Type Password</label>
            <input id="password21" type="password" placeholder="Re-Type Password" autocorrect="off" autocapitalize="off" value="${this.password2}" @input="${this.updatePassword2}" />
          </div>
          <button style="--color: ${this.backgroundColor};" @click="${this.resetPassword}">Reset Password</button>
          ` : html`<p>Your account has been reset successfully.</p><a style="--color: ${this.backgroundColor};" href="/login">Click here to go to the login page.</a>`
      }
      </div> `;
  }
}

customElements.define('reset-password-form', ResetPasswordFormComponent)
