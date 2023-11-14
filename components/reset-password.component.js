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

    #hibpWarning {
      display: flex;
      flex-direction: column;
      text-align: center;
      gap: 0.5rem;
      align-items: center;
    }

    #hibpWarning p {
      margin: 0;
      padding: 0;
      max-width: 36ch;
    }

    #hibpWarning p#warning {
      font-weight: 800;
      font-size: 1.25rem;
      color: var(--error);
    }

    #hibpWarning strong {
      color: var(--error);
    }

    #hibpWarning button {
      width: 100%;
    }

    hr {
      border: 1px solid rgb(225, 225, 225);
      width: 100%;
    }

    #hibpWarning a {
      color: #476ade;
      text-decoration: underline;
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
    minimumPasswordLength: { type: Number },
    hibp: { type: String },
    passwordSecurityLink: { type: String }
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
    this.minimumPasswordLength = 0
    this.hibp = "loose"
    this.passwordSecurityLink = "#"
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

  firstUpdated() {
    this.stage = window.location.search == "?hibp=true" ? -1 : this.stage
  }

  async resetPassword(hasSeenPwn) {
    if (this.isSending) { return }

    if (this.password1.length == 0 || this.password2.length == 0) {
      alert("Please confirm both password fields are filled.")
      return
    }

    if (this.password1 != this.password2) {
      alert("Passwords do not match. Please try again.")
      return
    }

    if (this.minimumPasswordLength > this.password1.length) {
      alert("Passwords must be at least " + this.minimumPasswordLength + " characters long.")
      return
    }

    this.isSending = true

    const body = {
      password: this.password1,
    }

    if (hasSeenPwn === true) {
      body["pwnok"] = true
    }

    const options = {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body)
    };

    const response = await fetch(`/api/reset-password`, options)
    const status = response.status

    if (status == 409) {
      var json = await response.json()
      if (json["pwned"] !== undefined) {
        if (json["pwned"] === true) {
          this.stage = -2;
          this.isSending = false
          return
        }
      }
    }

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

  renderHIBPResetWarning() {
    return html`<section id="hibpWarning">
  <p>The password you tried to use has been used before and is no longer safe because it was shared publicly online after a theft. <strong>Don't worry, this didn't happen on our website.</strong> To keep your account safe, we ${this.hibp != "strict" ? "really " : ""}need you to pick a new password. Thank you for understanding and helping us keep your information secure.</p>
  <button style="--color: ${this.backgroundColor};" @click=${() => {
        this.password1 = ""
        this.password2 = ""
        this.stage = 2;
      }
      }>Choose a new Password</button>
    ${this.hibp == "loose" ? html`
  <button id="unsafe" style="--color: ${this.backgroundColor};" @click=${() => {
          this.resetPassword(true)
        }}>Or, continue with insecure password.</button>` : ""}

    <hr>

    <a href="${this.passwordSecurityLink}">Learn more about how we protect your account.</a>
  </section>`
  }

  render() {
    return html`<div id="root">
      ${this.stage == -2 ? this.renderHIBPResetWarning() : this.stage == -1 ? html`
        <div id="hibpWarning">
        <p id="warning">Account Security Alert</p>
        <p>We've found that your password has been shared in a public leak <strong>(don't worry, this didn't happen from our site)</strong>. To keep your account safe, please change your password by clicking "Continue" and entering your email address. It's best to choose a new password that you haven't used before. If you need any assistance, we're here to help! Your security is our priority. Thank you for your prompt attention to this matter.</p>

        <button style="--color: ${this.backgroundColor};" @click="${() => this.stage = 0}">Continue</button>

        <hr>

        <a href="${this.passwordSecurityLink}">Learn more about how we protect your account.</a>
        </div>
        ` :
        this.stage == 0 ? html`
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
          <button style="--color: ${this.backgroundColor};" @click="${this.resetPassword}">${this.isSending ? "Resetting.." : "Reset Password"}</button>
          ` : html`<p>Your account has been reset successfully.</p><a style="--color: ${this.backgroundColor};" href="/login">Click here to go to the login page.</a>`
      }
      </div> `;
  }
}

customElements.define('reset-password-form', ResetPasswordFormComponent)
