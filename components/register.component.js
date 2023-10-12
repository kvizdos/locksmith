import { LitElement, html, css } from 'https://cdn.jsdelivr.net/gh/lit/dist@2/core/lit-core.min.js';

export class RegisterFormComponent extends LitElement {
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
    `;

  static properties = {
    backgroundColor: { type: String },
    username: { type: String },
    password: { type: String },
    confirmedPassword: { type: String },
    email: { type: String },
    registrationError: { type: Number },
    registrationSuccess: { type: Boolean },
    emailDisabled: { type: Boolean },
    code: { type: this.toString },
    emailAsUsername: { type: Boolean },
    registering: { type: Boolean },
    hasOnboarding: { type: String },
  };

  constructor() {
    super();
    this.backgroundColor = "#565b66"
    this.username = ""
    this.password = ""
    this.confirmedPassword = ""
    this.email = ""
    this.registrationError = 0
    this.emailDisabled = false;
    this.code = ""
    this.emailAsUsername = false
    this.hasOnboarding = ""
    // 0 = none
    // 1 = password confirmation error
    // 2 = username taken
    this.registrationSuccess = false;
    this.registering = false;
  }

  updateUsername(e) {
    this.username = e.srcElement.value;
    e.srcElement.parentElement.classList.remove("error")
    e.srcElement.setCustomValidity("")

    if (this.registrationError == 2) {
      const el = this.shadowRoot.querySelector("#email")
      el.parentElement.classList.remove("error")
      el.parentElement.setCustomValidity("")
    }
  }

  updatePassword(e) {
    this.password = e.srcElement.value;
    e.srcElement.parentElement.classList.remove("error")
    e.srcElement.setCustomValidity("")
  }

  updateConfirmedPassword(e) {
    this.confirmedPassword = e.srcElement.value;
    e.srcElement.parentElement.classList.remove("error")
    e.srcElement.setCustomValidity("")
  }

  updateEmail(e) {
    this.email = e.srcElement.value;
    e.srcElement.parentElement.classList.remove("error")
    e.srcElement.setCustomValidity("")
  }

  isValidUsername(input) {
    const pattern = /^[a-zA-Z0-9]+$/;
    return pattern.test(input);
  }

  isValidEmail(email) {
    const pattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return pattern.test(email);
  }

  register() {
    let fail = false;
    if (this.regitering == true) {
      return
    }
    this.registering = true;

    if (this.emailAsUsername) {
      this.username = this.email
    }

    if (!this.emailAsUsername) {
      this.username = this.username.trim()
      if (this.username.length == 0) {
        const input = this.shadowRoot.getElementById("username")
        input.setCustomValidity("Please enter a username")
        input.reportValidity()
        input.parentElement.classList.add("error")
        fail = true;
      } else if (!this.isValidUsername(this.username)) {
        const input = this.shadowRoot.getElementById("username")
        input.setCustomValidity(this.emailAsUsername ? "Username must be a valid email address." : "Username can only contain alphanumerical characters.")
        input.reportValidity()
        input.parentElement.classList.add("error")
        fail = true;
      }
    }

    this.email = this.email.trim()

    if (this.email.length == 0) {
      const input = this.shadowRoot.getElementById("email")
      if (this.username.length != 0) {
        input.setCustomValidity("Please enter your email")
        input.reportValidity()
      }
      input.parentElement.classList.add("error")
      fail = true;
    } else if (!this.isValidEmail(this.email)) {
      const input = this.shadowRoot.getElementById("email")
      input.setCustomValidity("Invalid email")
      input.reportValidity()
      input.parentElement.classList.add("error")
      fail = true;
    }

    if (this.password != this.confirmedPassword) {
      this.registrationError = 1;
      this.registering = false;
      return
    }

    if (this.password.length == 0) {
      const input = this.shadowRoot.getElementById("password")
      if (this.email.length != 0) {
        input.setCustomValidity("Please enter a password")
        input.reportValidity()
      }
      input.parentElement.classList.add("error")
      fail = true;
    }

    if (this.confirmedPassword.length == 0) {
      const input = this.shadowRoot.getElementById("confPassword")
      if (this.username.length != 0 && this.password != 0) {
        input.setCustomValidity("Please enter a password")
        input.reportValidity()
      }
      input.parentElement.classList.add("error")
      fail = true;
    }

    if (fail) {
      this.registering = false;
      return
    }

    const options = {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: `{"username":"${this.username}","password":"${this.password}","email":"${this.email}","code":"${this.code}"}`
    };

    fetch('/api/register', options)
      .then(response => this.handleAPIResponse(response))
      .catch(err => {
        this.registering = false;
        console.error(err)
      });
  }

  handleAPIResponse(response) {
    switch (response.status) {
      case 200:
        console.log(response)
        this.registrationSuccess = true;
        setTimeout(() => {
          window.location.href = `/login${this.hasOnboarding == "true" ? "?onboard=true" : ""}`
        }, 1000)
        break;
      case 400:
        this.registering = false;
        alert("Email does not match invitation email. Please reload and try again.")
        break;
      case 404:
        this.registering = false;
        alert("Public registration is disabled.")
        break;
      case 409:
        this.registering = false;
        this.registrationError = 2;
        break;
      case 500:
        this.registering = false;
        alert("Something went wrong, please try again later.")
        break;
    }
  }

  getRegistrationErrorMessage() {
    switch (this.registrationError) {
      case 0:
        return ""
      case 1:
        return "Passwords do not match."
      case 2:
        return "Username or email taken."
    }
  }

  render() {
    return html`<div id="root">
    ${!this.emailAsUsername ? html`
      <div class="input${this.registrationError == 2 ? " error" : ''}">
        <label for="username">Username</label>
        <input id="username" type="text" placeholder="Username" autocorrect="off" autocapitalize="off" value="${this.username}" @input="${this.updateUsername}" />
      </div>
      ` : ""}
      <div class="input${this.registrationError == 2 ? " error" : ''}">
        <label for="email">Email</label>
        <input id="email" type="email" placeholder="Email" autocorrect="off" autocapitalize="off" value="${this.email}" @input="${this.updateEmail}" ?disabled=${this.emailDisabled} />
      </div>
      <div class="input">
          <label for="password">Password</label>
          <input id="password" type="password" placeholder="Password" autocorrect="off" autocapitalize="off" value="${this.password}" @input="${this.updatePassword}" />
      </div>
      <div class="input${this.registrationError == 1 ? " error" : ''}">
          <label for="confPassword">Confirm Password</label>
          <input id="confPassword" type="password" placeholder="Confirm Password" autocorrect="off" autocapitalize="off" value="${this.confirmedPassword}" @input="${this.updateConfirmedPassword}" />
      </div>

      <button style="--color: ${this.backgroundColor};" @click="${this.register}">${this.registering ? "Registering..." : "Register"}</button>

      <p id="error">${this.getRegistrationErrorMessage()}</p>
      ${this.registrationSuccess ? html`<p id="success">User registered, redirecting to login page..</p>` : ""}
      </div>`;
  }
}

customElements.define('register-form', RegisterFormComponent)
