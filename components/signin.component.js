import { LitElement, html, css } from 'https://cdn.jsdelivr.net/gh/lit/dist@2/core/lit-core.min.js';

export class SignInComponent extends LitElement {
  static styles = css`div#root {
      display: flex;
      gap: 0.65rem;
      background-color: var(--color);
      border-radius: 0.25rem;
      padding: 0.65rem 1rem 0.65rem 1rem;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      transition: 200ms;
    }

    div#root:hover {
      filter: saturate(1.25);
      transition: 200ms;
    }

    p {
      color: white;
      margin: 0;
    }

    img {
      width: 24px;
    }

    .section {
      display: flex;
      align-items: center;
      gap: 0.65rem;
      color: white;
    }

    .section p {
      user-select: none;
    }

    .section:not(.active) {
    display: none;
    }

    p#fallback {
      text-align: center;
      color: var(--color);
      opacity: 0.75;
      margin-top: 0.5rem;
      cursor: pointer;
    }
    `;

  static properties = {
    backgroundColor: { type: String },
    stage: { type: Number },
    signInText: { type: String },
  };

  getDeviceType() {
    var userAgent = navigator.userAgent;

    if (/iPhone/i.test(userAgent)) {
      return 'iPhone';
    } else if (/iPad/i.test(userAgent)) {
      return 'iPad';
    } else if (/iPod/i.test(userAgent)) {
      return 'iPod';
    } else if (/Mac/i.test(userAgent)) {
      return 'Mac';
    } else if (/Android/i.test(userAgent)) {
      return 'Android';
    } else if (/Windows/i.test(userAgent)) {
      return 'Windows';
    }

    return 'Unknown Device';
  }

  constructor() {
    super();
    this.device = this.getDeviceType()
    this.backgroundColor = "#565b66"
    this.stage = 0
    this.signInText = "Sign In"
  }

  setStage(stage) {
    if (stage == 1 && !window.PublicKeyCredential) {
      this.stage = 2;

      this.dispatchEvent(new CustomEvent("next-stage", {
        bubbles: true,
        detail: this.stage
      }));

      return;
    }

    this.stage = stage;

    this.dispatchEvent(new CustomEvent("next-stage", {
      bubbles: true,
      detail: this.stage
    }));
  }

  continue() {
    this.setStage(1)
  }

  fallbackPassword() {
    this.setStage(2)
  }

  render() {
    return html`
      <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:opsz,wght,FILL,GRAD@48,400,0,0" />
      <div>
        <div id="root" style="--color: ${this.backgroundColor};">
          <div class="section${this.stage == 0 ? " active" : ""}" id="continue">
            <p id="passkey">Continue..</p>
          </div>
          <div class="section${this.stage == 1 ? " active" : ""}" id="passkey">
            <p id="passkey">Sign in with ${this.device}</p>
            <img src="https://passkeys.dev/images/fido-passkey-white.svg" id="passkeyicon" />
          </div>
          <div class="section${this.stage == 2 ? " active" : ""}" id="password">
            <p id="passkey">${this.signInText}</p>
          </div>
        </div>
        ${this.stage == 1 ? html`<p id="fallback" @click="${this.fallbackPassword}" style="--color: ${this.backgroundColor};">Continue with Password</p>` : ""}

      </div>`;
  }
}

export class LoginFormComponent extends LitElement {
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

    sign-in {
      margin-top: 0.5rem;
    }
    `;

  static properties = {
    backgroundColor: { type: String },
    stage: { type: Number },
    username: { type: String },
    password: { type: String },
    loginError: { type: Number },
    emailAsUsername: { type: Boolean },
    signingIn: { type: Boolean },
  };

  constructor() {
    super();
    this.backgroundColor = "#565b66"
    this.stage = 0
    this.username = ""
    this.password = ""
    this.loginError = 0
    this.emailAsUsername = false
    this.signingIn = false
    // 0 = none
    // 1 = invalid username
    // 2 = invalid password
  }

  updateUsername(e) {
    this.username = e.srcElement.value;
    e.srcElement.parentElement.classList.remove("error")
    e.srcElement.setCustomValidity("")
  }

  updatePassword(e) {
    this.password = e.srcElement.value;
    e.srcElement.parentElement.classList.remove("error")
    e.srcElement.setCustomValidity("")
  }

  stageChange({ detail: stage }) {
    this.stage = stage
  }


  signin() {
    let fail = false;
    if (this.signingIn == true) {
      return
    }
    this.signingIn = true;
    this.username = this.username.trim()
    if (this.username.length == 0) {
      const input = this.shadowRoot.getElementById("username")
      input.setCustomValidity("Please enter a username")
      input.reportValidity()
      input.parentElement.classList.add("error")
      fail = true;
    }

    if (this.password.length == 0) {
      const input = this.shadowRoot.getElementById("password")
      if (this.username.length != 0) {
        input.setCustomValidity("Please enter a password")
        input.reportValidity()
      }
      input.parentElement.classList.add("error")
      fail = true;
    }

    if (fail) {
      this.signingIn = false;
      return
    }

    const options = {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: `{"username":"${this.username}","password":"${this.password}"}`
    };

    fetch('/api/login', options)
      .then(response => this.handleAPIResponse(response))
      .catch(err => console.error(err));
  }

  handleAPIResponse(response) {
    switch (response.status) {
      case 200:
        window.location.href = "/app"
        break;
      case 404:
        this.loginError = 1;
        this.signingIn = false;
        break;
      case 401:
        this.loginError = 2;
        this.signingIn = false;
        break;
      case 500:
        this.signingIn = false;
        alert("Something went wrong, please try again later.")
        break;
    }
  }

  getLoginErrorMessage() {
    switch (this.loginError) {
      case 0:
        return ""
      default:
        return "Invalid username or password."
    }
  }

  render() {
    return html`<div id="root">
      <div class="input">
        <label for="username">${this.emailAsUsername ? "Email" : "Username"}</label>
        <input id="username" type="text" placeholder="${this.emailAsUsername ? "Email" : "Username"}" autocorrect="off" autocapitalize="off" value="${this.username}" @input="${this.updateUsername}" />
      </div>
      ${this.stage == 0 ? html`
      <div class="input">
          <label for="password">Password</label>
          <input id="password" type="password" placeholder="Password" autocorrect="off" autocapitalize="off" value="${this.password}" @input="${this.updatePassword}" />
        </div>
        ` : ""}
      <sign-in backgroundColor="${this.backgroundColor}" stage="2" @next-stage=${this.stageChange} .signInText=${this.signingIn ? "Signing In" : "Sign In"} @click=${this.signin}></sign-in>

      <p id="error">${this.getLoginErrorMessage()}</p>

      </div>`;
  }
}

customElements.define('sign-in', SignInComponent)
customElements.define('login-form', LoginFormComponent)
