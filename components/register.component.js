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
      border-radius: 0.35rem;
      border: 2px solid #c9ccd4;
    }

    div.input.error input {
      border: 2px solid var(--error);
    }

    button {
      border: 0;
      color: white;
      display: flex;
      gap: 0.65rem;
      background-color: var(--color);
      border-radius: 0.35rem;
      padding: 0.65rem 1rem 0.65rem 1rem;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      transition: 200ms;
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
  };

  constructor() {
    super();
    this.backgroundColor = "#565b66"
    this.username = ""
    this.password = ""
    this.confirmedPassword = ""
    this.email = ""
    this.registrationError = 0
    // 0 = none
    // 1 = password confirmation error
    // 2 = username taken
    this.registrationSuccess = false;
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

  updateConfirmedPassword(e) {
    this.confirmedPassword = e.srcElement.value;
    e.srcElement.parentElement.classList.remove("error")
    e.srcElement.setCustomValidity("")
  }

  updateEmail(e) {
    this.email = e.srcElement.value;
  }

  register() {
    this.username = this.username.trim()
    if (this.username.length == 0) {
      const input = this.shadowRoot.getElementById("username")
      input.setCustomValidity("Please enter a username")
      input.reportValidity()
      input.parentElement.classList.add("error")
    }
    this.password = this.password.trim()
    this.confirmedPassword = this.confirmedPassword.trim()

    if (this.password != this.confirmedPassword) {
      this.registrationError = 1;
      return
    }

    let fail = false;

    if (this.password.length == 0) {
      const input = this.shadowRoot.getElementById("password")
      if (this.username.length != 0) {
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
      return
    }

    const options = {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: `{"username":"${this.username}","password":"${this.password}"}`
    };

    fetch('/api/register', options)
      .then(response => this.handleAPIResponse(response))
      .catch(err => console.error(err));
  }

  handleAPIResponse(response) {
    switch (response.status) {
      case 200:
        console.log(response)
        this.registrationSuccess = true;
        setTimeout(() => {
          window.location.href = "/login"
        }, 4000)
        break;
      case 409:
        this.registrationError = 2;
        break;
      case 500:
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
        return "Username taken."
    }
  }

  render() {
    return html`<div id="root">
      <div class="input${this.registrationError == 2 ? " error" : ''}">
        <label for="username">Username</label>
        <input id="username" type="text" placeholder="Username" value="${this.username}" @input="${this.updateUsername}" />
      </div>
      <div class="input">
          <label for="password">Password</label>
          <input id="password" type="password" placeholder="Password" value="${this.password}" @input="${this.updatePassword}" />
      </div>
      <div class="input${this.registrationError == 1 ? " error" : ''}">
          <label for="confPassword">Confirm Password</label>
          <input id="confPassword" type="password" placeholder="Confirm Password" value="${this.confirmedPassword}" @input="${this.updateConfirmedPassword}" />
      </div>

      <button style="--color: ${this.backgroundColor};" @click="${this.register}">Register</button>

      <p id="error">${this.getRegistrationErrorMessage()}</p>
      ${this.registrationSuccess ? html`<p id="success">User registered, redirecting to login page..</p>` : ""}
      </div>`;
  }
}

customElements.define('register-form', RegisterFormComponent)
